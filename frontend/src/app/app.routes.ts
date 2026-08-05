import { Routes } from '@angular/router';
import { HomePage } from './pages/home/home-page/home-page';
import { AuthPage } from './pages/auth/auth-page/auth-page';
import { guestGuard } from './core/guards/guest.guard';
import { TicketPage } from './pages/ticket/ticket-page';
import { authGuard } from './core/guards/auth-guard';

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
    },
    {
        path: 'tickets',
        title: 'TaskForge - Ticket',
        component: TicketPage,
        canActivate: [authGuard]
    }
];
