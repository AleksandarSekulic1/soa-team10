import { Component, OnInit, Inject, PLATFORM_ID } from '@angular/core';
import { isPlatformBrowser, CommonModule } from '@angular/common';
import { ActivatedRoute } from '@angular/router';
import { TourService } from '../../services/tour.service';
import { forkJoin, of } from 'rxjs';
import { catchError } from 'rxjs/operators';
import { Tour } from '../../models/tour.model';
import { TouristPosition } from '../../models/tourist-position.model';

// Svi 'import' za leaflet su uklonjeni sa vrha

@Component({
  selector: 'app-position-simulator',
  standalone: true,
  imports: [CommonModule],
  templateUrl: './position-simulator.component.html',
  styleUrls: ['./position-simulator.component.scss']
})
export class PositionSimulatorComponent implements OnInit {
  private map: any;
  private touristMarker: any;
  private routeToStartLine: any;
  public tour: Tour | undefined;
  public touristPosition: TouristPosition | undefined;

  constructor(
    private route: ActivatedRoute,
    private tourService: TourService,
    @Inject(PLATFORM_ID) private platformId: Object
  ) {}

  ngOnInit(): void {
    const tourId = this.route.snapshot.paramMap.get('id');
    if (tourId) {
      forkJoin({
        tour: this.tourService.getTourById(tourId),
        position: this.tourService.getTouristPosition().pipe(
          catchError(error => of(undefined))
        )
      }).subscribe(({ tour, position }) => {
        this.tour = tour;
        this.touristPosition = position;
        
        // Inicijalizacija mape se poziva tek kada stignu podaci
        setTimeout(() => this.initMap(), 0);
      });
    }
  }

  private async initMap(): Promise<void> {
    // Proveravamo da li se kod izvršava u browseru
    if (isPlatformBrowser(this.platformId)) {
      if (!this.tour) return;
    
      // Dinamički uvozimo Leaflet i njegove dodatke
      const L = await import('leaflet');
      (window as any).L = L; // Činimo L globalno dostupnim za dodatke
      await import('leaflet-defaulticon-compatibility');
      await import('leaflet-routing-machine');

      this.map = L.map('simulator-map');
      L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png').addTo(this.map);

      // Iscrtaj glavnu rutu ture ako ima ključnih tačaka
      if (this.tour.keyPoints && this.tour.keyPoints.length > 1) {
        const waypoints = this.tour.keyPoints.map(kp => L.latLng(kp.latitude, kp.longitude));
        (L as any).Routing.control({
          waypoints: waypoints,
          routeWhileDragging: false,
          addWaypoints: false,
          show: false,
          lineOptions: { styles: [{ color: '#007bff', weight: 5, opacity: 0.8 }] },
          createMarker: () => { return null; }
        }).addTo(this.map);

        this.tour.keyPoints.forEach(kp => {
          L.marker([kp.latitude, kp.longitude]).addTo(this.map).bindPopup(`<b>${kp.name}</b>`);
        });
      }
      
      this.updateTouristVisuals(L);
      
      if(this.tour.keyPoints && this.tour.keyPoints.length > 0){
          const tourBounds = L.latLngBounds(this.tour.keyPoints.map(kp => [kp.latitude, kp.longitude]));
          this.map.fitBounds(tourBounds, { padding: [50, 50] });
      }

      this.map.on('click', (e: L.LeafletMouseEvent) => {
        const { lat, lng } = e.latlng;
        this.tourService.updateTouristPosition({ latitude: lat, longitude: lng }).subscribe(
          updatedPosition => {
            this.touristPosition = updatedPosition;
            this.updateTouristVisuals(L);
          }
        );
      });
    }
  }
  
    private updateTouristVisuals(L: any): void {
    if (!this.map || !this.touristPosition) return;
    
    const touristLatLng = L.latLng(this.touristPosition.Latitude, this.touristPosition.Longitude);
    
    // Ažuriraj marker turiste (ovaj deo ostaje isti)
    if (!this.touristMarker) {
      const touristIcon = L.icon({
        iconUrl: 'assets/images/default-avatar.png',
        iconSize: [40, 40],
        iconAnchor: [20, 40],
        popupAnchor: [0, -40],
        className: 'tourist-marker'
      });
      this.touristMarker = L.marker(touristLatLng, { icon: touristIcon, zIndexOffset: 1000 }).addTo(this.map);
    } else {
      this.touristMarker.setLatLng(touristLatLng);
    }
    
    // AŽURIRANA LOGIKA: Koristimo Routing Machine umesto Polyline
    if (this.tour?.keyPoints && this.tour.keyPoints.length > 0) {
      const startPointLatLng = L.latLng(this.tour.keyPoints[0].latitude, this.tour.keyPoints[0].longitude);
      
      // Ako linija (sada ruter) već postoji, samo joj ažuriramo tačke
      if (this.routeToStartLine) {
        this.routeToStartLine.setWaypoints([touristLatLng, startPointLatLng]);
      } else {
        // Ako ne postoji, kreiramo novi ruter za pomoćnu liniju
        this.routeToStartLine = (L as any).Routing.control({
          waypoints: [touristLatLng, startPointLatLng],
          // Opcije da se sakriju svi nepotrebni elementi
          addWaypoints: false,
          fitSelectedRoutes: false, // Ne želimo da mapa zumira na ovu liniju
          show: false,
          // Podešavanje stila linije da bude isprekidana i crvena
          lineOptions: {
            styles: [{ color: '#eb0b22ff', weight: 5, opacity: 0.8, dashArray: '10, 10' }]
          },
          // Ne želimo da ovaj ruter kreira svoje pinove
          createMarker: () => null 
        }).addTo(this.map);
      }
    }
  }
}