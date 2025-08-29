import { Component, OnInit, OnDestroy, Inject, PLATFORM_ID } from '@angular/core';
import { isPlatformBrowser, CommonModule, DecimalPipe } from '@angular/common';
import { ActivatedRoute, Router, RouterLink } from '@angular/router';
import { TourService } from '../../services/tour.service';
import { forkJoin, of, Subscription, interval } from 'rxjs';
import { catchError } from 'rxjs/operators';
import { Tour } from '../../models/tour.model';
import { TouristPosition } from '../../models/tourist-position.model';
import { TourExecution } from '../../models/tour-execution.model';

@Component({
  selector: 'app-position-simulator',
  standalone: true,
  imports: [CommonModule, DecimalPipe, RouterLink],
  templateUrl: './position-simulator.component.html',
  styleUrls: ['./position-simulator.component.scss']
})
export class PositionSimulatorComponent implements OnInit, OnDestroy {
  private map: any;
  private touristMarker: any;
  private routeToStartLine: any;
  public tour: Tour | undefined;
  public touristPosition: TouristPosition | undefined;
  public tourExecution: TourExecution | undefined;
  private pollingSubscription: Subscription | undefined;
  private keyPointMarkers: any[] = [];
  private leaflet: any;

  private routePoints: any[] = [];
  private toStartRoutePoints: any[] = [];
  private movementSimulationTimer: any;
  
  // Svojstva za tajmer proteklog vremena
  public elapsedTime: string = '00:00:00';
  private elapsedTimeTimer: any;


  constructor(
    private route: ActivatedRoute,
    private router: Router,
    private tourService: TourService,
    @Inject(PLATFORM_ID) private platformId: Object
  ) {}

  public getKeyPointName(keyPointId: string): string {
  // Proverite da li tura i ključne tačke postoje
  if (!this.tour || !this.tour.keyPoints) {
    return `ID: ${keyPointId}`; // Vraća ID ako podaci o turi još nisu učitani
  }

  // Pronađite ključnu tačku po ID-u
  // VAŽNO: Proverite da li se polja u vašem modelu za KeyPoint zovu 'id' i 'name'
  const keyPoint = this.tour.keyPoints.find(kp => kp.id === keyPointId);

  // Vratite ime tačke ako je pronađena, u suprotnom vratite njen ID
  return keyPoint ? keyPoint.name : `ID: ${keyPointId}`;
} 

  ngOnInit(): void {
    const tourId = this.route.snapshot.paramMap.get('id');
    if (tourId) {
      forkJoin({
        tour: this.tourService.getTourById(tourId),
        position: this.tourService.getTouristPosition().pipe(catchError(() => of(undefined))),
        execution: this.tourService.getActiveExecutionForUser()
      }).subscribe(({ tour, position, execution }) => {
        this.tour = tour;
        this.touristPosition = position;

        if (execution && execution.TourId !== tour.id) {
          alert(`You have another active tour. Please complete or abandon it first.`);
          this.router.navigate(['/tours', execution.TourId, 'simulate']);
          return;
        }

        this.tourExecution = execution || undefined;
        if (this.tourExecution) {
          this.startPositionPolling();
          this.startElapsedTimeTimer();
        }
        setTimeout(() => this.initMap(), 0);
      });
    }
  }

  ngOnDestroy(): void {
    if (this.pollingSubscription) this.pollingSubscription.unsubscribe();
    if (this.movementSimulationTimer) clearInterval(this.movementSimulationTimer);
    if (this.elapsedTimeTimer) clearInterval(this.elapsedTimeTimer);
    if (this.map) this.map.remove();
  }

  startTour(): void {
    if (!this.tour) return;
    this.tourService.startTour(this.tour.id).subscribe({
      next: (execution) => {
        this.tourExecution = execution;
        alert('Tour started successfully!');
        if (this.map) this.map.off('click');
        this.startMovementSimulation(true);
        this.startPositionPolling();
        this.startElapsedTimeTimer();
      },
      error: (err) => alert(`Error starting tour: ${err.error.error || 'Unknown error'}`)
    });
  }

  // Tajmer koji na svakih 10 sekundi šalje poziciju na backend
  private startPositionPolling(): void {
    if (!isPlatformBrowser(this.platformId)) return;
    
    this.pollingSubscription = interval(10000).subscribe(() => {
      if (this.touristPosition && this.tourExecution && this.tourExecution.Status === 'Active') {
        this.tourService.checkPosition(this.tourExecution.ID, {
          latitude: this.touristPosition.Latitude,
          longitude: this.touristPosition.Longitude
        }).subscribe(updatedExecution => {
          const oldCompletedCount = this.tourExecution?.CompletedKeyPoints?.length || 0;
          updatedExecution.CompletedKeyPoints = updatedExecution.CompletedKeyPoints || [];
          this.tourExecution = updatedExecution;

          if ((this.tourExecution.CompletedKeyPoints?.length || 0) > oldCompletedCount) {
            console.log(`A key point was completed! Total: ${this.tourExecution.CompletedKeyPoints.length}`);
            this.updateKeyPointMarkers();
          }
        });
      }
    });
  }

  // Tajmer koji samo pomera pin na mapi radi vizuelne simulacije
  private startMovementSimulation(isMovingToStart: boolean): void {
    if (!isPlatformBrowser(this.platformId)) return;

    const currentRoute = isMovingToStart ? this.toStartRoutePoints : this.routePoints;
    if (currentRoute.length === 0) {
      if (isMovingToStart) this.startMovementSimulation(false);
      return;
    }

    if (isMovingToStart && this.routeToStartLine) {
      this.map.removeControl(this.routeToStartLine);
      this.routeToStartLine = null;
    }

    if (this.movementSimulationTimer) clearInterval(this.movementSimulationTimer);
    
    let routeIndex = 0;
    const intervalTime = 1000;

    this.movementSimulationTimer = setInterval(() => {
      if (routeIndex >= currentRoute.length - 1) {
        clearInterval(this.movementSimulationTimer);
        if (isMovingToStart) {
          this.startMovementSimulation(false);
        } else {
          alert("Simulation finished! You can now complete the tour.");
        }
        return;
      }

      routeIndex = Math.min(routeIndex + 1, currentRoute.length - 1);
      const newPosition = currentRoute[routeIndex];
      
      this.touristPosition = {
        ...(this.touristPosition!),
        Latitude: newPosition.lat,
        Longitude: newPosition.lng
      };
      
      this.updateTouristVisuals(true);
    }, intervalTime);
  }

  private startElapsedTimeTimer(): void {
    if (!isPlatformBrowser(this.platformId) || !this.tourExecution) return;

    const startTime = new Date(this.tourExecution.StartTime).getTime();

    this.elapsedTimeTimer = setInterval(() => {
      const now = new Date().getTime();
      const difference = now - startTime;

      const hours = Math.floor((difference % (1000 * 60 * 60 * 24)) / (1000 * 60 * 60));
      const minutes = Math.floor((difference % (1000 * 60 * 60)) / (1000 * 60));
      const seconds = Math.floor((difference % (1000 * 60)) / 1000);

      this.elapsedTime = 
        `${hours.toString().padStart(2, '0')}:${minutes.toString().padStart(2, '0')}:${seconds.toString().padStart(2, '0')}`;
    }, 1000);
  }

  private updateKeyPointMarkers(): void {
    if (this.leaflet) {
      this.keyPointMarkers.forEach(marker => {
          const keyPointId = marker.options.keyPointId;
          const completedPoints = this.tourExecution?.CompletedKeyPoints || [];
          const isCompleted = completedPoints.some(kp => kp.KeyPointId === keyPointId);
          
          if (isCompleted && !marker.isCompleted) {
              marker.setIcon(this.leaflet.icon({
                  iconUrl: 'assets/images/default-avatar.png',
                  iconSize: [25, 41],
                  iconAnchor: [12, 41],
              }));
              marker.isCompleted = true;
          }
      });
    }
  }

  abandonTour(): void {
    if (!this.tourExecution) return;
    clearInterval(this.movementSimulationTimer);
    clearInterval(this.elapsedTimeTimer);
    this.pollingSubscription?.unsubscribe();
    this.tourService.abandonTour(this.tourExecution.ID).subscribe(() => {
      alert('Tour abandoned.');
      //this.router.navigate(['/tours', this.tour?.id]);
    });
  }

  completeTour(): void {
    if (!this.tourExecution) return;
    
    const totalKeyPoints = this.tour?.keyPoints?.length || 0;
    const completedKeyPointsCount = this.tourExecution.CompletedKeyPoints?.length || 0;

    if (completedKeyPointsCount < totalKeyPoints) {
      alert('You have not completed all key points yet!');
      return;
    }
 
    clearInterval(this.movementSimulationTimer);
    clearInterval(this.elapsedTimeTimer);
    this.pollingSubscription?.unsubscribe();
    this.tourService.completeTour(this.tourExecution.ID).subscribe(() => {
      alert('Tour completed successfully!');
      //this.router.navigate(['/tours', this.tour?.id]);
    });
  }

  private async initMap(): Promise<void> {
    if (isPlatformBrowser(this.platformId)) {
      if (!this.tour) return;
    
      this.leaflet = await import('leaflet');
      const L = this.leaflet;
      (window as any).L = L;
      await import('leaflet-defaulticon-compatibility');
      await import('leaflet-routing-machine');

      this.map = L.map('simulator-map');
      L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png').addTo(this.map);

      if (this.tour.keyPoints && this.tour.keyPoints.length > 1) {
        const waypoints = this.tour.keyPoints.map(kp => L.latLng(kp.latitude, kp.longitude));
        
        const routingControl = (L as any).Routing.control({
          waypoints: waypoints,
          routeWhileDragging: false, addWaypoints: false, show: false,
          lineOptions: { styles: [{ color: '#007bff', weight: 5, opacity: 0.8 }] },
          createMarker: () => null
        }).addTo(this.map);
        
        routingControl.on('routesfound', (e: any) => {
          if (e.routes && e.routes.length > 0) {
            this.routePoints = e.routes[0].coordinates;
            // Ako je tura već aktivna (nastavljamo je), odmah pokreni simulaciju
            if (this.tourExecution) {
              this.startMovementSimulation(false);
            }
          }
        });
      }

      this.keyPointMarkers = []; 
      if (this.tour.keyPoints) {
        this.tour.keyPoints.forEach(kp => {
          const marker = L.marker([kp.latitude, kp.longitude], { keyPointId: kp.id } as L.MarkerOptions)
            .addTo(this.map).bindPopup(`<b>${kp.name}</b>`);
          this.keyPointMarkers.push(marker);
        });
      }
      
      this.updateTouristVisuals();
      
      if (this.tour.keyPoints && this.tour.keyPoints.length > 0) {
          const tourBounds = L.latLngBounds(this.tour.keyPoints.map(kp => [kp.latitude, kp.longitude]));
          if (this.touristPosition) tourBounds.extend([this.touristPosition.Latitude, this.touristPosition.Longitude]);
          this.map.fitBounds(tourBounds, { padding: [50, 50] });
      }

      this.map.on('click', (e: any) => {
        if (!this.tourExecution) {
          const { lat, lng } = e.latlng;
          this.tourService.updateTouristPosition({ latitude: lat, longitude: lng }).subscribe(
            updatedPosition => {
              this.touristPosition = updatedPosition;
              this.updateTouristVisuals();
            }
          );
        }
      });
    }
  }
  
  private updateTouristVisuals(isSimulating: boolean = false): void {
    if (!this.map || !this.touristPosition || !this.leaflet) return;
    
    const L = this.leaflet;
    const touristLatLng = L.latLng(this.touristPosition.Latitude, this.touristPosition.Longitude);
    
    if (!this.touristMarker) {
      const touristIcon = L.icon({
        iconUrl: 'assets/images/default-avatar.png',
        iconSize: [40, 40], iconAnchor: [20, 40], popupAnchor: [0, -40],
        className: 'tourist-marker'
      });
      this.touristMarker = L.marker(touristLatLng, { icon: touristIcon, zIndexOffset: 1000 }).addTo(this.map);
    } else {
      this.touristMarker.setLatLng(touristLatLng);
    }
    
    if (this.tour?.keyPoints && this.tour.keyPoints.length > 0) {
      const startPointLatLng = L.latLng(this.tour.keyPoints[0].latitude, this.tour.keyPoints[0].longitude);
      
      if (this.routeToStartLine) {
        this.routeToStartLine.setWaypoints([touristLatLng, startPointLatLng]);
      } else if (!this.tourExecution) {
        this.routeToStartLine = (L as any).Routing.control({
          waypoints: [touristLatLng, startPointLatLng],
          addWaypoints: false, fitSelectedRoutes: false, show: false,
          lineOptions: { styles: [{ color: '#dc3545', weight: 4, opacity: 0.8, dashArray: '10, 10' }] },
          createMarker: () => null
        }).addTo(this.map);

        this.routeToStartLine.on('routesfound', (e: any) => {
            if (e.routes && e.routes.length > 0) {
                this.toStartRoutePoints = e.routes[0].coordinates;
            }
        });
      }
    }
  }
}