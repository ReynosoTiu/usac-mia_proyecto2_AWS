package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func Crear_Mkdir(path_crear string, hay_p bool) string {
	var result int = crearCarpeta(path_crear, hay_p)
	if result == 0 {
		fmt.Println("ERROR LA CARPETA YA EXISTE")
		return "La carpeta ya existe"
	} else if result == 1 {
		fmt.Println("CARPETA CREADA CON EXITO")
		return "CARPETA CREADA CON EXITO"
	} else if result == 2 {
		fmt.Println("ERROR NO TIENE PERMISO DE ESCRITURA PAR CREAR CARPETA")
		return "No tiene permisos de escritura"
	} else if result == 3 {
		fmt.Println("ERROR: NO EXISTE EL DIRECTORIO Y  NO ESTA EL PARAMETRO -P")
		return "No existe el directorio, sug (>r)"
	}
	return "Ha ocurrido un problema al crear la carpeta"
}

func crearCarpeta(path string, p bool) int {

	fp, err := os.OpenFile(actualSesion.Path, os.O_RDWR, 0777)
	defer fp.Close()

	if err != nil {
		return -1
	}

	var auxPath [500]byte
	var existe int = buscarCarpetaArchivo(fp, path)

	copy(auxPath[:], []byte(path))
	var response int = -1

	if existe != -1 {
		response = 0
	} else {
		response = nuevaCarpeta(fp, actualSesion.Fit[:], p, path, 0)
	}

	return response
}

func buscarCarpetaArchivo(fp *os.File, path string) int {
	var super SUPER_BLOQUE
	var inodo TABLA_INODOS
	var carpeta BLOQUE_CARPETA

	var cont int = 0
	var numInodo int = 0
	var lista []string = separa_path(path)
	cont = len(lista)

	fp.Seek(int64(actualSesion.InicioSuper), 0)
	super = Leer_SuperBloque(fp, tamanoSuper)
	numInodo = int(super.S_inode_start) //Byte donde inicia el inodo

	for cont2 := 0; cont2 < cont; cont2++ {
		fp.Seek(int64(numInodo), 0)
		inodo = Leer_TablaInodo(fp, tamanoInodo)

		var siguiente int32 = 0
		for i := 0; i < 16; i++ {
			if inodo.I_block[i] != -1 { //Apuntadores directos
				var byteBloque int = byteInodoBloque(fp, int(inodo.I_block[i]), '2')
				fp.Seek(int64(byteBloque), 0)
				carpeta = Leer_BloqueCarpeta(fp, tamanoCarpeta)
				for j := 0; j < 4; j++ {
					if (cont2 == cont-1) && (strings.ToLower(aString(carpeta.B_content[j].B_name[:])) == strings.ToLower(lista[cont2])) { //Tendria que ser la carpeta
						return int(carpeta.B_content[j].B_inodo)
					} else if (cont2 != cont-1) && (strings.ToLower(aString(carpeta.B_content[j].B_name[:])) == strings.ToLower(lista[cont2])) {
						numInodo = byteInodoBloque(fp, int(carpeta.B_content[j].B_inodo), '1')
						siguiente = 1
						break
					}
				}

				if siguiente == 1 {
					break
				}

			}
		}
	}

	return -1
}

func byteInodoBloque(fp *os.File, pos int, tipo byte) int {
	var super SUPER_BLOQUE
	fp.Seek(int64(actualSesion.InicioSuper), 0)
	super = Leer_SuperBloque(fp, tamanoSuper)
	if tipo == '1' {
		return (int(super.S_inode_start) + int(tamanoInodo)*pos)
	} else if tipo == '2' {
		return (int(super.S_block_start) + int(tamanoCarpeta)*pos)
	}
	return 0
}

func nuevaCarpeta(fp *os.File, fit []byte, flagP bool, path string, index int) int {
	var super SUPER_BLOQUE
	var inodo, inodoNuevo TABLA_INODOS
	var carpeta, carpetaNueva, carpetaAux BLOQUE_CARPETA

	var copiaPath [500]byte
	var directorio [500]byte
	var nombreCarpeta [80]byte

	copy(copiaPath[:], []byte(path))
	copy(directorio[:], []byte(filepath.Dir(path))) //borrar por si las moscas
	copy(copiaPath[:], []byte(path))
	copy(nombreCarpeta[:], []byte(filepath.Base(path)))

	var cont int = 0
	var numInodo int = index
	var response int = 0

	var lista []string = separa_path(path)
	cont = len((lista))

	fp.Seek(int64(actualSesion.InicioSuper), 0)
	super = Leer_SuperBloque(fp, tamanoSuper)

	if cont == 1 { //Solo es una carpeta '/home' | '/archivos'
		var content int = 0
		var bloque int = 0
		var libre int = buscarContentLibre(fp, numInodo, &inodo, &carpeta, &content, &bloque)
		if libre == 1 {
			//Apuntadores directos
			var permisos bool = permisosDeEscritura(int(inodo.I_perm), (inodo.I_uid == actualSesion.Id_user), (inodo.I_gid == actualSesion.Id_grp))
			if permisos || actualSesion.Id_user == 1 && actualSesion.Id_grp == 1 {
				var buffer byte = '1'
				var bitInodo int = int(buscarBit1(fp, 'I', fit[:]))
				//Agregamos la carpeta al espacio libre en el bloque
				carpeta.B_content[content].B_inodo = int32(bitInodo)
				copy(carpeta.B_content[content].B_name[:], nombreCarpeta[:])
				fp.Seek(int64(super.S_block_start+tamanoCarpeta*inodo.I_block[bloque]), 0)
				Escribir_BloqueCarpeta(fp, &carpeta)

				//Creamos el nuevo inodo
				inodoNuevo = crearInodo(0, '0', 664)
				var bitBloque int = int(buscarBit1(fp, 'B', fit[:]))
				inodoNuevo.I_block[0] = int32(bitBloque)
				fp.Seek(int64(super.S_inode_start+tamanoInodo*int32(bitInodo)), 0)
				Escribir_TablaInodo(fp, &inodoNuevo)

				//Creamos el nuevo bloque carpeta
				carpetaNueva = crearBloqueCarpeta()
				carpetaNueva.B_content[0].B_inodo = int32(bitInodo)
				carpetaNueva.B_content[1].B_inodo = int32(numInodo)
				copy(carpetaNueva.B_content[0].B_name[:], ".")
				copy(carpetaNueva.B_content[1].B_name[:], "..")
				fp.Seek(int64(super.S_block_start+tamanoCarpeta*int32(bitBloque)), 0)
				Escribir_BloqueCarpeta(fp, &carpetaNueva)
				//Guardamos los bits en los bitmaps
				fp.Seek(int64(super.S_bm_inode_start+int32(bitInodo)), 0)
				Escribir_Byte(fp, &buffer)
				fp.Seek(int64(super.S_bm_block_start+int32(bitBloque)), 0)
				Escribir_Byte(fp, &buffer)
				//Sobreescribimos el super bloque
				super.S_free_inodes_count = super.S_free_inodes_count - 1
				super.S_free_blocks_count = super.S_free_blocks_count - 1
				super.S_first_ino = super.S_first_ino + 1
				super.S_first_blo = super.S_first_blo + 1
				fp.Seek(int64(actualSesion.InicioSuper), 0)
				Escribir_SuperBloque(fp, &super)
				return 1
			} else {
				return 2
			}

		} else if libre == 0 { //Todos bloques estan llenos

			fp.Seek(int64(super.S_inode_start+tamanoInodo*int32(numInodo)), 0)
			inodo = Leer_TablaInodo(fp, tamanoInodo)
			for i := 0; i < 16; i++ {

				if inodo.I_block[i] == -1 {
					bloque = i
					break
				}

			}

			//Apuntadores directos
			var permissions bool = permisosDeEscritura(int(inodo.I_perm), (inodo.I_uid == actualSesion.Id_user), (inodo.I_gid == actualSesion.Id_grp))
			if permissions || actualSesion.Id_user == 1 && actualSesion.Id_grp == 1 {

				var buffer byte = '1'
				var bitBloque int = int(buscarBit1(fp, 'B', fit[:]))

				inodo.I_block[bloque] = int32(bitBloque)
				//Sobreescribimos el inodo
				fp.Seek(int64(super.S_inode_start+tamanoInodo*int32(numInodo)), 0)
				//inodo = Leer_TablaInodo(fp, tamanoInodo)
				Escribir_TablaInodo(fp, &inodo)
				//Bloque carpeta auxiliar
				var bitInodo int = int(buscarBit1(fp, 'I', fit[:]))
				carpetaAux = crearBloqueCarpeta()
				carpetaAux.B_content[0].B_inodo = int32(bitInodo)
				copy(carpetaAux.B_content[0].B_name[:], nombreCarpeta[:])
				fp.Seek(int64(super.S_block_start+tamanoCarpeta*int32(bitBloque)), 0)
				Escribir_BloqueCarpeta(fp, &carpetaAux)
				//Escribimos el bit en el bitmap de blqoues
				fp.Seek(int64(super.S_bm_block_start+int32(bitBloque)), 0)
				Escribir_Byte(fp, &buffer)

				//Creamos el nuevo inodo
				inodoNuevo = crearInodo(0, '0', 664)
				bitBloque = int(buscarBit1(fp, 'B', fit[:]))
				inodoNuevo.I_block[0] = int32(bitBloque)
				fp.Seek(int64(super.S_inode_start+tamanoInodo*int32(bitInodo)), 0)
				Escribir_TablaInodo(fp, &inodoNuevo)
				//Escribimos el bit en el bitmap de inodos
				fp.Seek(int64(super.S_bm_inode_start+int32(bitInodo)), 0)
				Escribir_Byte(fp, &buffer)
				//Creamos el nuevo bloque carpeta
				carpetaNueva = crearBloqueCarpeta()
				carpetaNueva.B_content[0].B_inodo = int32(bitInodo)
				carpetaNueva.B_content[1].B_inodo = int32(numInodo)
				copy(carpetaNueva.B_content[0].B_name[:], ".")
				copy(carpetaNueva.B_content[1].B_name[:], "..")
				fp.Seek(int64(super.S_block_start+tamanoCarpeta*int32(bitBloque)), 0)
				Escribir_BloqueCarpeta(fp, &carpetaNueva)

				//Guardamos el bit en el bitmap de bloques
				fp.Seek(int64(super.S_bm_block_start+int32(bitBloque)), 0)
				Escribir_Byte(fp, &buffer)

				//Sobreescribimos el super bloque
				super.S_free_inodes_count = super.S_free_inodes_count - 1
				super.S_free_blocks_count = super.S_free_blocks_count - 2
				super.S_first_ino = super.S_first_ino + 1
				super.S_first_blo = super.S_first_blo + 2
				fp.Seek(int64(actualSesion.InicioSuper), 0)
				Escribir_SuperBloque(fp, &super)
				return 1
			} else {
				return 2
			}

		}
	} else { //Es un directorio '/home/usac/archivos'
		//Verificar que exista el directorio
		var existe int = buscarCarpetaArchivo(fp, aString(directorio[:]))
		if existe == -1 {
			if flagP {
				var index int = 0
				var aux string = ""
				//Crear posibles carpetas inexistentes
				for i := 0; i < cont; i++ {
					aux += "/" + lista[i]
					var dir [500]byte
					var auxDir [500]byte
					copy(dir[:], []byte(aux))
					copy(auxDir[:], []byte(aux))
					var carpeta int = buscarCarpetaArchivo(fp, aString(dir[:]))
					if carpeta == -1 {
						response = nuevaCarpeta(fp, fit[:], false, aString(auxDir[:]), index)
						if response == 2 {
							break
						}
						copy(auxDir[:], []byte(aux))
						index = buscarCarpetaArchivo(fp, aString(auxDir[:]))
					} else {
						index = carpeta
					}

				}
			} else {
				return 3
			}

		} else { //Solo crear la carpeta en el directorio
			var dir [100]byte
			copy(dir[:], []byte("/"))
			concatenar(dir[:], nombreCarpeta[:], len(aString(dir[:])))
			return nuevaCarpeta(fp, fit[:], false, aString(dir[:]), existe)
		}
	}

	return response
}

/*__________________________________________________________*/

func buscarContentLibre(fp *os.File, numInodo int, inodo *TABLA_INODOS, carpeta *BLOQUE_CARPETA, content *int, bloque *int) int {
	var libre int = 0
	var super SUPER_BLOQUE
	fp.Seek(int64(actualSesion.InicioSuper), 0)
	super = Leer_SuperBloque(fp, tamanoSuper)
	fp.Seek(int64(super.S_inode_start+tamanoInodo*int32(numInodo)), 0)
	*inodo = Leer_TablaInodo(fp, tamanoInodo)

	//Buscamos un espacio libre en el bloque carpeta
	for i := 0; i < 16; i++ {
		if inodo.I_block[i] != -1 {
			//Apuntadores directos
			fp.Seek(int64(super.S_block_start+tamanoCarpeta*inodo.I_block[i]), 0)
			*carpeta = Leer_BloqueCarpeta(fp, tamanoCarpeta)
			for j := 0; j < 4; j++ {
				if carpeta.B_content[j].B_inodo == -1 {
					libre = 1
					*bloque = i
					*content = j
					break
				}
			}
		}

		if libre == 1 {
			break
		}
	}

	return libre
}

/*________________________________________________________________________*/
func permisosDeEscritura(permisos int, flagUser bool, flagGroup bool) bool {
	var aux string = strconv.Itoa(permisos)
	var propietario byte = aux[0]
	var grupo byte = aux[1]
	var otros byte = aux[2]

	if (propietario == '2' || propietario == '3' || propietario == '6' || propietario == '7') && flagUser {
		return true
	} else if (grupo == '2' || grupo == '3' || grupo == '6' || grupo == '7') && flagGroup {
		return true
	} else if otros == '2' || otros == '3' || otros == '6' || otros == '7' {
		return true
	}

	return false
}

func crearInodo(size int, type_ byte, perm int) TABLA_INODOS {
	var inodo TABLA_INODOS
	inodo.I_uid = actualSesion.Id_user
	inodo.I_gid = actualSesion.Id_grp
	inodo.I_size = int32(size)
	copy(inodo.I_atime[:], []byte(Hora_fecha()))
	copy(inodo.I_ctime[:], []byte(Hora_fecha()))
	copy(inodo.I_mtime[:], []byte(Hora_fecha()))
	for i := 0; i < 16; i++ {
		inodo.I_block[i] = -1
	}

	inodo.I_type = type_
	inodo.I_perm = int32(perm)
	return inodo
}

func crearBloqueCarpeta() BLOQUE_CARPETA {
	var carpeta BLOQUE_CARPETA

	for i := 0; i < 4; i++ {
		copy(carpeta.B_content[i].B_name[:], "")
		carpeta.B_content[i].B_inodo = -1
	}

	return carpeta
}
