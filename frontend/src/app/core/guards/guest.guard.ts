import { inject } from '@angular/core';
import { CanActivateFn, Router } from '@angular/router';
import { Auth } from '../auth/auth';

export const guestGuard: CanActivateFn = (route, state) => {
  const authService = inject(Auth);
  const router = inject(Router);
  console.log('auth', authService.isAuthenticated());
  if (!authService.isAuthenticated()) {
    return true;
  }

  return router.parseUrl('/');
};
