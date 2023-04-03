package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func Reconocer_parametros(texto_parametros string, nombre_comando string) []Param {
	//string comando_analizar = "mkdisk -size=10  -unit=m -fit=ff  -path=\"/home/dark/documentos/nuevacarpeta/archivo.txt\" #comentario de una linea";
	var param string = ""
	var valor_param string = ""
	var hay_comentario bool = false
	var hay_parametro bool = true
	var hay_valor bool = false
	var hay_comilla bool = false

	//variables comando

	var parametros [6]Param
	for k := 0; k < 6; k++ {
		parametros[k].nombre = ""
		parametros[k].valor = ""
	}

	//iniciando analizar parametros
	var indice int = 0
	for i := 0 + len(nombre_comando); i < len(texto_parametros); i++ {

		if !hay_comentario {

			if texto_parametros[i] == 32 && !hay_comilla {
				//espacio
				hay_parametro = true
				hay_valor = false
			}
			if texto_parametros[i] == 35 {
				//comentario

				hay_comentario = true
				hay_valor = false
				hay_parametro = false

			}
			if hay_valor {
				//valor parametro
				valor_param += string(texto_parametros[i])
				if texto_parametros[i] == 34 {

					if hay_comilla == true {
						hay_comilla = false

					} else {
						hay_comilla = true

					}
				}

				if len(texto_parametros) == i+1 {
					parametros[indice].nombre = param
					parametros[indice].valor = valor_param
					param = ""
					valor_param = ""
					break
				}

				if texto_parametros[i+1] == 32 && !hay_comilla {
					parametros[indice].nombre = param
					parametros[indice].valor = valor_param
					param = ""
					valor_param = ""
					indice++

				}

			}
			if hay_parametro {
				//nombre parametro
				/*if len(texto_parametros)-1 == i {
					param += string(texto_parametros[i])
					parametros[indice].nombre = param
					parametros[indice].valor = ""
					param = ""
					valor_param = ""
					hay_parametro = false
					hay_valor = false
				} else if texto_parametros[i+1] == 32 && i < len(texto_parametros)-1 {
					param += string(texto_parametros[i])
					parametros[indice].nombre = param
					parametros[indice].valor = ""
					indice++
					param = ""
					valor_param = ""
					hay_parametro = false
					hay_valor = false

				} else if texto_parametros[i] != 32 {
					param += string(texto_parametros[i])
				}*/

				if texto_parametros[i] != 32 {
					param += string(texto_parametros[i])
				}
				if strings.ToLower(param) == "-p" {
					if i < len(texto_parametros)-1 {
						if texto_parametros[i+1] == 32 {
							parametros[indice].nombre = param
							parametros[indice].valor = ""
							indice++
							param = ""
							valor_param = ""
							hay_parametro = false
							hay_valor = false
						}
					}

					if i == (len(texto_parametros) - 1) {
						parametros[indice].nombre = param
						parametros[indice].valor = ""
						param = ""
						valor_param = ""
						hay_parametro = false
						hay_valor = false
					}

				}
				if strings.ToLower(param) == "-r" {
					if i < len(texto_parametros)-1 {
						if texto_parametros[i+1] == 32 {
							parametros[indice].nombre = param
							parametros[indice].valor = ""
							indice++
							param = ""
							valor_param = ""
							hay_parametro = false
							hay_valor = false
						}
					}

					if i == (len(texto_parametros) - 1) {
						parametros[indice].nombre = param
						parametros[indice].valor = ""
						param = ""
						valor_param = ""
						hay_parametro = false
						hay_valor = false
					}
				}
			}
			if texto_parametros[i] == 61 {
				//=
				hay_parametro = false
				hay_valor = true

			}

		}

	} //analizando parametros
	return parametros[:]

} // fin de reconocer parametros

func Reconocer_Comando(texto_comando string) int {
	fmt.Println(texto_comando)
	if texto_comando[0] != 35 {
		var ins_aux string = ""
		for i := 0; i < len(texto_comando); i++ {
			if texto_comando[i] == 32 {
				break
			}
			ins_aux += string(texto_comando[i])
		}

		switch strings.ToLower(ins_aux) {
		case "mkdisk":
			Reconocer_mkdisk(texto_comando, ins_aux)
			break
		case "fdisk":
			Reconocer_Fdisk(texto_comando, ins_aux)
			break
		case "mkfs":
			Reconocer_Mkfs(texto_comando, ins_aux)
			break
		case "rmdisk":
			Reconocer_Rmdisk(texto_comando, ins_aux)
			break
		case "mount":
			Reconocer_Mount(texto_comando, ins_aux)
			break
		case "login":
			Reconocer_Login(texto_comando, ins_aux)
			break
		case "logout":
			if Logout() == 0 {
				fmt.Println("CERRANDO SESION CORRECTAMENTE")
			} else if Logout() == 1 {
				fmt.Println("ERROR NO HAY SESION A CERRAR")
			}
			break
		case "mkgrp":
			Reconocer_Mkgrp(texto_comando, ins_aux)
			break
		case "rmgrp":
			Reconocer_Rmgrp(texto_comando, ins_aux)
			break
		case "mkusr":
			Reconocer_Mkusr(texto_comando, ins_aux)
			break
		case "rmusr":
			Reconocer_Rmusr(texto_comando, ins_aux)
			break
		case "mkdir":
			Reconocer_Mkdir(texto_comando, ins_aux)
			break
		case "mkfile":
			Reconocer_Mkfile(texto_comando, ins_aux)
			break
		case "comentario":
			break
		case "exec":

			parametros := Reconocer_parametros(texto_comando, ins_aux)
			if strings.ToLower(parametros[0].nombre) == "-path=" {
				EjecutarExec(parametros[0].valor)
			} else {
				fmt.Println("ERROR PARAMETRO DESCONOCIDO", parametros[0].nombre)
			}

			break

		case "rep":
			Reconocer_Rep(texto_comando, ins_aux)
			break
		case "pause":
			reader := bufio.NewReader(os.Stdin)
			texto, _ := reader.ReadString('\n')
			_ = texto

			break
		default:
			return -1

		}

	} else {
		//AQUI ES UN COMENTARIO
		//fmt.Println(texto_comando)
	}

	return 0
}

func Reconocer_mkdisk(lista_param string, comando_aux string) int {
	var size int32 = 0
	var unit string = ""
	var path string = ""
	var fit string = ""
	var hay_size bool = false
	var hay_path bool = false

	//inicio de comando para parametros
	if "mkdisk" == strings.ToLower(comando_aux) {
		parametros := Reconocer_parametros(lista_param, comando_aux)
		for i := 0; i < len(parametros); i++ {
			if parametros[i].nombre != "" && parametros[i].valor != "" {

				if "-size=" == strings.ToLower(parametros[i].nombre) {
					s, err := strconv.Atoi(parametros[i].valor)
					if err == nil && s > 0 {
						size = int32(s)
						hay_size = true

					} else {
						return -1
					}

				} else if "-unit=" == strings.ToLower(parametros[i].nombre) {
					if strings.ToLower(parametros[i].valor) == "" {
						fmt.Println("ERROR UNIT SIN VALOR EN PARAMETRO")
						return -1
					} else if strings.ToLower(parametros[i].valor) == "k" {
						unit = strings.ToLower(parametros[i].valor)

					} else if strings.ToLower(parametros[i].valor) == "m" {
						unit = strings.ToLower(parametros[i].valor)

					} else if strings.ToLower(parametros[i].valor) == "b" {
						unit = strings.ToLower(parametros[i].valor)

					} else {
						unit = strings.ToLower(parametros[i].valor)
						fmt.Println("ERROR UNIT PARAMETRO INCORRECTO: " + unit)
						return -1
					}
				} else if "-path=" == strings.ToLower(parametros[i].nombre) {
					var path_temp string = strings.Replace(parametros[i].valor, filepath.Ext(parametros[i].valor), ".dk", 1)
					path = strings.ReplaceAll(path_temp, "\"", "")

					if path == "" {
						fmt.Println("ERROR PATH SIN VALOR EN PARAMETRO")
						return -1
					}
					hay_path = true

				} else if "-fit=" == strings.ToLower(parametros[i].nombre) {

					if strings.ToLower(parametros[i].valor) == "" {
						fmt.Println("ERROR FIT SIN VALOR EN PARAMETRO")
						return -1
					} else if strings.ToLower(parametros[i].valor) == "wf" {
						fit = strings.ToLower(parametros[i].valor)

					} else if strings.ToLower(parametros[i].valor) == "ff" {
						fit = strings.ToLower(parametros[i].valor)

					} else if strings.ToLower(parametros[i].valor) == "bf" {
						fit = strings.ToLower(parametros[i].valor)

					} else {
						fit = strings.ToLower(parametros[i].valor)
						fmt.Println("ERROR FIT PARAMETRO INCORRECTO: " + fit)
						return -1
					}
				} else {
					fmt.Println("ERROR NINGUN PARAMETRO COICIDE: " + parametros[i].nombre)
					return -1
				}

			}

		} //fin for

		/*__________________________ CREAR DISCO ______________________*/
		if hay_size {
			if hay_path {
				Crear_carpetas(path)
				if Crear_archivo(size, path, fit, unit) == 0 {
					fmt.Println("ARCHIVO CREADO CORRECTAMENTE EN MKDISK")
				} else {
					fmt.Println("ERROR NO SE PUDO CREAR ARCHIVO EN MKDISK")
				}
			} else {
				fmt.Println("ERROR NO HAY PATH EN MKDISK")
			}
		} else {
			fmt.Println("ERROR NO HAY SIZE MKDISK")
		}

		return 0
	}
	return -1
} /*______________________ FIN DE ANALIZAR COMANDO MKDISK_____________________*/

func Reconocer_Rmdisk(lista_comando string, comando_aux string) int {
	var hay_path bool = false
	var path string = ""
	if "rmdisk" == strings.ToLower(comando_aux) {

		parametros := Reconocer_parametros(lista_comando, comando_aux)

		for i := 0; i < len(parametros); i++ {
			if parametros[i].valor != "" && parametros[i].nombre != "" {

				if "-path=" == strings.ToLower(parametros[i].nombre) {
					path = strings.ReplaceAll(parametros[i].valor, "\"", "")
					hay_path = true
				}

			}

		} //fin de for

		/*__________________________ ELIMINAR ARCHIVO ______________________*/
		if hay_path {
			var desicion int = Ejecutar_rmdisk(path)
			if desicion == 0 {
				fmt.Println("ARCHIVO ELIMINADO EXITOSAMENTE EN RMDISK")
			} else if desicion == 1 {
				fmt.Println("ELIMINAR ARCHIVO CANCELADO EN RMDISK")
			} else if desicion == 2 {
				fmt.Println("VALOR NO VALIDO PARA ELIMARAR ARCHIVO EN RMDISK")
			} else if desicion == -1 {
				fmt.Println("ARCHIVO NO EXISTE PARA ELIMARAR EN RMDISK")
			}
		} else {
			fmt.Println("ERROR NO HAY PATH EN RMDISK")
		}

		return 0
	}
	return -1
}

/*________________________ INICIO DE COMANDO FDISK ____________________________*/
func Reconocer_Fdisk(lista_comando string, comando_aux string) int {
	var size int32 = 0
	var unit string = ""
	var path string = ""
	var type_ string = ""
	var fit string = ""
	var name string = ""
	var hay_size bool = false
	var hay_path bool = false
	var hay_name bool = false
	//_ = type_
	//_ = name
	//_ = hay_name

	if "fdisk" == strings.ToLower(comando_aux) {

		parametros := Reconocer_parametros(lista_comando, comando_aux)

		for i := 0; i < len(parametros); i++ {
			if parametros[i].valor != "" && parametros[i].nombre != "" {

				if "-size=" == strings.ToLower(parametros[i].nombre) {
					s, err := strconv.Atoi(parametros[i].valor)
					if err == nil && s > 0 {
						size = int32(s)
						hay_size = true

					} else {
						return -1
					}

				} else if "-unit=" == strings.ToLower(parametros[i].nombre) {
					if strings.ToLower(parametros[i].valor) == "" {
						fmt.Println("ERROR UNIT SIN VALOR EN PARAMETRO")
						return -1
					} else if strings.ToLower(parametros[i].valor) == "k" {
						unit = strings.ToLower(parametros[i].valor)

					} else if strings.ToLower(parametros[i].valor) == "m" {
						unit = strings.ToLower(parametros[i].valor)

					} else if strings.ToLower(parametros[i].valor) == "b" {
						unit = strings.ToLower(parametros[i].valor)

					} else {
						unit = strings.ToLower(parametros[i].valor)
						fmt.Println("ERROR UNIT PARAMETRO INCORRECTO: " + unit)
						return -1
					}
				} else if "-path=" == strings.ToLower(parametros[i].nombre) {
					var path_temp string = strings.Replace(parametros[i].valor, filepath.Ext(parametros[i].valor), ".dk", 1)
					path = strings.ReplaceAll(path_temp, "\"", "")

					if path == "" {
						fmt.Println("ERROR PATH SIN VALOR EN PARAMETRO")
						return -1
					}
					hay_path = true

				} else if "-fit=" == strings.ToLower(parametros[i].nombre) {

					if strings.ToLower(parametros[i].valor) == "" {
						fmt.Println("ERROR FIT SIN VALOR EN PARAMETRO")
						return -1
					} else if strings.ToLower(parametros[i].valor) == "wf" {
						fit = strings.ToLower(parametros[i].valor)

					} else if strings.ToLower(parametros[i].valor) == "ff" {
						fit = strings.ToLower(parametros[i].valor)

					} else if strings.ToLower(parametros[i].valor) == "bf" {
						fit = strings.ToLower(parametros[i].valor)

					} else {
						//fit = strings.ToLower(parametros[i].valor)
						fmt.Println("ERROR FIT PARAMETRO INCORRECTO: " + strings.ToLower(parametros[i].valor))
						return -1
					}
				} else if "-type=" == strings.ToLower(parametros[i].nombre) {
					//type
					if strings.ToLower(parametros[i].valor) == "" {
						fmt.Println("ERROR TYPE SIN VALOR EN PARAMETRO")
						return -1
					} else if strings.ToLower(parametros[i].valor) == "p" {
						type_ = strings.ToLower(parametros[i].valor)

					} else if strings.ToLower(parametros[i].valor) == "l" {
						type_ = strings.ToLower(parametros[i].valor)

					} else if strings.ToLower(parametros[i].valor) == "e" {
						type_ = strings.ToLower(parametros[i].valor)

					} else {
						fmt.Println("ERROR TYPE PARAMETRO INCORRECTO: " + strings.ToLower(parametros[i].valor))
						return -1
					}

				} else if "-name=" == strings.ToLower(parametros[i].nombre) {
					//name
					if parametros[i].valor == "" {
						fmt.Println("ERROR NAME SIN VALOR EN PARAMETRO")
						return -1
					}
					name = parametros[i].valor
					hay_name = true

				} else {
					fmt.Println("ERROR NINGUN PARAMETRO COICIDE: " + parametros[i].nombre)
					return -1
				}

			}

		} //fin de for

		/*__________________________ CREAR PARTICION ______________________*/
		if hay_size {
			if hay_path {
				if hay_name {
					if Crear_particion(size, unit, path, type_, fit, name) == 0 {
						//fmt.Println("PARTICION  CREADA CORRECTAMENTE")
					} else {
						fmt.Println("ERROR NO SE PUDO CREAR PARTICION")
					}
				} else {
					fmt.Println("ERROR NO HAY NAME EN FDISK")
				}

			} else {
				fmt.Println("ERROR NO HAY PATH EN FDISK")
			}
		} else {
			fmt.Println("ERROR NO HAY SIZE FDISK")
		}

		return 0
	}
	return -1
}

/*________________________ FIN DE COMANDO FDISK _______________________________*/

/*________________________ INICIO DE COMANDO MOUNT ____________________________*/
func Reconocer_Mount(lista_comando string, comando_aux string) int {
	var path string = ""
	var name string = ""
	var hay_path bool = false
	var hay_name bool = false

	if "mount" == strings.ToLower(comando_aux) {

		parametros := Reconocer_parametros(lista_comando, comando_aux)
		for i := 0; i < len(parametros); i++ {

			if parametros[i].valor != "" && parametros[i].nombre != "" {
				if "-path=" == strings.ToLower(parametros[i].nombre) {
					var path_temp string = strings.Replace(parametros[i].valor, filepath.Ext(parametros[i].valor), ".dk", 1)
					path = strings.ReplaceAll(path_temp, "\"", "")

					if path == "" {
						fmt.Println("ERROR PATH SIN VALOR EN PARAMETRO")
						return -1
					}
					hay_path = true

				} else if "-name=" == strings.ToLower(parametros[i].nombre) {
					//name
					if parametros[i].valor == "" {
						fmt.Println("ERROR NAME SIN VALOR EN PARAMETRO")
						return -1
					}
					name = parametros[i].valor
					hay_name = true

				} else {
					fmt.Println("ERROR NINGUN PARAMETRO COICIDE: " + parametros[i].nombre)
					return -1
				}

			}

		} //fin de for

		/*__________________________ CREAR MOUNT ______________________*/

		if hay_path {
			if hay_name {
				if !listaS.buscar_nombre_path(path, name) {
					if listaS.insertar_Nodo(path, name) == 0 {
						//fmt.Println("PARTICION  CREADA CORRECTAMENTE")
					} else {
						fmt.Println("ERROR NO SE PUDO MONTAR")
						return -1
					}
				} else {
					fmt.Println("ERROR PARTICION YA MONTADA....")
				}

			} else {
				fmt.Println("ERROR NO HAY NAME EN MOUNT")
				return -1
			}

		} else {
			fmt.Println("ERROR NO HAY PATH EN MOUNT")
			return -1
		}
		return 0
	}
	return -1
}

/*____________________________________ FIN DE COMANDO MOUNT _____________________________*/

/*____________________________________ INICIO DE COMANDO MKFS ____________________________*/

func Reconocer_Mkfs(lista_comando string, comando_aux string) int {
	var id string = ""
	var type_ string = ""
	var hay_id bool = false

	if "mkfs" == strings.ToLower(comando_aux) {
		parametros := Reconocer_parametros(lista_comando, comando_aux)
		for i := 0; i < len(parametros); i++ {

			if parametros[i].valor != "" && parametros[i].nombre != "" {
				if "-id=" == strings.ToLower(parametros[i].nombre) {
					//var path_temp string = strings.Replace(parametros[i].valor, filepath.Ext(parametros[i].valor), ".dk", 1)
					id = strings.ReplaceAll(parametros[i].valor, "\"", "")

					if id == "" {
						fmt.Println("ERROR ID SIN VALOR EN PARAMETRO")
						return -1
					}
					hay_id = true

				} else if "-type=" == strings.ToLower(parametros[i].nombre) {
					//name
					if parametros[i].valor == "" {
						fmt.Println("ERROR TYPE SIN VALOR EN PARAMETRO")
						return -1
					}
					type_ = strings.ToLower(parametros[i].valor)

				} else {
					fmt.Println("ERROR NINGUN PARAMETRO COICIDE: " + parametros[i].nombre)
					return -1
				}

			}

		} //fin de for

		/*__________________________ CREAR MKFS ______________________*/

		if hay_id {

			if EjecutarMkfs(id, type_) == 0 {
				fmt.Println("EXT2... FORMATEADO CON EXITO...MKFS  EJECUTADA EXITOSAMENTE")
			} else {
				fmt.Println("ERROR NO SE PUDO MKFS")
				return -1
			}

		} else {
			fmt.Println("ERROR NO HAY ID EN MKFS")
			return -1
		}
		return 0
	}
	return -1
}

/*____________________________________ FIN DE COMANDO MKFS _______________________________*/

/*___________________________________ INCIO DE LOGIN _____________________________________*/
func Reconocer_Login(lista_comando string, comando_aux string) int {
	var id string = ""
	var password string = ""
	var user string = ""
	var hay_id bool = false
	var hay_password = false
	var hay_user = false

	if "login" == strings.ToLower(comando_aux) {
		parametros := Reconocer_parametros(lista_comando, comando_aux)
		for i := 0; i < len(parametros); i++ {

			if parametros[i].valor != "" && parametros[i].nombre != "" {
				if "-id=" == strings.ToLower(parametros[i].nombre) {
					//var path_temp string = strings.Replace(parametros[i].valor, filepath.Ext(parametros[i].valor), ".dk", 1)
					id = strings.ReplaceAll(parametros[i].valor, "\"", "")

					if id == "" {
						fmt.Println("ERROR ID SIN VALOR EN PARAMETRO")
						return -1
					}
					hay_id = true

				} else if "-password=" == strings.ToLower(parametros[i].nombre) {
					//name
					if parametros[i].valor == "" {
						fmt.Println("ERROR PASSWORD SIN VALOR EN PARAMETRO")
						return -1
					}
					//password = strings.ToLower(parametros[i].valor)
					password = parametros[i].valor
					hay_password = true

				} else if "-usuario=" == strings.ToLower(parametros[i].nombre) {
					//name
					if parametros[i].valor == "" {
						fmt.Println("ERROR USUARIO SIN VALOR EN PARAMETRO")
						return -1
					}
					//user = strings.ToLower(parametros[i].valor)
					user = parametros[i].valor
					hay_user = true

				} else {
					fmt.Println("ERROR NINGUN PARAMETRO COICIDE: " + parametros[i].nombre)
					return -1
				}

			}

		} //fin de for

		/*__________________________ INICIAR LOGIN ______________________*/

		if hay_id {
			if hay_user {
				if hay_password {
					var aux_nodo *NODO = listaS.obtenerNodo(id)
					if aux_nodo != nil {
						if Log_in(aux_nodo.Path, id, user, password) == 0 {
							fmt.Println("LOGIN EJECUTADA EXITOSAMENTE")
							return 0

						} else {
							fmt.Println("ERROR USUARIO NO SE PUDO LOGEAR")
							return -1
						}
					} else {
						fmt.Println("ERROR ID NO SE ENCOTRO EN LOGIN: ", id)
						return -1
					}

				} else {
					fmt.Println("ERROR NO HAY PASSWORD EN LOGIN")
					return -1
				}
			} else {
				fmt.Println("ERROR NO HAY USER EN LOGIN")
				return -1
			}

		} else {
			fmt.Println("ERROR NO HAY ID EN LOGIN")
			return -1
		}
	}

	return -1
}

/*___________________________________ FIN DE LOGIN _______________________________________*/

/*__________________________________  INICIO DE MKGRP ___________________________________*/
func Reconocer_Mkgrp(lista_comando string, comando_aux string) int {
	var name string = ""
	var hay_name bool = false

	if "mkgrp" == strings.ToLower(comando_aux) {
		parametros := Reconocer_parametros(lista_comando, comando_aux)
		for i := 0; i < len(parametros); i++ {

			if parametros[i].valor != "" && parametros[i].nombre != "" {

				if "-name=" == strings.ToLower(parametros[i].nombre) {
					//var path_temp string = strings.Replace(parametros[i].valor, filepath.Ext(parametros[i].valor), ".dk", 1)
					name = strings.ReplaceAll(parametros[i].valor, "\"", "")

					if name == "" {
						fmt.Println("ERROR NAME SIN VALOR EN PARAMETRO")
						return -1
					}
					hay_name = true

				} else {
					fmt.Println("ERROR NINGUN PARAMETRO COICIDE: " + parametros[i].nombre)
					return -1
				}

			}

		} //fin de for

		/*__________________________ CREAR MKGRP ______________________*/

		if hay_name {
			if Crear_Mkgrp(name) == 0 {
				//fmt.Println("LOGIN EJECUTADA EXITOSAMENTE")
				return 0
			} else {
				fmt.Println("ERROR USUARIO NO SE PUDO LOGEAR")
				return -1
			}
		} else {
			fmt.Println("ERROR NO HAY NAME EN MKGRP")
			return -1
		}
	}

	return -1
}

/*__________________________________  FIN DE MKGRP ______________________________________*/

/*__________________________________ INICO DE RMGRP ____________________________________*/
func Reconocer_Rmgrp(lista_comando string, comando_aux string) int {
	var name string = ""
	var hay_name bool = false

	if "rmgrp" == strings.ToLower(comando_aux) {
		parametros := Reconocer_parametros(lista_comando, comando_aux)
		for i := 0; i < len(parametros); i++ {

			if parametros[i].valor != "" && parametros[i].nombre != "" {

				if "-name=" == strings.ToLower(parametros[i].nombre) {
					//var path_temp string = strings.Replace(parametros[i].valor, filepath.Ext(parametros[i].valor), ".dk", 1)
					name = strings.ReplaceAll(parametros[i].valor, "\"", "")

					if name == "" {
						fmt.Println("ERROR NAME SIN VALOR EN PARAMETRO")
						return -1
					}
					hay_name = true

				} else {
					fmt.Println("ERROR NINGUN PARAMETRO COICIDE: " + parametros[i].nombre)
					return -1
				}

			}

		} //fin de for

		/*__________________________ CREAR RMGRP ______________________*/

		if hay_name {
			if Eliminar_Grupo(name, actualSesion.Path) == 0 {
				//fmt.Println("LOGIN EJECUTADA EXITOSAMENTE")
				return 0
			} else {
				fmt.Println("ERROR USUARIO NO SE PUDO ELIMINAR")
				return -1
			}
		} else {
			fmt.Println("ERROR NO HAY NAME EN RMGRP")
			return -1
		}
	}

	return -1
}

/*__________________________________ FIN DE RMGRP _____________________________________*/

/*__________________________________  INICIO DE MKUSR ___________________________________*/
func Reconocer_Mkusr(lista_comando string, comando_aux string) int {
	var usuario string = ""
	var pass_word = ""
	var grupo_incluir = ""
	var hay_usuario bool = false
	var hay_pass bool = false
	var hay_grupo bool = false

	if "mkusr" == strings.ToLower(comando_aux) {
		parametros := Reconocer_parametros(lista_comando, comando_aux)
		for i := 0; i < len(parametros); i++ {

			if parametros[i].valor != "" && parametros[i].nombre != "" {

				if "-usuario=" == strings.ToLower(parametros[i].nombre) {
					//var path_temp string = strings.Replace(parametros[i].valor, filepath.Ext(parametros[i].valor), ".dk", 1)
					usuario = strings.ReplaceAll(parametros[i].valor, "\"", "")

					if usuario == "" {
						fmt.Println("ERROR USUARIO SIN VALOR EN PARAMETRO MKUSR")
						return -1
					}
					hay_usuario = true

				} else if "-pwd=" == strings.ToLower(parametros[i].nombre) {
					//var path_temp string = strings.Replace(parametros[i].valor, filepath.Ext(parametros[i].valor), ".dk", 1)
					pass_word = strings.ReplaceAll(parametros[i].valor, "\"", "")

					if pass_word == "" {
						fmt.Println("ERROR PWD SIN VALOR EN PARAMETRO EN MKUSR")
						return -1
					}
					hay_pass = true

				} else if "-grp=" == strings.ToLower(parametros[i].nombre) {
					//var path_temp string = strings.Replace(parametros[i].valor, filepath.Ext(parametros[i].valor), ".dk", 1)
					grupo_incluir = strings.ReplaceAll(parametros[i].valor, "\"", "")

					if grupo_incluir == "" {
						fmt.Println("ERROR GRUPO SIN VALOR EN PARAMETRO MKUSR")
						return -1
					}
					hay_grupo = true

				} else {
					fmt.Println("ERROR NINGUN PARAMETRO COICIDE: "+parametros[i].nombre, "EN MKUSER")
					return -1
				}

			}

		} //fin de for

		/*__________________________ CREAR MKUSR ______________________*/

		if !hay_usuario {
			fmt.Println("ERROR NO HAY USUARIO EN MKUSR")
			return -1
		}
		if !hay_pass {
			fmt.Println("ERROR NO HAY PASS EN MKUSR")
			return -1
		}
		if !hay_grupo {
			fmt.Println("ERROR NO HAY GRUPO EN MKUSR")
			return -1
		}
		if actualSesion.Id_user == 1 && actualSesion.Id_grp == 1 {
			if Crear_Mkusr(usuario, pass_word, grupo_incluir) == 0 {
				fmt.Println("USUARIO CREADO CORRENTAMENTE")
				return 0
			}
		}

	}

	return -1
}

/*__________________________________  FIN DE MKUSR ______________________________________*/

/*__________________________________ INICO DE RMUSR ____________________________________*/
func Reconocer_Rmusr(lista_comando string, comando_aux string) int {
	var name string = ""
	var hay_name bool = false

	if "rmusr" == strings.ToLower(comando_aux) {
		parametros := Reconocer_parametros(lista_comando, comando_aux)
		for i := 0; i < len(parametros); i++ {

			if parametros[i].valor != "" && parametros[i].nombre != "" {

				if "-usuario=" == strings.ToLower(parametros[i].nombre) {
					name = strings.ReplaceAll(parametros[i].valor, "\"", "")

					if name == "" {
						fmt.Println("ERROR NAME SIN VALOR EN PARAMETRO")
						return -1
					}
					hay_name = true

				} else {
					fmt.Println("ERROR NINGUN PARAMETRO COICIDE: " + parametros[i].nombre)
					return -1
				}

			}

		} //fin de for

		/*__________________________ CREAR RMGRP ______________________*/

		if hay_name {
			if Eliminar_Usurio(name, actualSesion.Path) == 0 {
				fmt.Println("USUARIO ELIMINADO CORRECTAMENTE")
				return 0
			} else {
				fmt.Println("ERROR USUARIO NO SE PUDO ELIMINAR")
				return -1
			}
		} else {
			fmt.Println("ERROR NO HAY NAME EN RMGRP")
			return -1
		}
	}

	return -1
}

/*__________________________________ FIN DE RMUSR _______________________________________*/

/*__________________________________  INICIO DE MKDIR___________________________________*/
func Reconocer_Mkdir(lista_comando string, comando_aux string) int {
	var path_crear string = ""
	var hay_path_crear bool = false
	var hay_p bool = false

	if "mkdir" == strings.ToLower(comando_aux) {
		parametros := Reconocer_parametros(lista_comando, comando_aux)
		for i := 0; i < len(parametros); i++ {

			if parametros[i].valor != "" && parametros[i].nombre != "" || parametros[i].valor == "" && parametros[i].nombre != "" {

				if "-path=" == strings.ToLower(parametros[i].nombre) {
					//var path_temp string = strings.Replace(parametros[i].valor, filepath.Ext(parametros[i].valor), ".dk", 1)
					path_crear = strings.ReplaceAll(parametros[i].valor, "\"", "")

					if path_crear == "" {
						fmt.Println("ERROR PATH SIN VALOR EN PARAMETRO MKDIR")
						return -1
					}
					hay_path_crear = true

				} else if "-p" == strings.ToLower(parametros[i].nombre) {
					if hay_p {
						fmt.Println("ERROR P SE REPITE EN PARAMETRO EN MKDIR")
						return -1
					}
					hay_p = true

				} else {
					fmt.Println("ERROR NINGUN PARAMETRO COICIDE: "+parametros[i].nombre, "EN MKDIR")
					return -1
				}

			}

		} //fin de for

		/*__________________________ CREAR MKDIR ______________________*/

		if !hay_path_crear {
			fmt.Println("ERROR NO HAY PATH EN MKDIR")
			return -1
		}
		if actualSesion.Id_user == 1 && actualSesion.Id_grp == 1 {
			if Crear_Mkdir(path_crear, hay_p) == 0 {
				//fmt.Println("CARPETAS CREADA CORRENTAMENTE")
				return 0
			}
		}

	}

	return -1
}

/*___________________________________ FIN DE MKDIR _____________________________*/

/*____________________________________ INICIO DE MKFILE _________________________*/
func Reconocer_Mkfile(lista_comando string, comando_aux string) int {
	var path_crear_archivo string = ""
	var hay_path_crear bool = false
	var size_archivo int = 0
	var cont_archivo string = ""

	var hay_r bool = false

	if "mkfile" == strings.ToLower(comando_aux) {
		parametros := Reconocer_parametros(lista_comando, comando_aux)
		for i := 0; i < len(parametros); i++ {

			if parametros[i].valor != "" && parametros[i].nombre != "" || parametros[i].valor == "" && parametros[i].nombre != "" {

				if "-path=" == strings.ToLower(parametros[i].nombre) {
					//var path_temp string = strings.Replace(parametros[i].valor, filepath.Ext(parametros[i].valor), ".dk", 1)
					path_crear_archivo = strings.ReplaceAll(parametros[i].valor, "\"", "")

					if path_crear_archivo == "" {
						fmt.Println("ERROR PATH SIN VALOR EN PARAMETRO MKFILE")
						return -1
					}
					hay_path_crear = true

				} else if "-r" == strings.ToLower(parametros[i].nombre) {
					if hay_r {
						fmt.Println("ERROR R SE REPITE EN PARAMETRO EN MKFILE")
						return -1
					}
					hay_r = true

				} else if "-size=" == strings.ToLower(parametros[i].nombre) {
					s, err := strconv.Atoi(parametros[i].valor)
					if err == nil && s > 0 {
						size_archivo = int(s)

					} else {
						fmt.Println("ERROR SIZE  EN PARAMETRO EN MKFILE")
						return -1
					}

				} else if "-cont=" == strings.ToLower(parametros[i].nombre) {
					cont_archivo = strings.ReplaceAll(parametros[i].valor, "\"", "")

					if cont_archivo == "" {
						fmt.Println("ERROR CONT SIN VALOR EN PARAMETRO MKFILE")
						return -1
					}

				} else {
					fmt.Println("ERROR NINGUN PARAMETRO COICIDE: "+parametros[i].nombre, "EN MKFILE")
					return -1
				}

			}

		} //fin de for

		/*__________________________ CREAR MKDIR ______________________*/
		if !hay_path_crear {
			fmt.Println("ERROR NO HAY PATH EN MKFILE")
			return -1
		}
		if actualSesion.Id_user == 1 && actualSesion.Id_grp == 1 {
			if Crear_Mkfile(path_crear_archivo, size_archivo, cont_archivo, hay_r) == 0 {
				//fmt.Println("CARPETAS CREADA CORRENTAMENTE")
				return 0
			}
		}

	}

	return -1
}

/*____________________________________ FIN DE MKFILE ____________________________*/

/*__________________________________ INICIO REPORTE _____________________________________*/
func Reconocer_Rep(lista_param string, comando_aux string) int {
	var name string = ""
	var id string = ""
	var path string = ""
	var path_reporte = ""
	var extension_rep = ""
	var ruta string = ""
	var hay_name bool = false
	var hay_path bool = false
	var hay_id bool = false

	//inicio de comando para parametros
	if "rep" == strings.ToLower(comando_aux) {
		parametros := Reconocer_parametros(lista_param, comando_aux)
		for i := 0; i < len(parametros); i++ {
			if parametros[i].nombre != "" && parametros[i].valor != "" {

				if "-name=" == strings.ToLower(parametros[i].nombre) {
					name = strings.ReplaceAll(parametros[i].valor, "\"", "")

					if name == "" {
						fmt.Println("ERROR NAME SIN VALOR EN PARAMETRO REPORTE")
						return -1
					}
					hay_name = true

				} else if "-path=" == strings.ToLower(parametros[i].nombre) {
					path_reporte = strings.ReplaceAll(parametros[i].valor, "\"", "")
					var tem_ext = strings.ReplaceAll(filepath.Ext(parametros[i].valor), ".", "")
					extension_rep = strings.ReplaceAll(tem_ext, "\"", "")

					var path_temp string = strings.Replace(parametros[i].valor, filepath.Ext(parametros[i].valor), ".dot", 1)
					path = strings.ReplaceAll(path_temp, "\"", "")

					if path == "" {
						fmt.Println("ERROR ID SIN VALOR EN PARAMETRO REPORTE")
						return -1
					}
					hay_path = true
				} else if "-id=" == strings.ToLower(parametros[i].nombre) {
					//var path_temp string = strings.Replace(parametros[i].valor, filepath.Ext(parametros[i].valor), ".dk", 1)
					id = strings.ReplaceAll(parametros[i].valor, "\"", "")

					if id == "" {
						fmt.Println("ERROR ID SIN VALOR EN PARAMETRO REPORTE")
						return -1
					}
					hay_id = true

				} else if "-ruta=" == strings.ToLower(parametros[i].nombre) {
					ruta = strings.ReplaceAll(parametros[i].valor, "\"", "")

					if ruta == "" {
						fmt.Println("ERROR RUTA SIN VALOR EN PARAMETRO REPORTE")
						return -1
					}
					//hay_id = true
				} else {
					fmt.Println("ERROR NINGUN PARAMETRO COICIDE: " + parametros[i].nombre)
					return -1
				}

			}

		} //fin for

		/*__________________________ CREAR DISCO ______________________*/
		if hay_name {
			if hay_path {
				if hay_id {
					//Crear_carpetas(path)
					if listaS.buscarParticion(id) {
						nodito := listaS.obtenerNodo(id)
						if name == "mbr" {
							if crear_rep_mbr(nodito.Path, path, path_reporte, extension_rep) == 0 {
								fmt.Println("REPORTE MBR CREADO EXITOSAMENTE EN REP")
							} else {
								fmt.Println("ERROR NO SE PUDO CREAR REPORTE MBR EN REP")
							}
						} else if name == "disk" {
							if crear_rep_disk(nodito.Path, path, path_reporte, extension_rep) == 0 {
								fmt.Println("REPORTE DISK CREADO EXITOSAMENTE EN REP")
							} else {
								fmt.Println("ERROR NO SE PUDO CREAR REPORTE DISK EN REP")
							}
						} else if name == "tree" {
							if crear_rep_tree(nodito.Path, path, path_reporte, extension_rep) == 0 {
								fmt.Println("REPORTE TREE CREADO EXITOSAMENTE EN REP")
							} else {
								fmt.Println("ERROR NO SE PUDO CREAR REPORTE TREE EN REP")
							}
						} else if name == "sb" {
							if crear_rep_super_bloque(nodito.Path, path, path_reporte, extension_rep) == 0 {
								fmt.Println("REPORTE SB CREADO EXITOSAMENTE EN REP")
							} else {
								fmt.Println("ERROR NO SE PUDO CREAR REPORTE SB EN REP")
							}
						} else if name == "file" {
							if crear_rep_file(nodito.Path, path, path_reporte, extension_rep, ruta) == 0 {
								fmt.Println("REPORTE FILE CREADO EXITOSAMENTE EN REP")
							} else {
								fmt.Println("ERROR NO SE PUDO CREAR REPORTE FILE EN REP")
							}
						}

					} else {
						fmt.Println("ERROR DISCO NO MONTADO PARA CREAR REPORTE")
					}

				} else {
					fmt.Println("ERROR NO HAY ID EN REPORTE")
				}
			} else {
				fmt.Println("ERROR NO HAY PATH EN REPORTE")
			}
		} else {
			fmt.Println("ERROR NO HAY  NAME REPORTE")
		}

		return 0
	}
	return -1
} /*______________________ FIN DE REPORTE _____________________*/
