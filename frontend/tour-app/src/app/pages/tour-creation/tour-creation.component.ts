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
  // IZMENA: Model 'tour' sada direktno sadrži 'tags' kao niz stringova
  tour = {
    name: '',
    description: '',
    difficulty: 3,
    tags: [] as string[]
  };

  // UKLONJENO: 'tagsInput' više nije potreban
  // tagsInput: string = '';

  message: string = '';

  // NOVO: Definisana lista dostupnih tagova koji će se prikazati kao checkbox-ovi
  availableTags = [
    { name: 'Planinarenje', value: 'hiking' },
    { name: 'Priroda', value: 'nature' },
    { name: 'Istorija', value: 'history' },
    { name: 'Avantura', value: 'adventure' },
    { name: 'Kultura', value: 'cultural' },
    { name: 'Opuštanje', value: 'relaxing' }
  ];

  constructor(private tourService: TourService, private router: Router) { }

  // NOVO: Metoda koja se poziva svaki put kad se klikne na checkbox
  onTagChange(event: any): void {
    const tagName = event.target.value;
    const isChecked = event.target.checked;

    if (isChecked) {
      // Ako je tag odabran, dodajemo ga u niz 'tour.tags'
      this.tour.tags.push(tagName);
    } else {
      // Ako je odabir poništen, uklanjamo tag iz niza
      const index = this.tour.tags.indexOf(tagName);
      if (index > -1) {
        this.tour.tags.splice(index, 1);
      }
    }
  }

  onSubmit(): void {
    // IZMENA: Nema više potrebe za kreiranjem 'tourData' objekta i parsiranjem stringa.
    // Objekat 'this.tour' već ima ispravan format sa nizom tagova.
    const tourData = {
      ...this.tour,
      difficulty: Number(this.tour.difficulty) // Osiguravamo da je težina broj
    };

    this.tourService.createTour(tourData).subscribe({
      next: (response) => {
        console.log('Tura je uspešno kreirana!', response);
        alert('Tura je uspešno kreirana! Bićete preusmereni na početnu stranicu.');
        
        // Preusmeravamo korisnika
        this.router.navigate(['/home']);
      },
      error: (error) => {
        console.error('Došlo je do greške!', error);
        this.message = 'Greška prilikom kreiranja ture. Proverite podatke.';
      }
    });
  }
}