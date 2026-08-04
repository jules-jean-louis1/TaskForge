import { HttpClient } from '@angular/common/http';
import { Inject, Injectable, PLATFORM_ID } from '@angular/core';
import { isPlatformBrowser } from '@angular/common';
import { ConfigService } from '../../services/config.service';
import { tap } from 'rxjs';

@Injectable({
  providedIn: 'root',
})
export class Auth {
  constructor(
    private http: HttpClient,
    private config: ConfigService,
    @Inject(PLATFORM_ID) private platformId: object,
  ) {}

  login(email: string, password: string) {
    return this.http.post<any>(`${this.config.apiUrl}/auth/login`, { email, password }).pipe(
      tap((res) => {
        this.setToken(res.token);
      }),
    );
  }

  register(firstname: string, lastname: string, email: string, password: string) {
    return this.http.post(`${this.config.apiUrl}/auth/register`, {
      firstname,
      lastname,
      email,
      password,
    });
  }

  getRefreshToken() {
    return this.http
      .post<any>(`${this.config.apiUrl}/auth/refresh`, {}, { withCredentials: true })
      .pipe(
        tap((res) => {
          this.setToken(res.token);
        }),
      );
  }

  logout() {
    if (isPlatformBrowser(this.platformId)) {
      localStorage.removeItem('token');
    }
  }

  getToken(): string | null {
    if (isPlatformBrowser(this.platformId)) {
      return localStorage.getItem('token');
    }
    return null;
  }

  isAuthenticated(): boolean {
    return !!this.getToken();
  }

  private setToken(token: string) {
    if (isPlatformBrowser(this.platformId)) {
      localStorage.setItem('token', token);
    }
  }
}
