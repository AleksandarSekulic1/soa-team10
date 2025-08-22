import { Component } from '@angular/core';
import { Router, RouterLink } from '@angular/router';
import { UserService } from '../../services/user.service';
import { AuthService } from '../../services/auth.service';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';

@Component({
  selector: 'app-login',
  standalone: true,
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
    // Ovaj log je najvažniji za proveru odgovora sa servera
    console.log('CEO ODGOVOR SA SERVERA:', response);
    this.authService.login(response);
    this.router.navigate(['/home']);
  },
      error: (error) => {
        console.error('Došlo je do greške!', error);
        this.errorMessage = 'Neispravno korisničko ime ili lozinka.';
      }
    });
  }
}