package main

import (
	"fmt"
	"strconv"
)

func (nodo *NODO) incialiar(path string, name string, letra byte, numero int32, id string, contador int32) {
	nodo.Path = path
	nodo.Name = name
	nodo.Letra = letra
	nodo.Numero = numero
	nodo.Contador = contador
	nodo.Id = id
	nodo.Siguiente = nil
}

func inicialiarList() *LISTASIMPLE {
	//ls.Primero = nil
	//ls.Ultimo = nil
	//ls.Letra_temp = 48
	//ls.Numero_temp = 0
	//ls.Contador_temp = 0
	return &LISTASIMPLE{nil, nil, 0, 0, 0, "", ""}
}

func (ls *LISTASIMPLE) inicialiarLista() {
	ls.Primero = nil
	ls.Ultimo = nil
	ls.Letra_temp = 48
	ls.Numero_temp = 0
	ls.Contador_temp = 0
}

func (ls *LISTASIMPLE) Esta_Vacio() bool {
	return ls.Primero == nil
}

func (ls *LISTASIMPLE) insertar_Nodo(path string, name string) int64 {
	if existeArchivo(path) {
		if Obtener_extension_archivo(path) == ".dk" {
			if Buscar_Nombre_P_E_L(path, name) == 1 {
				ls.obtener_letra(path, name)
				var new_Node NODO //= nil //Nodo_mount(path,name,letra_temp,numero_temp,obtener_id(),contador_temp);
				new_Node.incialiar(path, name, ls.Letra_temp, ls.Numero_temp, ls.obtener_id(), ls.Contador_temp)
				if ls.Esta_Vacio() {
					ls.Primero = &new_Node
					ls.Primero.Siguiente = nil
					ls.Ultimo = ls.Primero

				} else {
					ls.Ultimo.Siguiente = &new_Node
					new_Node.Siguiente = nil
					ls.Ultimo = &new_Node
				}
				ls.mostrar_particiones()

			} else {
				fmt.Println("ERROR AL MONTAR, PARTICION NO EXISTE")
				return -1

			}

		} else {
			fmt.Println("ERROR EN MONTRAR EXTENSION DESCONOCIDO")
			return -1
		}

	} else {
		fmt.Println("ERROR EN MONTRAR ARCHIVO NO EXISTE")
		return -1

	}
	return 0
}

/*____________ FIN DE INSERTAR NODO ________*/

func (ls *LISTASIMPLE) borrar_nodo(id string) {
	var nodo_actual *NODO = ls.Primero
	var nodo_anterior *NODO = nil

	var es_encontrado bool = false
	if ls.Primero != nil {
		for nodo_actual != nil && es_encontrado != true {
			if nodo_actual.Id == id {
				if nodo_actual == ls.Primero {
					ls.Primero = ls.Primero.Siguiente

				} else if nodo_actual == ls.Ultimo {
					ls.Ultimo = nodo_anterior
					ls.Ultimo.Siguiente = nil
				} else {
					nodo_anterior.Siguiente = nodo_actual.Siguiente
				}

				es_encontrado = true
			}

			nodo_anterior = nodo_actual
			nodo_actual = nodo_actual.Siguiente

			if es_encontrado == true {
				break
			}

		}

	} else {
		fmt.Println("\tLa Lista esta vacia!!!! >:v")
	}

} /*____________ fin de borrar_nodo _______*/

func (ls *LISTASIMPLE) mostrar_particiones() {
	fmt.Println("\t---------------------------------")
	fmt.Println("\t|       Particiones Montadas    |")
	fmt.Println("\t---------------------------------")
	fmt.Println("\t|        Nombre     |   ID      |")
	fmt.Println("\t---------------------------------")
	var aux *NODO = ls.Primero
	for aux != nil {
		fmt.Println("\t|\t", aux.Name, "\t|\t", aux.Id, "\t|")
		fmt.Println("\t---------------------------------")
		aux = aux.Siguiente
	}
}

func (ls *LISTASIMPLE) obtener_letra(path string, name string) {
	ls.Letra_temp = 48
	ls.Numero_temp = 0
	ls.Contador_temp = 0

	var temp *NODO = ls.Primero
	var aumentarLetra int64 = 0

	for temp != nil {

		if path == temp.Path {
			ls.Letra_temp = temp.Letra
			ls.Numero_temp = temp.Numero
			ls.Contador_temp = temp.Contador
		} else {

			ls.Contador_temp = temp.Contador
			//cout<<"\tVALOR CONTADOR:"<<contador_temp<<endl;
		}
		temp = temp.Siguiente
	}

	if ls.Letra_temp == 48 && ls.Numero_temp == 0 {
		ls.Letra_temp = 97 //aqui era A
		//cout<<"\tVALOR CONTADOR:"<<contador_temp<<endl;
		if ls.Contador_temp != 0 {
			ls.Contador_temp++
			ls.Numero_temp = ls.Contador_temp
		} else {
			ls.Contador_temp++
			ls.Numero_temp = 1
		}

	} else {
		aumentarLetra = int64(ls.Letra_temp)
		aumentarLetra++
		ls.Letra_temp = byte(aumentarLetra)
		aumentarLetra = 0
	}

}

func (ls *LISTASIMPLE) obtener_id() string {
	//201504115
	//SE UTILIZO EL CARNET DEL AUX 06
	//fmt.Println("Numero Tem: ",ls.Numero_temp)
	//strconv.Itoa(int(ls.Numero_temp))
	return "15" + string(strconv.Itoa(int(ls.Numero_temp))) + string(ls.Letra_temp)
}

//mount -path=/home/dia/disco.dk -name=Particion1

func (ls *LISTASIMPLE) buscarParticion(id string) bool {
	var temp *NODO = ls.Primero
	for temp != nil {
		if temp.Id == id {
			return true
		}
		temp = temp.Siguiente
	}

	return false
}

func (ls *LISTASIMPLE) buscar_nombre_path(path string, nombr string) bool {
	var temp *NODO = ls.Primero
	for temp != nil {
		if temp.Path == path && temp.Name == nombr {
			return true
		}
		temp = temp.Siguiente
	}

	return false
}

func (ls *LISTASIMPLE) obtenerNodo(id string) *NODO {
	var temp *NODO = ls.Primero
	for temp != nil {
		if temp.Id == id {
			return temp
		}
		temp = temp.Siguiente
	}

	return nil
}
