import { Component, OnInit } from '@angular/core';
import { CommonModule, TitleCasePipe, DecimalPipe } from '@angular/common';
import { TourService } from '../../services/tour.service';
import { Tour } from '../../models/tour.model';
import { RouterLink } from '@angular/router';

@Component({
  selector: 'app-home',
  standalone: true,
  imports: [CommonModule, RouterLink, DecimalPipe],
  templateUrl: './home.component.html',
  styleUrls: ['./home.component.scss']
})
export class HomeComponent implements OnInit {
  publishedTours: Tour[] = [];
  isLoading = true;

  constructor(private tourService: TourService) {}

  ngOnInit(): void {
    // Pozivamo metodu koja vraća samo objavljene ture
    this.tourService.getPublishedTours().subscribe({
      next: (tours) => {
        this.publishedTours = tours;
        this.isLoading = false;
        console.log('Objavljene ture:', tours);
      },
      error: (err) => {
        console.error('Greška pri preuzimanju objavljenih tura:', err);
        this.isLoading = false;
      }
    });
  }
}