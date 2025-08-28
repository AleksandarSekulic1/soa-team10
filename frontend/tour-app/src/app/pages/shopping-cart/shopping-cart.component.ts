import { Component, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { ShoppingCartService } from '../../services/shopping-cart.service';
import { Observable } from 'rxjs';

@Component({
  selector: 'app-shopping-cart',
  standalone: true,
  imports: [CommonModule],
  templateUrl: './shopping-cart.component.html',
  styleUrls: ['./shopping-cart.component.scss']
})
export class ShoppingCartComponent implements OnInit {
  cart$: Observable<any>;

  constructor(private shoppingCartService: ShoppingCartService) {
    this.cart$ = this.shoppingCartService.cart$;
  }

  ngOnInit(): void {
    this.shoppingCartService.getCart().subscribe();
  }

  removeItem(itemId: string): void {
    this.shoppingCartService.removeItemFromCart(itemId).subscribe({
      next: () => console.log(`Stavka ${itemId} uklonjena.`),
      error: (err) => console.error('Greška pri uklanjanju stavke:', err)
    });
  }

  checkout(): void {
    this.shoppingCartService.checkout().subscribe({
      next: (response) => {
        alert('Kupovina je uspešno obavljena!');
        console.log('Dobijeni tokeni:', response.tokens);
      },
      error: (err) => {
        alert('Došlo je do greške prilikom kupovine.');
        console.error(err);
      }
    });
  }
}
