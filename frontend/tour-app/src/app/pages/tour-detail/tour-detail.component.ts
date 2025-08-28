import { Component, OnInit, Inject, PLATFORM_ID, OnDestroy } from '@angular/core';
import { isPlatformBrowser, CommonModule, CurrencyPipe, DatePipe } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { ActivatedRoute, RouterLink } from '@angular/router';
import { TourService } from '../../services/tour.service';
import { Tour, TourKeyPoint } from '../../models/tour.model';
import { AuthService } from '../../services/auth.service';
import { ShoppingCartService } from '../../services/shopping-cart.service';
import { Subscription } from 'rxjs';

@Component({
  selector: 'app-tour-detail',
  standalone: true,
  imports: [CommonModule, RouterLink, CurrencyPipe, DatePipe, FormsModule],
  templateUrl: './tour-detail.component.html',
  styleUrls: ['./tour-detail.component.scss']
})
export class TourDetailComponent implements OnInit, OnDestroy {
  tour: Tour | undefined;
  isLoading = true;
  error: string | null = null;

  // Stanja za prikaz
  isTourist: boolean = false;
  isAuthor: boolean = false;
  isPurchased: boolean = false;
  isInCart: boolean = false;
  displayedKeyPoints: TourKeyPoint[] = [];

  private subscriptions = new Subscription();

  // Propertiji za modal i mape
  isEditModalVisible = false;
  currentKeyPointToEdit: TourKeyPoint | null = null;
  private editMap: any;
  private editMarker: any;
  private routeMap: any;
  availableImages: string[] = [
    'assets/images/default-avatar.png', 'assets/images/men2.png', 'assets/images/men3.png',
    'assets/images/men4.png', 'assets/images/men5.png', 'assets/images/women1.png',
    'assets/images/women2.png', 'assets/images/women3.png', 'assets/images/women4.png',
    'assets/images/women5.png'
  ];

  constructor(
    private route: ActivatedRoute,
    private tourService: TourService,
    private authService: AuthService,
    private shoppingCartService: ShoppingCartService,
    @Inject(PLATFORM_ID) private platformId: Object
  ) {}

  ngOnInit(): void {
    this.isTourist = this.authService.isTourist();
    const currentUsername = this.authService.getUsername();

    const tourId = this.route.snapshot.paramMap.get('id');
    if (tourId) {
      // Pratimo promene u korpi
      this.subscriptions.add(
        this.shoppingCartService.cart$.subscribe(cart => {
          this.isInCart = cart?.items?.some((item: any) => item.tourId === tourId) || false;
        })
      );
      // Pratimo promene u kupljenim turama
      this.subscriptions.add(
        this.shoppingCartService.purchasedTours$.subscribe(purchasedIds => {
          this.isPurchased = purchasedIds.includes(tourId);
          this.updateVisibleKeyPoints();
        })
      );

      this.tourService.getTourById(tourId).subscribe({
        next: (fetchedTour) => {
          this.tour = fetchedTour;
          this.isLoading = false;
          this.isAuthor = !!currentUsername && fetchedTour.authorId === currentUsername;
          this.updateVisibleKeyPoints();
          setTimeout(() => this.initRouteMap(), 0);
        },
        error: (err) => { this.error = 'Tour not found.'; this.isLoading = false; }
      });
    }
  }

  ngOnDestroy(): void {
    this.subscriptions.unsubscribe();
  }

  private updateVisibleKeyPoints(): void {
    if (!this.tour) return;

    if (this.isAuthor || this.isPurchased) {
      this.displayedKeyPoints = this.tour.keyPoints;
    } else {
      this.displayedKeyPoints = this.tour.keyPoints.slice(0, 1);
    }
  }

  addToCart(): void {
    if (!this.tour) return;
    this.shoppingCartService.addItemToCart(this.tour).subscribe({
      next: () => alert(`Tura "${this.tour?.name}" je dodata u korpu!`),
      error: (err) => alert('Greška: ' + (err.error?.message || 'Pokušajte ponovo.'))
    });
  }

  private async initRouteMap(): Promise<void> {
    if (isPlatformBrowser(this.platformId)) {
      if (!this.tour || !this.tour.keyPoints || this.tour.keyPoints.length < 2) {
        return;
      }
      const L = await import('leaflet');
      (window as any).L = L;
      await import('leaflet-defaulticon-compatibility');
      await import('leaflet-routing-machine');
      this.setupMapWithRouting(L);
    }
  }

  private setupMapWithRouting(L: any): void {
    if (this.routeMap) {
      this.routeMap.remove();
    }
    this.routeMap = L.map('tour-route-map');
    L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png').addTo(this.routeMap);
    const waypoints = this.tour!.keyPoints.map(kp => L.latLng(kp.latitude, kp.longitude));
    (L.Routing as any).control({
      waypoints: waypoints,
      routeWhileDragging: false,
      addWaypoints: false,
      show: false,
      lineOptions: {
        styles: [{ color: '#007bff', weight: 5, opacity: 0.8 }],
        extendToWaypoints: true,
        missingRouteTolerance: 1
      },
      createMarker: (i: number, waypoint: any, n: number) => {
        const keyPoint = this.tour!.keyPoints[i];
        return L.marker(waypoint.latLng).bindPopup(`<b>${keyPoint.name}</b><br>${keyPoint.description}`);
      }
    }).addTo(this.routeMap);
  }

  async openEditModal(keyPoint: TourKeyPoint): Promise<void> {
    this.currentKeyPointToEdit = { ...keyPoint };
    this.isEditModalVisible = true;

    if (isPlatformBrowser(this.platformId)) {
      setTimeout(async () => {
        const L = await import('leaflet');
        await import('leaflet-defaulticon-compatibility');
        this.initEditMap(L);
      }, 0);
    }
  }

  closeEditModal(): void {
    if (this.editMap) { this.editMap.remove(); }
    this.isEditModalVisible = false;
    this.currentKeyPointToEdit = null;
  }

  private initEditMap(L: any): void {
    if (!this.currentKeyPointToEdit) return;
    const kp = this.currentKeyPointToEdit;

    this.editMap = L.map('map-edit').setView([kp.latitude, kp.longitude], 13);
    L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png').addTo(this.editMap);
    this.editMarker = L.marker([kp.latitude, kp.longitude]).addTo(this.editMap);

    this.editMap.on('click', (e: any) => {
      this.currentKeyPointToEdit!.latitude = e.latlng.lat;
      this.currentKeyPointToEdit!.longitude = e.latlng.lng;
      this.editMarker!.setLatLng(e.latlng);
    });
  }

  getStarArray(rating: number): any[] { return new Array(Math.round(rating)); }

  selectImageForEdit(imagePath: string): void {
    if (this.currentKeyPointToEdit) { this.currentKeyPointToEdit.imageUrl = imagePath; }
  }

  onUpdateKeyPoint(): void {
    if (!this.currentKeyPointToEdit || !this.tour) return;
    this.tourService.updateKeyPoint(this.tour.id, this.currentKeyPointToEdit.id, this.currentKeyPointToEdit)
      .subscribe({
        next: (updatedKeyPoint) => {
          const index = this.tour!.keyPoints.findIndex(kp => kp.id === updatedKeyPoint.id);
          if (index !== -1) { this.tour!.keyPoints[index] = updatedKeyPoint; }
          this.closeEditModal();
        },
        error: (err) => console.error('Error updating key point:', err)
      });
  }

  deleteKeyPoint(keyPointId: string): void {
    if (!this.tour) return;
    if (confirm('Are you sure you want to delete this key point?')) {
      this.tourService.deleteKeyPoint(this.tour.id, keyPointId).subscribe({
        next: () => {
          this.tour!.keyPoints = this.tour!.keyPoints.filter(kp => kp.id !== keyPointId);
        },
        error: (err) => console.error('Error deleting key point:', err)
      });
    }
  }
}
