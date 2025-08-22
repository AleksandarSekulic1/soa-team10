import { Routes } from '@angular/router';
import { RegistrationComponent } from './pages/registration/registration.component';
import { LoginComponent } from './pages/login/login.component';
import { UserListComponent } from './pages/user-list/user-list.component';

export const routes: Routes = [
  { path: 'register', component: RegistrationComponent },
  { path: 'login', component: LoginComponent },
  { path: 'users', component: UserListComponent },
  { path: '', redirectTo: '/login', pathMatch: 'full' } // Podrazumevana ruta je sada login
];
