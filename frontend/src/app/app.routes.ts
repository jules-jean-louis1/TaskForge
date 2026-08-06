import { Routes } from '@angular/router';
import { HomePage } from './pages/home/home-page/home-page';
import { AuthPage } from './pages/auth/auth-page/auth-page';
import { guestGuard } from './core/guards/guest.guard';
import { TicketsPage } from './pages/tickets/tickets-page';
import { authGuard } from './core/guards/auth-guard';
import { TicketPage } from './pages/ticket/ticket-page/ticket-page';

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
    canActivate: [guestGuard],
  },
  {
    path: 'tickets',
    title: 'TaskForge - Ticket',
    component: TicketsPage,
    canActivate: [authGuard],
  },
  {
    path: 'ticket/:id',
    title: 'TaskForge - Ticket',
    component: TicketPage,
    canActivate: [authGuard],
  },
];
