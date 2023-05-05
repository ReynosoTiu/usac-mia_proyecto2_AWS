import { Injectable } from '@angular/core';
import { Subject } from 'rxjs';

@Injectable({
  providedIn: 'root'
})
export class HeaderService {

  constructor() { }

  private enviarLogin = new Subject<boolean>();
  obtenerLogin = this.enviarLogin.asObservable();

  enviarValorLogin(value: boolean) {
    this.enviarLogin.next(value);
  }

}
