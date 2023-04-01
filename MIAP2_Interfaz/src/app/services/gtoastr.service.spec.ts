import { TestBed } from '@angular/core/testing';

import { GtoastrService } from './gtoastr.service';

describe('ToastrService', () => {
  let service: GtoastrService;

  beforeEach(() => {
    TestBed.configureTestingModule({});
    service = TestBed.inject(GtoastrService);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });
});
