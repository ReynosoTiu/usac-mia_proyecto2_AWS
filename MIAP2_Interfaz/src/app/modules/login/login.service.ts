import { Injectable } from '@angular/core';
import { CommandlineService } from '../commandline/commandline.service';
import { GtoastrService } from 'src/app/services/gtoastr.service';

@Injectable({
  providedIn: 'root'
})
export class LoginService {

  constructor(private command: CommandlineService,
    private toastr: GtoastrService,) { }

  async cerrarSesion(mensaje = false){
    await this.command.enviarContenidoEEA("logout")
    .then(
      res=>{
        if(mensaje){
          this.toastr.info(res.Mensaje);
        }
        localStorage.setItem("login", "false");
      }
    )
    .catch()
  }
}
