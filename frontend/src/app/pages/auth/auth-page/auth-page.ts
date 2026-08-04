import { Component, signal } from '@angular/core';
import { LoginForm } from '../components/login-form/login-form';
import { RegisterForm } from '../components/register-form/register-form';
import { CommonModule } from '@angular/common';

@Component({
  selector: 'app-auth-page',
  imports: [CommonModule, LoginForm, RegisterForm],
  templateUrl: './auth-page.html',
  styleUrl: './auth-page.css',
  standalone: true,
})
export class AuthPage {
  isLoginMode = signal(true);
  isRegisterSuccessful = signal<boolean>(false);

  toggleMode() {
    this.isLoginMode.update((current) => !current);
  }
  onRegisterSuccess(isSuccess: boolean) {
    this.isRegisterSuccessful.set(isSuccess);
  }
}
