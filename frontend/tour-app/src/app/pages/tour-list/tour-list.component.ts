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
    visitDate: new Date().toISOString().split('T')[0],
    imageUrlsInput: '' // Polje za unos URL-ova kao string
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
    // Resetujemo formu svaki put kad se otvori
    this.review = {
      rating: 5,
      comment: '',
      visitDate: new Date().toISOString().split('T')[0],
      imageUrlsInput: ''
    };
  }

  submitReview(): void {
    if (!this.selectedTour) return;

    const reviewData = {
      rating: this.review.rating,
      comment: this.review.comment,
      visitDate: new Date(this.review.visitDate).toISOString(),
      // Pretvaramo string sa URL-ovima u niz stringova
      imageUrls: this.review.imageUrlsInput.split(',').map(url => url.trim()).filter(url => url)
    };

    this.tourService.addReview(this.selectedTour.id, reviewData).subscribe({
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
