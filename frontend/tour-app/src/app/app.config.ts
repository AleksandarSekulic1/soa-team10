import { ApplicationConfig, provideZoneChangeDetection } from '@angular/core';
import { provideRouter } from '@angular/router';
import { routes } from './app.routes';
import { provideClientHydration } from '@angular/platform-browser';

// 1. Uvozimo 'withFetch' iz '@angular/common/http'
import { provideHttpClient, withInterceptors, withFetch } from '@angular/common/http';
import { authInterceptor } from './services/auth.interceptor';

export const appConfig: ApplicationConfig = {
  providers: [
    provideZoneChangeDetection({ eventCoalescing: true }),
    provideRouter(routes),
    provideClientHydration(),

    // 2. Dodajemo 'withFetch()' kao drugi argument u 'provideHttpClient'
    provideHttpClient(
      withInterceptors([authInterceptor]),
      withFetch()
    )
  ]
};