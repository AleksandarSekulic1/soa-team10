import { ComponentFixture, TestBed } from '@angular/core/testing';

import { PositionSimulatorComponent } from './position-simulator.component';

describe('PositionSimulatorComponent', () => {
  let component: PositionSimulatorComponent;
  let fixture: ComponentFixture<PositionSimulatorComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [PositionSimulatorComponent]
    })
    .compileComponents();

    fixture = TestBed.createComponent(PositionSimulatorComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
