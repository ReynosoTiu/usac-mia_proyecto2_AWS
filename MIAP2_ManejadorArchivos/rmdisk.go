package main

import (
	"fmt"
	"os"
)

func Ejecutar_rmdisk(path string) string {

	archivo := path
	fmt.Println(archivo)
	if archivoExiste(archivo) {
		e := os.Remove(path)
		if e != nil {
			return "Ocurrio un error al intentar eliminar el disco"
		}
		return "DISCO ELIMINADO CORRECTAMENTE"
	} else {
		return "DISCO NO ENCONTRADO"
	}
}
