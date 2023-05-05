package main

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"unsafe"
)

func crear_rep_mbr(direccion_disco string, destino_reporte_dot string, destino_reporte_grafico string, extension string) Respuesta {
	Crear_carpetas(destino_reporte_dot) //para crear las carpetas que no existen
	var contenidoDot string = ""
	//______________ LEYENDO ARCHIVO ____________________
	fp, err := os.OpenFile(direccion_disco, os.O_RDWR, 0777)
	defer fp.Close()
	if err != nil {
		//log.Fatal("ERROR AL ABRIR ARCHIVO", err)
		fmt.Println("ERROR AL ABRIR ARCHIVO:", direccion_disco)
		return Respuesta{
			Tipo:    -1,
			Mensaje: "",
			Data:    "",
		}
	}

	if err == nil {

		contenidoDot += "digraph G{ \n"
		contenidoDot += "\ntbl[shape=box,label=<\n"
		contenidoDot += "<table border='0' cellborder='1' cellspacing='0' width='300'  height='200' >\n"
		var tamanoMBR int32 = int32(unsafe.Sizeof(MBR{}))
		var tamanoEBR int32 = int32(unsafe.Sizeof(EBR{}))
		var masterBoot MBR = ObtenerInforMbr(fp, tamanoMBR)

		contenidoDot += "<tr >  <td  width='300' bgcolor=\"purple\" colspan=\"4\" >REPORTE MBR</td> </tr>"
		contenidoDot += "<tr>  <td  width='150'><b>mbr_tamaño</b></td><td  width='150'>" + strconv.Itoa(int(masterBoot.Mbr_tamano)) + "</td></tr>\n"
		contenidoDot += "<tr>  <td bgcolor=\"purple\"><b>mbr_fecha_creacion</b></td> <td>" + aString(masterBoot.Mbr_fecha_creacion[:]) + "</td> </tr>\n"
		contenidoDot += "<tr>  <td width='150'><b>mbr_disk_signature</b></td> <td width='150'>" + strconv.Itoa(int(masterBoot.Mbr_disk_signature)) + "</td>  </tr>\n"
		contenidoDot += "<tr>  <td bgcolor=\"purple\" width='150'><b>Disk_fit</b></td> <td width='150'>" + aString(masterBoot.Disk_fit[:]) + "</td>  </tr>\n"

		var index_Extendida int = -1
		for i := 0; i < 4; i++ {
			if masterBoot.Mbr_partition[i].Part_start != -1 && masterBoot.Mbr_partition[i].Part_status != 0 {
				if masterBoot.Mbr_partition[i].Part_type == 101 {
					index_Extendida = i
				}

				contenidoDot += "<tr><td colspan=\"4\" bgcolor=\"purple\">PARTICION</td> </tr>"
				contenidoDot += "<tr>  <td bgcolor=\"#FF00FF\"><b>part_status</b></td> <td>" + string(masterBoot.Mbr_partition[i].Part_status) + "</td></tr>\n"
				contenidoDot += "<tr>  <td><b>part_type</b></td> <td>" + string(masterBoot.Mbr_partition[i].Part_type) + "</td></tr>\n"
				contenidoDot += "<tr>  <td bgcolor=\"#FF00FF\"><b>part_fit</b></td> <td>" + aString(masterBoot.Mbr_partition[i].Part_fit[:]) + "</td></tr>\n"
				contenidoDot += "<tr>  <td><b>part_start</b></td> <td>" + strconv.Itoa(int(masterBoot.Mbr_partition[i].Part_start)) + "</td></tr>\n"
				contenidoDot += "<tr>  <td bgcolor=\"#FF00FF\"><b>part_size</b></td> <td>" + strconv.Itoa(int(masterBoot.Mbr_partition[i].Part_size)) + "</td></tr>\n"
				contenidoDot += "<tr>  <td><b>part_name</b></td> <td>" + aString(masterBoot.Mbr_partition[i].Part_name[:]) + "</td></tr>\n"
			}
		}

		contenidoDot += "</table>\n"
		contenidoDot += ">];"

		if index_Extendida != -1 {
			fp.Seek(int64(masterBoot.Mbr_partition[index_Extendida].Part_start), 0)
			var extendedBoot EBR = ObtenerInfoEbr(fp, tamanoEBR)
			contenidoDot += "\ntbl_1[shape=box, label=<\n "
			contenidoDot += "<table border='0' cellborder='1' cellspacing='0'  width='300' height='160' >\n "

			for {
				if extendedBoot.Part_status != 0 {

					contenidoDot += "<tr>  <td colspan=\"4\" bgcolor=\"purple\" width='300'>REPORTE EBR</td> </tr>"
					contenidoDot += "<tr>  <td width='150'><b>part_status</b></td> <td width='150'>" + string(extendedBoot.Part_status) + "</td>  </tr>\n"
					contenidoDot += "<tr>  <td><b>part_fit</b></td> <td>" + aString(extendedBoot.Part_fit[:]) + "</td>  </tr>\n"
					contenidoDot += "<tr>  <td><b>part_size</b></td> <td>" + strconv.Itoa(int(extendedBoot.Part_size)) + "</td>  </tr>\n"
					contenidoDot += "<tr>  <td><b>part_start</b></td> <td>" + strconv.Itoa(int(extendedBoot.Part_start)) + "</td>  </tr>\n"
					contenidoDot += "<tr>  <td><b>part_next</b></td> <td>" + strconv.Itoa(int(extendedBoot.Part_next)) + "</td>  </tr>\n"
					contenidoDot += "<tr>  <td><b>part_name</b></td> <td>" + aString(extendedBoot.Part_name[:]) + "</td>  </tr>\n"

				}

				if extendedBoot.Part_next == -1 {
					break
				}
				fp.Seek(int64(extendedBoot.Part_next), 0)
				extendedBoot = ObtenerInfoEbr(fp, tamanoEBR) //ruta quemada por el momento xD
			}
			contenidoDot += "</table>\n"
			contenidoDot += ">];\n"
		}
		contenidoDot += "}\n"
		crear_archivo_reporte(destino_reporte_dot, contenidoDot)
		contenidoDot = ""
		//var comando string = "dot -T" + extension + " " + destino_reporte_dot + " -o " + destino_reporte_grafico
		fmt.Println("REPORTASO")
		fmt.Println(destino_reporte_dot)
		fmt.Println(destino_reporte_grafico)
		fmt.Println(extension)

		res := systema_comando(destino_reporte_dot, destino_reporte_grafico, extension)
		if res == 0 {
			data := getReporteBase64(destino_reporte_grafico)
			return Respuesta{
				Tipo:      0,
				Mensaje:   "",
				Data:      data,
				Extension: extension,
			}
		}

		return Respuesta{
			Tipo:    -1,
			Mensaje: "",
			Data:    "",
		}

	} else {
		return Respuesta{
			Tipo:    -1,
			Mensaje: "",
			Data:    "",
		}
	}
}

/*____________________________________________ FUNCIONES IMPORTANTES PARA REPORTE __________________________________*/
func crear_archivo_reporte(path string, contenido_archivo string) {
	crear_archivo, err := os.Create(path)

	if err != nil {
		log.Fatal(err)
	}

	defer crear_archivo.Close()

	_, err2 := crear_archivo.WriteString(contenido_archivo)

	if err2 != nil {
		//log.Fatal(err2)
		fmt.Println("ERROR EN FUNCION CREAR ARCHIVO REPORTE")
	}

}

func aString(c []byte) string {
	n := -1
	for i, b := range c {
		if b == '\x00' {
			break
		}
		n = i
	}
	return string(c[:n+1])
}

func systema_comando(dir_dot string, dir_rep string, extension_rep string) int {
	fmt.Println("extension_rep", extension_rep)
	var arg1 string = "dot"
	var arg2 string = "-T" + extension_rep
	var arg3 string = dir_dot
	var arg4 string = "-o"
	var arg5 string = dir_rep
	//"dot -Tjpg /home/dark/Documentos/PruebaArchivo/parte1/particiones/d1.dot -o /home/dark/Documentos/PruebaArchivo/parte1/particiones/d1.jpg"
	cmd := exec.Command(arg1, arg2, arg3, arg4, arg5)

	err := cmd.Run()

	if err != nil {
		//log.Fatal(err)
		fmt.Println("ERROR AL CREAR REPORTE EN FUNCION SISTEMA COMANDO......")
		return -1
	}
	return 0
}

func calcular_porcentaje_usado(tamano_particion int, tamano_disco int) int {

	return int(math.Round(((float64(tamano_particion) * 100) / float64(tamano_disco))))
}

func obtener_espacio_libre(inicio_part int, fin_part int, tamano_part int, tipo string) int {
	if tipo == "normal" {
		return fin_part - (inicio_part + tamano_part)
	} else if tipo == "extendida" {
		return fin_part - (inicio_part + tamano_part + int(unsafe.Sizeof(EBR{})))
	}
	return 0
}

/*____________________________________________ FIN FUNCIONES IMPORTANTES PARA REPORTE __________________________________*/
func crear_rep_disk(direccion_disco string, destino_reporte_dot string, destino_reporte_grafico string, extension string) Respuesta {
	Crear_carpetas(destino_reporte_dot)

	var contenidoDot string = ""

	archivo_rep, err := os.OpenFile(direccion_disco, os.O_RDWR, 0777)
	defer archivo_rep.Close()
	if err != nil {
		//log.Fatal("ERROR AL ABRIR ARCHIVO", err)
		fmt.Println("ERROR AL ABRIR ARCHIVO:", direccion_disco)
		return Respuesta{
			Tipo:    -1,
			Mensaje: "",
			Data:    "",
		}
	}

	if err == nil {
		if true {
			contenidoDot += "digraph{\n\tnode [shape=plaintext]\n\tReporte_Disco[label=<\n\t<TABLE cellpadding=\"0\">\n\t<TR>\n\t\t<TD BGCOLOR=\"#27AE60\">MBR<br/></TD>\n"
			var tamanoMBR int32 = int32(unsafe.Sizeof(MBR{}))
			var tamanoEBR int32 = int32(unsafe.Sizeof(EBR{}))

			var mbr_temp MBR = ObtenerInforMbr(archivo_rep, tamanoMBR)

			var i int = 0
			var indice_usado_final int = 0

			for i = 0; i < 4; i++ {
				if mbr_temp.Mbr_partition[i].Part_status != '\x00' {
					indice_usado_final = i
					if mbr_temp.Mbr_partition[i].Part_type == 112 {

						if mbr_temp.Mbr_partition[i].Part_status == '\x00' {
							contenidoDot = contenidoDot + "<TD BGCOLOR=\"#00FF00\">" + "LIBRE<br/>" + strconv.Itoa(calcular_porcentaje_usado(int(mbr_temp.Mbr_partition[i].Part_size), int(mbr_temp.Mbr_tamano))) + "%</TD>"
						} else {
							contenidoDot = contenidoDot + "<TD BGCOLOR=\"#00FF00\">" + "PRIMARIA<br/>" + strconv.Itoa(calcular_porcentaje_usado(int(mbr_temp.Mbr_partition[i].Part_size), int(mbr_temp.Mbr_tamano))) + "%</TD>"
						}

						if i < 3 {
							var espacio_no_usado int = obtener_espacio_libre(int(mbr_temp.Mbr_partition[i].Part_start), int(mbr_temp.Mbr_partition[i+1].Part_start), int(mbr_temp.Mbr_partition[i].Part_size), "normal")
							if espacio_no_usado > 0 {
								contenidoDot = contenidoDot + "<TD BGCOLOR=\"#00FF00\">" + "LIBRE<br/>" + strconv.Itoa(calcular_porcentaje_usado(espacio_no_usado, int(mbr_temp.Mbr_tamano))) + "%</TD>"
							}
						}

					}

					if mbr_temp.Mbr_partition[i].Part_type == 101 {
						//crear metodo para obtener el string porque si esta antes o en medio F para reporte
						archivo_rep.Seek(int64(mbr_temp.Mbr_partition[i].Part_start), 0)
						var extendedBoot EBR = ObtenerInfoEbr(archivo_rep, tamanoEBR)
						var extend_temp EBR = EBR{}
						var espacio_libre int

						contenidoDot = contenidoDot + "<TD BORDER=\"0\">\n\t<TABLE BORDER=\"0\" CELLBORDER=\"1\" CELLSPACING=\"0\" CELLPADDING=\"4\">\n\t\t<TR>\n\t\t\t<TD BGCOLOR=\"#33FF9F\" COLSPAN=\"23\">EXTENDIDA</TD>\n\t\t</TR>\n\t\t<TR>"

						for {

							if extendedBoot.Part_next == -1 {
								contenidoDot = contenidoDot + "\n\t\t\t\t<TD bgcolor=\"#339CFF\">EBR</TD>\n\t\t\t\t<TD bgcolor=\"#D1FF33\">LOGICA<br/>" + strconv.Itoa(calcular_porcentaje_usado(int(extendedBoot.Part_size), int(mbr_temp.Mbr_tamano))) + "%</TD>"
								var termina_ultima_extendida int32 = extendedBoot.Part_start + extendedBoot.Part_size
								var total_no_usado = mbr_temp.Mbr_partition[i].Part_size - termina_ultima_extendida
								if total_no_usado > 0 {
									//contenidoDot = contenidoDot + "\n\t\t\t\t<TD bgcolor=\"#339CFF\">EBR</TD>\n\t\t\t\t<TD bgcolor=\"#D1FF33\">LIBRE<br/>" + strconv.Itoa(calcular_porcentaje_usado(int(total_no_usado), int(mbr_temp.Mbr_tamano))) + "%</TD>"
									contenidoDot = contenidoDot + "<TD BGCOLOR=\"#00FF00\">" + "LIBRE<br/>" + strconv.Itoa(calcular_porcentaje_usado(int(total_no_usado), int(mbr_temp.Mbr_tamano))) + "%</TD>"
								}

								break
							} else {
								extend_temp = extendedBoot
								archivo_rep.Seek(int64(extendedBoot.Part_next), 0)
								extendedBoot = ObtenerInfoEbr(archivo_rep, tamanoEBR)
								espacio_libre = obtener_espacio_libre(int(extend_temp.Part_start), int(extendedBoot.Part_start), int(extend_temp.Part_size), "extendida")
								if espacio_libre > 0 {

									contenidoDot = contenidoDot + "\n\t\t\t\t<TD bgcolor=\"#339CFF\">EBR</TD>\n\t\t\t\t<TD bgcolor=\"#D1FF33\">LOGICA<br/>" + strconv.Itoa(calcular_porcentaje_usado(int(extend_temp.Part_size), int(mbr_temp.Mbr_tamano))) + "%</TD>"
									contenidoDot = contenidoDot + "\n\t\t\t\t<TD bgcolor=\"#339CFF\">EBR</TD>\n\t\t\t\t<TD bgcolor=\"#D1FF33\">LIBRE<br/>" + strconv.Itoa(calcular_porcentaje_usado(espacio_libre, int(mbr_temp.Mbr_tamano))) + "%</TD>"

								} else {

									contenidoDot = contenidoDot + "\n\t\t\t\t<TD bgcolor=\"#339CFF\">EBR</TD>\n\t\t\t\t<TD bgcolor=\"#D1FF33\">LOGICA<br/>" + strconv.Itoa(calcular_porcentaje_usado(int(extend_temp.Part_size), int(mbr_temp.Mbr_tamano))) + "%</TD>"
								}
							}
						} //fin de while
						contenidoDot = contenidoDot + "\n\t\t</TR>\n\t</TABLE>\n</TD>"
					}
				}

			} //fin de for
			//____________________________________________________________PARA ESPACCIO FINAL EN  DE DISCO _______
			var termina_ultima_particion int32 = mbr_temp.Mbr_partition[indice_usado_final].Part_start + mbr_temp.Mbr_partition[indice_usado_final].Part_size
			var espacio_no_usado = mbr_temp.Mbr_tamano - termina_ultima_particion
			if espacio_no_usado > 0 {
				contenidoDot = contenidoDot + "<TD BGCOLOR=\"#00FF00\">" + "LIBRE<br/>" + strconv.Itoa(calcular_porcentaje_usado(int(espacio_no_usado), int(mbr_temp.Mbr_tamano))) + "%</TD>"
			}

			//____________________________________________________________ FIN PARA ESPACION FINAL EN DISCO

			contenidoDot += "\n\t</TR>\n\t</TABLE>\n\t>];}"
			crear_archivo_reporte(destino_reporte_dot, contenidoDot)
			contenidoDot = ""

			res := systema_comando(destino_reporte_dot, destino_reporte_grafico, extension)

			if res == 0 {
				data := getReporteBase64(destino_reporte_grafico)
				return Respuesta{
					Tipo:      0,
					Mensaje:   "",
					Data:      data,
					Extension: extension,
				}
			}

			return Respuesta{
				Tipo:    -1,
				Mensaje: "",
				Data:    "",
			}

		} else {

			fmt.Println("No se puede crear archivo dot en reporte disk")
		}

	} else {
		fmt.Println("Error al abrir archivo para crear reporte disk")
	}
	return Respuesta{
		Tipo:    -1,
		Mensaje: "",
		Data:    "",
	}
} //fin de metodo crear_rep_disk

/*_____________________________REPORTE TREE ___________________*/
func crear_rep_tree(direccion_disco string, destino_reporte_dot string, destino_reporte_grafico string, extension string) Respuesta {
	Crear_carpetas(destino_reporte_dot)

	var contenidoDot string = ""

	fp, err := os.OpenFile(direccion_disco, os.O_RDWR, 0777)
	defer fp.Close()
	if err != nil {
		//log.Fatal("ERROR AL ABRIR ARCHIVO", err)
		fmt.Println("ERROR AL ABRIR ARCHIVO:", direccion_disco)
		return Respuesta{
			Tipo:    -1,
			Mensaje: "",
			Data:    "",
		}
	}

	var super SUPER_BLOQUE
	var inodo TABLA_INODOS
	var carpeta BLOQUE_CARPETA
	var archivo BLOQUE_ARCHIVO

	fp.Seek(int64(actualSesion.InicioSuper), 0)
	super = Leer_SuperBloque(fp, tamanoSuper)

	var aux int32 = super.S_bm_inode_start //super.s_bm_inode_start
	var i int32 = 0

	var buffer byte

	contenidoDot = "digraph G{\n\n"
	contenidoDot = contenidoDot + "    rankdir=\"LR\" \n"
	//Creamos los inodos
	for aux < super.S_bm_block_start {
		//fmt.Println("INicio de super bloque:",supe)

		fp.Seek(int64(super.S_bm_inode_start+i), 0)
		buffer = Leer_Byte(fp, 1)
		aux++
		var port int = 0
		if buffer == '1' {
			fp.Seek(int64(super.S_inode_start+tamanoInodo*i), 0)
			inodo = Leer_TablaInodo(fp, tamanoInodo)
			contenidoDot = contenidoDot + "    inodo_" + strconv.Itoa(int(i)) + "[ shape=plaintext  label=<\n"
			contenidoDot = contenidoDot + "   <table bgcolor=\"royalblue\" border='0' >"
			contenidoDot = contenidoDot + "    <tr> <td colspan='2'><b>Inode " + strconv.Itoa(int(i)) + "</b></td></tr>\n"
			contenidoDot = contenidoDot + "    <tr> <td bgcolor=\"lightsteelblue\"> i_uid </td> <td bgcolor=\"white\">" + strconv.Itoa(int(inodo.I_uid)) + "</td>  </tr>\n"
			contenidoDot = contenidoDot + "    <tr> <td bgcolor=\"lightsteelblue\"> i_gid </td> <td bgcolor=\"white\">" + strconv.Itoa(int(inodo.I_gid)) + "</td>  </tr>\n"
			contenidoDot = contenidoDot + "    <tr> <td bgcolor=\"lightsteelblue\"> i_size </td><td bgcolor=\"white\">" + strconv.Itoa(int(inodo.I_size)) + "</td> </tr>\n"

			contenidoDot = contenidoDot + "    <tr> <td bgcolor=\"lightsteelblue\"> i_atime </td> <td bgcolor=\"white\">" + aString(inodo.I_atime[:]) + "</td> </tr>\n"
			contenidoDot = contenidoDot + "    <tr> <td bgcolor=\"lightsteelblue\"> i_ctime </td> <td bgcolor=\"white\">" + aString(inodo.I_ctime[:]) + "</td> </tr>\n"
			contenidoDot = contenidoDot + "    <tr> <td bgcolor=\"lightsteelblue\"> i_mtime </td> <td bgcolor=\"white\">" + aString(inodo.I_mtime[:]) + "</td> </tr>\n"

			for b := 0; b < 16; b++ {
				contenidoDot = contenidoDot + "    <tr> <td bgcolor=\"lightsteelblue\"> i_block_" + strconv.Itoa(port) + "</td> <td bgcolor=\"white\" port=\"f" + strconv.Itoa(b) + "\"> " + strconv.Itoa(int(inodo.I_block[b])) + "</td></tr>\n"
				port++
			}
			contenidoDot = contenidoDot + "    <tr> <td bgcolor=\"lightsteelblue\"> i_type </td> <td bgcolor=\"white\">" + string(inodo.I_type) + "</td>  </tr>\n"
			contenidoDot = contenidoDot + "    <tr> <td bgcolor=\"lightsteelblue\"> i_perm </td> <td bgcolor=\"white\">" + strconv.Itoa(int(inodo.I_perm)) + "</td>  </tr>\n"
			contenidoDot = contenidoDot + "   </table>>]\n\n"
			//Creamos los bloques relacionados al inodo
			for j := 0; j < 16; j++ {
				port = 0
				if inodo.I_block[j] != -1 {
					fp.Seek(int64(super.S_bm_block_start+inodo.I_block[j]), 0)
					buffer = Leer_Byte(fp, 1)
					if buffer == '1' { //Bloque carpeta
						fp.Seek(int64(super.S_block_start+tamanoCarpeta*inodo.I_block[j]), 0)
						carpeta = Leer_BloqueCarpeta(fp, tamanoCarpeta)
						contenidoDot = contenidoDot + "    bloque_" + strconv.Itoa(int(inodo.I_block[j])) + "[shape=plaintext  label=< \n"
						contenidoDot = contenidoDot + "   <table bgcolor=\"seagreen\" border='0'>\n"
						contenidoDot = contenidoDot + "    <tr> <td colspan='2'><b>Bloque Carpeta " + strconv.Itoa(int(inodo.I_block[j])) + "</b></td></tr>\n"
						contenidoDot = contenidoDot + "    <tr> <td bgcolor=\"mediumseagreen\"> b_name </td> <td bgcolor=\"mediumseagreen\"> b_inode </td></tr>\n"
						for c := 0; c < 4; c++ {
							contenidoDot = contenidoDot + "    <tr> <td bgcolor=\"white\" >" + aString(carpeta.B_content[c].B_name[:]) + "</td> <td bgcolor=\"white\"  port=\"f" + strconv.Itoa(port) + "\">" + strconv.Itoa(int(carpeta.B_content[c].B_inodo)) + "</td></tr>\n"
							port++
						}
						contenidoDot = contenidoDot + "   </table>>]\n\n"
						//Relacion de bloques a inodos
						for c := 0; c < 4; c++ {
							if carpeta.B_content[c].B_inodo != -1 {
								if aString(carpeta.B_content[c].B_name[:]) != "." && aString(carpeta.B_content[c].B_name[:]) != ".." {
									contenidoDot = contenidoDot + "    bloque_" + strconv.Itoa(int(inodo.I_block[j])) + ":f" + strconv.Itoa(c) + "-> inodo_" + strconv.Itoa(int(carpeta.B_content[c].B_inodo)) + ";\n"
								}

							}
						}
					} else if buffer == '2' { //Bloque archivo
						fp.Seek(int64(super.S_block_start+tamanoArchivo*inodo.I_block[j]), 0)
						archivo = Leer_BloqueArchivo(fp, tamanoArchivo)
						contenidoDot = contenidoDot + "    bloque_" + strconv.Itoa(int(inodo.I_block[j])) + "[shape=plaintext  label=< \n"
						contenidoDot = contenidoDot + "   <table border='0' bgcolor=\"sandybrown\">\n"
						contenidoDot = contenidoDot + "    <tr> <td> <b>Bloque Archivo " + strconv.Itoa(int(inodo.I_block[j])) + "</b></td></tr>\n"
						contenidoDot = contenidoDot + "    <tr> <td bgcolor=\"white\"> " + aString(archivo.B_content[:]) + "</td></tr>\n"
						contenidoDot = contenidoDot + "   </table>>]\n\n"
					}
					//Relacion de inodos a bloques
					contenidoDot = contenidoDot + "    inodo_" + strconv.Itoa(int(i)) + ":f" + strconv.Itoa(j) + " -> bloque_" + strconv.Itoa(int(inodo.I_block[j])) + "; \n"
				}
			}
		}
		i++
	}

	contenidoDot = contenidoDot + "\n\n}"
	crear_archivo_reporte(destino_reporte_dot, contenidoDot)
	contenidoDot = ""

	res := systema_comando(destino_reporte_dot, destino_reporte_grafico, extension)
	if res == 0 {
		data := getReporteBase64(destino_reporte_grafico)
		return Respuesta{
			Tipo:      0,
			Mensaje:   "",
			Data:      data,
			Extension: extension,
		}
	}

	return Respuesta{
		Tipo:    -1,
		Mensaje: "",
		Data:    "",
	}

	//Original----->string comando = "dot -T"+extension.toStdString()+" grafica.dot -o "+destino.toStdString();
	//Original---->system(comando.c_str());

	//string comando = "dot -T"+extension.toStdString()+" "+obtener_ruta_carpeta(destino_reporte.toStdString())+"/"+obtener_nombre_archivo(destino_reporte).toStdString()+".dot -o "+destino_reporte.toStdString();
	//system(comando.c_str());
	//cout << "\n\tReporte Tree generado con exito " << endl;
	//return 0
}

func crear_rep_super_bloque(direccion_disco string, destino_reporte_dot string, destino_reporte_grafico string, extension string) Respuesta {
	Crear_carpetas(destino_reporte_dot)

	var contenidoDot string = ""

	fp, err := os.OpenFile(direccion_disco, os.O_RDWR, 0777)
	defer fp.Close()
	if err != nil {
		//log.Fatal("ERROR AL ABRIR ARCHIVO", err)
		fmt.Println("ERROR AL ABRIR ARCHIVO:", direccion_disco)
		return Respuesta{
			Tipo:    -1,
			Mensaje: "",
			Data:    "",
		}
	}
	var super SUPER_BLOQUE

	fp.Seek(int64(actualSesion.InicioSuper), 0)
	super = Leer_SuperBloque(fp, tamanoSuper)

	contenidoDot = contenidoDot + "digraph G{\n"
	contenidoDot = contenidoDot + "    nodo [shape=none,  label=<"
	contenidoDot = contenidoDot + "   <table border='0' cellborder='1' cellspacing='0' bgcolor=\"cornflowerblue\">"
	contenidoDot = contenidoDot + "    <tr> <td COLSPAN='2'> <b>SUPER_BLOQUE</b> </td></tr>\n"
	contenidoDot = contenidoDot + "    <tr> <td bgcolor=\"lightsteelblue\"> s_inodes_count </td> <td bgcolor=\"white\"> " + strconv.Itoa(int(super.S_inodes_count)) + "</td> </tr>\n"
	contenidoDot = contenidoDot + "    <tr> <td bgcolor=\"lightsteelblue\"> s_blocks_count </td> <td bgcolor=\"white\"> " + strconv.Itoa(int(super.S_blocks_count)) + "</td> </tr>\n"
	contenidoDot = contenidoDot + "    <tr> <td bgcolor=\"lightsteelblue\"> s_free_block_count </td> <td bgcolor=\"white\"> " + strconv.Itoa(int(super.S_free_blocks_count)) + "</td> </tr>\n"
	contenidoDot = contenidoDot + "    <tr> <td bgcolor=\"lightsteelblue\"> s_free_inodes_count </td> <td bgcolor=\"white\"> " + strconv.Itoa(int(super.S_free_inodes_count)) + " </td> </tr>\n"

	contenidoDot = contenidoDot + "    <tr> <td bgcolor=\"lightsteelblue\"> s_mtime </td> <td bgcolor=\"white\">" + aString(super.S_mtime[:]) + " </td></tr>\n"
	//contenidoDot = contenidoDot + "    <tr> <td bgcolor=\"lightsteelblue\"> s_umtime </td> <td bgcolor=\"white\"> " + aString(super.S_umtime[:]) + "</td> </tr>\n"
	contenidoDot = contenidoDot + "    <tr> <td bgcolor=\"lightsteelblue\"> s_mnt_count </td> <td bgcolor=\"white\"> " + strconv.Itoa(int(super.S_mnt_count)) + "</td> </tr>\n"
	contenidoDot = contenidoDot + "    <tr> <td bgcolor=\"lightsteelblue\"> s_magic </td> <td bgcolor=\"white\"> " + strconv.Itoa(int(super.S_magic)) + " </td> </tr>\n"
	contenidoDot = contenidoDot + "    <tr> <td bgcolor=\"lightsteelblue\"> s_inode_size </td> <td bgcolor=\"white\"> " + strconv.Itoa(int(super.S_inode_size)) + " </td> </tr>\n"
	contenidoDot = contenidoDot + "    <tr> <td bgcolor=\"lightsteelblue\"> s_block_size </td> <td bgcolor=\"white\"> " + strconv.Itoa(int(super.S_block_size)) + " </td> </tr>\n"
	contenidoDot = contenidoDot + "    <tr> <td bgcolor=\"lightsteelblue\"> s_first_ino </td> <td bgcolor=\"white\"> " + strconv.Itoa(int(super.S_first_ino)) + " </td> </tr>\n"
	contenidoDot = contenidoDot + "    <tr> <td bgcolor=\"lightsteelblue\"> s_first_blo </td> <td bgcolor=\"white\"> " + strconv.Itoa(int(super.S_first_blo)) + " </td> </tr>\n"
	contenidoDot = contenidoDot + "    <tr> <td bgcolor=\"lightsteelblue\"> s_bm_inode_start </td> <td bgcolor=\"white\"> " + strconv.Itoa(int(super.S_bm_inode_start)) + " </td></tr>\n"
	contenidoDot = contenidoDot + "    <tr> <td bgcolor=\"lightsteelblue\"> s_bm_block_start </td> <td bgcolor=\"white\"> " + strconv.Itoa(int(super.S_bm_block_start)) + " </td> </tr>\n"
	contenidoDot = contenidoDot + "    <tr> <td bgcolor=\"lightsteelblue\"> s_inode_start </td> <td bgcolor=\"white\"> " + strconv.Itoa(int(super.S_inode_start)) + " </td> </tr>\n"
	contenidoDot = contenidoDot + "    <tr> <td bgcolor=\"lightsteelblue\"> s_block_start </td> <td bgcolor=\"white\"> " + strconv.Itoa(int(super.S_block_start)) + " </td> </tr>\n"
	contenidoDot = contenidoDot + "   </table>>]\n"
	contenidoDot = contenidoDot + "\n}"

	crear_archivo_reporte(destino_reporte_dot, contenidoDot)
	contenidoDot = ""
	res := systema_comando(destino_reporte_dot, destino_reporte_grafico, extension)
	if res == 0 {
		data := getReporteBase64(destino_reporte_grafico)
		return Respuesta{
			Tipo:      0,
			Mensaje:   "",
			Data:      data,
			Extension: extension,
		}
	}

	return Respuesta{
		Tipo:    -1,
		Mensaje: "",
		Data:    "",
	}

}

func crear_rep_file(direccion_disco string, destino_reporte_dot string, destino_reporte_grafico string, extension string, ruta_archivo string) Respuesta {
	Crear_carpetas(destino_reporte_dot)

	var contenidoDot string = ""

	fp, err := os.OpenFile(direccion_disco, os.O_RDWR, 0777)
	defer fp.Close()
	if err != nil {
		//log.Fatal("ERROR AL ABRIR ARCHIVO", err)
		fmt.Println("ERROR AL ABRIR ARCHIVO:", direccion_disco)
		return Respuesta{
			Tipo:    -1,
			Mensaje: "",
			Data:    "",
		}
	}
	var super SUPER_BLOQUE
	var inodo TABLA_INODOS
	var archivo BLOQUE_ARCHIVO

	fp.Seek(int64(actualSesion.InicioSuper), 0)
	super = Leer_SuperBloque(fp, tamanoSuper)
	var numero_inodo int = buscarCarpetaArchivo(fp, ruta_archivo)
	if numero_inodo == -1 {
		fmt.Println("ERROR EL ARCHIVO:", ruta_archivo, "NO EXISTE PARA REPORTE")
		return Respuesta{
			Tipo:    -1,
			Mensaje: "",
			Data:    "",
		}
	}
	var name string = filepath.Base(ruta_archivo)
	fp.Seek(int64(super.S_inode_start+tamanoInodo*int32(numero_inodo)), 0)
	inodo = Leer_TablaInodo(fp, tamanoInodo)

	//FILE *graph = fopen("grafica.dot","w");
	contenidoDot = contenidoDot + "digraph G{\n"
	contenidoDot = contenidoDot + "    nodo [shape=none,  label=<"
	contenidoDot = contenidoDot + "   <table border='0' cellborder='1' cellspacing='0' bgcolor=\"lightsteelblue\">"
	contenidoDot = contenidoDot + "    <tr><td align=\"left\"> <b>" + name + "</b> </td></tr>\n"
	contenidoDot = contenidoDot + "    <tr><td bgcolor=\"white\">"
	for i := 0; i < 15; i++ {
		if inodo.I_block[i] != -1 {
			//Apuntadores directos
			fp.Seek(int64(super.S_block_start+tamanoArchivo*inodo.I_block[i]), 0)
			archivo = Leer_BloqueArchivo(fp, tamanoArchivo)
			contenidoDot = contenidoDot + aString(archivo.B_content[:]) + "<br/>"

		}
	}
	contenidoDot = contenidoDot + "    </td></tr>\n"
	contenidoDot = contenidoDot + "   </table>>]\n"
	contenidoDot = contenidoDot + "\n}"
	//cout << "Reporte file generado con exito " << endl;
	crear_archivo_reporte(destino_reporte_dot, contenidoDot)
	contenidoDot = ""
	res := systema_comando(destino_reporte_dot, destino_reporte_grafico, extension)
	if res == 0 {
		data := getReporteBase64(destino_reporte_grafico)
		return Respuesta{
			Tipo:      0,
			Mensaje:   "",
			Data:      data,
			Extension: extension,
		}
	}

	return Respuesta{
		Tipo:    -1,
		Mensaje: "",
		Data:    "",
	}

}

func getReporteBase64(destino_reporte_grafico string) string {
	fmt.Println("Archivo lectura", destino_reporte_grafico)
	file, err := ioutil.ReadFile(destino_reporte_grafico)
	if err != nil {
		fmt.Println("Error al leer el archivo:", err)
		return ""
	}
	encoded := base64.StdEncoding.EncodeToString(file)
	return encoded
}
