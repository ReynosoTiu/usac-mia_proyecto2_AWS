package main

import (
	"fmt"
	"os"
	"strconv"
)

func Eliminar_Grupo(name string, path string) string {

	if actualSesion.Id_user != 1 || actualSesion.Id_grp != 1 {
		return "RMGRP solo puede ser ejecutato por el usuario root"
	}

	file, err := os.OpenFile(path, os.O_RDWR, 0777)
	defer file.Close()

	if err != nil {
		//log.Fatal("ERROR AL ABRIR ARCHIVO,ARCHIVO NO EXISTE", err)
		return "Error al abrir el archivo"
	}
	var super SUPER_BLOQUE
	var inodo TABLA_INODOS
	var archivo BLOQUE_ARCHIVO

	var col int32 = 1
	var actual byte
	var posicion int32 = 0
	var numBloque int32 = 0
	var id int32 = -1
	var tipo byte = '\x00'
	var grupo string = ""
	var palabra string = ""
	var flag bool = false

	file.Seek(int64(actualSesion.InicioSuper), 0)
	super = Leer_SuperBloque(file, tamanoSuper)

	//Nos posicionamos en el inodo del archivo users.txt
	file.Seek(int64(super.S_inode_start+tamanoInodo), 0)
	inodo = Leer_TablaInodo(file, tamanoInodo)

	for i := 0; i < 16; i++ {
		if inodo.I_block[i] != -1 {

			file.Seek(int64(super.S_block_start+tamanoArchivo*inodo.I_block[i]), 0)
			archivo = Leer_BloqueArchivo(file, tamanoArchivo)
			for j := 0; j < 64; j++ {
				actual = archivo.B_content[j]
				if actual == '\n' {
					if tipo == 'G' {
						grupo = palabra
						if grupo == name {
							file.Seek(int64(super.S_block_start+tamanoArchivo*numBloque), 0)
							archivo = Leer_BloqueArchivo(file, tamanoArchivo)
							//fread(&archivo,sizeof(BloqueCarpeta),1,fp);
							archivo.B_content[posicion] = '0'
							file.Seek(int64(super.S_block_start+tamanoArchivo*numBloque), 0)
							Escribir_BloqueArchivo(file, &archivo)
							//fwrite(&archivo,sizeof(BloqueArchivo),1,fp);
							fmt.Println("Grupo eliminado con exito")
							flag = true
							break
						}
					}
					col = 1
					palabra = ""
				} else if actual != ',' {
					palabra += string(actual)
					col++
				} else if actual == ',' {
					if col == 2 {
						idd, _ := strconv.Atoi(palabra)
						id = int32(idd)
						_ = id
						posicion = int32(j - 1)
						numBloque = inodo.I_block[i]
					} else if col == 4 {
						tipo = palabra[0]
					}

					col++
					palabra = ""
				}
			}
			if flag {
				return "GRUPO REMOVIDO EXITOSAMENTE"
				//break

			}

		}
	}

	return "OCURRIO UN ERROR AL TRATAR DE REMOVER EL GRUPO"
}
