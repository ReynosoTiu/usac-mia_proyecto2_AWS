export interface Respuesta {
    Tipo: number,
    Mensaje: string,
    Data: string,
    Ruta: string,
    Extension: string
}

export interface Reporte {
    Ruta: string,
    Data: string,
    NombreSave: string,
    Extension: string
}