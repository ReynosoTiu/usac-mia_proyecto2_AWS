import { Component, OnInit } from '@angular/core';
import { Router } from '@angular/router';
import { graphviz } from 'd3-graphviz';
import { Reporte } from 'src/app/models/servicios.model';
import { saveAs } from 'file-saver'
import { LoginService } from '../login/login.service';
import { HeaderService } from 'src/app/services/app/header.service';
@Component({
  selector: 'app-report',
  templateUrl: './report.component.html',
  styleUrls: ['./report.component.scss']
})
export class ReportComponent implements OnInit {

  reportes: Reporte[] = [];
  reporteActual:Reporte | undefined;
  generado = false;

  constructor(private route: Router,
    private loginS: LoginService,
    private headerS: HeaderService) {
    let grafosString = localStorage.getItem("reportes");
    if (grafosString) {
      this.reportes = JSON.parse(grafosString);
    }
  }

  ngOnInit(): void {
    this.headerS.cambiarLogin(true);
    let login = localStorage.getItem("login");
    if (login != "true") {
      this.route.navigate(["login"]);
    }
  }

  cerrarSesion(){
    this.loginS.cerrarSesion(true).then();
  }

  generar(grafo: Reporte) {
    this.reporteActual = grafo;
    console.log(grafo.Grafo);
    graphviz("#graph")
      .dot(grafo.Grafo)
      .render();

    setTimeout(() => {
      let svg_all = document.getElementsByTagName("svg");
      let svg = svg_all[0];
      svg.setAttribute("width", "100%");
      this.generado = true;
    }, 100);
  }

  abrir() {
    let svg_all = document.getElementsByTagName("svg");
    let svg = svg_all[0];
    const svgCode = new XMLSerializer().serializeToString(svg);
    const blob = new Blob([svgCode], { type: 'image/svg+xml' });
    var url = window.URL.createObjectURL(blob);
    window.open(url);
  }

  descargar() {
    let svg_all = document.getElementsByTagName("svg");
    let svg = svg_all[0];
    const svgCode = new XMLSerializer().serializeToString(svg);
    const blob = new Blob([svgCode], { type: 'image/svg+xml' });
    saveAs(blob, this.reporteActual?.NombreSave);
  }

  regresar(){
    this.reporteActual = undefined;
    this.generado = false;
  }
}
