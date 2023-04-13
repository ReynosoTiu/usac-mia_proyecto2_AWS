package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"math"
	"os"
	"unsafe"
)

func EjecutarMkfs(id string, type_ string) string {
	if type_ == "full" || type_ == "" {
	} else {
		fmt.Println("ERROR  VALOR DESCONOCIDO DE TYPE: ", type_)
		return "Valor del parametro TYPE incorrecto"
	}
	var aux *NODO = listaS.obtenerNodo(id)
	if aux != nil {
		masterBoot := leer_archivo(aux.Path)
		for _, a := range masterBoot.Mbr_partition {
			if string(bytes.Trim(a.Part_name[:], "\x00")) == aux.Name {
				if a.Part_type == 101 {
					return "No se permite formatear una particion extendida"
				}
			}
		}
		index, inicio, tamano, es_logica, _ := Buscar_Indice_P_E_L(aux.Path, aux.Name)
		if index != -1 {
			if es_logica {
				inicio += int32(unsafe.Sizeof(EBR{}))
			}
			formatearEXT2(inicio, tamano, aux.Path)
			return "PARTICION FORMATEADA EXITOSAMENTE CON EXT2"
		} else {
			fmt.Println("ERROR  NO COINCIDE EN MKFS")
			return "No se encuentra una particion montada con el ID ingresado"
		}
	} else {
		fmt.Println("ERROR ID NO MONTADA")
		return "No se encuentra una particion montada con el ID ingresado"
	}
}

func Buscar_Indice_P_E_L(path string, name string) (int32, int32, int32, bool, []byte) {

	//os.Open sirve para leer el archivo nada mas y no para modificar
	file, err := os.OpenFile(path, os.O_RDWR, 0777)
	defer file.Close()

	if err != nil {
		log.Fatal("ERROR AL ABRIR ARCHIVO,ARCHIVO NO EXISTE", err)
		os.Exit(4)
	}
	masterBoot := MBR{}

	var tamano_masterBoot int32 = int32(unsafe.Sizeof(masterBoot))
	masterBoot = ObtenerInforMbr(file, tamano_masterBoot)

	var indice_extendia int32 = -1
	for i := 0; i < 4; i++ {
		if string(bytes.Trim(masterBoot.Mbr_partition[i].Part_name[:], "\x00")) == name {
			file.Close()

			return int32(i), masterBoot.Mbr_partition[i].Part_start, masterBoot.Mbr_partition[i].Part_size, false, masterBoot.Mbr_partition[i].Part_fit[:]
		}
		if string(masterBoot.Mbr_partition[i].Part_type) == "e" {
			indice_extendia = int32(i)
		}
	}

	if indice_extendia != -1 {
		//aqui va lo de la logica
		extendBoot := EBR{}
		file.Seek(int64(masterBoot.Mbr_partition[indice_extendia].Part_start), 0) //(cantidad_bytes_desplazar,origen_de_archivo)
		extendBoot = ObtenerInfoEbr(file, int32(unsafe.Sizeof(EBR{})))
		for {
			//fmt.Println("Nombre Particion: ",string(extendBoot.Part_name[:])," NOMBRE COMPARADO: ",name)
			//fmt.Println("Tamano Partcicion: ",len(string(bytes.Trim(extendBoot.Part_name[:],"\x00")))," Tamano Comparado: ",len(name))
			//se utiliza bytes.Trimp para eliminar los caracteres nulos y comparar el nombre xD
			if string(bytes.Trim(extendBoot.Part_name[:], "\x00")) == name {
				return indice_extendia, extendBoot.Part_start, extendBoot.Part_size, true, extendBoot.Part_fit[:]
			}

			if extendBoot.Part_next == -1 {
				break
			}
			file.Seek(int64(extendBoot.Part_next), 0)
			extendBoot = ObtenerInfoEbr(file, int32(unsafe.Sizeof(EBR{}))) //ruta quemada por el momento xD
		}
	}
	file.Close()
	var temp []byte
	return -1, 0, 0, false, temp
}

/*___________________________________ 	INICIO DE EXT2 _______________________________________*/
func formatearEXT2(inicio int32, tamano int32, direccion string) {
	var n float64 = (float64(tamano) - float64(unsafe.Sizeof(SUPER_BLOQUE{}))) / (4 + float64(unsafe.Sizeof(TABLA_INODOS{})) + 3*float64(unsafe.Sizeof(BLOQUE_ARCHIVO{})))
	var num_estructuras int32 = int32(math.Floor(n)) //Numero de inodos
	var num_bloques int32 = 3 * num_estructuras      //Numero de bloques

	//fmt.Println("NUMERO DE BITMAP INODO", num_estructuras)
	//fmt.Println("NUMERO DE BITMAP BLOQUE", num_bloques)

	var super SUPER_BLOQUE
	super.S_filesystem_type = 2
	super.S_inodes_count = num_estructuras
	super.S_blocks_count = num_bloques
	super.S_free_blocks_count = num_bloques - 2
	super.S_free_inodes_count = num_estructuras - 2
	copy(super.S_mtime[:], Hora_fecha())
	//copy(super.S_umtime[:], "")
	super.S_mnt_count = 0
	super.S_magic = 0xEF53
	super.S_inode_size = int32(unsafe.Sizeof(TABLA_INODOS{}))
	super.S_block_size = int32(unsafe.Sizeof(BLOQUE_ARCHIVO{}))
	super.S_first_ino = 2
	super.S_first_blo = 2
	super.S_bm_inode_start = inicio + int32(unsafe.Sizeof(SUPER_BLOQUE{}))
	super.S_bm_block_start = inicio + int32(unsafe.Sizeof(SUPER_BLOQUE{})) + num_estructuras
	super.S_inode_start = inicio + int32(unsafe.Sizeof(SUPER_BLOQUE{})) + num_estructuras + num_bloques
	super.S_block_start = inicio + int32(unsafe.Sizeof(SUPER_BLOQUE{})) + num_estructuras + num_bloques + (int32(unsafe.Sizeof(TABLA_INODOS{})) * num_estructuras)

	var buffer_Cero byte = 48              //Cero para bitmap no usado
	var buffer_Carpeta_Inodo_uno byte = 49 //Uno para representar Bloque Carpeta
	var buffer_Archivo byte = 50           //Dos para representar Bloque Archivo

	file, err := os.OpenFile(direccion, os.O_RDWR, 0777)
	defer file.Close()

	if err != nil {
		log.Fatal("ERROR AL ABRIR ARCHIVO,ARCHIVO NO EXISTE", err)
		os.Exit(4)
	}

	/*___________________ SUPERBLOQUE ____________________*/
	file.Seek(int64(inicio), 0)
	Escribir_SuperBloque(file, &super)
	/*___________________ BITMAP DE INODOS _______________*/
	for i := 0; i < int(num_estructuras); i++ {
		file.Seek(int64(super.S_bm_inode_start)+int64(i), 0)
		Escribir_Byte(file, &buffer_Cero)
	}
	/*___________________ bit para / y users.txt en BM ____________*/
	file.Seek(int64(super.S_bm_inode_start), 0)
	Escribir_Byte(file, &buffer_Carpeta_Inodo_uno) //Para Root
	Escribir_Byte(file, &buffer_Carpeta_Inodo_uno) //Para archivo usuario

	/*_____________________ BITMAP DE BLOQUES ____________________*/
	for i := 0; i < int(num_bloques); i++ {
		file.Seek(int64(super.S_bm_block_start)+int64(i), 0)
		Escribir_Byte(file, &buffer_Cero)
	}
	/*____________________ bit para / y users.txt en BM ___________*/
	file.Seek(int64(super.S_bm_block_start), 0)
	Escribir_Byte(file, &buffer_Carpeta_Inodo_uno) //
	Escribir_Byte(file, &buffer_Archivo)

	/*_____________________ inodo para carpeta root ________________*/
	var inodoTabla TABLA_INODOS

	inodoTabla.I_uid = 1
	inodoTabla.I_gid = 1
	inodoTabla.I_size = 0
	copy(inodoTabla.I_atime[:], Hora_fecha())
	copy(inodoTabla.I_ctime[:], Hora_fecha())
	copy(inodoTabla.I_mtime[:], Hora_fecha())
	inodoTabla.I_block[0] = 0
	for i := 1; i < 16; i++ {
		inodoTabla.I_block[i] = -1
	}
	inodoTabla.I_type = 48
	inodoTabla.I_perm = 664
	file.Seek(int64(super.S_inode_start), 0)
	Escribir_TablaInodo(file, &inodoTabla)
	/*_____________________Bloque para carpeta root___________________*/
	var bloqueCarpeta BLOQUE_CARPETA

	copy(bloqueCarpeta.B_content[0].B_name[:], ".") //Actual (el mismo) poisble error con copy ya qye sino F haber que pex
	bloqueCarpeta.B_content[0].B_inodo = 0

	copy(bloqueCarpeta.B_content[1].B_name[:], "..") //Padre //verificar que copy no copie todo el contenido en el arreglo de bytes sino F
	bloqueCarpeta.B_content[1].B_inodo = 0

	copy(bloqueCarpeta.B_content[2].B_name[:], "users.txt")
	bloqueCarpeta.B_content[2].B_inodo = 1

	copy(bloqueCarpeta.B_content[3].B_name[:], ".")
	bloqueCarpeta.B_content[3].B_inodo = -1

	file.Seek(int64(super.S_block_start), 0)
	Escribir_BloqueCarpeta(file, &bloqueCarpeta)

	/*____________________ inodo para users.txt ____________________*/
	inodoTabla.I_uid = 1
	inodoTabla.I_gid = 1
	inodoTabla.I_size = 27
	copy(inodoTabla.I_atime[:], Hora_fecha())
	copy(inodoTabla.I_ctime[:], Hora_fecha())
	copy(inodoTabla.I_mtime[:], Hora_fecha())
	inodoTabla.I_block[0] = 1
	for i := 1; i < 16; i++ {
		inodoTabla.I_block[i] = -1
	}
	inodoTabla.I_type = 49
	inodoTabla.I_perm = 664
	file.Seek(int64(super.S_inode_start)+int64(unsafe.Sizeof(TABLA_INODOS{})), 0)
	Escribir_TablaInodo(file, &inodoTabla)
	/*___________________ Bloque para users.txt _____________________*/
	var archivo BLOQUE_ARCHIVO
	copy(archivo.B_content[:], "\x00")
	//memset(archivo.B_content, 0, sizeof(archivo.B_content))
	copy(archivo.B_content[:], "1,G,root\n1,U,root,root,123\n")
	file.Seek(int64(super.S_block_start)+int64(unsafe.Sizeof(BLOQUE_CARPETA{})), 0)
	Escribir_BloqueArchivo(file, &archivo)

	//fmt.Println("EXT2...")
	//fmt.Println("DISCO FORMATEADO CON EXITO")
	file.Close()
}

/*___________________________________ FIN DE EXT2 ____________________________________________*/

/*___________________________________ INICIO DE MODIFICAR ____________________________________*/

func Escribir_SuperBloque(file *os.File, superBloque *SUPER_BLOQUE) {
	var master_buffer bytes.Buffer
	binary.Write(&master_buffer, binary.BigEndian, superBloque)
	EscribirBytes(file, master_buffer.Bytes())
}

/*___________________________________  FIN DE MODIFICAR SUPER BLOQUE _________________________*/

/*___________________________________ INICIO DE MODIFICAR BYE ____________________________________*/

func Escribir_Byte(file *os.File, mibyte *byte) {
	var master_buffer bytes.Buffer
	binary.Write(&master_buffer, binary.BigEndian, mibyte)
	EscribirBytes(file, master_buffer.Bytes())
}

/*___________________________________  FIN DE MODIFICAR BYTE _________________________*/

/*___________________________________ INICIO DE MODIFICAR TABLA INODOS____________________________________*/

func Escribir_TablaInodo(file *os.File, tablaInodo *TABLA_INODOS) {
	var master_buffer bytes.Buffer
	binary.Write(&master_buffer, binary.BigEndian, tablaInodo)
	EscribirBytes(file, master_buffer.Bytes())
}

/*___________________________________  FIN DE MODIFICAR TABLA INODOS _________________________*/

/*___________________________________ INICIO DE MODIFICAR TABLA INODOS____________________________________*/

func Escribir_BloqueCarpeta(file *os.File, bloqueCarpeta *BLOQUE_CARPETA) {
	var master_buffer bytes.Buffer
	binary.Write(&master_buffer, binary.BigEndian, bloqueCarpeta)
	EscribirBytes(file, master_buffer.Bytes())
}

/*___________________________________  FIN DE MODIFICAR TABLA INODOS _________________________*/

/*___________________________________ INICIO DE MODIFICAR TABLA INODOS____________________________________*/

func Escribir_BloqueArchivo(file *os.File, bloqueArchiv *BLOQUE_ARCHIVO) {
	var master_buffer bytes.Buffer
	binary.Write(&master_buffer, binary.BigEndian, bloqueArchiv)
	EscribirBytes(file, master_buffer.Bytes())
}

/*___________________________________  FIN DE MODIFICAR TABLA INODOS _________________________*/

/*##############################################################-############################*/
/*___________________________________ INICIO DE LEER SUPERBLOQUE ____________________________________*/
func Leer_SuperBloque(file *os.File, tamano_leer int32) SUPER_BLOQUE {

	superBloque := SUPER_BLOQUE{}
	data := leerBytes(file, tamano_leer)
	buffer := bytes.NewBuffer(data)
	err := binary.Read(buffer, binary.BigEndian, &superBloque)
	if err != nil {
		log.Fatal("ERROR AL LEER ARCHIVO BINARIO", err)
		return SUPER_BLOQUE{}
	}
	return superBloque
}

/*___________________________________ FIN DE LEER SUPERBLOQUE ____________________________________*/

/*___________________________________ INICIO DE LEER BYE ____________________________________*/

func Leer_Byte(file *os.File, tamano_leer int32) byte {
	var leerByt byte
	data := leerBytes(file, tamano_leer)
	buffer := bytes.NewBuffer(data)
	err := binary.Read(buffer, binary.BigEndian, &leerByt)
	if err != nil {
		log.Fatal("ERROR AL LEER ARCHIVO BINARIO", err)
		return 0
	}
	return leerByt
}

/*___________________________________  FIN DE LEER BYTE _________________________*/

/*___________________________________ INICIO DE LEER TABLA INODOS____________________________________*/

func Leer_TablaInodo(file *os.File, tamano_leer int32) TABLA_INODOS {
	tablaInodo := TABLA_INODOS{}
	data := leerBytes(file, tamano_leer)
	buffer := bytes.NewBuffer(data)
	err := binary.Read(buffer, binary.BigEndian, &tablaInodo)
	if err != nil {
		log.Fatal("ERROR AL LEER ARCHIVO BINARIO", err)
		return TABLA_INODOS{}
	}
	return tablaInodo
}

/*___________________________________  FIN DE LEER TABLA INODOS _________________________*/

/*___________________________________ INICIO DE LEER TABLA INODOS____________________________________*/

func Leer_BloqueCarpeta(file *os.File, tamano_leer int32) BLOQUE_CARPETA {
	bloqueCarpeta := BLOQUE_CARPETA{}
	data := leerBytes(file, tamano_leer)
	buffer := bytes.NewBuffer(data)
	err := binary.Read(buffer, binary.BigEndian, &bloqueCarpeta)
	if err != nil {
		log.Fatal("ERROR AL LEER ARCHIVO BINARIO", err)
		return BLOQUE_CARPETA{}
	}
	return bloqueCarpeta
}

/*___________________________________  FIN DE LEER TABLA INODOS _________________________*/

/*___________________________________ INICIO DE LEER TABLA INODOS____________________________________*/

func Leer_BloqueArchivo(file *os.File, tamano_leer int32) BLOQUE_ARCHIVO {
	bloqueArchivo := BLOQUE_ARCHIVO{}
	data := leerBytes(file, tamano_leer)
	buffer := bytes.NewBuffer(data)
	err := binary.Read(buffer, binary.BigEndian, &bloqueArchivo)
	if err != nil {
		log.Fatal("ERROR AL LEER ARCHIVO BINARIO", err)
		return BLOQUE_ARCHIVO{}
	}
	return bloqueArchivo
}

/*___________________________________  FIN DE LEER TABLA INODOS _________________________*/
