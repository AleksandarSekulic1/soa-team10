import { Component } from '@angular/core';
import { AuthService } from '../../services/auth.service';
import { CommonModule } from '@angular/common';   // <-- DODAJEMO CommonModule
import { FormsModule } from '@angular/forms';     // <-- DODAJEMO FormsModule
import { RouterModule } from '@angular/router'; // <-- 1. Uvezite RouterModule


@Component({
  selector: 'app-navbar',
  standalone: true, 
  imports: [CommonModule, FormsModule,RouterModule],
  templateUrl: './navbar.component.html',
  styleUrls: ['./navbar.component.scss']
})
export class NavbarComponent {
  // Injektujemo AuthService da bismo mogli da koristimo njegove metode u HTML-u
  constructor(public authService: AuthService) {}

  logout(): void {
    this.authService.logout();
  }
}