import { Component } from '@angular/core';
import { HtmlInputEvent } from './models/global.model';
import { Router } from '@angular/router';

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.scss']
})
export class AppComponent {
  title = 'MIAP2_Interfaz';

  disableBtn = false;
  constructor(private routes: Router) { }

  ngOnInit(): void {
  }

  cambiar(){
    this.disableBtn = !this.disableBtn;
    if(this.disableBtn){
      this.routes.navigate(['login'])
    }else{
      this.routes.navigate(['commandline'])
    }
  }

  textoArchivoLeido = '';
  seleccionar(event: HtmlInputEvent){
    console.log(event);
    if (event.target.files && event.target.files[0]) {
      let file: File = <File>event.target.files[0];
      console.log(file.name);
      if(file.name.split('.').pop() === 'eea'){
        const reader = new FileReader();
        reader.onload = () => {
          this.textoArchivoLeido = String(reader.result);
        }
      reader.readAsText(file);
      }
    }
  }

  ejecutar(contenido:string){
    for(let linea of contenido.split('\n')){
      if(linea){
        console.log(linea);
      }
    }
  }
}
