import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';

const routes: Routes = [
  {
    path: '',
    loadChildren: () => import('./modules/commandline/commandline-module').then(m => m.CommandlineModule)
  },
  {
    path: 'login',
    loadChildren: () => import('./modules/login/login-module').then(m => m.LoginModule)
  },
  {
    path: 'reports',
    loadChildren: () => import('./modules/report/report-routing.module').then(m => m.ReportRoutingModule)
  },
  {
    path: '**', redirectTo: 'commandline'
  }
];

@NgModule({
  imports: [RouterModule.forRoot(routes)],
  exports: [RouterModule]
})
export class AppRoutingModule { }
