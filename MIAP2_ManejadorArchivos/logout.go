package main

import "bytes"

func Logout() int32 {
	if actualSesion.hay_Sesion {
		actualSesion.Id_user = -1
		actualSesion.Id_grp = -1
		actualSesion.InicioSuper = -1
		actualSesion.Tipo_sistema = -1
		actualSesion.Path = ""
		//temp := bytes.Repeat([]byte{0}, len(actualSesion.Fit[:]))
		copy(actualSesion.Fit[:], bytes.Repeat([]byte{0}, len(actualSesion.Fit[:])))
		actualSesion.hay_Sesion = false
		return 0

	} else {
		return 1

	}
	//return -1
}
