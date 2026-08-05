import { Component, Input, Output, EventEmitter } from '@angular/core';

@Component({
  selector: 'app-modal',
  standalone: true,
  imports: [],
  templateUrl: './modal.html',
  styleUrl: './modal.css',
})
export class ModalComponent {
  // Propriété pour contrôler l'ouverture/fermeture
  @Input() isOpen = false;

  // Titre optionnel de la modal
  @Input() title = '';

  // Événement déclenché quand l'utilisateur demande la fermeture
  @Output() close = new EventEmitter<void>();

  onClose() {
    this.close.emit();
  }

}
