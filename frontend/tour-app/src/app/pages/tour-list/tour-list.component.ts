import { Component, OnInit } from '@angular/core';
import { CommonModule, DatePipe } from '@angular/common'; // Dodajemo DatePipe
import { FormsModule } from '@angular/forms';
import { TourService } from '../../services/tour.service';
import { AuthService } from '../../services/auth.service';

@Component({
  selector: 'app-tour-list',
  standalone: true,
  // DatePipe se mora dodati i u imports
  imports: [CommonModule, FormsModule, DatePipe], 
  templateUrl: './tour-list.component.html',
  styleUrls: ['./tour-list.component.scss']
})
export class TourListComponent implements OnInit {
  allTours: any[] = [];
  selectedTour: any = null;
  
  availableReviewImages: string[] = [
    'assets/images/tura1.png',
    'assets/images/tura2.png',
    'assets/images/tura3.png',
    'assets/images/tura4.png',
    'assets/images/tura5.png',
    'assets/images/tura6.png',
    'assets/images/tura7.png',
    'assets/images/tura8.png',
    'assets/images/tura9.png',
    'assets/images/tura10.png'
  ];

  review = {
    rating: 5,
    comment: '',
    visitDate: new Date().toISOString().split('T')[0],
    imageUrls: [] as string[]
  };

  constructor(public authService: AuthService, private tourService: TourService) {}

  ngOnInit(): void {
    this.loadTours();
  }

  loadTours(): void {
    this.tourService.getAllTours().subscribe(tours => {
      // IZMENA: Svakoj turi dodajemo property za praćenje slajda
      this.allTours = tours.map(tour => ({
        ...tour,
        currentReviewIndex: 0 // Početni slajd je uvek prvi (index 0)
      }));
    });
  }

  // NOVO: Metode za navigaciju slajdera
  nextReview(tour: any): void {
    if (tour.reviews && tour.reviews.length > 0) {
      tour.currentReviewIndex = (tour.currentReviewIndex + 1) % tour.reviews.length;
    }
  }

  prevReview(tour: any): void {
    if (tour.reviews && tour.reviews.length > 0) {
      tour.currentReviewIndex = (tour.currentReviewIndex - 1 + tour.reviews.length) % tour.reviews.length;
    }
  }

  // NOVO: Pomoćna funkcija za prikaz zvezdica
  getStarArray(rating: number): any[] {
    return new Array(Math.round(rating));
  }
  
  // Ostatak koda (selectTourForReview, toggleImageSelection, submitReview) ostaje isti...
  selectTourForReview(tour: any): void {
    this.selectedTour = tour;
    this.review = {
      rating: 5,
      comment: '',
      visitDate: new Date().toISOString().split('T')[0],
      imageUrls: []
    };
  }

  toggleImageSelection(url: string): void {
    const index = this.review.imageUrls.indexOf(url);
    if (index > -1) {
      this.review.imageUrls.splice(index, 1);
    } else {
      this.review.imageUrls.push(url);
    }
  }

  submitReview(): void {
    if (!this.selectedTour) return;
    const reviewData = {
      rating: this.review.rating,
      comment: this.review.comment,
      visitDate: new Date(this.review.visitDate).toISOString(),
      imageUrls: this.review.imageUrls
    };
    this.tourService.addReview(this.selectedTour.id, reviewData).subscribe({
      next: () => {
        alert('Recenzija uspešno poslata!');
        this.selectedTour = null;
        this.loadTours();
      },
      error: (err) => {
        alert('Greška pri slanju recenzije.');
        console.error(err);
      }
    });
  }
}