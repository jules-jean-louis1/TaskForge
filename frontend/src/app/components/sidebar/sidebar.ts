import { Component, inject } from '@angular/core';
import { RouterLink, RouterLinkActive } from '@angular/router';
import { ModalService } from '../modal/modal.service';
import { Auth } from '../../core/auth/auth';

@Component({
  selector: 'app-sidebar',
  standalone: true,
  imports: [RouterLink, RouterLinkActive],
  templateUrl: './sidebar.html',
  styleUrl: './sidebar.css',
})
export class SidebarComponent {
  private modalService = inject(ModalService);
  authService = inject(Auth);

  openTicketModal() {
    this.modalService.open();
  }
}
