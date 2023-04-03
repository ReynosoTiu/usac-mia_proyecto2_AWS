package main

import (
	"os"
	"strconv"
)

func Eliminar_Usurio(name string, path_disco string) int {
	if actualSesion.Id_user == 1 && actualSesion.Id_grp == 1 {
		if buscarUsuario(name) {
			eliminarUsuario(name)
			return 0
		}
	}

	return -1
}

func eliminarUsuario(name string) {
	fp, err := os.OpenFile(actualSesion.Path, os.O_RDWR, 0777)
	defer fp.Close()

	if err != nil {

	}

	var super SUPER_BLOQUE
	var inodo TABLA_INODOS
	var archivo BLOQUE_ARCHIVO

	var col int32 = 1
	var actual byte
	var palabra string = ""
	var posicion int32 = 0
	var numBloque int32 = 0
	var id int32 = -1
	var tipo byte = '\x00'
	var grupo string = ""
	var usuario string = ""
	var flag bool = false

	fp.Seek(int64(actualSesion.InicioSuper), 0)
	super = Leer_SuperBloque(fp, tamanoSuper)
	//Nos posicionamos en el inodo del archivo users.txt
	fp.Seek(int64(super.S_inode_start+tamanoInodo), 0)
	inodo = Leer_TablaInodo(fp, tamanoInodo)

	for i := 0; i < 16; i++ {
		if inodo.I_block[i] != -1 {
			fp.Seek(int64(super.S_block_start+tamanoArchivo*inodo.I_block[i]), 0)
			archivo = Leer_BloqueArchivo(fp, tamanoArchivo)
			for j := 0; j < 64; j++ {
				actual = archivo.B_content[j]
				if actual == '\n' {
					if tipo == 'U' {
						if usuario == name {
							fp.Seek(int64(super.S_block_start+tamanoArchivo*numBloque), 0)
							archivo = Leer_BloqueArchivo(fp, tamanoArchivo)

							archivo.B_content[posicion] = '0'
							fp.Seek(int64(super.S_block_start+tamanoArchivo*numBloque), 0)
							Escribir_BloqueArchivo(fp, &archivo)
							flag = true
							//return 0
							break

						}
						usuario = ""
						grupo = ""
					}
					col = 1
					palabra = ""
				} else if actual != ',' {
					palabra += string(actual)
					col++
				} else if actual == ',' {
					if col == 2 {
						tem_id, _ := strconv.Atoi(palabra)
						id = int32(tem_id)
						_ = id
						posicion = int32(j) - 1
						numBloque = inodo.I_block[i]
					} else if col == 4 {
						tipo = palabra[0]
					} else if grupo == "" {
						grupo = palabra
					} else if usuario == "" {
						usuario = palabra
					}
					col++
					palabra = ""
				}
			}
			if flag {
				//return 0
				break
			}

		}
	}

}
