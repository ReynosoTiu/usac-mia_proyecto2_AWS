import { Component } from '@angular/core';
import { HtmlInputEvent } from './models/global.model';
import { Router } from '@angular/router';
import { HeaderService } from './services/app/header.service';

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.scss']
})
export class AppComponent {
  title = 'MIAP2_Interfaz';
  disableBtn = false;

  constructor(private routes: Router,
    private headerS: HeaderService) {
  }

  ngOnInit(): void {
    this.headerS.obtenerLogin.subscribe(res => {
      this.disableBtn = res;
    });
  }

  cambiar() {
    let valor = !this.disableBtn;
    if (valor) {
      this.headerS.enviarValorLogin(true);
      this.routes.navigate(['login'])
    } else {
      this.headerS.enviarValorLogin(false);
      this.routes.navigate(['commandline'])
    }
  }

  textoArchivoLeido = '';
  seleccionar(event: HtmlInputEvent) {
    if (event.target.files && event.target.files[0]) {
      let file: File = <File>event.target.files[0];
      if (file.name.split('.').pop() === 'eea') {
        const reader = new FileReader();
        reader.onload = () => {
          this.textoArchivoLeido = String(reader.result);
        }
        reader.readAsText(file);
      }
    }
  }
}
