import { Component, OnInit } from '@angular/core';
import { CommonModule, TitleCasePipe } from '@angular/common';
import { TourService } from '../../services/tour.service';
import { Tour } from '../../models/tour.model';
import { RouterLink } from '@angular/router';

@Component({
  selector: 'app-my-tours',
  standalone: true,
  imports: [CommonModule, RouterLink, TitleCasePipe], // Dodajemo TitleCasePipe
  templateUrl: './my-tours.component.html',
  styleUrls: ['./my-tours.component.scss']
})
export class MyToursComponent implements OnInit {
  myTours: Tour[] = []; // Koristimo jak tip

  constructor(private tourService: TourService) {}

  ngOnInit(): void {
    this.loadTours();
  }

  loadTours(): void {
    this.tourService.getMyTours().subscribe({
      next: (tours) => {
        this.myTours = tours;
      },
      error: (err) => {
        console.error('Error fetching tours:', err);
      }
    });
  }

  // --- NOVE METODE ZA UPRAVLJANJE STATUSIMA ---

  publishTour(tourId: string, index: number): void {
    this.tourService.publishTour(tourId).subscribe({
      next: (updatedTour) => {
        // Ažuriramo samo turu koja je promenjena u listi
        this.myTours[index] = updatedTour;
        alert('Tour successfully published!');
      },
      error: (err) => {
        // Prikazujemo grešku sa backenda (npr. "nema dovoljno tačaka")
        alert(`Error publishing tour: ${err.error.error}`);
        console.error(err);
      }
    });
  }

  archiveTour(tourId: string, index: number): void {
    this.tourService.archiveTour(tourId).subscribe({
      next: (updatedTour) => {
        this.myTours[index] = updatedTour;
        alert('Tour successfully archived!');
      },
      error: (err) => {
        alert(`Error archiving tour: ${err.error.error}`);
        console.error(err);
      }
    });
  }

  reactivateTour(tourId: string, index: number): void {
    this.tourService.reactivateTour(tourId).subscribe({
      next: (updatedTour) => {
        this.myTours[index] = updatedTour;
        alert('Tour successfully reactivated!');
      },
      error: (err) => {
        alert(`Error reactivating tour: ${err.error.error}`);
        console.error(err);
      }
    });
  }
}