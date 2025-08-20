import { bootstrapApplication } from '@angular/platform-browser';
import { provideRouter } from '@angular/router';
import { importProvidersFrom } from '@angular/core';
import { HttpClientModule } from '@angular/common/http';
import { FormsModule } from '@angular/forms';

import { AppComponent } from './app/app.component';
import { routes } from './app/app-routing.module'; // Uvozimo samo rute

bootstrapApplication(AppComponent, {
  providers: [
    provideRouter(routes), // Obezbeđujemo rute za celu aplikaciju
    importProvidersFrom(HttpClientModule, FormsModule) // Obezbeđujemo module globalno
  ]
})
  .catch(err => console.error(err));
