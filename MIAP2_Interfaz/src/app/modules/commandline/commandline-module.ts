import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';

import { CommandlineRoutingModule } from './commandline-routing.module';
import { CommandlineComponent } from './commandline/commandline.component';
import { FileInputComponent } from 'src/app/shared/file-input/file-input.component';
import { CommandlineService } from './commandline.service';

@NgModule({
  declarations: [
    CommandlineComponent,
    FileInputComponent
  ],
  imports: [
    CommonModule,
    CommandlineRoutingModule
  ],
  providers: [
    CommandlineService
  ]
})
export class CommandlineModule { }
