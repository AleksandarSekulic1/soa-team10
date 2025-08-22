// src/app/services/auth.service.ts
import { Injectable } from '@angular/core';
import { Router } from '@angular/router';
import { jwtDecode } from 'jwt-decode';

@Injectable({
  providedIn: 'root'
})
export class AuthService {

  constructor(private router: Router) { }

  // Metoda za čuvanje tokena nakon logovanja
  login(token: string): void {
    localStorage.setItem('jwt_token', token);
  }

  // Metoda za proveru da li je korisnik ulogovan
  isLoggedIn(): boolean {
    return !!localStorage.getItem('jwt_token');
  }

  // Metoda za odjavu
  logout(): void {
    localStorage.removeItem('jwt_token');
    this.router.navigate(['/login']);
  }

  // Metoda koja čita ulogu korisnika iz tokena
  getUserRole(): string | null {
    const token = localStorage.getItem('jwt_token');
    if (!token) {
      return null;
    }

    try {
      const decodedToken: { role: string } = jwtDecode(token);
      return decodedToken.role;
    } catch (error) {
      console.error("Invalid token", error);
      return null;
    }
  }

  // Pomoćna metoda za laku proveru da li je korisnik admin
  isAdmin(): boolean {
    return this.getUserRole() === 'administrator';
  }
}
