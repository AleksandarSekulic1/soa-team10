import { Component } from '@angular/core';
import { UserService } from '../../services/user.service';
import { FormsModule } from '@angular/forms';   // <-- Uvozimo FormsModule
import { CommonModule } from '@angular/common'; // <-- Uvozimo CommonModule (za *ngIf)

@Component({
  selector: 'app-registration',
  standalone: true, // <-- KAŽEMO DA JE KOMPONENTA SAMOSTALNA
  imports: [FormsModule, CommonModule], // <-- UVOZIMO ŠTA JOJ TREBA
  templateUrl: './registration.component.html',
  styleUrls: ['./registration.component.scss']
})
export class RegistrationComponent {
  // ... ostatak vaše klase ostaje isti ...
  user = {
    username: '',
    email: '',
    password: '',
    role: 'turista'
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
