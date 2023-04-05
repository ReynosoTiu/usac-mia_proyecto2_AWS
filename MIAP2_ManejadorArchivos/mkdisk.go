package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	_ "fmt"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func Hora_fecha() string {
	tiemp_hora := time.Now()
	return tiemp_hora.Format("01-02-2006 15:04:05")
}

func Crear_archivo(size int32, path string, fit string, unit string) string {
	var size_archivo int32 = 0
	var fit_archivo string = ""
	if Obtener_extension_archivo(path) == ".dsk" {
		//_ = size_archivo
		//_ = fit_archivo

		//validando el unit del archivo
		if unit == "k" {
			size_archivo = size * 1024
		} else if unit == "m" {
			size_archivo = size * 1024 * 1024
		} else if unit == "" {
			size_archivo = size * 1024 * 1024
		}

		//validando fit del archivo
		if fit == "bf" {
			fit_archivo = fit
		} else if fit == "ff" {
			fit_archivo = fit
		} else if fit == "wf" {
			fit_archivo = fit
		} else if fit == "" {
			fit_archivo = "ff"
		}

		var buffer [1024]byte
		for i := 0; i < 1024; i++ {
			buffer[i] = 0
		}
		binarios := &buffer
		var data_buffer bytes.Buffer
		binary.Write(&data_buffer, binary.BigEndian, binarios)

		mi_archivo, err := os.Create(path)
		defer mi_archivo.Close()

		if err != nil {
			log.Fatal(err)
			return "Error del SO al crear el disco"
		}
		var j int32
		for j = 0; j < (size_archivo / 1024); j++ {
			EscribirBytes(mi_archivo, data_buffer.Bytes())
		}

		masterBoot := &MBR{}
		masterBoot.Mbr_tamano = int32(size_archivo)
		copy(masterBoot.Mbr_fecha_creacion[:], Hora_fecha())
		masterBoot.Mbr_disk_signature = int32(Numero_Random())
		copy(masterBoot.Disk_fit[:], fit_archivo)

		for i := 0; i < 4; i++ {
			masterBoot.Mbr_partition[i].Part_status = 0
			masterBoot.Mbr_partition[i].Part_type = 0
			copy(masterBoot.Mbr_partition[i].Part_fit[:], "")
			masterBoot.Mbr_partition[i].Part_start = -1
			masterBoot.Mbr_partition[i].Part_size = -1
			copy(masterBoot.Mbr_partition[i].Part_name[:], "")

		}

		mi_archivo.Seek(0, 0)
		var master_buffer bytes.Buffer
		binary.Write(&master_buffer, binary.BigEndian, masterBoot)
		EscribirBytes(mi_archivo, master_buffer.Bytes())

		mi_archivo.Close()

	} else {
		fmt.Println("extension")
		return "Extension no admitida en el parametro PATH"
	}

	return "ARCHIVO CREADO EXITOSAMENTE"
}

func EscribirBytes(file *os.File, bytes []byte) {
	_, err := file.Write(bytes)

	if err != nil {
		log.Fatal(err)
	}

}

func Numero_Random() int32 {
	var minimo int32 = 1
	var maximo int32 = 200

	return int32(rand.Intn(int(maximo) - int(minimo) + int(minimo)))
}

func Crear_carpetas(direccion string) {
	path_temp := filepath.Dir(direccion)

	if _, err := os.Stat(path_temp); os.IsNotExist(err) {
		// path/to/whatever does not exist
		if existeError := os.MkdirAll(path_temp, os.ModePerm); existeError != nil {
			log.Fatal(existeError)
		}
	}

}

func Obtener_extension_archivo(path string) string {
	fmt.Println("Extension encontrada " + strings.ToLower(filepath.Ext(path)))
	return strings.ToLower(filepath.Ext(path))
}

func Obtener_nombe_archivo(path string, extension string) string {

	return strings.TrimSuffix(filepath.Base(path), extension)
}
