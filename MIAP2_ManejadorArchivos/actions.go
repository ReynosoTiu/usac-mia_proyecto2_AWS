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
	ejecutar_contenido_json(content)
	w.Header().Set("Content-type", "application/json")
	json.NewEncoder(w).Encode("mkdir exitosamente")
	//Borrar esto xD fmt.Println(content)

}

func ejecutar_contenido_json(contenido string) {
	var temp string = strings.ReplaceAll(contenido, "\r", "")

	res1 := strings.Split(temp, "\n")

	for i := 4; i < len(res1)-2; i++ {
		if res1[i] != "" {
			//var temp = res1[i]
			//fmt.Println("MI COntenido[", i, "]:", res1[i])
			//fmt.Println("MI COntenido[", i, "]:", []byte(temp))
			if Reconocer_Comando(res1[i]) == 0 {
				fmt.Println("RECONOCIO COMANDO CORRECTATEMTE")
			} else {
				fmt.Println("ERROR RETURN -1 EN ANALIZAR PARAMETROS")
			}

		}
	}

}
