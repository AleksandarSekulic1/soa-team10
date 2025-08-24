import { Component } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { Router } from '@angular/router';
import { TourService } from '../../services/tour.service';

@Component({
  selector: 'app-tour-creation',
  standalone: true,
  imports: [CommonModule, FormsModule],
  templateUrl: './tour-creation.component.html',
  styleUrls: ['./tour-creation.component.scss']
})
export class TourCreationComponent {
  tour: any = {
    name: '',
    description: '',
    difficulty: 3,
  };
  tagsInput: string = '';
  message: string = '';

  constructor(private tourService: TourService, private router: Router) { }

  onSubmit(): void {
    const tourData = {
      ...this.tour,
      tags: this.tagsInput.split(',').map(tag => tag.trim()).filter(tag => tag)
    };
    tourData.difficulty = Number(tourData.difficulty);

    this.tourService.createTour(tourData).subscribe({
      next: (response) => {
        console.log('Tura je uspešno kreirana!', response);
        // Možemo prikazati kratku poruku pre preusmeravanja, ali nije obavezno
        alert('Tura je uspešno kreirana! Vraćate se na početnu stranicu.');

        // Preusmeravamo korisnika na 'home' stranicu
        this.router.navigate(['/home']);
      },
      error: (error) => {
        console.error('Došlo je do greške!', error);
        this.message = 'Greška prilikom kreiranja ture. Proverite podatke.';
      }
    });
  }
}
