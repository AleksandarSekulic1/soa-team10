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

  checkout(): void {
    alert('Checkout funkcionalnost će biti implementirana uskoro!');
  }
}
