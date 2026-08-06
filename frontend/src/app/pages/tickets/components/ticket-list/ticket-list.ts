import { Component, inject, Signal } from '@angular/core';
import { DatePipe } from '@angular/common';
import { toSignal } from '@angular/core/rxjs-interop';
import { Subject, switchMap, startWith } from 'rxjs';
import { TicketService } from '../../../../core/services/ticket-service';
import { Ticket } from '../../../../core/models/models';
import { Router } from '@angular/router';

@Component({
  selector: 'app-ticket-list',
  standalone: true,
  imports: [DatePipe],
  templateUrl: './ticket-list.html',
  styleUrl: './ticket-list.css',
})
export class TicketList {
  private _ticketService = inject(TicketService);
  private router = inject(Router);

  // Subject pour déclencher le rafraîchissement des données
  private _refresh$ = new Subject<void>();

  // Signal réactif qui contient la liste des tickets
  // Se met à jour au démarrage ET à chaque fois que `refresh()` est appelé
  tickets: Signal<Ticket[]> = toSignal(
    this._refresh$.pipe(
      startWith(void 0), 
      switchMap(() => this._ticketService.getAll())
    ),
    { initialValue: [] }
  );

  /**
   * Méthode publique pour forcer le rafraîchissement
   */
  refresh(): void {
    this._refresh$.next();
  }

  navigateToTicket(id: string) {
    this.router.navigate(["/ticket",id])
  }
}