package main

import (
	"fmt"
	"math"
	"os"
	"path/filepath"
)

func Crear_Mkfile(path string, size_archivo int, cont_archivo string, hay_r bool) string {
	var result int = crearArchivo(path, size_archivo, cont_archivo, hay_r)
	if result == 1 {
		fmt.Println("ARCHIVO CREADA CORRECTAMENTE")
		return "ARCHIVO CREADO EXITOSAMENTE"
	} else if result == 2 {
		fmt.Println("ERROR USUARIO NO TIEME PERMISO PARA ESCRITURA")
		return "Usuario no tiene permisos de escritura"
	} else if result == 3 {
		fmt.Println("ERROR EL ARCHIVO CONTENIDO NO EXISTE")
		return "Archivo contenido no existe"
	} else if result == 4 {
		fmt.Println("ERROR NO EXISTE LA DIRECCION Y NO TIENE PARAMETRO >R")
		return "No existe la ruta, sug (>R)"
	}
	return "Hubo un error al chear el archivo"
}

func crearArchivo(path string, size int, cont string, p bool) int {
	fp, err := os.OpenFile(actualSesion.Path, os.O_RDWR, 0777)
	defer fp.Close()

	if err != nil {
		return -1
	}

	var auxPath [500]byte
	var auxPath2 [500]byte
	copy(auxPath[:], []byte(path))
	copy(auxPath2[:], []byte(path))
	var existe int = buscarCarpetaArchivo(fp, aString(auxPath[:]))
	copy(auxPath[:], []byte(path))
	var response int = -1

	if existe != -1 {
		response = 0
	} else {
		response = nuevoArchivo(fp, actualSesion.Fit[:], p, aString(auxPath[:]), size, cont, 0, aString(auxPath2[:]))
	}
	return response
}

func nuevoArchivo(fp *os.File, fit []byte, flagP bool, path string, size int, contenido string, index int, auxPath string) int {
	var super SUPER_BLOQUE
	var inodo, inodoNuevo TABLA_INODOS
	var carpeta, carpetaNueva BLOQUE_CARPETA

	var lista = separa_path(path)
	var copiaPath [500]byte
	var directorio [500]byte
	var nombreCarpeta [80]byte
	var content string = ""
	var contentSize string = "0123456789"

	copy(copiaPath[:], []byte(path))
	copy(directorio[:], []byte((filepath.Dir(aString(copiaPath[:])))))
	copy(copiaPath[:], []byte(path))
	copy(nombreCarpeta[:], []byte(filepath.Base(aString(copiaPath[:]))))
	copy(copiaPath[:], []byte(path))

	var cont int32 = 0
	var numInodo int32 = int32(index)
	var finalSize int = size
	cont = int32(len(lista))

	if len(contenido) != 0 {
		content = Leer_Cualquier_Archivo(contenido)
		if content != "" {
			finalSize = len(content)
		} else {
			return 3
		}

	}

	fp.Seek(int64(actualSesion.InicioSuper), 0)
	super = Leer_SuperBloque(fp, tamanoSuper)

	if cont == 1 {
		var bloque int = 0
		var b_content int = 0
		var libre int = buscarContentLibre(fp, int(numInodo), &inodo, &carpeta, &b_content, &bloque)

		if libre == 1 {
			var permisos bool = permisosDeEscritura(int(inodo.I_perm), (inodo.I_uid == actualSesion.Id_user), (inodo.I_gid == actualSesion.Id_grp))
			if permisos || (actualSesion.Id_user == 1 && actualSesion.Id_grp == 1) {
				var buffer byte = '1'
				var buffer2 byte = '2'

				//Agregamos el archivo al bloque correspondiente
				var bitInodo int = int(buscarBit1(fp, 'I', fit[:]))
				carpeta.B_content[b_content].B_inodo = int32(bitInodo)
				copy(carpeta.B_content[b_content].B_name[:], nombreCarpeta[:])
				fp.Seek(int64(super.S_block_start+tamanoCarpeta*inodo.I_block[bloque]), 0)
				Escribir_BloqueCarpeta(fp, &carpeta)

				//Creamos el nuevo inodo archivo
				inodoNuevo = crearInodo(0, '1', 664)
				fp.Seek(int64(super.S_inode_start+tamanoInodo*int32(bitInodo)), 0)
				Escribir_TablaInodo(fp, &inodoNuevo)

				//Registramos el inodo en el bitmap
				fp.Seek(int64(super.S_bm_inode_start+int32(bitInodo)), 0)
				Escribir_Byte(fp, &buffer)

				//Si viene el parametro -size/-cont
				if finalSize != 0 {

					var n float64 = float64(finalSize) / float64(64)
					var numBloques int = int(math.Ceil(n))
					var caracteres int = finalSize
					var charNum int = 0  //size_t
					var contChar int = 0 //size_t
					numInodo = int32(buscarCarpetaArchivo(fp, auxPath))
					for i := 0; i < numBloques; i++ {
						var archivo BLOQUE_ARCHIVO
						memseT(archivo.B_content[:], 0)

						var bitBloque int32 = buscarBit1(fp, 'B', fit[:])
						//Registramos el bloque en el bitmap
						fp.Seek(int64(super.S_bm_block_start+bitBloque), 0)
						Escribir_Byte(fp, &buffer2)

						if caracteres > 64 {
							for j := 0; j < 64; j++ {
								if len(content) != 0 { //-cont
									archivo.B_content[j] = content[contChar]
									contChar++
								} else { //-size
									archivo.B_content[j] = contentSize[charNum]
									charNum++
									if charNum == 10 {
										charNum = 0
									}
								}
							}
							//Guardamos el bloque en el respectivo inodo archivo
							fp.Seek(int64(super.S_inode_start+tamanoInodo*numInodo), 0)
							inodo = Leer_TablaInodo(fp, tamanoInodo)

							inodo.I_block[i] = bitBloque
							inodo.I_size = int32(finalSize)
							fp.Seek(int64(super.S_inode_start+tamanoInodo*numInodo), 0)
							Escribir_TablaInodo(fp, &inodo)
							//Guardamos el bloque
							fp.Seek(int64(super.S_block_start+tamanoArchivo*bitBloque), 0)
							Escribir_BloqueArchivo(fp, &archivo)

							caracteres -= 64
						} else {
							for j := 0; j < caracteres; j++ {
								if len(content) != 0 {
									archivo.B_content[j] = content[contChar]
									contChar++
								} else {
									archivo.B_content[j] = contentSize[charNum]
									charNum++
									if charNum == 10 {
										charNum = 0
									}

								}
							}
							//Guardamos el bloque en el respectivo inodo archivo
							fp.Seek(int64(super.S_inode_start+tamanoInodo*numInodo), 0)
							inodo = Leer_TablaInodo(fp, tamanoInodo)
							inodo.I_block[i] = bitBloque
							inodo.I_size = int32(finalSize)
							fp.Seek(int64(super.S_inode_start+tamanoInodo*numInodo), 0)
							Escribir_TablaInodo(fp, &inodo)

							//Guardamos el bloque
							fp.Seek(int64(super.S_block_start+tamanoArchivo*bitBloque), 0)
							Escribir_BloqueArchivo(fp, &archivo)

						}

					}
					//Modificamos el super bloque
					super.S_free_blocks_count = super.S_free_blocks_count - int32(numBloques)
					super.S_free_inodes_count = super.S_free_inodes_count - 1
					super.S_first_ino = super.S_first_ino + 1
					super.S_first_blo = super.S_first_blo + int32(numBloques)
					fp.Seek(int64(actualSesion.InicioSuper), 0)
					Escribir_SuperBloque(fp, &super)
					//fwrite(&super,sizeof(SuperBloque),1,stream)
					return 1
				}
				super.S_free_inodes_count = super.S_free_inodes_count - 1
				super.S_first_ino = super.S_first_ino + 1
				fp.Seek(int64(actualSesion.InicioSuper), 0)
				Escribir_SuperBloque(fp, &super)
				return 1
			} else {
				return 2
			}

		} else { //Todos los bloques estan llenos
			fp.Seek(int64(super.S_inode_start+tamanoInodo*numInodo), 0)
			inodo = Leer_TablaInodo(fp, tamanoInodo)

			for i := 0; i < 16; i++ {
				if inodo.I_block[i] == -1 {
					bloque = i
					break
				}
			}
			//Apuntadores directos
			var permisos bool = permisosDeEscritura(int(inodo.I_perm), (inodo.I_uid == actualSesion.Id_user), (inodo.I_gid == actualSesion.Id_grp))
			if permisos || (actualSesion.Id_user == 1 && actualSesion.Id_grp == 1) {
				var buffer byte = '1'
				var buffer2 byte = '2'

				var bitBloque int32 = buscarBit1(fp, 'B', fit[:])
				//Guardamos el bloque en el inodo
				inodo.I_block[bloque] = bitBloque
				fp.Seek(int64(super.S_inode_start+tamanoInodo*numInodo), 0)
				Escribir_TablaInodo(fp, &inodo)

				//Creamos el nuevo bloque carpeta
				var bitInodo int32 = buscarBit1(fp, 'I', fit[:])
				carpetaNueva.B_content[0].B_inodo = bitInodo
				carpetaNueva.B_content[1].B_inodo = -1
				carpetaNueva.B_content[2].B_inodo = -1
				carpetaNueva.B_content[3].B_inodo = -1
				copy(carpetaNueva.B_content[0].B_name[:], nombreCarpeta[:])
				copy(carpetaNueva.B_content[1].B_name[:], "")
				copy(carpetaNueva.B_content[2].B_name[:], "")
				copy(carpetaNueva.B_content[3].B_name[:], "")
				fp.Seek(int64(super.S_block_start+tamanoCarpeta*bitBloque), 0)
				Escribir_BloqueCarpeta(fp, &carpetaNueva)

				//Registramos el bloque en el bitmap
				fp.Seek(int64(super.S_bm_block_start+bitBloque), 0)
				Escribir_Byte(fp, &buffer)

				//Creamos el nuevo inodo
				inodoNuevo = crearInodo(0, '1', 664)
				fp.Seek(int64(super.S_inode_start+tamanoInodo*bitInodo), 0)
				Escribir_TablaInodo(fp, &inodoNuevo)

				fp.Seek(int64(super.S_inode_start+tamanoInodo*bitInodo), 0)
				Escribir_TablaInodo(fp, &inodoNuevo)

				//Registramos el inodo en el bitmap
				fp.Seek(int64(super.S_bm_inode_start+bitInodo), 0)
				Escribir_Byte(fp, &buffer)

				//Si viene el parametro -size/-cont
				if finalSize != 0 {
					var n float64 = float64(finalSize) / float64(64)
					var numBloques int = int(math.Ceil(n))
					var caracteres int = finalSize
					var charNum int = 0
					var contChar int = 0
					numInodo = int32(buscarCarpetaArchivo(fp, auxPath))
					for i := 0; i < numBloques; i++ {
						var archivo BLOQUE_ARCHIVO
						memseT(archivo.B_content[:], 0)
						//Apuntadores simples
						var bitBloque int32 = buscarBit1(fp, 'B', fit[:])
						//Registramos el bloque en el bitmap
						fp.Seek(int64(super.S_bm_block_start+bitBloque), 0)
						Escribir_Byte(fp, &buffer2)

						if caracteres > 64 {
							for j := 0; j < 64; j++ {
								if len(content) != 0 { //-cont
									archivo.B_content[j] = content[contChar]
									contChar++
								} else { //-size
									archivo.B_content[j] = contentSize[charNum]
									charNum++
									if charNum == 10 {
										charNum = 0
									}
								}
							}
							//Guardamos el bloque en el respectivo inodo archivo
							fp.Seek(int64(super.S_inode_start+tamanoInodo*numInodo), 0)
							inodo = Leer_TablaInodo(fp, tamanoInodo)

							inodo.I_block[i] = bitBloque
							fp.Seek(int64(super.S_inode_start+tamanoInodo*numInodo), 0)
							Escribir_TablaInodo(fp, &inodo)

							//Guardamos el bloque
							fp.Seek(int64(super.S_block_start+tamanoArchivo*bitBloque), 0)
							Escribir_BloqueArchivo(fp, &archivo)

							caracteres -= 64
						} else {
							for j := 0; j < caracteres; j++ {
								if len(content) != 0 {
									archivo.B_content[j] = content[contChar]
									contChar++
								} else {
									archivo.B_content[j] = contentSize[charNum]
									charNum++
									if charNum == 10 {
										charNum = 0
									}
								}
							}
							//Guardamos el bloque en el respectivo inodo archivo
							fp.Seek(int64(super.S_inode_start+tamanoInodo*numInodo), 0)
							inodo = Leer_TablaInodo(fp, tamanoInodo)

							inodo.I_block[i] = bitBloque
							fp.Seek(int64(super.S_inode_start+tamanoInodo*numInodo), 0)
							Escribir_TablaInodo(fp, &inodo)

							//Guardamos el bloque
							fp.Seek(int64(super.S_block_start+tamanoArchivo*bitBloque), 0)
							Escribir_BloqueArchivo(fp, &archivo)

						}

					}
					//Modificamos el super bloque
					super.S_free_blocks_count = super.S_free_blocks_count - int32(numBloques)
					super.S_free_inodes_count = super.S_free_inodes_count - 1
					super.S_first_ino = super.S_first_ino + 1
					super.S_first_blo = super.S_first_blo + int32(numBloques)
					fp.Seek(int64(actualSesion.InicioSuper), 0)
					Escribir_SuperBloque(fp, &super)

					return 1
				}
				super.S_free_inodes_count = super.S_free_inodes_count - 1
				super.S_first_ino = super.S_first_ino + 1
				fp.Seek(int64(actualSesion.InicioSuper), 0)
				Escribir_SuperBloque(fp, &super)

				return 1
			} else {
				return 2
			}

		}
	} else { //Directorio
		var existe int = buscarCarpetaArchivo(fp, aString(directorio[:]))
		if existe == -1 {
			if flagP {
				var index int = 0
				var aux string = ""
				//Crear posibles carpetas inexistentes
				for i := 0; i < int(cont); i++ {
					if int32(i) == cont-1 {
						var dir [100]byte
						copy(dir[:], []byte("/"))
						concatenar(dir[:], nombreCarpeta[:], len(aString(dir[:])))
						return nuevoArchivo(fp, fit[:], false, aString(dir[:]), size, contenido, index, auxPath)
					} else {
						aux += "/" + lista[i]
						var dir [500]byte
						var auxDir [500]byte
						copy(dir[:], []byte(aux))
						copy(auxDir[:], []byte(aux))
						var carpeta int = buscarCarpetaArchivo(fp, aString(dir[:]))
						if carpeta == -1 {
							nuevaCarpeta(fp, fit[:], false, aString(auxDir[:]), index)
							copy(auxDir[:], []byte(aux))
							index = buscarCarpetaArchivo(fp, aString(auxDir[:]))
						} else {
							index = carpeta
						}
					}
				}
			} else {
				return 4
			}
		} else { //Crear el archivo en el directorio
			var dir [100]byte
			copy(dir[:], []byte("/"))
			concatenar(dir[:], nombreCarpeta[:], len(aString(dir[:])))
			return nuevoArchivo(fp, fit, false, aString(dir[:]), size, contenido, existe, auxPath)
		}
	}

	return 0
}
