import { Injectable } from '@angular/core';
import { HotToastService } from '@ngneat/hot-toast';

@Injectable({
  providedIn: 'root'
})
export class GtoastrService {

  constructor(private toast: HotToastService) { }

  error(text: string){
    this.toast.error(text);
  }

  info(text: string){
    this.toast.info(text);
  }

  warning(text: string){
    this.toast.warning(text);
  }

  success(text: string){
    this.toast.success(text);
  }

  loading(text: string, duration = 5){
    this.toast.loading(text, {duration});
  }

  show(text: string){
    this.toast.show(text);
  }

}
