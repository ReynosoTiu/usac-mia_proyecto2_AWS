package main

import (
	"fmt"
	"os"
	"strconv"
)

func Crear_Mkusr(usuario string, pass string, grupo string) string {
	if len(usuario) > 10 {
		return "Valor del parametro USER debe tener maximo 10 caracteres"
	}

	if len(pass) > 10 {
		return "Valor del parametro PWD debe tener maximo 10 caracteres"
	}

	if len(grupo) > 10 {
		return "Valor del parametro GRP debe tener maximo 10 caracteres"
	}

	idGrupo := buscarGrupo(grupo, actualSesion.Path)
	if idGrupo != -1 {
		if idGrupo == 0 {
			return "Grupo eliminado"
		}

		if !buscarUsuario(usuario) {
			var id = getID_usr()
			fmt.Println("ID Cliente", id)
			var datos string = strconv.Itoa(id) + ",U," + grupo + "," + usuario + "," + pass + "\n"
			agregarUsersTXT1(datos, actualSesion.Path)

			return "USUARIO CREADO EXITOSAMENTE"
		} else {
			return "El usuario ya existe"
		}

	}
	return "Grupo no encontrado"
}

func buscarUsuario(name string) bool {
	fp, err := os.OpenFile(actualSesion.Path, os.O_RDWR, 0777)
	defer fp.Close()

	if err != nil {
		return false
	}

	var cadena [400]byte
	var super SUPER_BLOQUE
	var inodo TABLA_INODOS

	fp.Seek(int64(actualSesion.InicioSuper), 0)
	super = Leer_SuperBloque(fp, tamanoSuper)
	//Nos posicionamos en el inodo del archivo users.txt
	fp.Seek(int64(super.S_inode_start+tamanoInodo), 0)
	inodo = Leer_TablaInodo(fp, tamanoInodo)

	for i := 0; i < 16; i++ {
		if inodo.I_block[i] != -1 {
			var archivo BLOQUE_ARCHIVO
			fp.Seek(int64(super.S_block_start+tamanoArchivo*inodo.I_block[i]), 0)
			//for j := 0; j < int(inodo.I_block[i]); j++ {
			//fread(&archivo,sizeof(BloqueArchivo),1,fp);
			archivo = Leer_BloqueArchivo(fp, tamanoArchivo)
			//}
			//strcat(cadena,archivo.b_content);
			concatenar(cadena[:], archivo.B_content[:], len(aString(cadena[:])))
		}
	}

	var contenido_cadena string = aString(cadena[:])
	var contenido_ptr *string = &contenido_cadena
	var token *string = strtok_r(contenido_ptr, "\n")

	for *token != "" {
		var id [2]byte
		var tipo [2]byte
		var user [10]byte
		//id,Tipo,NombreGrupo,nombre_usuario,contrasena
		//1,U,root,root,123\n
		var token_id *string = strtok_r(token, ",")
		copy(id[:], []byte(*token_id))
		if aString(id[:]) != "0" { //Verificar que no sea un U/G eliminado
			var token_tipo *string = strtok_r(token, ",")
			copy(tipo[:], []byte(*token_tipo))

			if aString(tipo[:]) == "U" {
				var token_nombre_grupo *string = strtok_r(token, ",")
				_ = token_nombre_grupo
				var token_nombre_usuario *string = strtok_r(token, ",")
				copy(user[:], []byte(*token_nombre_usuario))
				if aString(user[:]) == name {
					fmt.Println("USUARIO ENCONTRADO", aString(user[:]))
					return true
				}

			}
		}
		token = strtok_r(contenido_ptr, "\n")

	}

	return false
}

func getID_usr() int {
	fp, err := os.OpenFile(actualSesion.Path, os.O_RDWR, 0777)
	defer fp.Close()

	if err != nil {
		return -1
	}
	var cadena []byte
	var res int = 0
	var super SUPER_BLOQUE
	var inodo TABLA_INODOS

	fp.Seek(int64(actualSesion.InicioSuper), 0)
	super = Leer_SuperBloque(fp, tamanoSuper)
	//Nos posicionamos en el inodo del archivo users.txt
	fp.Seek(int64(super.S_inode_start+tamanoInodo), 0)
	inodo = Leer_TablaInodo(fp, tamanoInodo)

	for i := 0; i < 16; i++ {
		if inodo.I_block[i] != -1 {
			var archivo BLOQUE_ARCHIVO
			fp.Seek(int64(super.S_block_start), 0)
			fmt.Println("tamanio", int(inodo.I_block[i]))
			for j := 0; j < int(inodo.I_block[i]); j++ {
				archivo = Leer_BloqueArchivo(fp, tamanoArchivo)
			}
			//concatenar(cadena[:], archivo.B_content[:], len(aString(cadena[:])))
			cadena = append(cadena[:], archivo.B_content[:]...)
		}
	}

	var contenido_cadena string = aString(cadena[:])
	var contenido_ptr *string = &contenido_cadena

	fmt.Println(contenido_cadena)
	var token *string = strtok_r(contenido_ptr, "\n")
	for *token != "" {
		var id [2]byte
		var tipo [2]byte

		var token_id *string = strtok_r(token, ",")
		copy(id[:], *token_id)
		fmt.Println("IDUSUARIO", aString(id[:]))
		if aString(id[:]) != "0" { //Verificar que no sea un U/G eliminado
			var token_tipo *string = strtok_r(token, ",")
			copy(tipo[:], *token_tipo)
			if aString(tipo[:]) == "U" {
				res++
			}
		}
		token = strtok_r(contenido_ptr, "\n")
	}
	return res
}
