import { Component, inject, signal } from '@angular/core';
import { Auth } from '../../core/auth/auth';
import { RouterLink, Router } from '@angular/router';
import { FormsModule } from '@angular/forms';

@Component({
  selector: 'app-header',
  imports: [RouterLink, FormsModule],
  templateUrl: './header.html',
  styleUrl: './header.css',
})
export class Header {
  authService = inject(Auth);
  private router = inject(Router);
  searchTerm = signal('');

  search() {
    const term = this.searchTerm().trim();
    if (!term) {
      this.router.navigate(['/tickets']);
      return;
    }

    this.router.navigate(['/tickets'], {
      queryParams: { search: term },
    });
  }
}
