import { ComponentFixture, TestBed } from '@angular/core/testing';

import { KeypointFormComponent } from './keypoint-form.component';

describe('KeypointFormComponent', () => {
  let component: KeypointFormComponent;
  let fixture: ComponentFixture<KeypointFormComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [KeypointFormComponent]
    })
    .compileComponents();

    fixture = TestBed.createComponent(KeypointFormComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
