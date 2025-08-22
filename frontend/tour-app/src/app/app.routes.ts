import { Routes } from '@angular/router';
// Uvozimo sve potrebne komponente
import { RegistrationComponent } from './pages/registration/registration.component';
import { LoginComponent } from './pages/login/login.component';
import { UserListComponent } from './pages/user-list/user-list.component';
import { HomeComponent } from './pages/home/home.component';
import { LayoutComponent } from './layout/layout.component';
import { authGuard } from './guards/auth.guard'; // Uvezemo guard
import { ProfileComponent } from './pages/profile/profile.component';

export const routes: Routes = [
  // --- Rute bez layout-a (bez navigacionog bara) ---
  { path: 'register', component: RegistrationComponent },
  { path: 'login', component: LoginComponent },

  // --- Ruta koja koristi LayoutComponent kao okvir (sve unutar nje imaće nav bar) ---
  {
    path: '', // Prazan path znači da će se primeniti na sve 'child' rute
    canActivate: [authGuard], // <-- ČUVAR JE POSTAVLJEN OVDE
    component: LayoutComponent,
    children: [
      { path: 'home', component: HomeComponent },
      { path: 'users', component: UserListComponent },
      // Sve buduće stranice koje zahtevaju login i navigaciju dodajte ovde
      { path: 'profile', component: ProfileComponent }, // <-- DODATA RUTA

      // Ako korisnik dođe na praznu putanju unutar layout-a (npr. nakon logina),
      // preusmeri ga na 'home' stranicu.
      { path: '', redirectTo: 'home', pathMatch: 'full' }
    ]
  },

  // Podrazumevana ruta ako ništa drugo ne odgovara - vodi na login
  { path: '**', redirectTo: '/login', pathMatch: 'full' }
];