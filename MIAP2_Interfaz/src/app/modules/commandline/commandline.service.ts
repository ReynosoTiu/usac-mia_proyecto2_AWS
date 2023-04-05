import { Injectable } from '@angular/core';
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

  enviarContenidoEEA(contenido: string) {
    return this.http.post('Carga', contenido);
    
  }
}
