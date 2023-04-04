package main

import (
	"bufio"
	_ "bytes"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	_ "github.com/godror/godror"
	"github.com/rs/cors"
)

func main() {
	InicializarVariablesGlobales()
	//iniciarServidor("127.0.0.1", "4000")
	leerComando() //siempre descomentar esto porque siver

}

func leerComando() {

	var bienvenida string = "~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~\n~~~~~~~~~~~ BIENVENIDO INGRESE UN COMANDO ~~~~~~~~~~~\n~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~"
	for {
		fmt.Println(bienvenida)
		reader := bufio.NewReader(os.Stdin)
		text, _ := reader.ReadString('\n')
		// convert CRLF to LF
		text = strings.Replace(text, "\n", "", -1)
		if strings.ToLower(text) == "exit" {
			break
		}
		if Reconocer_Comando(text) == 0 {
			fmt.Println("RECONOCIO COMANDO CORRECTATEMTE")
		} else {
			fmt.Println("ERROR RETURN -1 EN ANALIZAR PARAMETROS")
		}
	}

}

func iniciarServidor(ip string, puerto string) {
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders: []string{"*"},
	})
	router := NewRouter()

	fmt.Println("El servidor esta corriendo en http://" + ip + ":" + puerto)
	log.Fatal(http.ListenAndServe(":"+puerto, c.Handler(router)))
}

//______________ FIN DE LEER COMANDO ____________
//mkdisk -Size=3000 -unit=m -path="/home/us/Disco1.dk" -fit=wf
//mkdisk -size=1 -path=/home/dark/Archivo_temp/Ejemplo_Go/Prueba_disco/Hola/Dark/Prueba.DK -unit=k

/*
mkdisk -Size=3000 -unit=m -path="/home/us/Disco1.dk" -fit=wf

 COMANDOS FUNCIONALES
 MKDISK
 RMDISK
 FDISK
 MOUNT
 MKFS
 LOGIN
 LOGOUT
 MKGRP
 RMGRP
 PAUSE
 REP
 EXEC

*/
