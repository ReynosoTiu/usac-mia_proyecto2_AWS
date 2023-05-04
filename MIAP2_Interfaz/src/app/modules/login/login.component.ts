import { Component, OnInit } from '@angular/core';
import { GtoastrService } from 'src/app/services/gtoastr.service';
import { CommandlineService } from '../commandline/commandline.service';
import { Router } from '@angular/router';
import { LoginService } from './login.service';
import { HeaderService } from 'src/app/services/app/header.service';

@Component({
  selector: 'app-login',
  templateUrl: './login.component.html',
  styleUrls: ['./login.component.scss']
})
export class LoginComponent implements OnInit {

  id_particion = '';
  usuario = '';
  password = '';

  constructor(private toastr: GtoastrService,
    private command: CommandlineService,
    private loginS: LoginService,
    private router: Router,
    private headerS: HeaderService) { 
    }

  ngOnInit(): void {
    let logeado = localStorage.getItem("login");
    if(logeado == "true"){
      this.router.navigate(["reports"]);
    }
  }

  ngAfterViewInit(){
    setTimeout(() => {
      this.headerS.enviarValorLogin(true);
    }, 100);
  }

  ingresar() {
    this.id_particion = this.id_particion.trim();
    this.usuario = this.usuario.trim();
    this.password = this.password.trim();

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
    let comando = `login >pwd=${this.password} >user=${this.usuario} >id=${this.id_particion}`;
    comando = comando.trim();
    this.command.enviarContenidoEEA(comando)
    .then(async (res)=>{
      if(res.Mensaje == "Un usuario ya se encuentra logueado"){
        this.toastr.warning("Se cerró la sesión anterior, intente nuevamente");
        await this.loginS.cerrarSesion();
      } else if(res.Mensaje == "INICIO DE SESION EXITOSO"){
        localStorage.setItem("login", "true");
        this.headerS.enviarValorLogin(true);
        this.router.navigate(["reports"]);
      }else{
        this.toastr.error(res.Mensaje);
      }
    })
    .catch();
  }
}
