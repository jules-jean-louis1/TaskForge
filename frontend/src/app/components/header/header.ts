import { Component } from '@angular/core';
import { Auth } from '../../core/auth/auth';

@Component({
  selector: 'app-header',
  imports: [],
  templateUrl: './header.html',
  styleUrl: './header.css',
})
export class Header {
  constructor(private authService: Auth) {}
  isLogged: boolean = false

  logged() {
    this.isLogged = this.authService.isAuthenticated()
  }
}
