import { inject } from '@angular/core';
import { CanActivateFn, Router } from '@angular/router';
import { AuthService } from '../services/auth.service';
import { map, take } from 'rxjs/operators';

export const authGuard: CanActivateFn = (route, state) => {
  const authService = inject(AuthService);
  const router = inject(Router);

  console.log(`[AuthGuard] Provera pristupa za rutu: ${state.url}`); // <-- LOG 1

  return authService.isLoggedIn$.pipe(
    take(1),
    map(isLoggedIn => {
      if (isLoggedIn) {
        console.log('[AuthGuard] Pristup DOZVOLJEN.'); // <-- LOG 2
        return true;
      } else {
        console.log('[AuthGuard] Pristup ZABRANJEN. Preusmeravanje na /login.'); // <-- LOG 3
        router.navigate(['/login']);
        return false;
      }
    })
  );
};