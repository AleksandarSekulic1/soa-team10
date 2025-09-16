import { HttpInterceptorFn } from '@angular/common/http';

export const authInterceptor: HttpInterceptorFn = (req, next) => {
  const token = localStorage.getItem('jwt_token');

  console.log(`[AuthInterceptor] Presretnut zahtev za: ${req.url}`); // <-- DODAJTE OVAJ LOG

  if (token) {
    console.log('[AuthInterceptor] Token pronađen, dodajem Authorization heder.'); // <-- DODAJTE OVAJ LOG
    const cloned = req.clone({
      headers: req.headers.set('Authorization', `Bearer ${token}`)
    });
    return next(cloned);
  }

  console.log('[AuthInterceptor] Token nije pronađen, preskačem.'); // <-- DODAJTE OVAJ LOG
  return next(req);
};