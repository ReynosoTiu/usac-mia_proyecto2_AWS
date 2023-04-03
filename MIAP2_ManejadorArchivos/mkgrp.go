package main

import (
	"fmt"
	"os"
	"strconv"
)

func Crear_Mkgrp(name string) int32 {
	if actualSesion.hay_Sesion {

		if actualSesion.Id_user == 1 && actualSesion.Id_grp == 1 { //Usuario root
			if len(name) <= 10 {
				//grupo, id_grupo := Existe_Grupo(name, actualSesion.Path)
				var grupo int = buscarGrupo(name, actualSesion.Path)
				if grupo == -1 {
					var id_grupo int = getID_grp()
					var nuevo_Grupo string = strconv.Itoa(int(id_grupo)) + ",G," + name + "\n"
					agregarUsersTXT1(nuevo_Grupo, actualSesion.Path)
					fmt.Println("GRUPO CREADO CON EXITO")
					return 0
				} else {
					fmt.Println("ERROR EXISTE GRUPO CON NOMBRE IGUAL")
					return -1
				}

			} else {
				fmt.Println("ERROR NOMBRE EXECEDE 10 CARACTERES")
				return -1
			}
		} else {
			fmt.Println("ERROR SOLO ROOT EJECUTA COMANDO MKGRP")
			return -1
		}

	} else {
		fmt.Println("ERROR DEBE EXISTIR SESION PARA CREAR GRUPO")
		return -1
	}
	//return 0

}

func buscarGrupo(name string, path_disco string) int {
	file, err := os.OpenFile(path_disco, os.O_RDWR, 0777)
	defer file.Close()

	if err != nil {
		return -1
	}
	var cadena []byte
	var super SUPER_BLOQUE
	var inodo TABLA_INODOS

	file.Seek(int64(actualSesion.InicioSuper), 0)
	super = Leer_SuperBloque(file, tamanoSuper)
	//Leemos el inodo del archivo users.txt
	file.Seek(int64(super.S_inode_start+tamanoInodo), 0)
	inodo = Leer_TablaInodo(file, tamanoInodo)

	for i := 0; i < 16; i++ {
		if inodo.I_block[i] != -1 {
			var archivo BLOQUE_ARCHIVO
			file.Seek(int64(super.S_block_start), 0)
			for j := 0; j <= int(inodo.I_block[i]); j++ {
				archivo = Leer_BloqueArchivo(file, tamanoArchivo)
			}
			//strcat(cadena, archivo.B_content)
			cadena = append(cadena[:], archivo.B_content[:]...)
		}
	}

	var contenido_cadena string = aString(cadena)
	var contenido_ptr *string = &contenido_cadena
	var token *string = strtok_r(contenido_ptr, "\n")

	for *token != "" {
		var id [2]byte
		var tipo [2]byte
		var group [10]byte
		var token_id *string = strtok_r(token, ",")
		copy(id[:], []byte(*token_id))
		if aString(id[:]) != "0" { //Verificar que no sea un U/G eliminado
			var token_grupo *string = strtok_r(token, ",")
			copy(tipo[:], []byte(*token_grupo))
			if aString(tipo[:]) == "G" {
				//var token_nombre_grupo *string = strtok_r(token, ",")
				var token_nombre_grupo *string = token
				copy(group[:], []byte(*token_nombre_grupo))
				if aString(group[:]) == name {
					id_num, _ := strconv.Atoi(aString(id[:]))

					return id_num
				}

			}
		}
		token = strtok_r(contenido_ptr, "\n")

	}

	return -1
}

func strtok_r(cadena *string, delimitador string) *string {
	var info *string = nil

	var temp_cadena []byte = []byte(*cadena)
	var temp_info [1024]byte
	var temp_fin_cadena [1024]byte

	for i := 0; i < len(*cadena); i++ {
		if (*cadena)[i] == (delimitador)[0] {
			copy(temp_info[:], temp_cadena[:i])
			copy(temp_fin_cadena[:], temp_cadena[(i+1):])
			break
		}
	}
	var n, b string
	n = aString(temp_fin_cadena[:])
	b = aString(temp_info[:])
	*cadena = n
	info = &b
	*info = b
	return info
}

func agregarUsersTXT1(datos string, path_disco string) int {
	fp, err := os.OpenFile(path_disco, os.O_RDWR, 0777)
	defer fp.Close()

	if err != nil {
		return -1
	}
	var super SUPER_BLOQUE
	var inodo TABLA_INODOS
	var archivo BLOQUE_ARCHIVO
	var blockIndex int32 = 0

	fp.Seek(int64(actualSesion.InicioSuper), 0)
	super = Leer_SuperBloque(fp, tamanoSuper)
	//Leemos el inodo del archivo users.txt
	fp.Seek(int64(super.S_inode_start+tamanoInodo), 0)
	inodo = Leer_TablaInodo(fp, tamanoInodo)

	for i := 0; i < 16; i++ {
		if inodo.I_block[i] != -1 {
			blockIndex = inodo.I_block[i] //Ultimo bloque utilizado del archivo
		}

	}

	fp.Seek(int64(super.S_block_start+tamanoArchivo*blockIndex), 0)
	archivo = Leer_BloqueArchivo(fp, tamanoArchivo)
	var enUso int32 = int32(len(aString(archivo.B_content[:])))
	var libre int32 = 64 - enUso

	if int32(len(datos)) < libre {
		concatenar(archivo.B_content[:], []byte(datos), len(aString(archivo.B_content[:])))
		fp.Seek(int64(super.S_block_start+tamanoArchivo*blockIndex), 0)
		Escribir_BloqueArchivo(fp, &archivo)

		fp.Seek(int64(super.S_inode_start+tamanoInodo), 0)
		inodo = Leer_TablaInodo(fp, tamanoInodo)
		inodo.I_size = inodo.I_size + int32(len(datos))
		copy(inodo.I_mtime[:], Hora_fecha())
		fp.Seek(int64(super.S_inode_start+tamanoInodo), 0)
		Escribir_TablaInodo(fp, &inodo)

	} else {
		var aux string = ""
		var aux2 string = ""
		var i int32 = 0

		for i = 0; i < libre; i++ {
			aux += string(datos[i])
		}

		for ; i < int32(len(datos)); i++ {
			aux2 += string(datos[i])
		}

		//Guardamos lo que cabe en el primer bloque
		concatenar(archivo.B_content[:], []byte(aux), len(aString(archivo.B_content[:])))
		fp.Seek(int64(super.S_block_start+tamanoArchivo*blockIndex), 0)
		Escribir_BloqueArchivo(fp, &archivo)
		var auxArchivo BLOQUE_ARCHIVO
		copy(auxArchivo.B_content[:], []byte(aux2))
		var bit int32 = buscarBit1(fp, 'B', actualSesion.Fit[:])
		/*Guardamos el bloque en el bitmap y en la tabla de bloques*/
		fp.Seek(int64(super.S_bm_block_start+bit), 0)
		var byte_archivo byte = '2'
		Escribir_Byte(fp, &byte_archivo)

		fp.Seek(int64(super.S_block_start+tamanoArchivo*bit), 0)
		Escribir_BloqueArchivo(fp, &auxArchivo)
		/*Guardamos el modificado del inodo*/
		fp.Seek(int64(super.S_inode_start+tamanoInodo), 0)
		inodo = Leer_TablaInodo(fp, tamanoInodo)
		inodo.I_size = inodo.I_size + int32(len(datos))
		copy(inodo.I_mtime[:], Hora_fecha())
		inodo.I_block[blockIndex] = bit
		fp.Seek(int64(super.S_inode_start+tamanoInodo), 0)
		Escribir_TablaInodo(fp, &inodo)
		/*Guardamos la nueva cantidad de bloques libres y el primer bloque libre*/
		super.S_first_blo = super.S_first_blo + 1
		super.S_free_blocks_count = super.S_free_blocks_count - 1
		fp.Seek(int64(actualSesion.InicioSuper), 0)
		Escribir_SuperBloque(fp, &super)
		return 0
	}
	return -1
}

func buscarBit1(fp *os.File, tipo byte, fit []byte) int32 {
	var super SUPER_BLOQUE
	var inicio_bm int32 = 0
	var tempBit byte = '0'
	var bit_libre int32 = -1
	var tam_bm int32 = 0

	fp.Seek(int64(actualSesion.InicioSuper), 0)
	super = Leer_SuperBloque(fp, tamanoSuper)

	if tipo == 'I' {
		tam_bm = super.S_inodes_count
		inicio_bm = super.S_bm_inode_start
	} else if tipo == 'B' {
		tam_bm = super.S_blocks_count
		inicio_bm = super.S_bm_block_start
	}

	/*----------------Tipo de ajuste a utilizar----------------*/
	if aString(fit) == "ff" { //Primer ajuste
		for i := 0; i < int(tam_bm); i++ {
			fp.Seek(int64(inicio_bm+int32(i)), 0)
			tempBit = Leer_Byte(fp, 1)
			if tempBit == '0' {
				bit_libre = int32(i)
				return bit_libre
			}
		}

		if bit_libre == -1 {
			return -1
		}

	} else if aString(fit[:]) == "bf" { //Mejor ajuste
		var libres int32 = 0
		var auxLibres int32 = -1

		for i := 0; i < int(tam_bm); i++ { //Primer recorrido
			fp.Seek(int64(inicio_bm+int32(i)), 0)
			tempBit = Leer_Byte(fp, 1)
			if tempBit == '0' {
				libres++
				if int32(i)+1 == tam_bm {
					if auxLibres == -1 || auxLibres == 0 {
						auxLibres = libres
					} else {
						if auxLibres > libres {
							auxLibres = libres
						}

					}
					libres = 0
				}
			} else if tempBit == '1' {
				if auxLibres == -1 || auxLibres == 0 {
					auxLibres = libres
				} else {
					if auxLibres > libres {
						auxLibres = libres
					}

				}
				libres = 0
			}
		}

		for i := 0; i < int(tam_bm); i++ {
			fp.Seek(int64(inicio_bm+int32(i)), 0)
			tempBit = Leer_Byte(fp, 1)
			if tempBit == '0' {
				libres++
				if int32(i)+1 == tam_bm {
					return ((int32(i) + 1) - libres)
				}
			} else if tempBit == '1' {
				if auxLibres == libres {
					return ((int32(i) + 1) - libres)
				}

				libres = 0
			}
		}

		return -1

	} else if aString(fit[:]) == "wf" { //Peor ajuste
		var libres int32 = 0
		var auxLibres int32 = -1

		for i := 0; i < int(tam_bm); i++ { //Primer recorrido
			fp.Seek(int64(inicio_bm+int32(i)), 0)
			tempBit = Leer_Byte(fp, 1)
			if tempBit == '0' {
				libres++
				if int32(i)+1 == tam_bm {
					if auxLibres == -1 || auxLibres == 0 {
						auxLibres = libres
					} else {
						if auxLibres < libres {
							auxLibres = libres
						}
					}
					libres = 0
				}
			} else if tempBit == '1' {
				if auxLibres == -1 || auxLibres == 0 {
					auxLibres = libres
				} else {
					if auxLibres < libres {
						auxLibres = libres
					}
				}
				libres = 0
			}
		}

		for i := 0; i < int(tam_bm); i++ {
			fp.Seek(int64(inicio_bm+int32(i)), 0)
			tempBit = Leer_Byte(fp, 1)
			if tempBit == '0' {
				libres++
				if int32(i)+1 == tam_bm {
					return ((int32(i) + 1) - libres)
				}
			} else if tempBit == '1' {
				if auxLibres == libres {
					return ((int32(i) + 1) - libres)
				}
				libres = 0
			}
		}

		return -1
	}

	return 0
}

func concatenar(contenido_original []byte, copiar []byte, inicio int) {
	var contador int = 0
	for i := inicio; i < inicio+len(copiar); i++ {
		contenido_original[i] = copiar[contador]
		contador++
	}

}

func getID_grp() int {

	fp, err := os.OpenFile(actualSesion.Path, os.O_RDWR, 0777)
	defer fp.Close()

	if err != nil {
		return -1
	}

	var cadena [400]byte
	var aux_id int = -1
	var super SUPER_BLOQUE
	var inodo TABLA_INODOS
	fp.Seek(int64(actualSesion.InicioSuper), 0)
	super = Leer_SuperBloque(fp, tamanoSuper)
	//Leemos el inodo del archivo users.txt
	fp.Seek(int64(super.S_inode_start+tamanoInodo), 0)
	inodo = Leer_TablaInodo(fp, tamanoInodo)

	for i := 0; i < 16; i++ {
		if inodo.I_block[i] != -1 {
			var archivo BLOQUE_ARCHIVO
			fp.Seek(int64(super.S_block_start+tamanoArchivo*inodo.I_block[i]), 0)
			//for j := 0; j < int(inodo.I_block[i]); j++ {
			archivo = Leer_BloqueArchivo(fp, tamanoArchivo)
			//}

			concatenar(cadena[:], archivo.B_content[:], len(aString(cadena[:])))
		}
	}

	//cadena := []byte("1,G,Root\n2,G,Dark\n3,G,Ghost\n4,G,Deep\n5,G,Camino\n6,G,Usuarios\n")
	var contenido_cadena string = aString(cadena[:])
	var contenido_ptr *string = &contenido_cadena

	var token *string = strtok_r(contenido_ptr, "\n")
	for *token != "" {
		var id [2]byte
		var tipo [2]byte
		var token_id *string = strtok_r(token, ",")
		copy(id[:], *token_id)
		if aString(id[:]) != "0" { //Verificar que no sea un U/G eliminado
			var token_grupo *string = strtok_r(token, ",")
			copy(tipo[:], *token_grupo)
			if aString(tipo[:]) == "G" {
				temp_id, _ := strconv.Atoi(aString(id[:]))
				aux_id = temp_id
			}

		}
		token = strtok_r(contenido_ptr, "\n")
	}
	aux_id++
	return aux_id
}
