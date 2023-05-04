import { Component, ElementRef, Input, OnInit, ViewChild } from '@angular/core';
import { HtmlInputEvent } from 'src/app/models/global.model';
import { CommandlineService } from '../commandline.service';
import { SweetalertService } from 'src/app/services/sweetalert.service';
import { ReportService } from '../../report/report.service';
import { HeaderService } from 'src/app/services/app/header.service';
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
    private sweet: SweetalertService,
    private reporteS: ReportService,
    private headerS: HeaderService
  ) { 
  }

  ngOnInit(): void {
    let item = localStorage.getItem("reportes");
    if(!item){
      localStorage.setItem("reportes", JSON.stringify([]));
    }
  }

  ngAfterViewInit() {
    this.textAr2.nativeElement.value = '';
    this.textAr.nativeElement.value = '';
    setTimeout(() => {
      this.headerS.enviarValorLogin(false);
    }, 100);
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
                      this.textoArchivoResult += res.Mensaje + "\n";
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
                if(res.Mensaje == "SESION CERRADA EXITOSAMENTE"){
                  localStorage.setItem("login", "false");
                }
                if(res.Mensaje == "INICIO DE SESION EXITOSO"){
                  localStorage.setItem("login", "true");
                }
                if(res.Tipo == 2){
                  this.reporteS.guardarReporte(res);
                }
                this.textoArchivoResult += res.Mensaje + "\n";
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
