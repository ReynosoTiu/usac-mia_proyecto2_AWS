import { Component, ElementRef, Input, OnInit, ViewChild } from '@angular/core';
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

  @ViewChild('textAr') textAr!: ElementRef<any>;
  @ViewChild('textAr2') textAr2!: ElementRef<any>;

  constructor(
    private commandlineService: CommandlineService,
    private toastService: GtoastrService,
    private sweet: SweetalertService
  ) { }

  ngOnInit(): void {

  }

  ngAfterViewInit() {
    this.textAr2.nativeElement.value = '';
    this.textAr.nativeElement.value = '';
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
          this.textAr.nativeElement.value = reader.result;
        }
        reader.readAsText(file);
      }
    }
  }

  async ejecutar() {
    let lineasComando = this.textoArchivoLeido.split("\n");
    for (let comando of lineasComando) {
      let c = comando.trim().replace(/\n/g, "").replace(/\r/g, "");
      if (c) {
        if (c.length > 5) {
          if (c.slice(0, 6) == "rmdisk") {
            await this.sweet.confirmAction("Confirmar", "¿Desea elimianar el disco?")
              .then(async(res:any) => {
                this.textoArchivoResult += "Desea eliminar el disco";
                this.textAr2.nativeElement.value = this.textoArchivoResult;
                if (res) {
                  this.textoArchivoResult += " SI\n";
                  this.textAr2.nativeElement.value = this.textoArchivoResult;
                  await this.commandlineService.enviarContenidoEEA(c)
                    .then(res => {
                      this.textoArchivoResult += res + "\n";
                      this.textAr2.nativeElement.value = this.textoArchivoResult;
                    });
                } else {
                  this.textoArchivoResult += " NO\n";
                  this.textAr2.nativeElement.value = this.textoArchivoResult;
                }
              })
          } else {
            await this.commandlineService.enviarContenidoEEA(c)
              .then(res => {
                this.textoArchivoResult += res + "\n";
                this.textAr2.nativeElement.value = this.textoArchivoResult;
              });
          }
        }
      }
    }
  }

  modificaTexto(event: any) {
    this.textoArchivoLeido = event.target.value;
  }

  onKeypressEvent(event:any){
    return false;
  }
}
