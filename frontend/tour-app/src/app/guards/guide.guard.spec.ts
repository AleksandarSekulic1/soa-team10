import { TestBed } from '@angular/core/testing';
import { CanActivateFn } from '@angular/router';

import { guideGuard } from './guide.guard';

describe('guideGuard', () => {
  const executeGuard: CanActivateFn = (...guardParameters) => 
      TestBed.runInInjectionContext(() => guideGuard(...guardParameters));

  beforeEach(() => {
    TestBed.configureTestingModule({});
  });

  it('should be created', () => {
    expect(executeGuard).toBeTruthy();
  });
});
