import { Routes } from '@angular/router';
// Uvozimo sve potrebne komponente
import { RegistrationComponent } from './pages/registration/registration.component';
import { LoginComponent } from './pages/login/login.component';
import { UserListComponent } from './pages/user-list/user-list.component';
import { HomeComponent } from './pages/home/home.component';
import { LayoutComponent } from './layout/layout.component';
import { ProfileComponent } from './pages/profile/profile.component';
import { TourCreationComponent } from './pages/tour-creation/tour-creation.component'; // <-- Uvozimo novu komponentu
import { TourListComponent } from './pages/tour-list/tour-list.component';
// Uvozimo sve guardove
import { authGuard } from './guards/auth.guard';
import { adminGuard } from './guards/admin.guard';
import { guideGuard } from './guards/guide.guard';
import { touristGuard } from './guards/tourist.guard';
import { MyToursComponent } from './pages/my-tours/my-tours.component';
export const routes: Routes = [
  // --- Rute bez layout-a (bez navigacionog bara) ---
  { path: 'register', component: RegistrationComponent },
  { path: 'login', component: LoginComponent },

  // --- Ruta koja koristi LayoutComponent kao okvir (sve unutar nje imaće nav bar) ---
  {
    path: '', // Prazan path znači da će se primeniti na sve 'child' rute
    canActivate: [authGuard], // <-- Glavni čuvar proverava da li je korisnik UOPŠTE ulogovan
    component: LayoutComponent,
    children: [
      { path: 'home', component: HomeComponent },
      {
        path: 'users',
        component: UserListComponent,
        canActivate: [adminGuard] // <-- Ovaj čuvar proverava da li je korisnik ADMIN
      },
      {
        path: 'create-tour',
        component: TourCreationComponent,
        canActivate: [guideGuard] // <-- Ovaj čuvar proverava da li je korisnik VODIČ
      },
      { path: 'profile', component: ProfileComponent },
      {
        path: 'my-tours',
        component: MyToursComponent,
        canActivate: [guideGuard] // Može i vodič da vidi svoje ture
      },
      {
        path: 'tours',
        component: TourListComponent,
        canActivate: [touristGuard] // Može i vodič da vidi svoje ture
      },
      // Ako korisnik dođe na praznu putanju unutar layout-a (npr. nakon logina),
      // preusmeri ga na 'home' stranicu.
      { path: '', redirectTo: 'home', pathMatch: 'full' }
    ]
  },

  // Podrazumevana ruta ako ništa drugo ne odgovara - vodi na login
  { path: '**', redirectTo: '/login', pathMatch: 'full' }
];
