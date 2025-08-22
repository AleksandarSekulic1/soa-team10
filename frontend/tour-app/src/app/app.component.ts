import { Component } from '@angular/core';
import { RouterOutlet, RouterLink } from '@angular/router'; // Dodajemo RouterLink
import { CommonModule } from '@angular/common'; // Dodajemo CommonModule za *ngIf
import { AuthService } from './services/auth.service'; // Uvozimo AuthService

@Component({
  selector: 'app-root',
  standalone: true,
  imports: [
    RouterOutlet,
    RouterLink, // Dodajemo
    CommonModule // Dodajemo
  ],
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.scss']
})
export class AppComponent {
  title = 'tour-app';

  // Činimo servis javnim da bismo mu pristupili iz HTML-a
  constructor(public authService: AuthService) {}
}
