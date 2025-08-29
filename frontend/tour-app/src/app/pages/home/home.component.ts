import { Component, OnInit } from '@angular/core';
import { CommonModule, DecimalPipe } from '@angular/common';
import { TourService } from '../../services/tour.service';
import { Tour } from '../../models/tour.model';
import { RouterLink } from '@angular/router';
import { forkJoin } from 'rxjs'; // <-- 1. Uvezite forkJoin

@Component({
  selector: 'app-home',
  standalone: true,
  imports: [CommonModule, RouterLink, DecimalPipe],
  templateUrl: './home.component.html',
  styleUrls: ['./home.component.scss']
})
export class HomeComponent implements OnInit {
  publishedTours: Tour[] = [];
  archivedTours: Tour[] = []; // <-- 2. Dodajte niz za arhivirane ture
  isLoading = true;

  constructor(private tourService: TourService) {}

  ngOnInit(): void {
    this.isLoading = true;

    // 3. Koristimo forkJoin da istovremeno dobavimo i objavljene i arhivirane ture
    forkJoin({
      published: this.tourService.getPublishedTours(),
      archived: this.tourService.getArchivedTours()
    }).subscribe({
      next: ({ published, archived }) => {
        this.publishedTours = published;
        this.archivedTours = archived;
        this.isLoading = false;
        console.log('Objavljene ture:', this.publishedTours);
        console.log('Arhivirane ture:', this.archivedTours);
      },
      error: (err) => {
        console.error('Greška pri preuzimanju tura:', err);
        this.isLoading = false;
      }
    });
  }
}