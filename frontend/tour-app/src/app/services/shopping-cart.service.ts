import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { BehaviorSubject, Observable, tap } from 'rxjs';
import { AuthService } from './auth.service';

@Injectable({
  providedIn: 'root'
})
export class ShoppingCartService {
  private apiUrl = 'http://localhost:8084/api/shopping-cart';

  private cartSubject = new BehaviorSubject<any>(null);
  cart$ = this.cartSubject.asObservable();

  constructor(private http: HttpClient, private authService: AuthService) { }

  getCart(): Observable<any> {
    const touristUsername = this.authService.getUsername();
    if (!touristUsername) {
      throw new Error('Korisnik nije ulogovan');
    }
    return this.http.get<any>(`${this.apiUrl}/${touristUsername}`).pipe(
      tap(cart => this.cartSubject.next(cart))
    );
  }

  addItemToCart(tour: any): Observable<any> {
    const touristUsername = this.authService.getUsername();
    if (!touristUsername) {
      throw new Error('Korisnik nije ulogovan');
    }

    const orderItem = {
      tourName: tour.name,
      price: tour.price > 0 ? tour.price : 50.0, // Podrazumevana cena ako je 0
      tourId: tour.id
    };

    return this.http.post<any>(`${this.apiUrl}/${touristUsername}/items`, orderItem).pipe(
      tap(cart => this.cartSubject.next(cart))
    );
  }
}
