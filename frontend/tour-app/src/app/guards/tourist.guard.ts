import { inject } from '@angular/core';
import { CanActivateFn, Router } from '@angular/router';
import { AuthService } from '../services/auth.service';
import { map, take } from 'rxjs/operators';

export const touristGuard: CanActivateFn = (route, state) => {
  const authService = inject(AuthService);
  const router = inject(Router);

  return authService.isLoggedIn$.pipe(
    take(1),
    map(isLoggedIn => {
      // Prvo proveravamo da li je korisnik uopšte ulogovan
      if (isLoggedIn && authService.isTourist()) {
        // Ako jeste ulogovan I ako je turista (ili admin), dozvoli pristup
        return true;
      } else {
        // Ako nije ulogovan ILI nije turista, preusmeri ga
        console.error("Pristup odbijen. Potrebna je uloga turiste.");
        router.navigate(['/home']); // Preusmeravamo na home ili login
        return false;
      }
    })
  );
};
