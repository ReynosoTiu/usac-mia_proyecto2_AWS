import { Component, OnInit } from '@angular/core';
import { Router } from '@angular/router';
import { Reporte } from 'src/app/models/servicios.model';
import { saveAs } from 'file-saver'
import { LoginService } from '../login/login.service';
import { HeaderService } from 'src/app/services/app/header.service';
import { DomSanitizer, SafeResourceUrl } from '@angular/platform-browser';

@Component({
  selector: 'app-report',
  templateUrl: './report.component.html',
  styleUrls: ['./report.component.scss']
})
export class ReportComponent implements OnInit {

  reportes: Reporte[] = [];
  reporteActual: Reporte | undefined;
  generado = false;

  constructor(private route: Router,
    private loginS: LoginService,
    private headerS: HeaderService,
    private sanitizer: DomSanitizer
  ) {
    let grafosString = localStorage.getItem("reportes");
    if (grafosString) {
      this.reportes = JSON.parse(grafosString);
    }
  }

  ngOnInit(): void {
    let login = localStorage.getItem("login");
    if (login != "true") {
      this.route.navigate(["login"]);
    }
  }

  ngAfterViewInit() {
    setTimeout(() => {
      this.headerS.enviarValorLogin(true);
    }, 100);
  }

  cerrarSesion() {
    this.loginS.cerrarSesion(true).then();
  }

  imageData = "";
  tipoArchivo = ""
  extension = "";
  generar(grafo: Reporte) {
    this.reporteActual = grafo;
    let extension = "data:image/png;base64,";
    this.tipoArchivo = "imagen";
    switch (grafo.Extension.toLowerCase()) {
      case "jpg":
        extension = "data:image/jpeg;base64,"
        break;
      case "pdf":
        this.tipoArchivo = "pdf";
        extension = "data:application/pdf;base64,"
        break;
    }
    this.imageData = `${extension}${grafo.Data}`;
    this.generado = true;
  }

  pdfUrl(): SafeResourceUrl {
    const url = this.imageData;
    return this.sanitizer.bypassSecurityTrustResourceUrl(url);
  }

  abrir() {
    let base64Image = this.imageData;
    const byteString = atob(base64Image.split(',')[1]);
    const mimeString = base64Image.split(',')[0].split(':')[1].split(';')[0];
    const ab = new ArrayBuffer(byteString.length);
    const ia = new Uint8Array(ab);
    for (let i = 0; i < byteString.length; i++) {
      ia[i] = byteString.charCodeAt(i);
    }
    const blob = new Blob([ab], { type: mimeString });
    const blobUrl = URL.createObjectURL(blob);
    window.open(blobUrl, '_blank');
  }

  descargar() {
    let base64Image = this.imageData;
    const byteString = atob(base64Image.split(',')[1]);
    const mimeString = base64Image.split(',')[0].split(':')[1].split(';')[0];
    const ab = new ArrayBuffer(byteString.length);
    const ia = new Uint8Array(ab);
    for (let i = 0; i < byteString.length; i++) {
      ia[i] = byteString.charCodeAt(i);
    }
    const blob = new Blob([ab], { type: mimeString });
    console.log(this.reporteActual?.NombreSave);
    saveAs(blob, this.reporteActual?.NombreSave);
  }

  regresar() {
    this.reporteActual = undefined;
    this.generado = false;
  }
}
