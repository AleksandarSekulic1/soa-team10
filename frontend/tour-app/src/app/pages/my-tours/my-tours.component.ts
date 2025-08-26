import { Component, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { TourService } from '../../services/tour.service';
import { Tour } from '../../models/tour.model'; // <-- Uvozimo model
import { RouterLink } from '@angular/router'; // <-- Uvozimo RouterLink za navigaciju


@Component({
  selector: 'app-my-tours',
  standalone: true,
  imports: [CommonModule,RouterLink],
  templateUrl: './my-tours.component.html',
  styleUrls: ['./my-tours.component.scss']
})
export class MyToursComponent implements OnInit {
  myTours: any[] = [];

  constructor(private tourService: TourService) {}

  ngOnInit(): void {
    this.tourService.getMyTours().subscribe({
      next: (tours) => {
        this.myTours = tours;
        console.log('Moje ture:', tours);
      },
      error: (err) => {
        console.error('Greška pri preuzimanju tura:', err);
      }
    });
  }
}
