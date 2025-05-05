import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Router } from '@angular/router';
import { Observable } from 'rxjs';

@Injectable({
  providedIn: 'root',
})
export class AuthService {
  private tokenKey = 'token';
  private apiUrl = 'http://localhost:8080/api';

  constructor(private http: HttpClient, private router: Router) {}

  login(email: string, password: string): Observable<{ token: string }> {
    const payload = { email, password };
    console.log('Sending login payload:', payload); // Логируем данные перед отправкой
    return this.http.post<{ token: string }>(`${this.apiUrl}/auth/login`, payload);
  }

  register(username: string, name: string, surname: string, email: string, password: string, role: string ): Observable<any> {
    return this.http.post(`${this.apiUrl}/auth/register`, { 
      username: username,
      name: name,
      surname: surname,
      email: email,
      password: password,
      role: role
    });
  }

  setToken(token: string): void {
    localStorage.setItem(this.tokenKey, token);
  }

  getToken(): string | null {
    return localStorage.getItem(this.tokenKey);
  }

  logout(): void {
    localStorage.removeItem(this.tokenKey);
    this.router.navigate(['/login']);
  }

  isLoggedIn(): boolean {
    return !!this.getToken();
  }
}