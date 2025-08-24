import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';

@Injectable({
  providedIn: 'root'
})
export class TourService {
  private apiUrl = 'http://localhost:8083/api/tours';

  constructor(private http: HttpClient) { }

  createTour(tourData: any): Observable<any> {
    return this.http.post<any>(this.apiUrl, tourData);
  }

  // NOVA METODA
  getMyTours(): Observable<any[]> {
    // AuthInterceptor će automatski dodati JWT token
    return this.http.get<any[]>(`${this.apiUrl}/my-tours`);
  }
}
