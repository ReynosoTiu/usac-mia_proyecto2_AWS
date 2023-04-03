package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"os"
	"unsafe"
)

func Print_fdisk() {
	fmt.Println("Estoy en Fdisk")
}

func Crear_particion(size int32, unit string, path string, type_ string, fit string, name string) int {

	if type_ == "p" {
		if buscar_indice_libre(path, "p") != -1 {
			if Buscar_Nombre_P_E_L(path, name) == -1 {
				particion_primaria(size, unit, path, type_, fit, name)
			} else {
				fmt.Println("ERROR AL CREAR PARTICION CON NOMBRE REPETIDO")
			}
		} else {
			fmt.Println("ERROR,YA EXISTEN 4 PARTICIONES")
		}

	} else if type_ == "e" {
		if buscar_indice_libre(path, "e") == -1 {
			if Buscar_Nombre_P_E_L(path, name) == -1 {
				particion_extendida(size, unit, path, type_, fit, name)
			} else {
				fmt.Println("ERROR AL CREAR PARTICION CON NOMBRE REPETIDO")
			}
		} else {
			fmt.Println("ERROR,YA EXISTE UNA PARTICION EXTENDIDA")
		}

	} else if type_ == "l" {
		if Buscar_Nombre_P_E_L(path, name) == -1 {
			particion_logica(size, unit, path, type_, fit, name)
		} else {
			fmt.Println("ERROR AL CREAR PARTICION CON NOMBRE REPETIDO")
		}

	} else {
		if Buscar_Nombre_P_E_L(path, name) == -1 {
			particion_primaria(size, unit, path, type_, fit, name)
		} else {
			fmt.Println("ERROR AL CREAR PARTICION CON NOMBRE REPETIDO")
		}
	}

	return 0
}

func particion_primaria(size int32, unit string, path string, type_ string, fit string, name string) {
	masterBoot := leer_archivo(path)
	var indice int32 = Posicion_arreglo_disponible(&masterBoot, size)

	if indice != -1 {
		if unit == "b" {
			size = size * 1
		} else if unit == "k" {
			size = size * 1024
		} else if unit == "m" {
			size = size * 1024 * 1024
		} else if unit == "" {
			size = size * 1024
		}

		if type_ == "" {
			type_ = "p"
		}

		if size <= EspacioDisponible(path) {
			if fit == "ff" {
				//primer ajuste(FF)
				if indice == 0 {

					masterBoot.Mbr_partition[indice].Part_status = 49
					masterBoot.Mbr_partition[indice].Part_type = ([]byte(type_))[0]

					copy(masterBoot.Mbr_partition[indice].Part_fit[:], fit[:])
					masterBoot.Mbr_partition[indice].Part_start = int32(unsafe.Sizeof(masterBoot))
					masterBoot.Mbr_partition[indice].Part_size = size
					copy(masterBoot.Mbr_partition[indice].Part_name[:], name[:])

				} else {
					masterBoot.Mbr_partition[indice].Part_status = 49
					masterBoot.Mbr_partition[indice].Part_type = ([]byte(type_))[0]
					copy(masterBoot.Mbr_partition[indice].Part_fit[:], fit[:])
					masterBoot.Mbr_partition[indice].Part_start = masterBoot.Mbr_partition[indice-1].Part_start + masterBoot.Mbr_partition[indice-1].Part_size
					masterBoot.Mbr_partition[indice].Part_size = size
					copy(masterBoot.Mbr_partition[indice].Part_name[:], name[:])
				}
				Modificar_archivo(&masterBoot, path)
				fmt.Println("PARTICION PRIMARIA CREADA CORRECTAMENTE")
			} else if fit == "bf" {
				//mejor ajuste(BF)
				var mejorIndice int32 = indice
				for i := 0; i < 4; i++ {
					if masterBoot.Mbr_partition[i].Part_start == -1 || (masterBoot.Mbr_partition[i].Part_status == 48 && masterBoot.Mbr_partition[i].Part_size >= size) {
						if int32(i) != indice {
							if masterBoot.Mbr_partition[mejorIndice].Part_size > masterBoot.Mbr_partition[i].Part_size {
								mejorIndice = int32(i)
								break
							}
						}
					}
				}
				if mejorIndice == 0 {
					if fit == "ff" {
						//primer ajuste(FF)
						if indice == 0 {

							masterBoot.Mbr_partition[indice].Part_status = 49
							masterBoot.Mbr_partition[indice].Part_type = ([]byte(type_))[0]

							copy(masterBoot.Mbr_partition[indice].Part_fit[:], fit[:])
							masterBoot.Mbr_partition[indice].Part_start = int32(unsafe.Sizeof(masterBoot))
							masterBoot.Mbr_partition[indice].Part_size = size
							copy(masterBoot.Mbr_partition[indice].Part_name[:], name[:])

						} else {
							masterBoot.Mbr_partition[indice].Part_status = 49
							masterBoot.Mbr_partition[indice].Part_type = ([]byte(type_))[0]
							copy(masterBoot.Mbr_partition[indice].Part_fit[:], fit[:])
							masterBoot.Mbr_partition[indice].Part_start = masterBoot.Mbr_partition[indice-1].Part_start + masterBoot.Mbr_partition[indice-1].Part_size
							masterBoot.Mbr_partition[indice].Part_size = size
							copy(masterBoot.Mbr_partition[indice].Part_name[:], name[:])
						}
						Modificar_archivo(&masterBoot, path)
						fmt.Println("PARTICION PRIMARIA CREADA CORRECTAMENTE")
					} else if fit == "bf" {
						//mejor ajuste(BF)
						var mejorIndice int32 = indice
						for i := 0; i < 4; i++ {
							if masterBoot.Mbr_partition[i].Part_start == -1 || (masterBoot.Mbr_partition[i].Part_status == 48 && masterBoot.Mbr_partition[i].Part_size >= size) {
								if int32(i) != indice {
									if masterBoot.Mbr_partition[mejorIndice].Part_size > masterBoot.Mbr_partition[i].Part_size {
										mejorIndice = int32(i)
										break
									}
								}
							}
						}
						if mejorIndice == 0 {

							masterBoot.Mbr_partition[mejorIndice].Part_status = 49
							masterBoot.Mbr_partition[mejorIndice].Part_type = ([]byte(type_))[0]

							copy(masterBoot.Mbr_partition[mejorIndice].Part_fit[:], fit[:])
							masterBoot.Mbr_partition[mejorIndice].Part_start = int32(unsafe.Sizeof(masterBoot))
							masterBoot.Mbr_partition[mejorIndice].Part_size = size
							copy(masterBoot.Mbr_partition[mejorIndice].Part_name[:], name[:])

						} else {
							masterBoot.Mbr_partition[mejorIndice].Part_status = 49
							masterBoot.Mbr_partition[mejorIndice].Part_type = ([]byte(type_))[0]
							copy(masterBoot.Mbr_partition[mejorIndice].Part_fit[:], fit[:])
							masterBoot.Mbr_partition[mejorIndice].Part_start = masterBoot.Mbr_partition[indice-1].Part_start + masterBoot.Mbr_partition[indice-1].Part_size
							masterBoot.Mbr_partition[mejorIndice].Part_size = size
							copy(masterBoot.Mbr_partition[mejorIndice].Part_name[:], name[:])
						}
						Modificar_archivo(&masterBoot, path)
						fmt.Println("PARTICION PRIMARIA CREADA CORRECTAMENTE")

					} else if fit == "wf" {
						//peor ajuste(WF)
						var peorIndice int32 = indice
						for i := 0; i < 4; i++ {
							if masterBoot.Mbr_partition[i].Part_start == -1 || (masterBoot.Mbr_partition[i].Part_status == 48 && masterBoot.Mbr_partition[i].Part_size >= size) {
								if int32(i) != indice {
									if masterBoot.Mbr_partition[peorIndice].Part_size < masterBoot.Mbr_partition[i].Part_size {
										peorIndice = int32(i)
										break
									}
								}
							}
						}
						if peorIndice == 0 {

							masterBoot.Mbr_partition[peorIndice].Part_status = 49
							masterBoot.Mbr_partition[peorIndice].Part_type = ([]byte(type_))[0]

							copy(masterBoot.Mbr_partition[peorIndice].Part_fit[:], fit[:])
							masterBoot.Mbr_partition[peorIndice].Part_start = int32(unsafe.Sizeof(masterBoot))
							masterBoot.Mbr_partition[peorIndice].Part_size = size
							copy(masterBoot.Mbr_partition[peorIndice].Part_name[:], name[:])

						} else {
							masterBoot.Mbr_partition[peorIndice].Part_status = 49
							masterBoot.Mbr_partition[peorIndice].Part_type = ([]byte(type_))[0]
							copy(masterBoot.Mbr_partition[peorIndice].Part_fit[:], fit[:])
							masterBoot.Mbr_partition[peorIndice].Part_start = masterBoot.Mbr_partition[indice-1].Part_start + masterBoot.Mbr_partition[indice-1].Part_size
							masterBoot.Mbr_partition[peorIndice].Part_size = size
							copy(masterBoot.Mbr_partition[peorIndice].Part_name[:], name[:])
						}
						Modificar_archivo(&masterBoot, path)
						fmt.Println("PARTICION PRIMARIA CREADA CORRECTAMENTE")

					} else if fit == "" {
						//peor ajuste(WF)
						var peorIndice int32 = indice
						for i := 0; i < 4; i++ {
							if masterBoot.Mbr_partition[i].Part_start == -1 || (masterBoot.Mbr_partition[i].Part_status == 48 && masterBoot.Mbr_partition[i].Part_size >= size) {
								if int32(i) != indice {
									if masterBoot.Mbr_partition[peorIndice].Part_size < masterBoot.Mbr_partition[i].Part_size {
										peorIndice = int32(i)
										break
									}
								}
							}
						}
						if peorIndice == 0 {

							masterBoot.Mbr_partition[peorIndice].Part_status = 49
							masterBoot.Mbr_partition[peorIndice].Part_type = ([]byte(type_))[0]

							copy(masterBoot.Mbr_partition[peorIndice].Part_fit[:], fit[:])
							masterBoot.Mbr_partition[peorIndice].Part_start = int32(unsafe.Sizeof(masterBoot))
							masterBoot.Mbr_partition[peorIndice].Part_size = size
							copy(masterBoot.Mbr_partition[peorIndice].Part_name[:], name[:])

						} else {
							masterBoot.Mbr_partition[peorIndice].Part_status = 49
							masterBoot.Mbr_partition[peorIndice].Part_type = ([]byte(type_))[0]
							copy(masterBoot.Mbr_partition[peorIndice].Part_fit[:], fit[:])
							masterBoot.Mbr_partition[peorIndice].Part_start = masterBoot.Mbr_partition[indice-1].Part_start + masterBoot.Mbr_partition[indice-1].Part_size
							masterBoot.Mbr_partition[peorIndice].Part_size = size
							copy(masterBoot.Mbr_partition[peorIndice].Part_name[:], name[:])
						}
						Modificar_archivo(&masterBoot, path)
						fmt.Println("PARTICION PRIMARIA CREADA CORRECTAMENTE")

					}
					masterBoot.Mbr_partition[mejorIndice].Part_status = 49
					masterBoot.Mbr_partition[mejorIndice].Part_type = ([]byte(type_))[0]

					copy(masterBoot.Mbr_partition[mejorIndice].Part_fit[:], fit[:])
					masterBoot.Mbr_partition[mejorIndice].Part_start = int32(unsafe.Sizeof(masterBoot))
					masterBoot.Mbr_partition[mejorIndice].Part_size = size
					copy(masterBoot.Mbr_partition[mejorIndice].Part_name[:], name[:])

				} else {
					masterBoot.Mbr_partition[mejorIndice].Part_status = 49
					masterBoot.Mbr_partition[mejorIndice].Part_type = ([]byte(type_))[0]
					copy(masterBoot.Mbr_partition[mejorIndice].Part_fit[:], fit[:])
					masterBoot.Mbr_partition[mejorIndice].Part_start = masterBoot.Mbr_partition[indice-1].Part_start + masterBoot.Mbr_partition[indice-1].Part_size
					masterBoot.Mbr_partition[mejorIndice].Part_size = size
					copy(masterBoot.Mbr_partition[mejorIndice].Part_name[:], name[:])
				}
				Modificar_archivo(&masterBoot, path)
				fmt.Println("PARTICION PRIMARIA CREADA CORRECTAMENTE")

			} else if fit == "wf" {
				//peor ajuste(WF)
				var peorIndice int32 = indice
				for i := 0; i < 4; i++ {
					if masterBoot.Mbr_partition[i].Part_start == -1 || (masterBoot.Mbr_partition[i].Part_status == 48 && masterBoot.Mbr_partition[i].Part_size >= size) {
						if int32(i) != indice {
							if masterBoot.Mbr_partition[peorIndice].Part_size < masterBoot.Mbr_partition[i].Part_size {
								peorIndice = int32(i)
								break
							}
						}
					}
				}
				if peorIndice == 0 {

					masterBoot.Mbr_partition[peorIndice].Part_status = 49
					masterBoot.Mbr_partition[peorIndice].Part_type = ([]byte(type_))[0]

					copy(masterBoot.Mbr_partition[peorIndice].Part_fit[:], fit[:])
					masterBoot.Mbr_partition[peorIndice].Part_start = int32(unsafe.Sizeof(masterBoot))
					masterBoot.Mbr_partition[peorIndice].Part_size = size
					copy(masterBoot.Mbr_partition[peorIndice].Part_name[:], name[:])

				} else {
					masterBoot.Mbr_partition[peorIndice].Part_status = 49
					masterBoot.Mbr_partition[peorIndice].Part_type = ([]byte(type_))[0]
					copy(masterBoot.Mbr_partition[peorIndice].Part_fit[:], fit[:])
					masterBoot.Mbr_partition[peorIndice].Part_start = masterBoot.Mbr_partition[indice-1].Part_start + masterBoot.Mbr_partition[indice-1].Part_size
					masterBoot.Mbr_partition[peorIndice].Part_size = size
					copy(masterBoot.Mbr_partition[peorIndice].Part_name[:], name[:])
				}
				Modificar_archivo(&masterBoot, path)
				fmt.Println("PARTICION PRIMARIA CREADA CORRECTAMENTE")

			} else if fit == "" {
				//peor ajuste(WF)
				var peorIndice int32 = indice
				for i := 0; i < 4; i++ {
					if masterBoot.Mbr_partition[i].Part_start == -1 || (masterBoot.Mbr_partition[i].Part_status == 48 && masterBoot.Mbr_partition[i].Part_size >= size) {
						if int32(i) != indice {
							if masterBoot.Mbr_partition[peorIndice].Part_size < masterBoot.Mbr_partition[i].Part_size {
								peorIndice = int32(i)
								break
							}
						}
					}
				}
				if peorIndice == 0 {

					masterBoot.Mbr_partition[peorIndice].Part_status = 49
					masterBoot.Mbr_partition[peorIndice].Part_type = ([]byte(type_))[0]

					copy(masterBoot.Mbr_partition[peorIndice].Part_fit[:], "wf")
					masterBoot.Mbr_partition[peorIndice].Part_start = int32(unsafe.Sizeof(masterBoot))
					masterBoot.Mbr_partition[peorIndice].Part_size = size
					copy(masterBoot.Mbr_partition[peorIndice].Part_name[:], name[:])

				} else {
					masterBoot.Mbr_partition[peorIndice].Part_status = 49
					masterBoot.Mbr_partition[peorIndice].Part_type = ([]byte(type_))[0]
					copy(masterBoot.Mbr_partition[peorIndice].Part_fit[:], "wf")
					masterBoot.Mbr_partition[peorIndice].Part_start = masterBoot.Mbr_partition[indice-1].Part_start + masterBoot.Mbr_partition[indice-1].Part_size
					masterBoot.Mbr_partition[peorIndice].Part_size = size
					copy(masterBoot.Mbr_partition[peorIndice].Part_name[:], name[:])
				}
				Modificar_archivo(&masterBoot, path)
				fmt.Println("PARTICION PRIMARIA CREADA CORRECTAMENTE")

			}

		} else {
			fmt.Println("ERROR AL CREAR PARTICION PRIMARIA ESPACIO INSUFICIENTE.   ESPACIO DISPONIBLE:", EspacioDisponible(path), "ESPACIO FALTANTE: ", (size - EspacioDisponible(path)))
		}

	} else {
		fmt.Println("ERROR,AL CREAR PARTICION YA EXISTE 4 PARTICIONES")
	}
}

func particion_extendida(size int32, unit string, path string, type_ string, fit string, name string) {
	masterBoot := leer_archivo(path)
	var indice int32 = Posicion_arreglo_disponible(&masterBoot, size)

	if indice != -1 {
		if unit == "b" {
			size = size * 1
		} else if unit == "k" {
			size = size * 1024
		} else if unit == "m" {
			size = size * 1024 * 1024
		} else if unit == "" {
			size = size * 1024
		}

		if size <= EspacioDisponible(path) {
			if fit == "ff" {
				//primer ajuste(FF)
				if indice == 0 {

					masterBoot.Mbr_partition[indice].Part_status = 49
					masterBoot.Mbr_partition[indice].Part_type = ([]byte(type_))[0]

					copy(masterBoot.Mbr_partition[indice].Part_fit[:], fit[:])
					masterBoot.Mbr_partition[indice].Part_start = int32(unsafe.Sizeof(masterBoot))
					masterBoot.Mbr_partition[indice].Part_size = size
					copy(masterBoot.Mbr_partition[indice].Part_name[:], name[:])

				} else {
					masterBoot.Mbr_partition[indice].Part_status = 49
					masterBoot.Mbr_partition[indice].Part_type = ([]byte(type_))[0]
					copy(masterBoot.Mbr_partition[indice].Part_fit[:], fit[:])
					masterBoot.Mbr_partition[indice].Part_start = masterBoot.Mbr_partition[indice-1].Part_start + masterBoot.Mbr_partition[indice-1].Part_size
					masterBoot.Mbr_partition[indice].Part_size = size
					copy(masterBoot.Mbr_partition[indice].Part_name[:], name[:])
				}
				Modificar_archivo(&masterBoot, path)
				//se crea el objeto ebr
				extendBoot := EBR{}
				extendBoot.Part_status = 0
				copy(extendBoot.Part_fit[:], fit)
				extendBoot.Part_start = masterBoot.Mbr_partition[indice].Part_start
				extendBoot.Part_size = -1
				extendBoot.Part_next = -1
				copy(extendBoot.Part_name[:], "")
				Modificar_archivo_ebr(&extendBoot, path, extendBoot.Part_start)
				fmt.Println("PARTICION EXTENDIDA CREADA CORRECTAMENTE")

			} else if fit == "bf" {
				//mejor ajuste(BF)
				var mejorIndice int32 = indice
				for i := 0; i < 4; i++ {
					if masterBoot.Mbr_partition[i].Part_start == -1 || (masterBoot.Mbr_partition[i].Part_status == 48 && masterBoot.Mbr_partition[i].Part_size >= size) {
						if int32(i) != indice {
							if masterBoot.Mbr_partition[mejorIndice].Part_size > masterBoot.Mbr_partition[i].Part_size {
								mejorIndice = int32(i)
								break
							}
						}
					}
				}
				if mejorIndice == 0 {

					masterBoot.Mbr_partition[mejorIndice].Part_status = 49
					masterBoot.Mbr_partition[mejorIndice].Part_type = ([]byte(type_))[0]

					copy(masterBoot.Mbr_partition[mejorIndice].Part_fit[:], fit[:])
					masterBoot.Mbr_partition[mejorIndice].Part_start = int32(unsafe.Sizeof(masterBoot))
					masterBoot.Mbr_partition[mejorIndice].Part_size = size
					copy(masterBoot.Mbr_partition[mejorIndice].Part_name[:], name[:])

				} else {
					masterBoot.Mbr_partition[mejorIndice].Part_status = 49
					masterBoot.Mbr_partition[mejorIndice].Part_type = ([]byte(type_))[0]
					copy(masterBoot.Mbr_partition[mejorIndice].Part_fit[:], fit[:])
					masterBoot.Mbr_partition[mejorIndice].Part_start = masterBoot.Mbr_partition[indice-1].Part_start + masterBoot.Mbr_partition[indice-1].Part_size
					masterBoot.Mbr_partition[mejorIndice].Part_size = size
					copy(masterBoot.Mbr_partition[mejorIndice].Part_name[:], name[:])
				}
				Modificar_archivo(&masterBoot, path)
				//se crea el objeto ebr
				extendBoot := EBR{}
				extendBoot.Part_status = 0
				copy(extendBoot.Part_fit[:], fit)
				extendBoot.Part_start = masterBoot.Mbr_partition[mejorIndice].Part_start
				extendBoot.Part_size = -1
				extendBoot.Part_next = -1
				copy(extendBoot.Part_name[:], "")
				Modificar_archivo_ebr(&extendBoot, path, extendBoot.Part_start)
				fmt.Println("PARTICION EXTENDIDA CREADA CORRECTAMENTE")

			} else if fit == "wf" {
				//peor ajuste(WF)
				var peorIndice int32 = indice
				for i := 0; i < 4; i++ {
					if masterBoot.Mbr_partition[i].Part_start == -1 || (masterBoot.Mbr_partition[i].Part_status == 48 && masterBoot.Mbr_partition[i].Part_size >= size) {
						if int32(i) != indice {
							if masterBoot.Mbr_partition[peorIndice].Part_size < masterBoot.Mbr_partition[i].Part_size {
								peorIndice = int32(i)
								break
							}
						}
					}
				}
				if peorIndice == 0 {

					masterBoot.Mbr_partition[peorIndice].Part_status = 49
					masterBoot.Mbr_partition[peorIndice].Part_type = ([]byte(type_))[0]

					copy(masterBoot.Mbr_partition[peorIndice].Part_fit[:], fit[:])
					masterBoot.Mbr_partition[peorIndice].Part_start = int32(unsafe.Sizeof(masterBoot))
					masterBoot.Mbr_partition[peorIndice].Part_size = size
					copy(masterBoot.Mbr_partition[peorIndice].Part_name[:], name[:])

				} else {
					masterBoot.Mbr_partition[peorIndice].Part_status = 49
					masterBoot.Mbr_partition[peorIndice].Part_type = ([]byte(type_))[0]
					copy(masterBoot.Mbr_partition[peorIndice].Part_fit[:], fit[:])
					masterBoot.Mbr_partition[peorIndice].Part_start = masterBoot.Mbr_partition[indice-1].Part_start + masterBoot.Mbr_partition[indice-1].Part_size
					masterBoot.Mbr_partition[peorIndice].Part_size = size
					copy(masterBoot.Mbr_partition[peorIndice].Part_name[:], name[:])
				}
				Modificar_archivo(&masterBoot, path)
				//se crea el objeto ebr
				extendBoot := EBR{}
				extendBoot.Part_status = 0
				copy(extendBoot.Part_fit[:], fit)
				extendBoot.Part_start = masterBoot.Mbr_partition[peorIndice].Part_start
				extendBoot.Part_size = -1
				extendBoot.Part_next = -1
				copy(extendBoot.Part_name[:], "")
				Modificar_archivo_ebr(&extendBoot, path, extendBoot.Part_start)
				fmt.Println("PARTICION EXTENDIDA CREADA CORRECTAMENTE")

			} else if fit == "" {
				//peor ajuste(WF)
				var peorIndice int32 = indice
				for i := 0; i < 4; i++ {
					if masterBoot.Mbr_partition[i].Part_start == -1 || (masterBoot.Mbr_partition[i].Part_status == 48 && masterBoot.Mbr_partition[i].Part_size >= size) {
						if int32(i) != indice {
							if masterBoot.Mbr_partition[peorIndice].Part_size < masterBoot.Mbr_partition[i].Part_size {
								peorIndice = int32(i)
								break
							}
						}
					}
				}
				if peorIndice == 0 {

					masterBoot.Mbr_partition[peorIndice].Part_status = 49
					masterBoot.Mbr_partition[peorIndice].Part_type = ([]byte(type_))[0]

					copy(masterBoot.Mbr_partition[peorIndice].Part_fit[:], "wf")
					masterBoot.Mbr_partition[peorIndice].Part_start = int32(unsafe.Sizeof(masterBoot))
					masterBoot.Mbr_partition[peorIndice].Part_size = size
					copy(masterBoot.Mbr_partition[peorIndice].Part_name[:], name[:])

				} else {
					masterBoot.Mbr_partition[peorIndice].Part_status = 49
					masterBoot.Mbr_partition[peorIndice].Part_type = ([]byte(type_))[0]
					copy(masterBoot.Mbr_partition[peorIndice].Part_fit[:], fit[:])
					masterBoot.Mbr_partition[peorIndice].Part_start = masterBoot.Mbr_partition[indice-1].Part_start + masterBoot.Mbr_partition[indice-1].Part_size
					masterBoot.Mbr_partition[peorIndice].Part_size = size
					copy(masterBoot.Mbr_partition[peorIndice].Part_name[:], name[:])
				}
				Modificar_archivo(&masterBoot, path)
				//se crea el objeto ebr
				extendBoot := EBR{}
				extendBoot.Part_status = 0
				copy(extendBoot.Part_fit[:], "wf")
				extendBoot.Part_start = masterBoot.Mbr_partition[peorIndice].Part_start
				extendBoot.Part_size = -1
				extendBoot.Part_next = -1
				copy(extendBoot.Part_name[:], "")
				Modificar_archivo_ebr(&extendBoot, path, extendBoot.Part_start)
				fmt.Println("PARTICION EXTENDIDA CREADA CORRECTAMENTE")
			}

		} else {
			fmt.Println("ERROR AL CREAR PARTICION EXTENDIDA ESPACIO INSUFICIENTE.   ESPACIO DISPONIBLE:", EspacioDisponible(path), "ESPACIO FALTANTE: ", (size - EspacioDisponible(path)))
		}

	} else {
		fmt.Println("ERROR,AL CREAR PARTICION YA EXISTE 4 PARTICIONES")
	}

}

func particion_logica(size int32, unit string, path string, type_ string, fit string, name string) {
	//os.Open sirve para leer el archivo nada mas y no para modificar
	file, err := os.OpenFile(path, os.O_RDWR, 0777)
	defer file.Close()

	if err != nil {
		log.Fatal("ERROR AL ABRIR ARCHIVO,ARCHIVO NO EXISTE", err)
	}

	var tamano_masterBoot int32 = int32(unsafe.Sizeof(MBR{}))
	file.Seek(0, 0)
	masterBoot := ObtenerInforMbr(file, tamano_masterBoot)

	//------------------------------------------
	var indice int32 = -1
	for i := 0; i < 4; i++ {
		if string(masterBoot.Mbr_partition[i].Part_type) == "e" {
			indice = int32(i)
			break
		}
	}
	if indice != -1 {
		if fit == "" {
			fit = "wf"
		}
		if unit == "b" {
			size = size * 1
		} else if unit == "k" || unit == "" {
			size = size * 1024
		} else if unit == "m" {
			size = size * 1024 * 1024
		}

		extendedBoot := EBR{}
		var cont int32 = masterBoot.Mbr_partition[indice].Part_start
		file.Seek(int64(cont), 0)
		extendedBoot = ObtenerInfoEbr(file, int32(unsafe.Sizeof(EBR{})))

		var espacio_disp int32 = EspacioDisponibleExtendida(path)

		if size <= espacio_disp {
			if extendedBoot.Part_size == -1 { //Si es la primera
				if masterBoot.Mbr_partition[indice].Part_size < size {
					//espacio insuficiente
					fmt.Println("ERROR AL CREAR LOGICA ESPACIO INSUFICIENTE")
				} else {
					extendedBoot.Part_status = 49
					copy(extendedBoot.Part_fit[:], fit)
					pos_actual, _ := file.Seek(0, 1)
					//extendedBoot.Part_start = masterBoot.Mbr_partition[indice].Part_start // ftell(fp) - sizeof(EBR); //Para regresar al inicio de la extendida
					extendedBoot.Part_start = int32(pos_actual) - int32(unsafe.Sizeof(EBR{}))
					extendedBoot.Part_size = size
					extendedBoot.Part_next = -1
					copy(extendedBoot.Part_name[:], name)
					//file.Seek(masterBoot.Mbr_partition[indice].Part_start,0);
					Modificar_archivo_ebr(&extendedBoot, path, masterBoot.Mbr_partition[indice].Part_start) //espremor que no de erro aqui sino mandar el archiv abierto de una
					fmt.Println("PARTICION LOGICA CREADA CORRECTAMENTE")
				}
			} else {
				for (extendedBoot.Part_next != -1) && (Ftell(file) < (masterBoot.Mbr_partition[indice].Part_size + masterBoot.Mbr_partition[indice].Part_start)) {
					file.Seek(int64(extendedBoot.Part_next), 0)
					extendedBoot = ObtenerInfoEbr(file, int32(unsafe.Sizeof(EBR{})))
					//fread(&extendedBoot,sizeof(EBR),1,fp);
				}

				var espacioNecesario int32 = extendedBoot.Part_start + extendedBoot.Part_size + size

				if espacioNecesario <= (masterBoot.Mbr_partition[indice].Part_size + masterBoot.Mbr_partition[indice].Part_start) {
					extendedBoot.Part_next = extendedBoot.Part_start + extendedBoot.Part_size
					//Escribimos el next del ultimo EBR
					mover, _ := file.Seek(int64(Ftell(file))-int64(unsafe.Sizeof(EBR{})), 0)
					Modificar_archivo_ebr(&extendedBoot, path, int32(mover)) //verificar aqui sino F para escribir el archivo
					//Escribimos el nuevo EBR
					file.Seek(int64(extendedBoot.Part_start+extendedBoot.Part_size), 0)
					extendedBoot.Part_status = 49
					copy(extendedBoot.Part_fit[:], fit)
					extendedBoot.Part_start = Ftell(file)
					extendedBoot.Part_size = size
					extendedBoot.Part_next = -1
					copy(extendedBoot.Part_name[:], name)
					//fwrite(&extendedBoot, sizeof(EBR), 1, fp)
					Modificar_archivo_ebr(&extendedBoot, path, Ftell(file))
					fmt.Println("PARTICION LOGICA CREADA CORRECTAMENTE")

				} else {
					fmt.Println("ERROR la particion logica a crear excede el")
					fmt.Println("espacio disponible de la particion extendida")
				}
			}

		} else {
			fmt.Println("ERROR AL CREAR PARTICION LOGICA ESPACIO INSUFICIENTE.   ESPACIO DISPONIBLE:", espacio_disp, "ESPACIO FALTANTE: ", (size - espacio_disp))
		}

	} else {
		fmt.Println("ERROR AL CREAR PARTICION LOGICA,NO EXISTE EXTENDIDA")
	}

	file.Close()
	//fin de particion logica
}

func leer_archivo(path string) MBR {
	//os.Open sirve para leer el archivo nada mas y no para modificar
	file, err := os.Open(path)
	defer file.Close()
	if err != nil {
		log.Fatal("ERROR AL ABRIR ARCHIVO,ARCHIVO NO EXISTE", err)
		os.Exit(4)
	}

	masterBoot := MBR{}

	var tamano_masterBoot int32 = int32(unsafe.Sizeof(masterBoot))
	data := leerBytes(file, tamano_masterBoot)
	buffer := bytes.NewBuffer(data)
	err = binary.Read(buffer, binary.BigEndian, &masterBoot)
	if err != nil {
		log.Fatal("ERROR AL LEER ARCHIVO BINARIO", err)
		os.Exit(4)
	}

	file.Close()

	return masterBoot
}

func leerBytes(file *os.File, number int32) []byte {
	bytes := make([]byte, number)
	_, err := file.Read(bytes)
	if err != nil {
		log.Fatal("ERROR  A LEER BYTES", err)
	}

	return bytes
}

func Posicion_arreglo_disponible(masterBoot *MBR, size int32) int32 {
	for i := 0; i < 4; i++ {
		if masterBoot.Mbr_partition[i].Part_start == -1 || (masterBoot.Mbr_partition[i].Part_status == 48 && masterBoot.Mbr_partition[i].Part_size >= size) {
			return int32(i)
		}
	}
	return -1
}

func Modificar_archivo(masterBoot *MBR, path string) {
	mi_archivo, err := os.OpenFile(path, os.O_RDWR, 0777)
	defer mi_archivo.Close()
	if err != nil {
		log.Fatal("ERROR AL ABRIR ARCHIVO", err)
		//return -1
	}

	mi_archivo.Seek(0, 0)
	var master_buffer bytes.Buffer
	binary.Write(&master_buffer, binary.BigEndian, masterBoot)
	EscribirBytes(mi_archivo, master_buffer.Bytes())
	mi_archivo.Close()
}

func Modificar_archivo_ebr(extendBoot *EBR, path string, moverPuntero int32) {
	mi_archivo, err := os.OpenFile(path, os.O_RDWR, 0777)
	defer mi_archivo.Close()
	if err != nil {
		log.Fatal("ERROR AL ABRIR ARCHIVO", err)
		//return -1
	}

	mi_archivo.Seek(int64(moverPuntero), 0)
	var master_buffer bytes.Buffer
	binary.Write(&master_buffer, binary.BigEndian, extendBoot)
	EscribirBytes(mi_archivo, master_buffer.Bytes())
	mi_archivo.Close()
}

/*___________________ PARTICION EXTENDIDA _____________*/
func Buscar_Nombre_P_E_L(path string, name string) int32 {

	//os.Open sirve para leer el archivo nada mas y no para modificar
	file, err := os.OpenFile(path, os.O_RDWR, 0777)
	defer file.Close()

	if err != nil {
		log.Fatal("ERROR AL ABRIR ARCHIVO,ARCHIVO NO EXISTE", err)
	}
	masterBoot := MBR{}

	var tamano_masterBoot int32 = int32(unsafe.Sizeof(masterBoot))
	masterBoot = ObtenerInforMbr(file, tamano_masterBoot)

	var indice_extendia int32 = -1
	for i := 0; i < 4; i++ {
		if string(bytes.Trim(masterBoot.Mbr_partition[i].Part_name[:], "\x00")) == name {
			file.Close()
			return 1
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
				return 1
			}

			if extendBoot.Part_next == -1 {
				break
			}
			file.Seek(int64(extendBoot.Part_next), 0)
			extendBoot = ObtenerInfoEbr(file, int32(unsafe.Sizeof(EBR{}))) //ruta quemada por el momento xD
		}
	}
	file.Close()

	return -1
}

func ObtenerInforMbr(file *os.File, tamano_leer int32) MBR {
	//file.Seek(0,0) //esto se agrego para mover el puntero en el origen del archivo sino borrar
	masterBoot := MBR{}
	data := leerBytes(file, tamano_leer)
	buffer := bytes.NewBuffer(data)
	err := binary.Read(buffer, binary.BigEndian, &masterBoot)
	if err != nil {
		log.Fatal("ERROR AL LEER ARCHIVO BINARIO", err)
		return MBR{}
	}
	return masterBoot
}

func ObtenerInfoEbr(file *os.File, tamano_leer int32) EBR {
	extendBoot := EBR{}
	data := leerBytes(file, tamano_leer)
	buffer := bytes.NewBuffer(data)
	err := binary.Read(buffer, binary.BigEndian, &extendBoot)
	if err != nil {
		log.Fatal("ERROR AL LEER ARCHIVO BINARIO", err)
		return EBR{}
	}
	return extendBoot
}

func Ftell(file *os.File) int32 {
	var posiconActual int64 = 0
	var err error = nil
	posiconActual, err = (file.Seek(0, 1))
	if err != nil {
		fmt.Println("EXISTEN ERRORE EN LA POSICION ACTUAL FTELL")
	}

	return int32(posiconActual)
}

func EspacioDisponible(path string) int32 {
	masterBoot := leer_archivo(path)
	var espacio_usado int32 = 0
	espacio_usado = int32(unsafe.Sizeof(MBR{}))
	for i := 0; i < 4; i++ {
		if masterBoot.Mbr_partition[i].Part_status == 49 {
			espacio_usado += masterBoot.Mbr_partition[i].Part_size
		}
	}

	return masterBoot.Mbr_tamano - espacio_usado
}

func EspacioDisponibleExtendida(path string) int32 {
	//os.Open sirve para leer el archivo nada mas y no para modificar
	file, err := os.OpenFile(path, os.O_RDWR, 0777)
	defer file.Close()

	if err != nil {
		log.Fatal("ERROR AL ABRIR ARCHIVO,ARCHIVO NO EXISTE", err)
		return -1 //aqui lo regresamos porque no encontro el archivo sino da error
	}

	var tamano_masterBoot int32 = int32(unsafe.Sizeof(MBR{}))
	file.Seek(0, 0)
	masterBoot := ObtenerInforMbr(file, tamano_masterBoot)

	var indice int32 = -1
	var espacio_usado int32 = 0
	for i := 0; i < 4; i++ {
		if masterBoot.Mbr_partition[i].Part_type == 101 && masterBoot.Mbr_partition[i].Part_status == 49 {
			indice = int32(i)
		}

	}

	extendedBoot := EBR{}
	var cont int32 = masterBoot.Mbr_partition[indice].Part_start
	file.Seek(int64(cont), 0)
	extendedBoot = ObtenerInfoEbr(file, int32(unsafe.Sizeof(EBR{})))

	if extendedBoot.Part_size >= 0 {
		espacio_usado += extendedBoot.Part_size
	}

	for (extendedBoot.Part_next != -1) && (Ftell(file) < (masterBoot.Mbr_partition[indice].Part_size + masterBoot.Mbr_partition[indice].Part_start)) {
		file.Seek(int64(extendedBoot.Part_next), 0)
		extendedBoot = ObtenerInfoEbr(file, int32(unsafe.Sizeof(EBR{})))
		espacio_usado += extendedBoot.Part_size
		//fread(&extendedBoot,sizeof(EBR),1,fp);
	}

	return (masterBoot.Mbr_partition[indice].Part_size - espacio_usado)
}

func Ver_particion_Primaria_Extendida_Logica(path string) {
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
	fmt.Println("*___________________________________________*")
	fmt.Println("Disco Fit: ", string(masterBoot.Disk_fit[:]))
	fmt.Println("Disco Asignatura: ", masterBoot.Mbr_disk_signature)
	fmt.Println("Fecha Creacion:", string(masterBoot.Mbr_fecha_creacion[:]))
	fmt.Println("Disco Tamano:", masterBoot.Mbr_tamano)

	var indice_extendia int32 = -1
	for i := 0; i < 4; i++ {
		fmt.Println("*___________________________________________*")
		fmt.Println("Particion[", i, "].Fit   :", string(masterBoot.Mbr_partition[i].Part_fit[:]))
		fmt.Println("Particion[", i, "].Name  :", string(masterBoot.Mbr_partition[i].Part_name[:]))
		fmt.Println("Particion[", i, "].Size  :", masterBoot.Mbr_partition[i].Part_size)
		fmt.Println("Particion[", i, "].Start :", masterBoot.Mbr_partition[i].Part_start)
		fmt.Println("Particion[", i, "].Status:", string(masterBoot.Mbr_partition[i].Part_status))
		fmt.Println("Particion[", i, "].Type  :", string(masterBoot.Mbr_partition[i].Part_type))

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
			fmt.Println("*___________________________________________*")
			fmt.Println("Particion[].Fit   :", string(extendBoot.Part_fit[:]))
			fmt.Println("Particion[].Name  :", string(extendBoot.Part_name[:]))
			fmt.Println("Particion[].Size  :", extendBoot.Part_next)
			fmt.Println("Particion[].Start :", extendBoot.Part_size)
			fmt.Println("Particion[].Status:", extendBoot.Part_start)
			fmt.Println("Particion[].Type  :", string(extendBoot.Part_status))

			if extendBoot.Part_next == -1 {
				break
			}
			file.Seek(int64(extendBoot.Part_next), 0)
			extendBoot = ObtenerInfoEbr(file, int32(unsafe.Sizeof(EBR{}))) //ruta quemada por el momento xD
		}
	}
	file.Close()
}

/*
Metodos

	Error al no crear la operacion soliciada
	por Espacio
	verificarExtendida para logica y verificar tamano
	Solo una extendida no se puede crear particion logica sino existe extendia
	No repetir nombre
	Eliminar_particion_debe existir
*/
func buscar_indice_libre(path string, type_ string) int {
	masterBoot := leer_archivo(path)
	if type_ == "p" {
		for i := 0; i < 4; i++ {
			if masterBoot.Mbr_partition[i].Part_status == 0 {
				return i
			}
		}
		return -1 //indica que ya no espacio para primaria
	} else if type_ == "e" {
		for i := 0; i < 4; i++ {
			if masterBoot.Mbr_partition[i].Part_type == 101 {
				return i
			}
		}
		return -1 //Inidica que no hay extendida
	}

	return -1
}
