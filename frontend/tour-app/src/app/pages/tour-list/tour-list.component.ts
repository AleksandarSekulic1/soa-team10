import { Component, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { TourService } from '../../services/tour.service';
import { AuthService } from '../../services/auth.service';

@Component({
  selector: 'app-tour-list',
  standalone: true,
  imports: [CommonModule, FormsModule],
  templateUrl: './tour-list.component.html',
  styleUrls: ['./tour-list.component.scss']
})
export class TourListComponent implements OnInit {
  allTours: any[] = [];
  selectedTour: any = null;
  review = {
    rating: 5,
    comment: '',
    visitDate: new Date().toISOString().split('T')[0] // Današnji datum
  };

  constructor(public authService: AuthService, private tourService: TourService) {}

  ngOnInit(): void {
    this.loadTours();
  }

  loadTours(): void {
    this.tourService.getAllTours().subscribe(tours => this.allTours = tours);
  }

  selectTourForReview(tour: any): void {
    this.selectedTour = tour;
  }

  submitReview(): void {
    if (!this.selectedTour) return;

    this.tourService.addReview(this.selectedTour.id, this.review).subscribe({
      next: () => {
        alert('Recenzija uspešno poslata!');
        this.selectedTour = null;
        this.loadTours(); // Ponovo učitaj ture da se vidi nova recenzija
      },
      error: (err) => {
        alert('Greška pri slanju recenzije.');
        console.error(err);
      }
    });
  }
}
