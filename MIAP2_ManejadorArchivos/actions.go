package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

func getHola(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode("HOLA MUNDO")
	//fmt.Fprintf(w, "Hola MUndo")
}

type inicioS struct {
	Username string `json:"Username"`
	Password string `json:"Password"`
}

type logn []inicioS

func Login(w http.ResponseWriter, r *http.Request) {
	var lgn = logn{}
	var inicio inicioS
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Error datos no validos")
	}
	json.Unmarshal(reqBody, &inicio) //La función Unmarshal() en la codificación del paquete / json se utiliza para descomprimir o decodificar los datos de JSON a la estructura
	lgn = append(lgn, inicio)
	fmt.Println(inicio)
	fmt.Println("USUARIO: ", inicio.Username, len(inicio.Username))
	fmt.Println("PASSWORD: ", inicio.Password, len(inicio.Password))
	w.Header().Set("Content-type", "application/json")
	json.NewEncoder(w).Encode(inicio.Username)

}

type Archivo struct {
	Contenido string `json:"Contenido"`
}

func Carga(w http.ResponseWriter, r *http.Request) {

	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Error datos no validos")
	}
	content := string(reqBody)
	resultado := ejecutar_contenido_json(content)
	w.Header().Set("Content-type", "application/json")
	json.NewEncoder(w).Encode(resultado)
	//Borrar esto xD fmt.Println(content)

}

func ejecutar_contenido_json(contenido string) string {
	var temp string = strings.ReplaceAll(contenido, "\r", "")
	return Reconocer_Comando(temp + "\n")
	// res1 := strings.Split(temp, "\n")
	// var result = ""
	// for i := 0; i < len(res1); i++ {
	// 	if res1[i] != "" {
	// 		result += Reconocer_Comando(res1[i]) + "\n"
	// 	}
	// }
	// return result
}
