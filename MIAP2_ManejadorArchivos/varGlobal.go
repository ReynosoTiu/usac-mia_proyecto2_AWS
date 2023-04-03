package main

import "unsafe"

/*_________________________ VARIABLES GLOBALES _______________*/
var listaS LISTASIMPLE
var actualSesion Sesion
var tamanoCarpeta int32 = int32(unsafe.Sizeof(BLOQUE_CARPETA{}))
var tamanoArchivo int32 = int32(unsafe.Sizeof(BLOQUE_ARCHIVO{}))
var tamanoSuper int32 = int32(unsafe.Sizeof(SUPER_BLOQUE{}))
var tamanoInodo int32 = int32(unsafe.Sizeof(TABLA_INODOS{}))

func InicializarVariablesGlobales() {
	listaS.inicialiarLista()

}
