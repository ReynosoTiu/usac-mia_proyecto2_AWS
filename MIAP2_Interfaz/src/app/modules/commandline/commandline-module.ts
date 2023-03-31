import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';

import { CommandlineRoutingModule } from './commandline-routing.module';
//import { ConfirmButtonComponent } from 'src/app/shared/buttons/confirm-button/confirm-button.component';


@NgModule({
  declarations: [
    //ConfirmButtonComponent
  ],
  imports: [
    CommonModule,
    CommandlineRoutingModule
  ]
})
export class CommandlineModule { }
