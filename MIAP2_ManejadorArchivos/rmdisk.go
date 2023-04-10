package main

import (
	"os"
)

func Ejecutar_rmdisk(path string) string {

	archivo := path
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

func archivoExiste(ruta string) bool {
	if _, err := os.Stat(ruta); os.IsNotExist(err) {
		return false
	}
	return true
}
