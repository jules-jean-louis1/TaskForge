import { Routes } from '@angular/router';
import { HomePage } from './pages/home/home-page/home-page';
import { AuthPage } from './pages/auth/auth-page/auth-page';
import { guestGuard } from './core/guards/guest.guard';

export const routes: Routes = [
    {
        path: '',
        title: 'TaskForge',
        component: HomePage,
    },
    {
        path: 'auth',
        title: 'TaskForge - Auth',
        component: AuthPage,
        canActivate: [guestGuard]
    }
];
