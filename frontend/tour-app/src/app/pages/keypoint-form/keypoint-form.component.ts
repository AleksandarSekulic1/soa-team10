// src/app/pages/keypoint-form/keypoint-form.component.ts

import { Component, OnInit, AfterViewInit, Inject, PLATFORM_ID } from '@angular/core';
import { isPlatformBrowser, CommonModule, DecimalPipe } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { ActivatedRoute, Router } from '@angular/router';
import { TourService } from '../../services/tour.service';
import { TourKeyPoint } from '../../models/tour.model';

// 1. Uvozimo SAMO slike, ne i celu 'leaflet' biblioteku
import iconUrl from 'leaflet/dist/images/marker-icon.png';
import iconRetinaUrl from 'leaflet/dist/images/marker-icon-2x.png';
import shadowUrl from 'leaflet/dist/images/marker-shadow.png';

@Component({
  selector: 'app-keypoint-form',
  standalone: true,
  imports: [CommonModule, FormsModule, DecimalPipe],
  templateUrl: './keypoint-form.component.html',
  styleUrls: ['./keypoint-form.component.scss']
})
export class KeypointFormComponent implements OnInit, AfterViewInit {
  tourId: string | null = null;
  keyPoint: Partial<TourKeyPoint> = {
    name: '',
    description: '',
    latitude: 44.787197, // Početna lokacija: Beograd
    longitude: 20.457273,
    imageUrl: ''
  };

  availableImages: string[] = [
    'assets/images/default-avatar.png', 'assets/images/men2.png', 'assets/images/men3.png',
    'assets/images/men4.png', 'assets/images/men5.png', 'assets/images/women1.png',
    'assets/images/women2.png', 'assets/images/women3.png', 'assets/images/women4.png',
    'assets/images/women5.png'
  ];

  private map: any;
  private marker: any;

  constructor(
    private route: ActivatedRoute,
    private router: Router,
    private tourService: TourService,
    @Inject(PLATFORM_ID) private platformId: Object
  ) {}

  ngOnInit(): void {
    this.tourId = this.route.snapshot.paramMap.get('tourId');
  }

  ngAfterViewInit(): void {
    if (isPlatformBrowser(this.platformId)) {
      // 2. Vraćamo se na dinamički import za SAMO Leaflet biblioteku
      import('leaflet').then(L => {
        // 3. Podešavanje ikonice sada ide UNUTAR ovog bloka
        const iconDefault = L.icon({
          iconUrl,
          iconRetinaUrl,
          shadowUrl,
          iconSize: [25, 41],
          iconAnchor: [12, 41],
          popupAnchor: [1, -34],
          shadowSize: [41, 41]
        });
        L.Marker.prototype.options.icon = iconDefault;

        this.initMap(L);
      });
    }
  }

  private initMap(L: any): void {
    this.map = L.map('map').setView([this.keyPoint.latitude!, this.keyPoint.longitude!], 13);
    L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
      maxZoom: 18,
      attribution: '© OpenStreetMap'
    }).addTo(this.map);

    this.marker = L.marker([this.keyPoint.latitude!, this.keyPoint.longitude!]).addTo(this.map);

    this.map.on('click', (e: any) => {
      const { lat, lng } = e.latlng;
      this.keyPoint.latitude = lat;
      this.keyPoint.longitude = lng;
      this.marker.setLatLng(e.latlng);
    });
  }

  selectImage(imagePath: string): void {
    this.keyPoint.imageUrl = imagePath;
  }

  onSubmit(): void {
    if (!this.tourId) {
      console.error('Tour ID not found!');
      return;
    }
    this.tourService.addKeyPoint(this.tourId, this.keyPoint).subscribe({
      next: () => {
        this.router.navigate(['/my-tours']);
      },
      error: (err) => console.error('Error adding key point:', err)
    });
  }
}