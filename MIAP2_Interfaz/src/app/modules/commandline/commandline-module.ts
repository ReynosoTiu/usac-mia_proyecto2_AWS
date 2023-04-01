import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';

import { CommandlineRoutingModule } from './commandline-routing.module';
import { CommandlineComponent } from './commandline/commandline.component';



@NgModule({
  declarations: [
    CommandlineComponent
  ],
  
  imports: [
    CommonModule,
    CommandlineRoutingModule
  ]
})
export class CommandlineModule { }
