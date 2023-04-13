package main

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os"
	"strings"
	"unsafe"
)

func existeArchivo(path string) bool {
	filiInfo, err := os.Stat(path)
	if err != nil {
		fmt.Println("ERROR EL ARCHIVO NO EXISTE")
		return false
	}
	_ = filiInfo // no se utilizo para nada solo para verificar que exist archivo
	return true
}

/*_________________________ INICIO SUPER BLOQUE ___________________*/
func inicio_super_bloque(path string, name string) int32 {

	//os.Open sirve para leer el archivo nada mas y no para modificar
	file, err := os.OpenFile(path, os.O_RDWR, 0777)
	defer file.Close()

	if err != nil {
		log.Fatal("ERROR AL ABRIR ARCHIVO,ARCHIVO NO EXISTE", err)
		return -1
	}
	masterBoot := MBR{}

	var tamano_masterBoot int32 = int32(unsafe.Sizeof(masterBoot))
	masterBoot = ObtenerInforMbr(file, tamano_masterBoot)

	var indice_extendia int32 = -1
	for i := 0; i < 4; i++ {
		if string(bytes.Trim(masterBoot.Mbr_partition[i].Part_name[:], "\x00")) == name {
			file.Close()
			return masterBoot.Mbr_partition[i].Part_start
		}
		if string(masterBoot.Mbr_partition[i].Part_type) == "e" {
			indice_extendia = int32(i)
		}
	}

	if indice_extendia != -1 {
		//aqui va lo de la logica
		extendBoot := EBR{}
		file.Seek(int64(masterBoot.Mbr_partition[indice_extendia].Part_start), 0) //(cantidad_bytes_desplazar,origen_de_archivo)
		extendBoot = ObtenerInfoEbr(file, int32(unsafe.Sizeof(EBR{})))
		for {

			if string(bytes.Trim(extendBoot.Part_name[:], "\x00")) == name {
				return extendBoot.Part_start
			}

			if extendBoot.Part_next == -1 {
				break
			}
			file.Seek(int64(extendBoot.Part_next), 0)
			extendBoot = ObtenerInfoEbr(file, int32(unsafe.Sizeof(EBR{}))) //ruta quemada por el momento xD
		}
	}
	file.Close()

	return -1
}

/*_______________________________ SEPARAR PATH CARPETA_____*/
func separa_path(path string) []string {
	var respuesta = []string{}
	var temp_path = path
	path_separado := strings.Split(temp_path, "/")
	for i := 0; i < len(path_separado); i++ {
		if path_separado[i] != "" {
			respuesta = append(respuesta, path_separado[i])
		}
	}
	return respuesta
}

/*_____________________ MEMSET RELLENA UN ARRAGLO DE BYTES CON UN CARACTER ESPECIFICO _______*/
func memseT(a []byte, v byte) {
	for i := range a {
		a[i] = v
	}
}

/*_____________________ LEER ARCHIVO ______________________________________*/
func Leer_Cualquier_Archivo(path string) string {

	f, err := os.Open(path)
	if err != nil {
		//log.Fatalf("open file error: %v", err)
		fmt.Println("ERROR NO SE PUDO LEER CONTENIDO EN CREAR MKFILE")
		return ""
	}
	defer f.Close()
	/*
		var contenido_retornar string = ""
		sc := bufio.NewScanner(f)
		for sc.Scan() {
			lineaTexto := sc.Text() // GET the line string
			//fmt.Println("Contenido Linea:")
			if len(lineaTexto) > 0 {
				contenido_retornar = contenido_retornar + lineaTexto
			}

		}
	*/

	_, err = f.Seek(0, 0)
	r4 := bufio.NewReader(f)
	b4, err := r4.Peek(1024)

	//return contenido_retornar
	return aString(b4[:])
}
