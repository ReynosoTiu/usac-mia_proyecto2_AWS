import { Injectable } from '@angular/core';
import { Reporte, Respuesta } from 'src/app/models/servicios.model';

@Injectable({
  providedIn: 'root'
})
export class ReportService {

  constructor() { }

  guardarReporte(res: Respuesta) {
    let grafosString = localStorage.getItem("reportes");
    let grafos: Reporte[] = []
    if (grafosString) {
      grafos = JSON.parse(grafosString);
    }

    let nombres = res.Ruta.split("/");
    let nuevoGrafo: Reporte = {
      Data: res.Data,
      Ruta: res.Ruta,
      NombreSave: nombres[nombres.length - 1],
      Extension: res.Extension
    }

    console.log(nuevoGrafo);

    let encontrado = false;
    for (let item of grafos) {
      if (item.Ruta == nuevoGrafo.Ruta) {
        item.Data = nuevoGrafo.Data;
        item.NombreSave = nuevoGrafo.NombreSave;
        encontrado = true;
        break;
      }
    }
    if (!encontrado) {
      grafos.push(nuevoGrafo);
    }

    localStorage.setItem("reportes", JSON.stringify(grafos));
  }
}
