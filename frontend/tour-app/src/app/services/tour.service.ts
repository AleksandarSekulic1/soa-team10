import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable, of } from 'rxjs'; // <-- 'of' je dodat ovde
import { catchError } from 'rxjs/operators'; // <-- Dodat je novi import za catchError
import { Tour, TourKeyPoint, TourTransport } from '../models/tour.model';
import { TouristPosition } from '../models/tourist-position.model';
import { TourExecution } from '../models/tour-execution.model';

@Injectable({
  providedIn: 'root'
})
export class TourService {
  private apiUrl = 'http://localhost:8083/api/tours';
// --- IZMENA OVE LINIJE ---
  private positionApiUrl = 'http://localhost:8085/api/tourist-position';
  private executionApiUrl = 'http://localhost:8085/api/tour-executions';
  constructor(private http: HttpClient) { }

  createTour(tourData: any): Observable<any> {
    return this.http.post<any>(this.apiUrl, tourData);
  }

  // NOVA METODA
  getMyTours(): Observable<any[]> {
    // AuthInterceptor će automatski dodati JWT token
    return this.http.get<any[]>(`${this.apiUrl}/my-tours`);
  }

  getAllTours(): Observable<any[]> {
    return this.http.get<any[]>(this.apiUrl);
  }

  // NOVA METODA
  addReview(tourId: string, reviewData: any): Observable<any> {
    return this.http.post<any>(`${this.apiUrl}/${tourId}/reviews`, reviewData);
  }

  getTourById(id: string): Observable<Tour> {
    return this.http.get<Tour>(`${this.apiUrl}/${id}`);
  }

  addKeyPoint(tourId: string, keyPointData: Partial<TourKeyPoint>): Observable<TourKeyPoint> {
    return this.http.post<TourKeyPoint>(`${this.apiUrl}/${tourId}/keypoints`, keyPointData);
  }

  updateKeyPoint(tourId: string, keypointId: string, keyPointData: TourKeyPoint): Observable<TourKeyPoint> {
    return this.http.put<TourKeyPoint>(`${this.apiUrl}/${tourId}/keypoints/${keypointId}`, keyPointData);
  }

  deleteKeyPoint(tourId: string, keypointId: string): Observable<any> {
    return this.http.delete<any>(`${this.apiUrl}/${tourId}/keypoints/${keypointId}`);
  }

  getTouristPosition(): Observable<TouristPosition> {
    return this.http.get<TouristPosition>(this.positionApiUrl);
  }

  updateTouristPosition(data: { latitude: number, longitude: number }): Observable<TouristPosition> {
    return this.http.post<TouristPosition>(this.positionApiUrl, data);
  }

  getPublishedTours(): Observable<Tour[]> {
    return this.http.get<Tour[]>(`${this.apiUrl}/published`);
  }

  addTransportInfo(tourId: string, transportInfo: TourTransport[]): Observable<Tour> {
    return this.http.post<Tour>(`${this.apiUrl}/${tourId}/transport-info`, transportInfo);
  }

  publishTour(tourId: string): Observable<Tour> {
    return this.http.post<Tour>(`${this.apiUrl}/${tourId}/publish`, {}); // Šaljemo prazan body
  }

  archiveTour(tourId: string): Observable<Tour> {
    return this.http.post<Tour>(`${this.apiUrl}/${tourId}/archive`, {});
  }

  reactivateTour(tourId: string): Observable<Tour> {
    return this.http.post<Tour>(`${this.apiUrl}/${tourId}/reactivate`, {});
  }

  startTour(tourId: string): Observable<TourExecution> {
    return this.http.post<TourExecution>(`${this.executionApiUrl}/start/${tourId}`, {});
  }

  checkPosition(executionId: string, position: { latitude: number, longitude: number }): Observable<TourExecution> {
    return this.http.post<TourExecution>(`${this.executionApiUrl}/check-position`, position);
  }
  
  completeTour(executionId: string): Observable<TourExecution> {
    return this.http.post<TourExecution>(`${this.executionApiUrl}/${executionId}/complete`, {});
  }
  
  abandonTour(executionId: string): Observable<TourExecution> {
    return this.http.post<TourExecution>(`${this.executionApiUrl}/${executionId}/abandon`, {});
  }

  getActiveExecutionForUser(): Observable<TourExecution | null> {
    return this.http.get<TourExecution>(`${this.executionApiUrl}/active`).pipe(
      catchError(() => of(null)) // Ako backend vrati 404, mi vratimo null
    );
  }

  getArchivedTours(): Observable<Tour[]> {
    return this.http.get<Tour[]>(`${this.apiUrl}/archived`);
  }
}
