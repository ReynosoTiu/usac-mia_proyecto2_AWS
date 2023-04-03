package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

func EjecutarExec(path string) {
	if Obtener_extension_archivo(path) == ".script" {
		f, err := os.Open(path)
		if err != nil {
			log.Fatalf("open file error: %v", err)
			return
		}
		defer f.Close()

		sc := bufio.NewScanner(f)
		for sc.Scan() {
			lineaTexto := sc.Text() // GET the line string

			if len(lineaTexto) > 0 {
				if Reconocer_Comando(lineaTexto) == 0 {
					//descomentar si est es el comentario xD
					//fmt.Println(lineaTexto)
				} else {
					fmt.Println("ERROR AL RECONOCER COMANDO:", lineaTexto)
				}

			}

		}
		fmt.Println("COMANDO EXEC EJECUTADO CORRECTAMENTE")

		if err := sc.Err(); err != nil {
			log.Fatalf("scan file error: %v", err)
			return
		}
	} else {
		fmt.Println("ERROR SOLO EJECUTAR ARCHIVOS .SCRIPT...EXTENSION ERRONEO:", Obtener_extension_archivo(path))
	}

}
