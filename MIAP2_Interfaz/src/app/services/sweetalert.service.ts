import { Injectable } from '@angular/core';
import Swal from 'sweetalert2'

@Injectable({
  providedIn: 'root'
})
export class SweetalertService {

  constructor() { }

  async confirmAction(titulo: string, texto: string) {
    return await Swal.fire({
      title: titulo,
      text: texto,
      icon: 'warning',
      showCancelButton: true,
      confirmButtonColor: '#3085d6',
      cancelButtonColor: '#d33',
      confirmButtonText: 'SI',
      cancelButtonText: 'NO'
    }).then((result:any) => {
      if (result.isConfirmed) {
        return true;
      }
      return false;
    }).catch((err:any) => { return false })
  }

  async informacion(titulo: string, texto: string, btn:string) {
    return await Swal.fire({
      title: titulo,
      text: texto,
      icon: 'info',
      confirmButtonColor: '#3085d6',
      confirmButtonText: btn,
    }).then((result:any) => {
      if (result.isConfirmed) {
        return true;
      }
      return false;
    }).catch((err:any) => { return false })
  }
}
