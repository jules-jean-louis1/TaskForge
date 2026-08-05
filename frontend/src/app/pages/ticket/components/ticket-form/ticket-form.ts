import { Component, inject, OnInit } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { ModalService } from '../../../../components/modal/modal.service';
import { TicketService } from '../../../../core/services/ticket-service';
import { CategoryService } from '../../../../core/services/category-service';
import { Category } from '../../../../core/models/models';

@Component({
  selector: 'app-ticket-form',
  standalone: true,
  imports: [FormsModule],
  templateUrl: './ticket-form.html',
  styleUrl: './ticket-form.css',
})
export class TicketForm implements OnInit {
  public categories: Category[] = [];

  private modalService = inject(ModalService);
  private _ticketService = inject(TicketService);
  private _CategoryService = inject(CategoryService);

  ticketData = {
    title: '',
    description: '',
    priority: 'LOW',
    category_id: '',
    status: '',
    assigned_to: '',
    created_by: '',
  };

  ngOnInit() {
    this._CategoryService.get().subscribe({
      next: (res: Category[]) => {
        this.categories = res;
      },
    });
  }

  submitTicket() {
    console.log('Envoi du ticket :', this.ticketData);
    this._ticketService
      .create(
        this.ticketData.title,
        this.ticketData.priority,
        this.ticketData.description,
        this.ticketData.assigned_to,
        parseInt(this.ticketData.category_id),
      )
      .subscribe({
        next: (res) => {
          this.closeModal();
        },
        error: (err) => {
          console.log(err);
        },
      });
  }

  closeModal() {
    this.modalService.close();
  }
}
