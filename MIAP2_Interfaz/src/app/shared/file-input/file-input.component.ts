import { Component, EventEmitter, Input, OnInit, Output } from '@angular/core';

@Component({
  selector: 'app-file-input',
  templateUrl: './file-input.component.html',
  styleUrls: ['./file-input.component.scss']
})
export class FileInputComponent implements OnInit {

  @Input() accept = '';
  @Output() changee = new EventEmitter<any>();

  constructor() { }

  ngOnInit(): void {
  }

  seleccionar(event:any){
    this.changee.emit(event);
  }

}
