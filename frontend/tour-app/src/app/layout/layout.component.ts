import { Component } from '@angular/core';
import { RouterModule } from '@angular/router'; // 1. Uvezite RouterModule
import { NavbarComponent } from '../components/navbar/navbar.component'; // 2. Uvezite NavbarComponent

@Component({
  selector: 'app-layout',
  standalone: true, // 3. Označite komponentu kao standalone
  imports: [
    RouterModule,    // 4. Dodajte RouterModule u imports niz
    NavbarComponent    // 5. Dodajte NavbarComponent u imports niz
  ],
  templateUrl: './layout.component.html',
  styleUrls: ['./layout.component.scss']
})
export class LayoutComponent {

}