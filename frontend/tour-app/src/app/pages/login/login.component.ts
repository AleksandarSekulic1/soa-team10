import { Component } from '@angular/core';
import { Router, RouterLink } from '@angular/router'; // <-- DODAJEMO RouterLink
import { UserService } from '../../services/user.service';
import { AuthService } from '../../services/auth.service';
import { CommonModule } from '@angular/common';   // <-- DODAJEMO CommonModule
import { FormsModule } from '@angular/forms';     // <-- DODAJEMO FormsModule

@Component({
  selector: 'app-login',
  standalone: true,
  // Dodajemo sve module koje templejt ove komponente koristi
  imports: [CommonModule, FormsModule, RouterLink],
  templateUrl: './login.component.html',
  styleUrls: ['./login.component.scss']
})
export class LoginComponent {
  credentials = {
    username: '',
    password: ''
  };
  errorMessage = '';

  constructor(
    private userService: UserService,
    private authService: AuthService,
    private router: Router
  ) { }

  onSubmit(): void {
    this.userService.login(this.credentials).subscribe({
      next: (response) => {
        console.log('Prijava uspešna!', response);
        this.authService.login(response.token);

        if (this.authService.isAdmin()) {
          this.router.navigate(['/users']);
        } else {
          // Za sada ga samo preusmeravamo na login sa porukom
          this.errorMessage = 'Prijava uspešna, ali nemate administratorski pristup.';
          // U pravoj aplikaciji bismo ga preusmerili na početnu stranicu za korisnike
          // npr. this.router.navigate(['/']);
        }
      },
      error: (error) => {
        console.error('Došlo je do greške!', error);
        this.errorMessage = 'Neispravno korisničko ime ili lozinka.';
      }
    });
  }
}
