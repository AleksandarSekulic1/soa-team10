import { Component } from '@angular/core';
import { Router } from '@angular/router';
import { UserService } from '../../services/user.service';
import { AuthService } from '../../services/auth.service';
import { CommonModule } from '@angular/common';   // <-- DODAJEMO IMPORT
import { FormsModule } from '@angular/forms';     // <-- DODAJEMO IMPORT

@Component({
  selector: 'app-login',
  standalone: true,
  // Dodajemo module koje komponenta koristi u svom templejtu
  imports: [CommonModule, FormsModule],
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
          this.errorMessage = 'Prijava uspešna, ali nemate administratorski pristup.';
        }
      },
      error: (error) => {
        console.error('Došlo je do greške!', error);
        this.errorMessage = 'Neispravno korisničko ime ili lozinka.';
      }
    });
  }
}
