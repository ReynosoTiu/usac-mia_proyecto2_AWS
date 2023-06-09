import { Injectable } from '@angular/core';
import { Respuesta } from 'src/app/models/servicios.model';
import { HttpRequestService } from 'src/app/services/http-request.service';

@Injectable({
  providedIn: 'root'
})
export class CommandlineService {

  constructor(private http: HttpRequestService) { }

  consumirInicial() {
    this.http.get('').then(res => {
      console.log(res);
    }).catch(err => {
      console.log(err);
    });
  }

  enviarContenidoEEA(contenido: string): Promise<Respuesta> {
    return this.http.post('Carga', contenido);
  }
}
