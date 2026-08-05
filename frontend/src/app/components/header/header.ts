import { Component, inject } from '@angular/core';
import { Auth } from '../../core/auth/auth';
import { RouterLink } from '@angular/router';

@Component({
  selector: 'app-header',
  imports: [RouterLink],
  templateUrl: './header.html',
  styleUrl: './header.css',
})
export class Header {
  authService = inject(Auth)
  test() {
    console.log(this.authService.currentUser()?.email)
  }
}
