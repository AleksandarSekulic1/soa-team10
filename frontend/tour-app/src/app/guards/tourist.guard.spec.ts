import { TestBed } from '@angular/core/testing';
import { CanActivateFn } from '@angular/router';

import { touristGuard } from './tourist.guard';

describe('touristGuard', () => {
  const executeGuard: CanActivateFn = (...guardParameters) => 
      TestBed.runInInjectionContext(() => touristGuard(...guardParameters));

  beforeEach(() => {
    TestBed.configureTestingModule({});
  });

  it('should be created', () => {
    expect(executeGuard).toBeTruthy();
  });
});
