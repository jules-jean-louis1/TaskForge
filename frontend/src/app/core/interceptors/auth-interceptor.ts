import { HttpErrorResponse, HttpInterceptorFn } from '@angular/common/http';
import { inject } from '@angular/core';
import { Router } from '@angular/router';
import { catchError, switchMap, throwError } from 'rxjs';
import { Auth } from '../auth/auth';
import { ConfigService } from '../../services/config.service';

export const authInterceptor: HttpInterceptorFn = (req, next) => {
  const authService = inject(Auth);
  const router = inject(Router);

  if (req.url.includes('/auth/login') || req.url.includes('/auth/register') || req.url.includes('/auth/refresh')) {
    return next(req);
  }

  const token = authService.getToken();
  let authReq = req;

  if (token) {
    authReq = req.clone({
      setHeaders: { Authorization: `Bearer ${token}` },
    });
  }

  return next(authReq).pipe(
    catchError((error) => {
      // Si on reçoit une 401
      if (error instanceof HttpErrorResponse && error.status === 401) {
        
        // 2. On appelle getRefreshToken() qui renvoie un Observable (SANS réinvoquer inject)
        return authService.getRefreshToken().pipe(
          switchMap((res: any) => {
            // On rejoue la requête initiale avec le nouveau token
            const newReq = req.clone({
              setHeaders: { Authorization: `Bearer ${res.token}` },
            });
            return next(newReq);
          }),
          catchError((refreshErr) => {
            // Si le refresh échoue (cookie expiré, etc.)
            authService.logout();
            router.navigate(['/auth']);
            return throwError(() => refreshErr);
          })
        );
      }

      return throwError(() => error);
    })
  );
};