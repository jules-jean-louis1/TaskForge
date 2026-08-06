import { HttpClient } from '@angular/common/http';
import { inject, Injectable } from '@angular/core';
import { ConfigService } from '../../services/config.service';
import { Category } from '../models/models';

@Injectable({
  providedIn: 'root',
})
export class CategoryService {
  config = inject(ConfigService);
  http = inject(HttpClient);

  get() {
    return this.http.get<Category[]>(`${this.config.apiUrl}/categories`);
  }

  getOne(id: number) {
    return this.http.get(`${this.config.apiUrl}/categories/${id}`);
  }

  create(name: string) {
    return this.http.post(`${this.config.apiUrl}/categories`, { name });
  }

  delete(id: number) {
    return this.http.delete(`${this.config.apiUrl}/categories/${id}`);
  }
}
