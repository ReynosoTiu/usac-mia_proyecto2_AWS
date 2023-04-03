package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"strconv"
	_ "strings"
	"unsafe"
)

func Log_in(path string, id string, usuario string, password string) int32 {
	var aux_nodo *NODO = listaS.obtenerNodo(id)
	if aux_nodo != nil {
		index, inicio, tamano, es_logica, fit_ := Buscar_Indice_P_E_L(aux_nodo.Path, aux_nodo.Name) //returna 4 variables
		_ = tamano

		file, err := os.OpenFile(path, os.O_RDWR, 0777)
		defer file.Close()

		if err != nil {
			log.Fatal("ERROR AL ABRIR ARCHIVO,ARCHIVO NO EXISTE", err)
			return -1
		}

		if index != -1 {
			if es_logica {
				inicio += int32(unsafe.Sizeof(EBR{}))
			}

			var super SUPER_BLOQUE
			var inodo TABLA_INODOS

			file.Seek(int64(inicio), 0)
			super = Leer_SuperBloque(file, int32(unsafe.Sizeof(SUPER_BLOQUE{})))

			file.Seek(int64(super.S_inode_start)+int64(unsafe.Sizeof(TABLA_INODOS{})), 0)
			inodo = Leer_TablaInodo(file, int32(unsafe.Sizeof(TABLA_INODOS{})))

			file.Seek(int64(super.S_inode_start)+int64(unsafe.Sizeof(TABLA_INODOS{})), 0)
			copy(inodo.I_atime[:], Hora_fecha())
			Escribir_TablaInodo(file, &inodo)
			file.Close()

			actualSesion.InicioSuper = int32(inicio)

			if !actualSesion.hay_Sesion {
				//LLenenado los datos de Usuario actual
				idGU, existeUsuario := verificarDatos(usuario, password, path)
				if existeUsuario {
					actualSesion.Id_user = int32(idGU)
					actualSesion.Id_grp = int32(idGU)
					//actualSesion.InicioSuper = inicio
					actualSesion.Tipo_sistema = super.S_filesystem_type
					actualSesion.Path = aux_nodo.Path
					copy(actualSesion.Fit[:], string(fit_)) //aqui arreglar a ver que pex
					actualSesion.hay_Sesion = existeUsuario
					return 0
				}
			} else {
				fmt.Println("ERROR DEBES CERRAR SESION PARA INICIAR OTRA SESION")
				return 0
			}
			//return -1
		}

	}

	return -1
}

func verificarDatos(user string, password string, direccion string) (int32, bool) {
	file, err := os.OpenFile(direccion, os.O_RDWR, 0777)
	defer file.Close()

	if err != nil {
		log.Fatal("ERROR AL ABRIR ARCHIVO,ARCHIVO NO EXISTE", err)
		//os.Exit(4)
		return -1, false
	}

	var cadena []byte = bytes.Repeat([]byte{0}, 400)
	var super SUPER_BLOQUE
	var inodo TABLA_INODOS

	file.Seek(int64(actualSesion.InicioSuper), 0)
	super = Leer_SuperBloque(file, int32(unsafe.Sizeof(SUPER_BLOQUE{})))

	//Leemos el inodo del archivo users.txt
	file.Seek(int64(super.S_inode_start)+int64(unsafe.Sizeof(TABLA_INODOS{})), 0)
	inodo = Leer_TablaInodo(file, int32(unsafe.Sizeof(TABLA_INODOS{})))

	for i := 0; i < 15; i++ {
		if inodo.I_block[i] != -1 {
			var archivo BLOQUE_ARCHIVO

			file.Seek(int64(super.S_block_start)+int64(unsafe.Sizeof(BLOQUE_CARPETA{})), 0)
			for j := 0; j < int(inodo.I_block[i]); j++ {
				archivo = Leer_BloqueArchivo(file, int32(unsafe.Sizeof(BLOQUE_ARCHIVO{})))
				//fread(&archivo,sizeof(BloqueArchivo),1,fp);
			}
			copy(cadena[:], string(archivo.B_content[:]))
		}
	}

	file.Close()

	return Usuario_Contrasena(cadena, user, password)

	//return -1,false
}
func Usuario_Contrasena(cadena []byte, user string, pass string) (int32, bool) {
	cadena_sinNull := bytes.Trim([]byte(cadena), "\x00")
	//fmt.Println("CADENA SIN NULL:", cadena_sinNull)
	grupo_usuario := bytes.Split(cadena_sinNull, []byte("\n"))
	//fmt.Println("GRUPO_USUARIO:", grupo_usuario)
	//fmt.Println("GRUPO_USUARIO TAMANO:", len(grupo_usuario))

	//grupo := bytes.Split([]byte(grupo_usuario[0]), []byte(","))
	////_=grupo
	//fmt.Println("GRUPO:", grupo)
	//usuario := bytes.Split([]byte(grupo_usuario[1]), []byte(","))
	//fmt.Println("USUARIO:", usuario)

	//GID, Tipo, Grupo
	//UID, Tipo, Grupo, Usuario, Contraseña

	for i := 0; i < len(grupo_usuario); i++ {
		if len(grupo_usuario[i]) > 0 {
			indice_grupo_usuario := bytes.Split([]byte(grupo_usuario[i]), []byte(","))
			if string(indice_grupo_usuario[1]) == "U" {
				if string(indice_grupo_usuario[3]) == user {
					if string(indice_grupo_usuario[4]) == pass {
						if string(indice_grupo_usuario[0]) != "0" {
							fmt.Println("CONVERTIENDOME EN STRING:", string(indice_grupo_usuario[0]))
							intID, _ := strconv.Atoi(string(indice_grupo_usuario[0]))
							var id int32 = int32(intID)
							return id, true
						} else {
							fmt.Println("ERROR EL USUARIO YA ELIMINADO PARA PODER INICIAR SESION")
							return -1, false
						}
					} else {
						fmt.Println("ERROR AUTENTIFICACION FALLIDA EN LOGIN")
						return -1, false
					}
					//return true
				}
			}
		}

	}

	fmt.Println("ERROR INCIAR LOGIN,USUARIO NO EXISTE")
	return -1, false
}

/*
func ByteArrayToint32(arr []byte) int32 {
	//fmt.Println("AQUIIII VIENTO ^^^^^^^^^^^^^^^^^", arr)
	val := int32(0)
	size := len(arr)
	for i := 0; i < size; i++ {
		*(*uint8)(unsafe.Pointer(uintptr(unsafe.Pointer(&val)) + uintptr(i))) = arr[i]
	}
	return val
}
*/
