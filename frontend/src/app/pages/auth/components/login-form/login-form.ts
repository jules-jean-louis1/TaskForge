import { Component, inject } from '@angular/core';
import { FormsModule, NgModel } from '@angular/forms';
import { Auth } from '../../../../core/auth/auth';
import { error } from 'console';
import { Router } from '@angular/router';

@Component({
  selector: 'app-login-form',
  imports: [FormsModule],
  templateUrl: './login-form.html',
  styleUrl: './login-form.css',
  standalone: true,
})
export class LoginForm {
  private router = inject(Router);

  email: string = '';
  password: string = '';

  isLoginSuccess: boolean = false;

  constructor(private authService: Auth) {}

  onSubmit() {
    if (!this.email && !this.password) {
      return;
    }
    this.authService.login(this.email, this.password).subscribe({
      next: (res) => {
        if (res.token.length) {
          this.isLoginSuccess = !this.isLoginSuccess;
          setTimeout(() => {
            this.router.navigate(['/']);
          }, 5000);
        }
      },
      error: (error) => {
        console.log(error);
      },
    });
  }
}
