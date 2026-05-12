package main

import "fmt"

// Mapa de funciones
// Llave: nombre
// Valor: variables, parametros, tipo retorno

type Funcion struct {
	Nombre string
	TipoRetorno Tipo
	Parametros []*Variable
	Variables *TablaVariables
}

type DirectorioFunciones struct {
	Funciones map[string]*Funcion
}

func NewDirectorioFunciones() *DirectorioFunciones {
	return &DirectorioFunciones{Funciones: make(map[string]*Funcion)}
}

func (d *DirectorioFunciones) Agregar(nombre string, tipo Tipo) (*Funcion, error) {
    if _, existe := d.Funciones[nombre]; existe {
        return nil, fmt.Errorf("función '%s' ya declarada", nombre)
    }
    f := &Funcion{
        Nombre:      nombre,
        TipoRetorno: tipo,
        Parametros:  []*Variable{},
        Variables:   NewTablaVariables(),
    }
    d.Funciones[nombre] = f
    return f, nil
}

func (d *DirectorioFunciones) Buscar(nombre string) (*Funcion, bool) {
    f, existe := d.Funciones[nombre]
    return f, existe
}