import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';
import { User } from '../models/user.model';

@Injectable({
  providedIn: 'root'
})
export class UserService {
  // ISPRAVKA: Vraćamo kosu crtu (/) na kraj, da odgovara vašem backendu.
  private apiUrl = '/api/stakeholders';

  constructor(private http: HttpClient) { }

  register(user: any): Observable<any> {
    // Putanja postaje /api/stakeholders/register
    return this.http.post<any>(`${this.apiUrl}register`, user);
  }

  login(credentials: any): Observable<any> {
    // Putanja postaje /api/stakeholders/login
    return this.http.post<any>(`${this.apiUrl}login`, credentials);
  }

  getAllUsers(): Observable<User[]> {
    // Pozivamo /api/stakeholders/
    return this.http.get<User[]>(this.apiUrl);
  }

  getProfile(): Observable<User> {
    // Putanja postaje /api/stakeholders/profile
    return this.http.get<User>(`${this.apiUrl}profile`);
  }

  updateProfile(user: User): Observable<User> {
    return this.http.put<User>(`${this.apiUrl}profile`, user);
  }

  blockUser(username: string): Observable<any> {
    // Putanja postaje npr. /api/stakeholders/pera/block
    return this.http.put(`${this.apiUrl}${username}/block`, {});
  }

  unblockUser(username: string): Observable<any> {
    // Putanja postaje npr. /api/stakeholders/pera/unblock
    return this.http.put(`${this.apiUrl}${username}/unblock`, {});
  }
}
