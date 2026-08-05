import { Injectable, signal } from '@angular/core';

@Injectable({
  providedIn: 'root'
})
export class ModalService {
  isTicketModalOpen = signal(false);

  open() {
    this.isTicketModalOpen.set(true);
  }

  close() {
    this.isTicketModalOpen.set(false);
  }
}