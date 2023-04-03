package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func Ejecutar_rmdisk(path string) int {

	if _, err := os.Stat(path); err == nil {
		// path/to/whatever exists
		var mensaje string = "Deseas Eliminar el archivo S o N?:"
		fmt.Println(mensaje)
		reader := bufio.NewReader(os.Stdin)
		text, _ := reader.ReadString('\n')
		// convert CRLF to LF
		text = strings.Replace(text, "\n", "", -1)
		if strings.ToLower(text) == "s" {
			e := os.Remove(path)
			if e != nil {
				log.Fatal(e)
			}
			return 0
		} else if strings.ToLower(text) == "n" {
			return 1
		} else {
			return 2
		}

	}

	return -1
}
