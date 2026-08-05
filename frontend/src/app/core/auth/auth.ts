import { HttpClient } from '@angular/common/http';
import { computed, Inject, Injectable, PLATFORM_ID, signal } from '@angular/core';
import { isPlatformBrowser } from '@angular/common';
import { ConfigService } from '../../services/config.service';
import { tap } from 'rxjs';
import { jwtDecode } from 'jwt-decode';

interface JWTPayload {
  id: string;
  firstname: string;
  lastname: string;
  email: string;
  role: string;
  exp: number;
  iat: number;
}

@Injectable({
  providedIn: 'root',
})
export class Auth {
  constructor(
    private http: HttpClient,
    private config: ConfigService,
    @Inject(PLATFORM_ID) private platformId: object,
  ) {}
  isLogged = signal<boolean>(false);
  currentUser = signal<JWTPayload | null>(this.getDecodedTokenFromStorage());

  isAuthenticated = computed(() => this.currentUser() !== null);
  userRole = computed(() => this.currentUser()?.role ?? null);
  isAdmin = computed(() => this.currentUser()?.role === 'ADMIN');

  login(email: string, password: string) {
    return this.http
      .post<any>(`${this.config.apiUrl}/auth/login`, { email, password }, { withCredentials: true })
      .pipe(
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
      this.currentUser.set(null);
    }
  }

  getToken(): string | null {
    if (isPlatformBrowser(this.platformId)) {
      return localStorage.getItem('token');
    }
    return null;
  }

  private setToken(token: string) {
    if (isPlatformBrowser(this.platformId)) {
      localStorage.setItem('token', token);
      const decoded = this.decodeToken(token);
      this.currentUser.set(decoded);
    }
  }

  private decodeToken(token: string): JWTPayload | null {
    try {
      return jwtDecode<JWTPayload>(token);
    } catch {
      return null;
    }
  }

  // Utilitaire privé exécuté uniquement au démarrage
  private getDecodedTokenFromStorage(): JWTPayload | null {
    const token = localStorage.getItem('access_token');
    if (!token) return null;
    return this.decodeToken(token);
  }
}
