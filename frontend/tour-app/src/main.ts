import { bootstrapApplication } from '@angular/platform-browser';
import { provideRouter } from '@angular/router';
import { importProvidersFrom } from '@angular/core';
import { HttpClientModule, provideHttpClient, withInterceptors } from '@angular/common/http';
import { FormsModule } from '@angular/forms';

import { AppComponent } from './app/app.component';
import { routes } from './app/app-routing.module';
import { authInterceptor } from './app/services/auth.interceptor'; // Uvezite interceptor

bootstrapApplication(AppComponent, {
  providers: [
    provideRouter(routes),
    // Obezbeđujemo HttpClient i registrujemo interceptor
    provideHttpClient(withInterceptors([authInterceptor])),
    importProvidersFrom(FormsModule) // FormsModule i dalje treba za ngModel
  ]
})
  .catch(err => console.error(err));
