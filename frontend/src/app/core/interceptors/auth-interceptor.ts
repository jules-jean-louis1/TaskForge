import { HttpErrorResponse, HttpInterceptorFn } from '@angular/common/http';
import { inject } from '@angular/core';
import { Router } from '@angular/router';
import { BehaviorSubject, catchError, switchMap, throwError } from 'rxjs';
import { filter, take } from 'rxjs/operators';
import { Auth } from '../auth/auth';
import { ConfigService } from '../../services/config.service';

export const authInterceptor: HttpInterceptorFn = (req, next) => {
  const authService = inject(Auth);
  const router = inject(Router);

  if (req.url.includes('/auth/login') || req.url.includes('/auth/register') || req.url.includes('/auth/refresh')) {
    return next(req);
  }

  // Prevent multiple simultaneous refresh calls which would revoke the same token
  // and cause "Refresh token invalide" for concurrent requests.
  // Shared state at module level is acceptable here since interceptor is singleton.
  if ((authInterceptor as any)._isInitialized !== true) {
    (authInterceptor as any)._isRefreshing = false;
    (authInterceptor as any)._refreshSubject = new BehaviorSubject<string | null>(null);
    (authInterceptor as any)._isInitialized = true;
  }
  const isRefreshing = (authInterceptor as any)._isRefreshing as boolean;
  const refreshSubject = (authInterceptor as any)._refreshSubject as BehaviorSubject<string | null>;

  const token = authService.getToken();
  let authReq = req;

  if (token) {
    authReq = req.clone({
      setHeaders: { Authorization: `Bearer ${token}` },
    });
  }

  return next(authReq).pipe(
    catchError((error) => {
      if (!(error instanceof HttpErrorResponse) || error.status !== 401) {
        return throwError(() => error);
      }

      // If a refresh is already in progress, wait for it to complete and retry
      if ((authInterceptor as any)._isRefreshing) {
        return refreshSubject.pipe(
          filter((token) => token != null),
          take(1),
          switchMap((token) => {
            const newReq = req.clone({
              setHeaders: { Authorization: `Bearer ${token}` },
            });
            return next(newReq);
          })
        );
      }

      // No refresh in progress -> start one
      (authInterceptor as any)._isRefreshing = true;
      refreshSubject.next(null);

      return authService.getRefreshToken().pipe(
        switchMap((res: any) => {
          (authInterceptor as any)._isRefreshing = false;
          refreshSubject.next(res.token);
          const newReq = req.clone({
            setHeaders: { Authorization: `Bearer ${res.token}` },
          });
          return next(newReq);
        }),
        catchError((refreshErr) => {
          (authInterceptor as any)._isRefreshing = false;
          refreshSubject.next(null);
          authService.logout();
          router.navigate(['/auth']);
          return throwError(() => refreshErr);
        })
      );
    })
  );
};