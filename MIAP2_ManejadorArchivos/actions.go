package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

type Payload struct {
	Carnet string `json:"carnet"`
	Nombre string `json:"nombre"`
}

func hojaTrabajo4(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	p := Payload{
		Carnet: "201345126",
		Nombre: "José Luis Reynoso Tiu",
	}
	json.NewEncoder(w).Encode(p)
}

type inicioS struct {
	Username string `json:"Username"`
	Password string `json:"Password"`
}
type Respuesta struct {
	Tipo      int32  `json:"Tipo"`
	Mensaje   string `json:"Mensaje"`
	Data      string `json:"Data"`
	Ruta      string `json:"Ruta"`
	Extension string `json:"Extension"`
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

func ejecutar_contenido_json(contenido string) Respuesta {
	return Reconocer_Comando(contenido)
}

func archivoExiste(ruta string) bool {
	if _, err := os.Stat(ruta); os.IsNotExist(err) {
		return false
	}
	return true
}
