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

func main() {
	InicializarVariablesGlobales()
	iniciarServidor("44.201.225.17", "4000")
	//leerComando() //siempre descomentar esto porque siver

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

		fmt.Println(Reconocer_Comando(text))
	}

}
