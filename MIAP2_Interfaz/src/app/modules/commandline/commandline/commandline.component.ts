import { Component, Input, OnInit } from '@angular/core';
import { HtmlInputEvent } from 'src/app/models/global.model';
import { CommandlineService } from '../commandline.service';
import { GtoastrService } from 'src/app/services/gtoastr.service';

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
    private toastService: GtoastrService
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
    this.commandlineService.enviarContenidoEEA(this.textoArchivoLeido)
    .then(res=>{
      this.textoArchivoResult += res;
    }).catch(err=>{
      this.toastService.error("Ha ocurrido un error");
    });
  }

  get() {
    this.commandlineService.consumirInicial();
  }

  modificaTexto(event:any){
    this.textoArchivoLeido = event.target.value;
  }
}
