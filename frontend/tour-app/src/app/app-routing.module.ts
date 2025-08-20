import { Routes } from '@angular/router';
import { RegistrationComponent } from './pages/registration/registration.component';

// Definišemo i eksportujemo samo konstantu sa rutama.
// Ne treba nam više @NgModule dekorator.
export const routes: Routes = [
  { path: 'register', component: RegistrationComponent },
  // Preusmeravanje na registraciju ako je putanja prazna
  { path: '', redirectTo: '/register', pathMatch: 'full' }
];
