import { Component, OnInit } from '@angular/core';

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
}
