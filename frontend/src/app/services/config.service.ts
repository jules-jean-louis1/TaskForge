import { Injectable, inject } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { firstValueFrom } from 'rxjs';

@Injectable({
  providedIn: 'root',
})
export class ConfigService {
  private http = inject(HttpClient);
  private config: any = {};

  loadConfig() {
    return firstValueFrom(this.http.get('/config.json')).then((config) => {
      this.config = config;
    });
  }

  get apiUrl(): string {
    return this.config.apiUrl;
  }
}
