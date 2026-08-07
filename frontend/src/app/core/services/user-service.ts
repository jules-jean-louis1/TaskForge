import { Injectable, inject } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { ConfigService } from '../../services/config.service';
import { User } from '../models/models';
import { Observable } from 'rxjs';

export interface CreateUserPayload {
  firstname: string;
  lastname: string;
  email: string;
  password: string;
  role: string;
}

export interface UpdateUserPayload {
  firstname?: string;
  lastname?: string;
  email?: string;
  password?: string;
  role?: string;
}

@Injectable({
  providedIn: 'root',
})
export class UserService {
  private http = inject(HttpClient);
  private config = inject(ConfigService);

  getAll(): Observable<User[]> {
    return this.http.get<User[]>(`${this.config.apiUrl}/users`);
  }

  create(payload: CreateUserPayload) {
    return this.http.post(`${this.config.apiUrl}/users`, payload);
  }

  update(id: string, payload: UpdateUserPayload) {
    return this.http.patch(`${this.config.apiUrl}/users/${id}`, payload);
  }

  delete(id: string) {
    return this.http.delete(`${this.config.apiUrl}/users/${id}`);
  }
}
