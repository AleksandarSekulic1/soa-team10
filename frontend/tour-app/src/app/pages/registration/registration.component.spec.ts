// src/app/pages/registration/registration.component.ts
import { Component } from '@angular/core';
import { UserService } from '../../services/user.service';

@Component({
  selector: 'app-registration',
  templateUrl: './registration.component.html',
  styleUrls: ['./registration.component.scss']
})
export class RegistrationComponent {
  user = {
    username: '',
    email: '',
    password: '',
    role: 'turista' // Podrazumevana vrednost
  };
  message = '';

  constructor(private userService: UserService) { }

  onSubmit(): void {
    this.userService.register(this.user).subscribe({
      next: (response) => {
        console.log('Registracija uspešna!', response);
        this.message = 'Uspešno ste se registrovali!';
      },
      error: (error) => {
        console.error('Došlo je do greške!', error);
        this.message = 'Greška prilikom registracije. Pokušajte ponovo.';
      }
    });
  }
}
