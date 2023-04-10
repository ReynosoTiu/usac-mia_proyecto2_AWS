import { Component, Input, OnInit } from '@angular/core';
import { HtmlInputEvent } from 'src/app/models/global.model';
import { CommandlineService } from '../commandline.service';
import { GtoastrService } from 'src/app/services/gtoastr.service';
import { SweetalertService } from 'src/app/services/sweetalert.service';
@Component({
  selector: 'app-commandline',
  templateUrl: './commandline.component.html',
  styleUrls: ['./commandline.component.scss']
})
export class CommandlineComponent implements OnInit {

  textoArchivoLeido = "";
  textoArchivoResult = "";
  disableBtn = false;

  constructor(
    private commandlineService: CommandlineService,
    private toastService: GtoastrService,
    private sweet: SweetalertService
  ) { }

  ngOnInit(): void {
  }

  cambiar() {
    this.disableBtn = !this.disableBtn;
  }

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

  ejecutar() {
    let lineasComando = this.textoArchivoLeido.split("\n");
    for (let comando of lineasComando) {
      let c = comando.trim();
      if (c) {
        if (c.length > 5) {
          if (c.slice(0, 6) == "rmdisk") {
            this.sweet.confirmAction("Confirmar", "¿Desea elimianar el disco?")
              .then(res => {
                this.textoArchivoResult += "Desea eliminar el disco";
                if (res) {
                  this.textoArchivoResult += "SI";
                  this.commandlineService.enviarContenidoEEA(c)
                    .then(res => {
                      this.textoArchivoResult += res;
                    });
                }
                this.textoArchivoResult += "NO";
              })
          } else {
            this.commandlineService.enviarContenidoEEA(c)
              .then(res => {
                this.textoArchivoResult += res;
              });
          }
        }
      }
    }
  }

  get() {
    this.commandlineService.consumirInicial();
  }

  modificaTexto(event: any) {
    this.textoArchivoLeido = event.target.value;
  }
}
