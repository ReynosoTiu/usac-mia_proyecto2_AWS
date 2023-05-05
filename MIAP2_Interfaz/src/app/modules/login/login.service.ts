import { Injectable } from '@angular/core';
import { CommandlineService } from '../commandline/commandline.service';
import { GtoastrService } from 'src/app/services/gtoastr.service';
import { HeaderService } from 'src/app/services/app/header.service';
import { Router } from '@angular/router';

@Injectable({
  providedIn: 'root'
})
export class LoginService {

  constructor(private command: CommandlineService,
    private toastr: GtoastrService,
    private headerS: HeaderService,
    private router: Router) { }

  async cerrarSesion(mensaje = false){
    await this.command.enviarContenidoEEA("logout")
    .then(
      res=>{
        if(mensaje){
          this.toastr.info(res.Mensaje);
        }
        localStorage.setItem("login", "false");
        this.headerS.enviarValorLogin(false);
        this.router.navigate(["login"]);
      }
    )
    .catch()
  }
}
