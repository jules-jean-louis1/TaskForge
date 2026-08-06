import { Component, output } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { Auth } from '../../../../core/auth/auth';

@Component({
  selector: 'app-register-form',
  imports: [FormsModule],
  templateUrl: './register-form.html',
  styleUrl: './register-form.css',
  standalone: true,
})
export class RegisterForm {
  firstname: string = '';
  lastname: string = '';
  email: string = '';
  password: string = '';

  registerSuccess = output<boolean>();

  constructor(private authService: Auth) {}

  onSubmit() {
    if (!this.email || !this.email.length) return;
    this.authService.register(this.firstname, this.lastname, this.email, this.password).subscribe({
      next: (res) => {
        if (res.status === 201) {
          this.registerSuccess.emit(true);
        }
      },
      error: (err) => {
        console.log(err);
      },
    });
  }
}
