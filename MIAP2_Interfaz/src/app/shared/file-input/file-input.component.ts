import { Component, ElementRef, EventEmitter, Input, OnInit, Output, ViewChild } from '@angular/core';

@Component({
  selector: 'app-file-input',
  templateUrl: './file-input.component.html',
  styleUrls: ['./file-input.component.scss']
})
export class FileInputComponent implements OnInit {

  @ViewChild('inputFile') inputFile!: ElementRef<any>;
  @Input() accept = '';
  @Output() changee = new EventEmitter<any>();

  constructor() { 
  }

  ngOnInit(): void {
  }

  abrirDialogo(){
    this.inputFile.nativeElement.value = '';
    this.inputFile.nativeElement.click();
  }

  seleccionar(event:any){
    this.changee.emit(event);
  }

}
