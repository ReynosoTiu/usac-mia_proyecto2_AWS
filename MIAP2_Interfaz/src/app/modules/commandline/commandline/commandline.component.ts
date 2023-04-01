import { Component, OnInit } from '@angular/core';
import { HtmlInputEvent } from 'src/app/models/global.model';

@Component({
  selector: 'app-commandline',
  templateUrl: './commandline.component.html',
  styleUrls: ['./commandline.component.scss']
})
export class CommandlineComponent implements OnInit {

  disableBtn = false;
  constructor() { }

  ngOnInit(): void {
  }

  cambiar(){
    this.disableBtn = !this.disableBtn;
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