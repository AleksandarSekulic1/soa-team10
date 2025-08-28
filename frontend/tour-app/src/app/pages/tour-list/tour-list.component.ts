import { Component, OnDestroy, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { TourService } from '../../services/tour.service';
import { AuthService } from '../../services/auth.service';
import { ShoppingCartService } from '../../services/shopping-cart.service';
import { Subscription } from 'rxjs';

@Component({
  selector: 'app-tour-list',
  standalone: true,
  imports: [CommonModule, FormsModule],
  templateUrl: './tour-list.component.html',
  styleUrls: ['./tour-list.component.scss']
})
export class TourListComponent implements OnInit, OnDestroy {
  allTours: any[] = [];
  cartItems: any[] = [];
  purchasedTourIds: string[] = [];
  private cartSubscription: Subscription | undefined;
  private purchasedSubscription: Subscription | undefined;

  selectedTour: any = null;
  review = {
    rating: 5,
    comment: '',
    visitDate: new Date().toISOString().split('T')[0],
    imageUrlsInput: ''
  };

  constructor(
    public authService: AuthService,
    private tourService: TourService,
    private shoppingCartService: ShoppingCartService
  ) {}

  ngOnInit(): void {
    this.loadTours();

    this.cartSubscription = this.shoppingCartService.cart$.subscribe(cart => {
      this.cartItems = cart?.items || [];
    });

    this.purchasedSubscription = this.shoppingCartService.purchasedTours$.subscribe(ids => {
      this.purchasedTourIds = ids;
    });

    this.shoppingCartService.getCart().subscribe();
  }

  ngOnDestroy(): void {
    this.cartSubscription?.unsubscribe();
    this.purchasedSubscription?.unsubscribe();
  }

  loadTours(): void {
    this.tourService.getAllTours().subscribe(tours => this.allTours = tours);
  }

  selectTourForReview(tour: any): void {
    this.selectedTour = tour;
  }

  submitReview(): void {
    if (!this.selectedTour) return;
    const reviewData = {
      rating: this.review.rating,
      comment: this.review.comment,
      visitDate: new Date(this.review.visitDate).toISOString(),
      imageUrls: this.review.imageUrlsInput.split(',').map(url => url.trim()).filter(url => url)
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

  addToCart(tour: any): void {
    this.shoppingCartService.addItemToCart(tour).subscribe({
      next: () => {
        alert(`Tura "${tour.name}" je dodata u korpu!`);
      },
      error: (err) => {
        alert('Greška prilikom dodavanja u korpu: ' + (err.error?.message || 'Pokušajte ponovo.'));
        console.error(err);
      }
    });
  }

  isTourInCart(tourId: string): boolean {
    return this.cartItems.some(item => item.tourId === tourId);
  }

  isTourPurchased(tourId: string): boolean {
    return this.purchasedTourIds.includes(tourId);
  }
}
