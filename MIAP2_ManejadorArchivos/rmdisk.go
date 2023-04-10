package main

import (
	"fmt"
	"log"
	"os"
)

func Ejecutar_rmdisk(path string) string {
	fmt.Println("path " + path)
	if _, err := os.Stat(path); err == nil {

		e := os.Remove(path)
		if e != nil {
			log.Fatal(e)
		}
		return "DISCO ELIMINADO CORRECTAMENTE"
	}

	return "DISCO NO ENCONTRADO"
}
