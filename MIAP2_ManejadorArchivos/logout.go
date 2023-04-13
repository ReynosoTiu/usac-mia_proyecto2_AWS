package main

import "bytes"

func Logout() string {
	if actualSesion.hay_Sesion {
		actualSesion.Id_user = -1
		actualSesion.Id_grp = -1
		actualSesion.InicioSuper = -1
		actualSesion.Tipo_sistema = -1
		actualSesion.Path = ""
		//temp := bytes.Repeat([]byte{0}, len(actualSesion.Fit[:]))
		copy(actualSesion.Fit[:], bytes.Repeat([]byte{0}, len(actualSesion.Fit[:])))
		actualSesion.hay_Sesion = false
		return "SESION CERRADA EXITOSAMENTE"

	} else {
		return "No existe una sesion activa"

	}
	//return -1
}
