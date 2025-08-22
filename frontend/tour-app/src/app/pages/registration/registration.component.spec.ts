import { Component } from '@angular/core';
import { UserService } from '../../services/user.service';
import { CommonModule } from '@angular/common'; // <-- DODAJEMO IMPORT
import { FormsModule } from '@angular/forms';   // <-- DODAJEMO IMPORT

@Component({
  selector: 'app-registration',
  standalone: true, // <-- ČINIMO GA SAMOSTALNIM
  imports: [CommonModule, FormsModule], // <-- DODAJEMO IMPORTE
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
        this.message = 'Uspešno ste se registrovali! Sada se možete prijaviti.';
      },
      error: (error) => {
        console.error('Došlo je do greške!', error);
        this.message = 'Greška prilikom registracije. Korisničko ime ili email možda već postoje.';
      }
    });
  }
}
