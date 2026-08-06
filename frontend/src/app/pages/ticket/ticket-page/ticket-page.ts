import { CommonModule, DatePipe } from '@angular/common';
import { Component, effect, inject, input, signal, Signal } from '@angular/core';
import { TicketService } from '../../../core/services/ticket-service';
import { Router } from '@angular/router';
import { toObservable, toSignal } from '@angular/core/rxjs-interop';
import { switchMap } from 'rxjs';
import { Ticket } from '../../../core/models/models';
import { Auth } from '../../../core/auth/auth';
import { FormsModule } from '@angular/forms';

@Component({
  selector: 'app-ticket-page',
  imports: [CommonModule, DatePipe, FormsModule],
  templateUrl: './ticket-page.html',
  styleUrl: './ticket-page.css',
})
export class TicketPage {
  private _ticketService = inject(TicketService);
  private _router = inject(Router);
  private _auth = inject(Auth);
  currentUser: any;
  isEditMode = false;
  isSaving = signal(false);
  draftTitle = signal('');
  draftDescription = signal('');
  draftPriority = signal<'low' | 'mid' | 'high' | 'critical'>('low');
  draftStatus = signal<'open' | 'in_progress' | 'resolved' | 'closed'>('open');

  constructor() {
    this.currentUser = this._auth.currentUser();

    effect(() => {
      const currentTicket = this.ticket();
      if (!currentTicket) {
        return;
      }

      this.draftTitle.set(currentTicket.title);
      this.draftDescription.set(currentTicket.description);
      this.draftPriority.set(currentTicket.priority);
      this.draftStatus.set(currentTicket.status);
    });
  }

  id = input.required<string>();

  ticket: Signal<Ticket | any> = toSignal(
    toObservable(this.id).pipe(
      switchMap((ticketId: string) => this._ticketService.getOne(ticketId)),
    ),
    { initialValue: null },
  );

  ButtonBackTicket() {
    this._router.navigate(['/tickets']);
  }

  onSubmit() {
    const currentTicket = this.ticket();

    if (!currentTicket || !this.isEditMode) {
      return;
    }

    this.isSaving.set(true);
    this._ticketService
      .update(
        currentTicket.id,
        this.draftTitle(),
        this.draftPriority(),
        this.draftDescription(),
        this.draftStatus(),
        currentTicket.assignee?.id,
        currentTicket.categoryId,
      )
      .subscribe({
        next: (updatedTicket) => {
          this.ticket = signal(updatedTicket);
          this.currentUser = this._auth.currentUser();
          this.isEditMode = false;
          this.isSaving.set(false);
        },
        error: () => {
          this.isSaving.set(false);
        },
      });
  }

  onEdit() {
    this.isEditMode = !this.isEditMode;
  }
}
