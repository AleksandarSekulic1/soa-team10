import { Injectable, Inject, PLATFORM_ID } from '@angular/core';
import { isPlatformBrowser } from '@angular/common';
import { Router } from '@angular/router';
import { BehaviorSubject } from 'rxjs';
import { jwtDecode } from 'jwt-decode';

export interface User {
  id: string;
  username: string;
  role: string;
}

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

  login(response: { token: string }): void {
    if (this.isBrowser) {
      const token = response.token;
      localStorage.setItem('jwt_token', token);
      const user: User = jwtDecode(token);
      localStorage.setItem('user', JSON.stringify(user));
      this.loggedInStatus.next(true);
      this.currentUser.next(user);
      console.log(`[AuthService] Korisnik '${user.username}' se uspešno prijavio. Rola: ${user.role}`);
    }
  }

  logout(): void {
    if (this.isBrowser) {
      const user = this.currentUser.getValue();
      console.log(`[AuthService] Korisnik '${user?.username}' se odjavio.`);
      localStorage.removeItem('jwt_token');
      localStorage.removeItem('user');
      this.loggedInStatus.next(false);
      this.currentUser.next(null);
      this.router.navigate(['/login']);
    }
  }

  private hasToken(): boolean {
    if (this.isBrowser) {
      return !!localStorage.getItem('jwt_token');
    }
    return false;
  }

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

  getUserRole(): string | null {
    return this.currentUser.getValue()?.role || null;
  }
  getUsername(): string | null {
    return this.currentUser.getValue()?.username || null;
  }

  isAdmin(): boolean {
    return this.getUserRole() === 'administrator';
  }

  isGuide(): boolean {
    return this.getUserRole() === 'vodic';
  }

  isTourist(): boolean {
    const role = this.getUserRole();
    return role === 'turista';
  }
}
