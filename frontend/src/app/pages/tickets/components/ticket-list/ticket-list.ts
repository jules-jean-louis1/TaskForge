import { Component, inject, Signal } from '@angular/core';
import { DatePipe } from '@angular/common';
import { toSignal } from '@angular/core/rxjs-interop';
import { Subject, switchMap, startWith } from 'rxjs';
import { TicketService } from '../../../../core/services/ticket-service';
import { Ticket } from '../../../../core/models/models';
import { ActivatedRoute, Router } from '@angular/router';

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
  private route = inject(ActivatedRoute);

  private _refresh$ = new Subject<void>();

  tickets: Signal<Ticket[]> = toSignal(
    this._refresh$.pipe(
      startWith(void 0),
      switchMap(() => {
        const search = this.route.snapshot.queryParamMap.get('search') ?? '';
        return this._ticketService.getAll(undefined, undefined, search || undefined);
      }),
    ),
    { initialValue: [] },
  );

  /**
   * Méthode publique pour forcer le rafraîchissement
   */
  refresh(): void {
    this._refresh$.next();
  }

  navigateToTicket(id: string) {
    this.router.navigate(['/ticket', id]);
  }
}
