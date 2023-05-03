import { Component, OnInit } from '@angular/core';
import { GtoastrService } from 'src/app/services/gtoastr.service';

@Component({
  selector: 'app-login',
  templateUrl: './login.component.html',
  styleUrls: ['./login.component.scss']
})
export class LoginComponent implements OnInit {

  id_particion = '';
  usuario = '';
  password = '';

  constructor(private toastr: GtoastrService) { }

  ngOnInit(): void {
  }

  ingresar() {
    if (!this.id_particion) {
      this.toastr.error("Ingrese el campo 'ID Partición'");
      return;
    }

    if (!this.usuario) {
      this.toastr.error("Ingrese el campo 'Usuario'");
      return;
    }

    if (!this.password) {
      this.toastr.error("Ingrese el campo 'password'");
      return;
    }
    console.log(this.id_particion, this.usuario, this.password);
  }
}
