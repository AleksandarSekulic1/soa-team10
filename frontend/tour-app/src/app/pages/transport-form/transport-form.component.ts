import { Component, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { ActivatedRoute, Router } from '@angular/router';
import { TourService } from '../../services/tour.service';
import { TourTransport } from '../../models/tour.model';

@Component({
  selector: 'app-transport-form',
  standalone: true,
  imports: [CommonModule, FormsModule],
  templateUrl: './transport-form.component.html',
  styleUrls: ['./transport-form.component.scss']
})
export class TransportFormComponent implements OnInit {
  tourId: string | null = null;
  transportInfo: TourTransport[] = [
    { type: 'walking', timeInMinutes: 0 },
    { type: 'bicycle', timeInMinutes: 0 },
    { type: 'car', timeInMinutes: 0 }
  ];

  constructor(
    private route: ActivatedRoute,
    private router: Router,
    private tourService: TourService
  ) {}

  ngOnInit(): void {
    this.tourId = this.route.snapshot.paramMap.get('id');
    if (this.tourId) {
      // Opciono: Dobaviti postojeće informacije ako postoje
      this.tourService.getTourById(this.tourId).subscribe(tour => {
        if (tour.transportInfo && tour.transportInfo.length > 0) {
          this.transportInfo = tour.transportInfo;
        }
      });
    }
  }

  onSubmit(): void {
    if (!this.tourId) return;
    // Filtriramo samo one unose gde je vreme veće od 0
    const validTransportInfo = this.transportInfo.filter(t => t.timeInMinutes > 0);

    this.tourService.addTransportInfo(this.tourId, validTransportInfo).subscribe({
      next: () => {
        this.router.navigate(['/my-tours']);
      },
      error: (err) => {
        console.error('Error updating transport info:', err);
      }
    });
  }
}