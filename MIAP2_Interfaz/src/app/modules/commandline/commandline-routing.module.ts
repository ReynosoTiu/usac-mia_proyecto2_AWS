import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';
import { CommandlineComponent } from './commandline/commandline.component';

const routes: Routes = [
    {
        path: '',
        redirectTo: 'commandline'
    },
    {
        path: 'commandline',
        component: CommandlineComponent
    }
];

@NgModule({
    imports: [RouterModule.forChild(routes)],
    exports: [RouterModule]
})
export class CommandlineRoutingModule { }