import { Component, inject, OnInit } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { ModalService } from '../../../../components/modal/modal.service';
import { TicketService } from '../../../../core/services/ticket-service';
import { CategoryService } from '../../../../core/services/category-service';
import { UserService } from '../../../../core/services/user-service';
import { Category, User } from '../../../../core/models/models';
import { Auth } from '../../../../core/auth/auth';

@Component({
  selector: 'app-ticket-form',
  standalone: true,
  imports: [FormsModule],
  templateUrl: './ticket-form.html',
  styleUrl: './ticket-form.css',
})
export class TicketForm implements OnInit {
  public categories: Category[] = [];
  public users: User[] = [];

  private modalService = inject(ModalService);
  private _ticketService = inject(TicketService);
  private _CategoryService = inject(CategoryService);
  private _userService = inject(UserService);
  authService = inject(Auth);

  ticketData = {
    title: '',
    description: '',
    priority: 'low',
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

    this._userService.getAll().subscribe({
      next: (res: User[]) => {
        this.users = res.filter((user) => user.role === 'tech' || user.role === 'admin');
      },
    });
  }

  submitTicket() {
    const canAssign =
      this.authService.currentUser()?.role === 'admin' ||
      this.authService.currentUser()?.role === 'tech';
    const assignedTo = canAssign ? this.ticketData.assigned_to : '';

    this._ticketService
      .create(
        this.ticketData.title,
        this.ticketData.priority,
        this.ticketData.description,
        assignedTo,
        parseInt(this.ticketData.category_id),
      )
      .subscribe({
        next: () => {
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
