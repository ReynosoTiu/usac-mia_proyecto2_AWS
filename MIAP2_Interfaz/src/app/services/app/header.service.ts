import { Injectable } from '@angular/core';

@Injectable({
  providedIn: 'root'
})
export class HeaderService {

  constructor() { }

  login = false;

  cambiarLogin(valor: boolean){
    this.login = valor;
  }

  getLogin(){
    return this.login;
  }

}
