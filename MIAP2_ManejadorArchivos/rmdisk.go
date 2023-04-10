package main

import (
	"log"
	"os"
)

func Ejecutar_rmdisk(path string) string {

	if _, err := os.Stat(path); err == nil {
		e := os.Remove(path)
		if e != nil {
			log.Fatal(e)
		}
		return "DISCO ELIMINADO CORRECTAMENTE"
	}

	return "DISCO NO ENCONTRADO"
}
