import { Component, inject, signal } from '@angular/core';
import { RouterOutlet } from '@angular/router';
import { SidebarComponent } from './components/sidebar/sidebar';
import { Header } from './components/header/header';
import { ModalComponent } from './components/modal/modal';
import { TicketForm } from './pages/tickets/components/ticket-form/ticket-form';
import { ModalService } from './components/modal/modal.service';

@Component({
  selector: 'app-root',
  standalone: true,
  imports: [RouterOutlet, SidebarComponent, Header, ModalComponent, TicketForm],
  templateUrl: './app.html',
  styleUrl: './app.css'
})
export class App {
  protected readonly title = signal('frontend');
  
  // Correction de l'injection ici !
  public modalService = inject(ModalService);
}