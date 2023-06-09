package main

import (
	"fmt"
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
				if strings.ToLower(param) == ">r" {
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

func Reconocer_Comando(texto_comando string) Respuesta {
	fmt.Println(texto_comando)
	mensajeRespuesta := ""
	respuestaStruct := Respuesta{
		Tipo:    1,
		Mensaje: "",
		Data:    "",
	}
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
			mensajeRespuesta = Reconocer_mkdisk(texto_comando, ins_aux)
		case "fdisk":
			mensajeRespuesta = Reconocer_Fdisk(texto_comando, ins_aux)
		case "mkfs":
			mensajeRespuesta = Reconocer_Mkfs(texto_comando, ins_aux)
		case "rmdisk":
			mensajeRespuesta = Reconocer_Rmdisk(texto_comando, ins_aux)
		case "mount":
			mensajeRespuesta = Reconocer_Mount(texto_comando, ins_aux)
		case "login":
			mensajeRespuesta = Reconocer_Login(texto_comando, ins_aux)
		case "logout":
			mensajeRespuesta = Logout()
		case "mkgrp":
			mensajeRespuesta = Reconocer_Mkgrp(texto_comando, ins_aux)
		case "rmgrp":
			mensajeRespuesta = Reconocer_Rmgrp(texto_comando, ins_aux)
		case "mkuser":
			mensajeRespuesta = Reconocer_Mkusr(texto_comando, ins_aux)
		case "rmusr":
			mensajeRespuesta = Reconocer_Rmusr(texto_comando, ins_aux)
		case "mkdir":
			mensajeRespuesta = Reconocer_Mkdir(texto_comando, ins_aux)
		case "mkfile":
			mensajeRespuesta = Reconocer_Mkfile(texto_comando, ins_aux)
		case "rep":
			respuestaStruct = Reconocer_Rep(texto_comando, ins_aux)
		case "pause":
			mensajeRespuesta = "PAUSE"
		default:
			mensajeRespuesta = "COMANDO NO RECONOCIDO0"
		}

	} else {
		mensajeRespuesta = texto_comando
	}

	if respuestaStruct.Mensaje != "" {
		return respuestaStruct
	}

	return Respuesta{
		Tipo:    0,
		Mensaje: mensajeRespuesta,
		Data:    "",
	}
}

func Reconocer_mkdisk(lista_param string, comando_aux string) string {
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

				if ">size=" == strings.ToLower(parametros[i].nombre) {
					s, err := strconv.Atoi(parametros[i].valor)
					if err == nil && s > 0 {
						size = int32(s)
						hay_size = true

					} else {
						return "El parametro SIZE debe ser mayor a cero"
					}

				} else if ">unit=" == strings.ToLower(parametros[i].nombre) {
					if strings.ToLower(parametros[i].valor) == "" {
						fmt.Println("ERROR UNIT SIN VALOR EN PARAMETRO")
						return "Valor del parametro UNIT requerido"
					} else if strings.ToLower(parametros[i].valor) == "k" {
						unit = strings.ToLower(parametros[i].valor)
					} else if strings.ToLower(parametros[i].valor) == "m" {
						unit = strings.ToLower(parametros[i].valor)
					} else if strings.ToLower(parametros[i].valor) == "b" {
						unit = strings.ToLower(parametros[i].valor)
					} else {
						unit = strings.ToLower(parametros[i].valor)
						fmt.Println("ERROR UNIT PARAMETRO INCORRECTO: " + unit)
						return "Valor del parametro UNIT incorrecto"
					}
				} else if ">path=" == strings.ToLower(parametros[i].nombre) {
					var path_temp string = strings.Replace(parametros[i].valor, filepath.Ext(parametros[i].valor), ".dsk", 1)
					path = strings.ReplaceAll(path_temp, "\"", "")
					if path == "" {
						fmt.Println("ERROR PATH SIN VALOR EN PARAMETRO")
						return "Valor del parametro PATH requerido"
					}
					hay_path = true

				} else if ">fit=" == strings.ToLower(parametros[i].nombre) {

					if strings.ToLower(parametros[i].valor) == "" {
						fmt.Println("ERROR FIT SIN VALOR EN PARAMETRO")
						return "Valor del parametro FIT requerido"
					} else if strings.ToLower(parametros[i].valor) == "wf" {
						fit = strings.ToLower(parametros[i].valor)
					} else if strings.ToLower(parametros[i].valor) == "ff" {
						fit = strings.ToLower(parametros[i].valor)
					} else if strings.ToLower(parametros[i].valor) == "bf" {
						fit = strings.ToLower(parametros[i].valor)
					} else {
						fit = strings.ToLower(parametros[i].valor)
						fmt.Println("ERROR FIT PARAMETRO INCORRECTO: " + fit)
						return "Valor del parametro FIT incorrecto"
					}
				} else {
					fmt.Println("ERROR NINGUN PARAMETRO COICIDE: " + parametros[i].nombre)
					return "No se reconoce el parametro " + parametros[i].nombre
				}

			}

		} //fin for

		/*__________________________ CREAR DISCO ______________________*/
		if hay_size {
			if hay_path {
				Crear_carpetas(path)
				return Crear_archivo(size, path, fit, unit)
			} else {
				fmt.Println("ERROR NO HAY PATH EN MKDISK")
				return "Parametro PATH requerido"
			}
		} else {
			fmt.Println("ERROR NO HAY SIZE MKDISK")
			return "Parametro SIZE requerido"
		}
	}
	return "Ocurrio un error al tratar de crear el disco"
} /*______________________ FIN DE ANALIZAR COMANDO MKDISK_____________________*/

func Reconocer_Rmdisk(lista_comando string, comando_aux string) string {
	var hay_path bool = false
	var path string = ""
	if "rmdisk" == strings.ToLower(comando_aux) {

		parametros := Reconocer_parametros(lista_comando, comando_aux)

		for i := 0; i < len(parametros); i++ {
			if parametros[i].valor != "" && parametros[i].nombre != "" {

				if ">path=" == strings.ToLower(parametros[i].nombre) {
					path = strings.ReplaceAll(parametros[i].valor, "\"", "")
					hay_path = true
				} else {
					return "No se reconoce el parametro " + parametros[i].nombre
				}

			}

		} //fin de for

		fmt.Println("MIPAT " + path)
		/*__________________________ ELIMINAR ARCHIVO ______________________*/
		if hay_path {
			return Ejecutar_rmdisk(path)
		} else {
			fmt.Println("ERROR NO HAY PATH EN RMDISK")
			return "Parametro PATH requerido"
		}
	}
	return "No se reconoce comando"
}

/*________________________ INICIO DE COMANDO FDISK ____________________________*/
func Reconocer_Fdisk(lista_comando string, comando_aux string) string {
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
				if ">size=" == strings.ToLower(parametros[i].nombre) {
					s, err := strconv.Atoi(parametros[i].valor)
					if err == nil && s > 0 {
						size = int32(s)
						hay_size = true
					} else {
						return "El parametro SIZE debe ser mayor a cero"
					}
				} else if ">unit=" == strings.ToLower(parametros[i].nombre) {
					if strings.ToLower(parametros[i].valor) == "" {
						fmt.Println("ERROR UNIT SIN VALOR EN PARAMETRO")
						return "Valor del parametro UNIT requerido"
					} else if strings.ToLower(parametros[i].valor) == "k" {
						unit = strings.ToLower(parametros[i].valor)

					} else if strings.ToLower(parametros[i].valor) == "m" {
						unit = strings.ToLower(parametros[i].valor)

					} else if strings.ToLower(parametros[i].valor) == "b" {
						unit = strings.ToLower(parametros[i].valor)

					} else {
						unit = strings.ToLower(parametros[i].valor)
						fmt.Println("ERROR UNIT PARAMETRO INCORRECTO: " + unit)
						return "Valor del parametro UNIT incorrecto"
					}
				} else if ">path=" == strings.ToLower(parametros[i].nombre) {
					var path_temp string = strings.Replace(parametros[i].valor, filepath.Ext(parametros[i].valor), ".dsk", 1)
					path = strings.ReplaceAll(path_temp, "\"", "")
					if path == "" {
						fmt.Println("ERROR PATH SIN VALOR EN PARAMETRO")
						return "Valor del parametro PATH requerido"
					}
					hay_path = true

				} else if ">fit=" == strings.ToLower(parametros[i].nombre) {

					if strings.ToLower(parametros[i].valor) == "" {
						fmt.Println("ERROR FIT SIN VALOR EN PARAMETRO")
						return "Valor del parametro FIT requerido"
					} else if strings.ToLower(parametros[i].valor) == "wf" {
						fit = strings.ToLower(parametros[i].valor)

					} else if strings.ToLower(parametros[i].valor) == "ff" {
						fit = strings.ToLower(parametros[i].valor)

					} else if strings.ToLower(parametros[i].valor) == "bf" {
						fit = strings.ToLower(parametros[i].valor)

					} else {
						//fit = strings.ToLower(parametros[i].valor)
						fmt.Println("ERROR FIT PARAMETRO INCORRECTO: " + strings.ToLower(parametros[i].valor))
						return "Valor del parametro FIT incorrecto"
					}
				} else if ">type=" == strings.ToLower(parametros[i].nombre) {
					//type
					if strings.ToLower(parametros[i].valor) == "" {
						fmt.Println("ERROR TYPE SIN VALOR EN PARAMETRO")
						return "Valor del parametro TYPE requerido"
					} else if strings.ToLower(parametros[i].valor) == "p" {
						type_ = strings.ToLower(parametros[i].valor)

					} else if strings.ToLower(parametros[i].valor) == "l" {
						type_ = strings.ToLower(parametros[i].valor)

					} else if strings.ToLower(parametros[i].valor) == "e" {
						type_ = strings.ToLower(parametros[i].valor)

					} else {
						fmt.Println("ERROR TYPE PARAMETRO INCORRECTO: " + strings.ToLower(parametros[i].valor))
						return "Valor del parametro TYPE incorrecto"
					}

				} else if ">name=" == strings.ToLower(parametros[i].nombre) {
					//name
					if parametros[i].valor == "" {
						fmt.Println("ERROR NAME SIN VALOR EN PARAMETRO")
						return "Valor del parametro NAME requerido"
					}
					name = parametros[i].valor
					hay_name = true

				} else {
					fmt.Println("ERROR NINGUN PARAMETRO COICIDE: " + parametros[i].nombre)
					return "No se reconoce el parametro " + parametros[i].nombre
				}

			}

		} //fin de for

		/*__________________________ CREAR PARTICION ______________________*/
		if hay_size {
			if hay_path {
				if hay_name {
					return Crear_particion(size, unit, path, type_, fit, name)
				} else {
					fmt.Println("ERROR NO HAY NAME EN FDISK")
					return "Parametro NAME requerido"
				}
			} else {
				fmt.Println("ERROR NO HAY PATH EN FDISK")
				return "Parametro PATH requerido"
			}
		} else {
			fmt.Println("ERROR NO HAY SIZE FDISK")
			return "Parametro SIZE requerido"
		}
	}

	return "Ocurrio un error en FDISK"
}

/*________________________ FIN DE COMANDO FDISK _______________________________*/

/*________________________ INICIO DE COMANDO MOUNT ____________________________*/
func Reconocer_Mount(lista_comando string, comando_aux string) string {
	var path string = ""
	var name string = ""
	var hay_path bool = false
	var hay_name bool = false

	if "mount" == strings.ToLower(comando_aux) {
		parametros := Reconocer_parametros(lista_comando, comando_aux)
		for i := 0; i < len(parametros); i++ {
			if parametros[i].valor != "" && parametros[i].nombre != "" {
				if ">path=" == strings.ToLower(parametros[i].nombre) {
					var path_temp string = strings.Replace(parametros[i].valor, filepath.Ext(parametros[i].valor), ".dsk", 1)
					path = strings.ReplaceAll(path_temp, "\"", "")

					if path == "" {
						fmt.Println("ERROR PATH SIN VALOR EN PARAMETRO")
						return "Valor del parametro PATH incorrecto"
					}
					hay_path = true

				} else if ">name=" == strings.ToLower(parametros[i].nombre) {
					//name
					if parametros[i].valor == "" {
						fmt.Println("ERROR NAME SIN VALOR EN PARAMETRO")
						return "Valor del parametro NAME incorrecto"
					}
					name = parametros[i].valor
					hay_name = true
				} else {
					fmt.Println("ERROR NINGUN PARAMETRO COICIDE: " + parametros[i].nombre)
					return "No se reconoce el parametro " + parametros[i].nombre
				}

			}

		} //fin de for

		/*__________________________ CREAR MOUNT ______________________*/

		if hay_path {
			if hay_name {
				if !listaS.buscar_nombre_path(path, name) {
					return listaS.insertar_Nodo(path, name)
				} else {
					fmt.Println("ERROR PARTICION YA MONTADA....")
					return "La particion " + name + " ya se encuentra montada"
				}
			} else {
				fmt.Println("ERROR NO HAY NAME EN MOUNT")
				return "Parametro NAME requerido"
			}

		} else {
			fmt.Println("ERROR NO HAY PATH EN MOUNT")
			return "Parametro PATH requerido"
		}
	}
	return "NO SE RECONOCE COMANDO: " + strings.ToLower(comando_aux)
}

/*____________________________________ FIN DE COMANDO MOUNT _____________________________*/

/*____________________________________ INICIO DE COMANDO MKFS ____________________________*/

func Reconocer_Mkfs(lista_comando string, comando_aux string) string {
	var id string = ""
	var type_ string = ""
	var hay_id bool = false

	if "mkfs" == strings.ToLower(comando_aux) {
		parametros := Reconocer_parametros(lista_comando, comando_aux)
		for i := 0; i < len(parametros); i++ {

			if parametros[i].valor != "" && parametros[i].nombre != "" {
				if ">id=" == strings.ToLower(parametros[i].nombre) {
					//var path_temp string = strings.Replace(parametros[i].valor, filepath.Ext(parametros[i].valor), ".dk", 1)
					id = strings.ReplaceAll(parametros[i].valor, "\"", "")

					if id == "" {
						fmt.Println("ERROR ID SIN VALOR EN PARAMETRO")
						return "Valor del parametro ID incorrecto"
					}
					hay_id = true

				} else if ">type=" == strings.ToLower(parametros[i].nombre) {
					if parametros[i].valor == "" {
						fmt.Println("ERROR TYPE SIN VALOR EN PARAMETRO")
						return "Valor del parametro TYPE incorrecto"
					}
					parametros[i].valor = strings.ReplaceAll(parametros[i].valor, "\"", "")
					type_ = strings.ToLower(parametros[i].valor)
				} else {
					fmt.Println("ERROR NINGUN PARAMETRO COICIDE: " + parametros[i].nombre)
					return "No se reconoce el parametro " + parametros[i].nombre
				}

			}

		} //fin de for

		/*__________________________ CREAR MKFS ______________________*/

		if hay_id {
			return EjecutarMkfs(id, type_)
		} else {
			fmt.Println("ERROR NO HAY ID EN MKFS")
			return "Parametro ID requerido"
		}
	}
	return "NO SE RECONOCE EL COMANDO"
}

/*____________________________________ FIN DE COMANDO MKFS _______________________________*/

/*___________________________________ INCIO DE LOGIN _____________________________________*/
func Reconocer_Login(lista_comando string, comando_aux string) string {
	var id string = ""
	var password string = ""
	var user string = ""
	var hay_id bool = false
	var hay_password = false
	var hay_user = false

	if "login" == strings.ToLower(comando_aux) {
		parametros := Reconocer_parametros(lista_comando, comando_aux)
		for i := 0; i < len(parametros); i++ {
			//parametros[i].valor = strings.ReplaceAll(parametros[i].valor, "\"", "")
			if parametros[i].valor != "" && parametros[i].nombre != "" {
				if ">id=" == strings.ToLower(parametros[i].nombre) {
					//var path_temp string = strings.Replace(parametros[i].valor, filepath.Ext(parametros[i].valor), ".dk", 1)
					id = strings.ReplaceAll(parametros[i].valor, "\"", "")

					if id == "" {
						fmt.Println("ERROR ID SIN VALOR EN PARAMETRO")
						return "Valor del parametro ID incorrecto"
					}
					hay_id = true

				} else if ">pwd=" == strings.ToLower(parametros[i].nombre) {
					//name
					if parametros[i].valor == "" {
						fmt.Println("ERROR PASSWORD SIN VALOR EN PARAMETRO")
						return "Valor del parametro PWD incorrecto"
					}
					//password = strings.ToLower(parametros[i].valor)
					password = parametros[i].valor
					hay_password = true

				} else if ">user=" == strings.ToLower(parametros[i].nombre) {
					if parametros[i].valor == "" {
						fmt.Println("ERROR USUARIO SIN VALOR EN PARAMETRO")
						return "Valor del parametro USER incorrecto"
					}
					//user = strings.ToLower(parametros[i].valor)
					user = parametros[i].valor
					hay_user = true

				} else {
					fmt.Println("ERROR NINGUN PARAMETRO COICIDE: " + parametros[i].nombre)
					return "Parametro no reconocido"
				}

			}

		} //fin de for

		/*__________________________ INICIAR LOGIN ______________________*/

		if hay_id {
			if hay_user {
				if hay_password {
					var aux_nodo *NODO = listaS.obtenerNodo(id)
					if aux_nodo != nil {
						return Log_in(aux_nodo.Path, id, user, password)
					} else {
						fmt.Println("ERROR ID NO SE ENCOTRO EN LOGIN: ", id)
						return "No se encuentra la particion con el ID ingresado"
					}

				} else {
					fmt.Println("ERROR NO HAY PASSWORD EN LOGIN")
					return "Parametro PWD requerido"
				}
			} else {
				fmt.Println("ERROR NO HAY USER EN LOGIN")
				return "Parametro USER requerido"
			}

		} else {
			fmt.Println("ERROR NO HAY ID EN LOGIN")
			return "Parametro ID requerido"
		}
	}
	return "NO SE RECONOCE EL COMANDO"
}

/*___________________________________ FIN DE LOGIN _______________________________________*/

/*__________________________________  INICIO DE MKGRP ___________________________________*/
func Reconocer_Mkgrp(lista_comando string, comando_aux string) string {
	var name string = ""
	var hay_name bool = false

	if "mkgrp" == strings.ToLower(comando_aux) {
		parametros := Reconocer_parametros(lista_comando, comando_aux)
		for i := 0; i < len(parametros); i++ {

			if parametros[i].valor != "" && parametros[i].nombre != "" {

				if ">name=" == strings.ToLower(parametros[i].nombre) {
					//var path_temp string = strings.Replace(parametros[i].valor, filepath.Ext(parametros[i].valor), ".dk", 1)
					name = strings.ReplaceAll(parametros[i].valor, "\"", "")

					if name == "" {
						fmt.Println("ERROR NAME SIN VALOR EN PARAMETRO")
						return "Valor del parametro NAME incorrecto"
					}
					hay_name = true

				} else {
					fmt.Println("ERROR NINGUN PARAMETRO COICIDE: " + parametros[i].nombre)
					return "Parametro no reconocido"
				}

			}

		} //fin de for

		/*__________________________ CREAR MKGRP ______________________*/

		if hay_name {
			return Crear_Mkgrp(name)
		} else {
			fmt.Println("ERROR NO HAY NAME EN MKGRP")
			return "Parametro NAME requerido"
		}
	}

	return "COMANDO NO RECONOCIDO"
}

/*__________________________________  FIN DE MKGRP ______________________________________*/

/*__________________________________ INICO DE RMGRP ____________________________________*/
func Reconocer_Rmgrp(lista_comando string, comando_aux string) string {
	var name string = ""
	var hay_name bool = false

	if "rmgrp" == strings.ToLower(comando_aux) {
		parametros := Reconocer_parametros(lista_comando, comando_aux)
		for i := 0; i < len(parametros); i++ {

			if parametros[i].valor != "" && parametros[i].nombre != "" {

				if ">name=" == strings.ToLower(parametros[i].nombre) {
					//var path_temp string = strings.Replace(parametros[i].valor, filepath.Ext(parametros[i].valor), ".dk", 1)
					name = strings.ReplaceAll(parametros[i].valor, "\"", "")

					if name == "" {
						fmt.Println("ERROR NAME SIN VALOR EN PARAMETRO")
						return "Valor del parametro NAME incorrecto"
					}
					hay_name = true

				} else {
					fmt.Println("ERROR NINGUN PARAMETRO COICIDE: " + parametros[i].nombre)
					return "Parametro no reconocido"
				}

			}

		} //fin de for

		/*__________________________ CREAR RMGRP ______________________*/

		if hay_name {
			return Eliminar_Grupo(name, actualSesion.Path)
		} else {
			fmt.Println("ERROR NO HAY NAME EN RMGRP")
			return "Parametro NAME requerido"
		}
	}

	return "COMANDO NO RECONOCIDO"
}

/*__________________________________ FIN DE RMGRP _____________________________________*/

/*__________________________________  INICIO DE MKUSR ___________________________________*/
func Reconocer_Mkusr(lista_comando string, comando_aux string) string {
	var usuario string = ""
	var pass_word = ""
	var grupo_incluir = ""
	var hay_usuario bool = false
	var hay_pass bool = false
	var hay_grupo bool = false

	if "mkuser" == strings.ToLower(comando_aux) {
		parametros := Reconocer_parametros(lista_comando, comando_aux)
		for i := 0; i < len(parametros); i++ {

			if parametros[i].valor != "" && parametros[i].nombre != "" {

				if ">user=" == strings.ToLower(parametros[i].nombre) {
					//var path_temp string = strings.Replace(parametros[i].valor, filepath.Ext(parametros[i].valor), ".dk", 1)
					usuario = strings.ReplaceAll(parametros[i].valor, "\"", "")

					if usuario == "" {
						fmt.Println("ERROR USUARIO SIN VALOR EN PARAMETRO MKUSR")
						return "Valor del parametro USER incorrecto"
					}
					hay_usuario = true

				} else if ">pwd=" == strings.ToLower(parametros[i].nombre) {
					//var path_temp string = strings.Replace(parametros[i].valor, filepath.Ext(parametros[i].valor), ".dk", 1)
					pass_word = strings.ReplaceAll(parametros[i].valor, "\"", "")

					if pass_word == "" {
						fmt.Println("ERROR PWD SIN VALOR EN PARAMETRO EN MKUSR")
						return "Valor del parametro PWD incorrecto"
					}
					hay_pass = true

				} else if ">grp=" == strings.ToLower(parametros[i].nombre) {
					//var path_temp string = strings.Replace(parametros[i].valor, filepath.Ext(parametros[i].valor), ".dk", 1)
					grupo_incluir = strings.ReplaceAll(parametros[i].valor, "\"", "")

					if grupo_incluir == "" {
						fmt.Println("ERROR GRUPO SIN VALOR EN PARAMETRO MKUSR")
						return "Valor del parametro GRP incorrecto"
					}
					hay_grupo = true

				} else {
					fmt.Println("ERROR NINGUN PARAMETRO COICIDE: "+parametros[i].nombre, "EN MKUSER")
					return "Parametro no reconocido " + parametros[i].nombre
				}

			}

		} //fin de for

		/*__________________________ CREAR MKUSR ______________________*/

		if !hay_usuario {
			fmt.Println("ERROR NO HAY USUARIO EN MKUSR")
			return "Parametro USER requerido"
		}
		if !hay_pass {
			fmt.Println("ERROR NO HAY PASS EN MKUSR")
			return "Parametro PWD requerido"
		}
		if !hay_grupo {
			fmt.Println("ERROR NO HAY GRUPO EN MKUSR")
			return "Parametro GRP requerido"
		}
		if actualSesion.Id_user == 1 && actualSesion.Id_grp == 1 {
			return Crear_Mkusr(usuario, pass_word, grupo_incluir)
		} else {
			return "MKUSER solo puede ser ejecutato por el usuario root"
		}

	}

	return "COMANDO NO RECONOCIDO"
}

/*__________________________________  FIN DE MKUSR ______________________________________*/

/*__________________________________ INICO DE RMUSR ____________________________________*/
func Reconocer_Rmusr(lista_comando string, comando_aux string) string {
	var name string = ""
	var hay_name bool = false

	if "rmusr" == strings.ToLower(comando_aux) {
		parametros := Reconocer_parametros(lista_comando, comando_aux)
		for i := 0; i < len(parametros); i++ {

			if parametros[i].valor != "" && parametros[i].nombre != "" {

				if ">user=" == strings.ToLower(parametros[i].nombre) {
					name = strings.ReplaceAll(parametros[i].valor, "\"", "")

					if name == "" {
						fmt.Println("ERROR NAME SIN VALOR EN PARAMETRO")
						return "Valor del parametro NAME incorrecto"
					}
					hay_name = true

				} else {
					fmt.Println("ERROR NINGUN PARAMETRO COICIDE: " + parametros[i].nombre)
					return "Parametro no reconocido"
				}

			}

		} //fin de for

		/*__________________________ CREAR RMGRP ______________________*/

		if hay_name {
			return Eliminar_Usurio(name, actualSesion.Path)
		} else {
			fmt.Println("ERROR NO HAY NAME EN RMGRP")
			return "Parametro NAME requerido"
		}
	}

	return "COMANDO NO RECONOCIDO1"
}

/*__________________________________ FIN DE RMUSR _______________________________________*/

/*__________________________________  INICIO DE MKDIR___________________________________*/
func Reconocer_Mkdir(lista_comando string, comando_aux string) string {
	var path_crear string = ""
	var hay_path_crear bool = false
	var hay_p bool = false

	if "mkdir" == strings.ToLower(comando_aux) {
		parametros := Reconocer_parametros(lista_comando, comando_aux)
		for i := 0; i < len(parametros); i++ {
			if parametros[i].valor != "" && parametros[i].nombre != "" || parametros[i].valor == "" && parametros[i].nombre != "" {

				if ">path=" == strings.ToLower(parametros[i].nombre) {
					//var path_temp string = strings.Replace(parametros[i].valor, filepath.Ext(parametros[i].valor), ".dk", 1)
					path_crear = strings.ReplaceAll(parametros[i].valor, "\"", "")

					if path_crear == "" {
						fmt.Println("ERROR PATH SIN VALOR EN PARAMETRO MKDIR")
						return "Valor del parametro PATH incorrecto"
					}
					hay_path_crear = true

				} else if ">r" == strings.ToLower(parametros[i].nombre) {
					if hay_p {
						fmt.Println("ERROR R SE REPITE EN PARAMETRO EN MKDIR")
						//RETURN
					}
					hay_p = true

				} else {
					fmt.Println("ERROR NINGUN PARAMETRO COICIDE: "+parametros[i].nombre, "EN MKDIR")
					return "No se reconoce el parametro"
				}

			}

		} //fin de for

		/*__________________________ CREAR MKDIR ______________________*/

		if !hay_path_crear {
			fmt.Println("ERROR NO HAY PATH EN MKDIR")
			return "Parametro PATH requerido"
		}
		return Crear_Mkdir(path_crear, hay_p)

	}

	return "COMANDO NO RECONOCIDO1"
}

/*___________________________________ FIN DE MKDIR _____________________________*/

/*____________________________________ INICIO DE MKFILE _________________________*/
func Reconocer_Mkfile(lista_comando string, comando_aux string) string {
	var path_crear_archivo string = ""
	var hay_path_crear bool = false
	var size_archivo int = 0
	var cont_archivo string = ""

	var hay_r bool = false

	if "mkfile" == strings.ToLower(comando_aux) {
		parametros := Reconocer_parametros(lista_comando, comando_aux)
		for i := 0; i < len(parametros); i++ {

			if parametros[i].valor != "" && parametros[i].nombre != "" || parametros[i].valor == "" && parametros[i].nombre != "" {

				if ">path=" == strings.ToLower(parametros[i].nombre) {
					//var path_temp string = strings.Replace(parametros[i].valor, filepath.Ext(parametros[i].valor), ".dk", 1)
					path_crear_archivo = strings.ReplaceAll(parametros[i].valor, "\"", "")

					if path_crear_archivo == "" {
						fmt.Println("ERROR PATH SIN VALOR EN PARAMETRO MKFILE")
						return "Valor del parametro PATH incorrecto"
					}
					hay_path_crear = true

				} else if ">r" == strings.ToLower(parametros[i].nombre) {
					hay_r = true

				} else if ">size=" == strings.ToLower(parametros[i].nombre) {
					s, err := strconv.Atoi(parametros[i].valor)
					if err == nil && s > 0 {
						size_archivo = int(s)

					} else {
						fmt.Println("ERROR SIZE  EN PARAMETRO EN MKFILE")
						return "Valor del parametro SIZE incorrecto"
					}

				} else if ">cont=" == strings.ToLower(parametros[i].nombre) {
					cont_archivo = strings.ReplaceAll(parametros[i].valor, "\"", "")

					if cont_archivo == "" {
						fmt.Println("ERROR CONT SIN VALOR EN PARAMETRO MKFILE")
						return "Valor del parametro COUNT incorrecto"
					}

				} else {
					fmt.Println("ERROR NINGUN PARAMETRO COICIDE: "+parametros[i].nombre, "EN MKFILE")
					return "No se reconoce el parametro"
				}

			}

		} //fin de for

		/*__________________________ CREAR MKDIR ______________________*/
		if !hay_path_crear {
			fmt.Println("ERROR NO HAY PATH EN MKFILE")
			return "Parametro PATH requerido"
		}
		if actualSesion.Id_user == 1 && actualSesion.Id_grp == 1 {
			return Crear_Mkfile(path_crear_archivo, size_archivo, cont_archivo, hay_r)
		}

	}

	return "NO SE RECONOCE EL COMANDO1"
}

/*____________________________________ FIN DE MKFILE ____________________________*/

/*__________________________________ INICIO REPORTE _____________________________________*/
func Reconocer_Rep(lista_param string, comando_aux string) Respuesta {
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

				if ">name=" == strings.ToLower(parametros[i].nombre) {
					name = strings.ReplaceAll(parametros[i].valor, "\"", "")

					if name == "" {
						fmt.Println("ERROR NAME SIN VALOR EN PARAMETRO REPORTE")
						return Respuesta{
							Tipo:    0,
							Mensaje: "Valor del parametro NAME incorrecto",
							Data:    "",
						}
					}
					hay_name = true

				} else if ">path=" == strings.ToLower(parametros[i].nombre) {
					path_reporte = strings.ReplaceAll(parametros[i].valor, "\"", "")
					var tem_ext = strings.ReplaceAll(filepath.Ext(parametros[i].valor), ".", "")
					extension_rep = strings.ReplaceAll(tem_ext, "\"", "")

					var path_temp string = strings.Replace(parametros[i].valor, filepath.Ext(parametros[i].valor), ".dot", 1)
					path = strings.ReplaceAll(path_temp, "\"", "")

					if path == "" {
						fmt.Println("ERROR ID SIN VALOR EN PARAMETRO REPORTE")
						return Respuesta{
							Tipo:    0,
							Mensaje: "Valor del parametro PATH incorrecto",
							Data:    "",
						}
					}
					hay_path = true
				} else if ">id=" == strings.ToLower(parametros[i].nombre) {
					//var path_temp string = strings.Replace(parametros[i].valor, filepath.Ext(parametros[i].valor), ".dk", 1)
					id = strings.ReplaceAll(parametros[i].valor, "\"", "")

					if id == "" {
						fmt.Println("ERROR ID SIN VALOR EN PARAMETRO REPORTE")
						return Respuesta{
							Tipo:    0,
							Mensaje: "Valor del parametro ID incorrecto",
							Data:    "",
						}
					}
					hay_id = true

				} else if ">ruta=" == strings.ToLower(parametros[i].nombre) {
					ruta = strings.ReplaceAll(parametros[i].valor, "\"", "")

					if ruta == "" {
						fmt.Println("ERROR RUTA SIN VALOR EN PARAMETRO REPORTE")
						return Respuesta{
							Tipo:    0,
							Mensaje: "Valor del parametro RUTA incorrecto",
							Data:    "",
						}
					}
					//hay_id = true
				} else {
					fmt.Println("ERROR NINGUN PARAMETRO COICIDE: " + parametros[i].nombre)
					return Respuesta{
						Tipo:    0,
						Mensaje: "No se reconoce el parametro",
						Data:    "",
					}
				}

			}

		} //fin for

		/*__________________________ CREAR DISCO ______________________*/
		if hay_name {
			if hay_path {
				if hay_id {
					//Crear_carpetas(path)
					if listaS.buscarParticion(id) {
						if !actualSesion.hay_Sesion {
							return Respuesta{
								Tipo:    0,
								Mensaje: "No hay una sesión activa REP-LOGIN",
							}
						}
						nodito := listaS.obtenerNodo(id)
						reporteRespuesta := Respuesta{
							Tipo:    0,
							Mensaje: "",
							Data:    "",
						}
						if name == "mbr" {
							reporteRespuesta = crear_rep_mbr(nodito.Path, path, path_reporte, extension_rep)
							if reporteRespuesta.Tipo == 0 {
								fmt.Println("REPORTE MBR CREADO EXITOSAMENTE EN REP")
							} else {
								fmt.Println("ERROR NO SE PUDO CREAR REPORTE MBR EN REP")
								return Respuesta{
									Tipo:    0,
									Mensaje: "No se pudo crear el reporte MBR",
									Data:    "",
								}
							}
						} else if name == "disk" {
							reporteRespuesta = crear_rep_disk(nodito.Path, path, path_reporte, extension_rep)
							if reporteRespuesta.Tipo == 0 {
								fmt.Println("REPORTE DISK CREADO EXITOSAMENTE EN REP")
							} else {
								fmt.Println("ERROR NO SE PUDO CREAR REPORTE DISK EN REP")
								return Respuesta{
									Tipo:    0,
									Mensaje: "No se pudo crear el reporte REP",
									Data:    "",
								}
							}
						} else if name == "tree" {
							reporteRespuesta = crear_rep_tree(nodito.Path, path, path_reporte, extension_rep)
							if reporteRespuesta.Tipo == 0 {
								fmt.Println("REPORTE TREE CREADO EXITOSAMENTE EN REP")
							} else {
								fmt.Println("ERROR NO SE PUDO CREAR REPORTE TREE EN REP")
								return Respuesta{
									Tipo:    0,
									Mensaje: "No se pudo crear el reporte TREE",
									Data:    "",
								}
							}
						} else if name == "sb" {
							reporteRespuesta = crear_rep_super_bloque(nodito.Path, path, path_reporte, extension_rep)
							if reporteRespuesta.Tipo == 0 {
								fmt.Println("REPORTE SB CREADO EXITOSAMENTE EN REP")
							} else {
								fmt.Println("ERROR NO SE PUDO CREAR REPORTE SB EN REP")
								return Respuesta{
									Tipo:    0,
									Mensaje: "No se pudo crear el reporte SUPERBLOQUE",
									Data:    "",
								}
							}
						} else if name == "file" {
							reporteRespuesta = crear_rep_file(nodito.Path, path, path_reporte, extension_rep, ruta)
							if reporteRespuesta.Tipo == 0 {
								fmt.Println("REPORTE FILE CREADO EXITOSAMENTE EN REP")
							} else {
								fmt.Println("ERROR NO SE PUDO CREAR REPORTE FILE EN REP")
								return Respuesta{
									Tipo:    0,
									Mensaje: "No se pudo crear el reporte FILE",
									Data:    "",
								}
							}
						}
						reporteRespuesta.Tipo = 2
						reporteRespuesta.Ruta = path_reporte
						reporteRespuesta.Mensaje = "REPORTE CREADO EXITOSAMENTE"
						return reporteRespuesta
					} else {
						fmt.Println("ERROR DISCO NO MONTADO PARA CREAR REPORTE")
						return Respuesta{
							Tipo:    0,
							Mensaje: "La particion indicada no se encuentra montada",
							Data:    "",
						}
					}

				} else {
					fmt.Println("ERROR NO HAY ID EN REPORTE")
					return Respuesta{
						Tipo:    0,
						Mensaje: "Parametro ID requerido",
						Data:    "",
					}
				}
			} else {
				fmt.Println("ERROR NO HAY PATH EN REPORTE")
				return Respuesta{
					Tipo:    0,
					Mensaje: "Parametro PATH requerido",
					Data:    "",
				}
			}
		} else {
			fmt.Println("ERROR NO HAY  NAME REPORTE")
			return Respuesta{
				Tipo:    0,
				Mensaje: "Parametro NAME requerido",
				Data:    "",
			}
		}
	}
	return Respuesta{
		Tipo:    0,
		Mensaje: "NO SE RECONOCE EL COMANDO1",
		Data:    "",
	}
} /*______________________ FIN DE REPORTE _____________________*/
