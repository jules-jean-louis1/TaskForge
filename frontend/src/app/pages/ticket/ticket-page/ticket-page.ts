import { CommonModule, DatePipe } from '@angular/common';
import { Component, effect, inject, input, signal, Signal } from '@angular/core';
import { TicketService } from '../../../core/services/ticket-service';
import { Router } from '@angular/router';
import { toObservable, toSignal } from '@angular/core/rxjs-interop';
import { BehaviorSubject, catchError, combineLatest, of, switchMap } from 'rxjs';
import { Ticket, User } from '../../../core/models/models';
import { Auth } from '../../../core/auth/auth';
import { FormsModule } from '@angular/forms';
import { UserService } from '../../../core/services/user-service';

@Component({
  selector: 'app-ticket-page',
  standalone: true,
  imports: [CommonModule, DatePipe, FormsModule],
  templateUrl: './ticket-page.html',
  styleUrl: './ticket-page.css',
})
export class TicketPage {
  private _ticketService = inject(TicketService);
  private _router = inject(Router);
  private _auth = inject(Auth);
  private _userService = inject(UserService);

  currentUser = this._auth.currentUser();
  isEditMode = false;
  isSaving = signal(false);

  id = input.required<string>();

  // Subject pour rafraîchir les données après une modification
  private _refresh$ = new BehaviorSubject<void>(undefined);

  // Ticket réactif qui se re-fetch au démarrage et à chaque _refresh$.next()
  ticket: Signal<Ticket | null> = toSignal(
    combineLatest([toObservable(this.id), this._refresh$]).pipe(
      switchMap(([ticketId]) => this._ticketService.getOne(ticketId)),
    ),
    { initialValue: null },
  );

  users: Signal<User[]> = toSignal(this._userService.getAll().pipe(catchError(() => of([]))), {
    initialValue: [],
  });

  draftTitle = signal('');
  draftDescription = signal('');
  draftPriority = signal<'low' | 'mid' | 'high' | 'critical'>('low');
  draftStatus = signal<'open' | 'in_progress' | 'resolved' | 'closed'>('open');
  draftAssigneeId = signal<string | undefined>(undefined);

  constructor() {
    effect(() => {
      const currentTicket = this.ticket();
      if (!currentTicket) return;

      this.draftTitle.set(currentTicket.title);
      this.draftDescription.set(currentTicket.description);
      this.draftPriority.set(currentTicket.priority);
      this.draftStatus.set(currentTicket.status);
      this.draftAssigneeId.set(currentTicket.assignee?.id);
    });
  }

  ButtonBackTicket() {
    this._router.navigate(['/tickets']);
  }

  onSubmit() {
    const currentTicket = this.ticket();
    if (!currentTicket) return;

    this.isSaving.set(true);

    // RÈGLES MÉTIER :
    // - Priorité : Seul l'Admin peut la changer. Pour un Tech, on garde la priorité actuelle.
    // - Assigné : Seul l'Admin peut assigner/réassigner.
    const finalPriority =
      this.currentUser?.role === 'admin' ? this.draftPriority() : currentTicket.priority;

    const finalAssigneeId =
      this.currentUser?.role === 'admin' ? this.draftAssigneeId() : currentTicket.assignee?.id;

    this._ticketService
      .update(
        currentTicket.id,
        this.draftTitle(),
        finalPriority,
        this.draftDescription(),
        this.draftStatus(),
        finalAssigneeId,
        currentTicket.categoryId,
      )
      .subscribe({
        next: () => {
          this._refresh$.next(); // Rafraîchit automatiquement le signal ticket
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

  // Permet de changer rapidement le statut via les boutons d'action
  updateStatusQuickly(newStatus: 'open' | 'in_progress' | 'resolved' | 'closed') {
    this.draftStatus.set(newStatus);
    if (!this.isEditMode) {
      this.onSubmit();
    }
  }
}
