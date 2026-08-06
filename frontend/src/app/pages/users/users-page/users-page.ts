import { CommonModule } from '@angular/common';
import { Component, inject, signal } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { UserService, CreateUserPayload, UpdateUserPayload } from '../../../core/services/user-service';
import { User } from '../../../core/models/models';

@Component({
  selector: 'app-users-page',
  standalone: true,
  imports: [CommonModule, FormsModule],
  templateUrl: './users-page.html',
  styleUrl: './users-page.css',
})
export class UsersPage {
  private userService = inject(UserService);
  users = signal<User[]>([]);
  loading = signal(true);
  form = signal({ firstname: '', lastname: '', email: '', password: '', role: 'standard' });
  editingId = signal<string | null>(null);

  ngOnInit() {
    this.loadUsers();
  }

  loadUsers() {
    this.loading.set(true);
    this.userService.getAll().subscribe({
      next: (response) => {
        this.users.set(Array.isArray(response) ? response : []);
        this.loading.set(false);
      },
      error: () => {
        this.users.set([]);
        this.loading.set(false);
      },
    });
  }

  submit() {
    const payload = this.form();
    const request = this.editingId()
      ? this.userService.update(this.editingId()!, payload as UpdateUserPayload)
      : this.userService.create(payload as CreateUserPayload);

    request.subscribe({
      next: () => {
        this.resetForm();
        this.loadUsers();
      },
    });
  }

  edit(user: User) {
    this.editingId.set(user.id);
    this.form.set({
      firstname: user.firstname,
      lastname: user.lastname,
      email: user.email,
      password: '',
      role: user.role,
    });
  }

  deleteUser(id: string) {
    this.userService.delete(id).subscribe({ next: () => this.loadUsers() });
  }

  resetForm() {
    this.editingId.set(null);
    this.form.set({ firstname: '', lastname: '', email: '', password: '', role: 'standard' });
  }
}
