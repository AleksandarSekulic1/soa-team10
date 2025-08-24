import { Injectable, Inject, PLATFORM_ID } from '@angular/core';
import { isPlatformBrowser } from '@angular/common';
import { Router } from '@angular/router';
import { BehaviorSubject } from 'rxjs';
import { User } from '../models/user.model';
import { jwtDecode } from 'jwt-decode'; // Uvozimo jwt-decode

@Injectable({
  providedIn: 'root'
})
export class AuthService {
  public currentUser = new BehaviorSubject<User | null>(null);
  private loggedInStatus = new BehaviorSubject<boolean>(false);
  isLoggedIn$ = this.loggedInStatus.asObservable();

  private isBrowser: boolean;

  constructor(
    private router: Router,
    @Inject(PLATFORM_ID) private platformId: Object
  ) {
    this.isBrowser = isPlatformBrowser(this.platformId);
    if (this.isBrowser) {
      const user = this.getUserFromStorage();
      const tokenExists = this.hasToken();
      if (user && tokenExists) {
        this.currentUser.next(user);
        this.loggedInStatus.next(true);
      }
    }
  }

// U AuthService klasi

  login(response: { token: string }): void {
    if (this.isBrowser) {
      const token = response.token;
      localStorage.setItem('jwt_token', token);
      const user: User = jwtDecode(token);
      localStorage.setItem('user', JSON.stringify(user));
      this.loggedInStatus.next(true);
      this.currentUser.next(user);

      console.log(`[AuthService] Korisnik '${user.username}' se uspešno prijavio. Rola: ${user.role}`); // <-- LOG
    }
  }

  logout(): void {
    if (this.isBrowser) {
      const user = this.currentUser.getValue();
      console.log(`[AuthService] Korisnik '${user?.username}' se odjavio.`); // <-- LOG

      localStorage.removeItem('jwt_token');
      localStorage.removeItem('user');
      this.loggedInStatus.next(false);
      this.currentUser.next(null);
      this.router.navigate(['/login']);
    }
  }
  // Pomoćna metoda koja proverava da li token postoji
  private hasToken(): boolean {
    if (this.isBrowser) {
      return !!localStorage.getItem('jwt_token');
    }
    return false;
  }

  // Metoda koja čita korisnika iz localStorage-a sa dodatnim proverama
  private getUserFromStorage(): User | null {
    if (this.isBrowser) {
      const userString = localStorage.getItem('user');
      if (userString && userString !== 'undefined') {
        try {
          return JSON.parse(userString) as User;
        } catch (e) {
          console.error('Greška pri parsiranju korisnika iz localStorage', e);
          return null;
        }
      }
    }
    return null;
  }

  // Pomoćna metoda za laku proveru da li je korisnik admin
  isAdmin(): boolean {
    return this.currentUser.getValue()?.role === 'administrator';
  }
  isGuide(): boolean {
    return this.getUserRole() === 'vodic';
  }
  // Metoda koja vraća ulogu trenutnog korisnika
  getUserRole(): string | null {
    return this.currentUser.getValue()?.role || null;
  }

    getUsername(): string | null {
    const user = this.currentUser.getValue();
    return user ? user.username : null;
  }
}
