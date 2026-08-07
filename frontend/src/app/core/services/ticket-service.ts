import { inject, Injectable } from '@angular/core';
import { ConfigService } from '../../services/config.service';
import { HttpClient, HttpParams } from '@angular/common/http';
import { Ticket } from '../models/models';
import { Observable } from 'rxjs';

@Injectable({
  providedIn: 'root',
})
export class TicketService {
  config = inject(ConfigService);
  http = inject(HttpClient);

  create(
    title: string,
    priority: string,
    description: string,
    assigned_to?: string,
    category_id?: number,
  ) {
    return this.http.post<any>(`${this.config.apiUrl}/tickets`, {
      title,
      description,
      priority,
      category_id,
      assigned_to,
    });
  }

  update(
    id: string,
    title: string,
    priority: string,
    description: string,
    status: string,
    assigned_to?: string,
    category_id?: number,
  ): Observable<Ticket> {
    return this.http.patch<Ticket>(`${this.config.apiUrl}/tickets/${id}`, {
      title,
      description,
      priority,
      status,
      assigned_to,
      category_id,
    });
  }

  delete(id: string) {
    return this.http.delete<any>(`${this.config.apiUrl}/tickets/${id}`);
  }

  getAll(
    priority?: string,
    status?: string,
    search?: string,
    sort?: string,
    order?: string,
    category_id?: number,
    created_by?: string,
    assigned_to?: string,
  ) {
    let params = new HttpParams();

    if (priority) params = params.set('priority', priority);
    if (status) params = params.set('status', status);
    if (search) params = params.set('search', search);
    if (sort) params = params.set('sort', sort);
    if (order) params = params.set('order', order);
    if (category_id) params = params.set('category_id', category_id.toString());
    if (created_by) params = params.set('created_by', created_by);
    if (assigned_to) params = params.set('assigned_to', assigned_to);

    return this.http.get<any>(`${this.config.apiUrl}/tickets`, { params });
  }

  getOne(id: string): Observable<Ticket> {
    return this.http.get<Ticket>(`${this.config.apiUrl}/tickets/${id}`);
  }
}
