import { CommonModule, DatePipe } from '@angular/common';
import { Component, inject, OnInit, signal } from '@angular/core';
import { TicketService } from '../../../core/services/ticket-service';
import { UserService } from '../../../core/services/user-service';
import { Ticket, User } from '../../../core/models/models';

@Component({
  selector: 'app-dashboard-page',
  standalone: true,
  imports: [CommonModule, DatePipe],
  templateUrl: './dashboard-page.html',
  styleUrl: './dashboard-page.css',
})
export class DashboardPage implements OnInit {
  private ticketService = inject(TicketService);
  private userService = inject(UserService);

  tickets = signal<Ticket[]>([]);
  users = signal<User[]>([]);
  loading = signal(true);

  ngOnInit() {
    this.loadData();
  }

  private loadData() {
    this.loading.set(true);
    this.ticketService.getAll().subscribe({
      next: (response) => {
        this.tickets.set(Array.isArray(response) ? response : []);
      },
      error: () => {
        this.tickets.set([]);
      },
    });

    this.userService.getAll().subscribe({
      next: (response) => {
        this.users.set(Array.isArray(response) ? response : []);
        this.loading.set(false);
      },
      error: () => {
        this.users.set([]);
        this.loading.set(false);
      },
    });
  }

  get openCount() {
    return this.tickets().filter(
      (ticket) => ticket.status === 'open' || ticket.status === 'in_progress',
    ).length;
  }

  get resolvedCount() {
    return this.tickets().filter(
      (ticket) => ticket.status === 'resolved' || ticket.status === 'closed',
    ).length;
  }

  get averageResolutionHours() {
    const resolved = this.tickets().filter((ticket) => ticket.resolvedAt);
    if (!resolved.length) {
      return 0;
    }

    const totalHours = resolved.reduce((sum, ticket) => {
      const createdAt = new Date(ticket.createdAt).getTime();
      const resolvedAt = new Date(ticket.resolvedAt!).getTime();
      return sum + (resolvedAt - createdAt) / (1000 * 60 * 60);
    }, 0);

    return Math.round((totalHours / resolved.length) * 10) / 10;
  }

  get priorityStats() {
    return [
      {
        label: 'Basse',
        value: this.tickets().filter((ticket) => ticket.priority === 'low').length,
      },
      {
        label: 'Moyenne',
        value: this.tickets().filter((ticket) => ticket.priority === 'mid').length,
      },
      {
        label: 'Haute',
        value: this.tickets().filter((ticket) => ticket.priority === 'high').length,
      },
      {
        label: 'Critique',
        value: this.tickets().filter((ticket) => ticket.priority === 'critical').length,
      },
    ];
  }
}
