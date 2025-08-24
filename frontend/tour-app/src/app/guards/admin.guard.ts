import { inject } from '@angular/core';
import { CanActivateFn, Router } from '@angular/router';
import { AuthService } from '../services/auth.service';
import { map, take } from 'rxjs/operators';

export const adminGuard: CanActivateFn = (route, state) => {
  const authService = inject(AuthService);
  const router = inject(Router);

  return authService.isLoggedIn$.pipe(
    take(1),
    map(isLoggedIn => {
      // Prvo proveravamo da li je korisnik uopšte ulogovan
      if (isLoggedIn && authService.isAdmin()) {
        // Ako jeste ulogovan I ako je admin, dozvoli pristup
        return true;
      } else {
        // Ako nije ulogovan ILI nije admin, preusmeri ga
        console.error("Pristup odbijen. Potrebna je uloga administratora.");
        router.navigate(['/home']); // Preusmeravamo na home ili login
        return false;
      }
    })
  );
};
