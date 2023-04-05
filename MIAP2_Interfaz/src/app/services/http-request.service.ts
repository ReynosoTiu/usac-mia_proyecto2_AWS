import { HttpClient, HttpHeaders } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { environment } from 'src/environments/environment';

@Injectable({
  providedIn: 'root'
})
export class HttpRequestService {

  constructor(private http: HttpClient) { }

  post(endpoint: string, params: any = {}): Promise<any> {
    let HEADERS = new HttpHeaders({
      'Content-Type': 'application/json',
      'responseType': 'json',
      'auth-token': ''
    });

    return new Promise((resolve, reject) => {
      this.http.post(`${environment.URI}/${endpoint}`, params, { headers: HEADERS }).subscribe({
        next(res) {
          resolve(res);
        },
        error(err) {
          reject(err);
        },
      });
    });
  }

  get(endpoint: string): Promise<any> {
    let HEADERS = new HttpHeaders({
      'Content-Type': 'application/json',
      'responseType': 'json',
      'auth-token': ''
    });

    return new Promise((resolve, reject) => {
      this.http.get(`${environment.URI}/${endpoint}`, { headers: HEADERS }).subscribe({
        next(res) {
          resolve(res);
        },
        error(err) {
          reject(err);
        },
      });
    });
  }
}
