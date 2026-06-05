package symbols

import (
	"fmt"

	"gustavocoutino.compilador/internal/types"
)

type Funcion struct {
	Nombre string // nombre de la función
	TipoRetorno types.Tipo // tipo de retorno
	Parametros []*Variable // parámetros de la función
	Variables *TablaVariables // todas las variables usadas
    DirInicio int // cuadruplo de inicio
    Recursos int // recursos utilizados por la función cada llamada
}

// DirectorioFunciones es una estructura que tiene un 
// atributo: Funciones, el cual es un mapa con llave string (nombre de la funcion)
// y valor Funcion (objeto de la funcion)
type DirectorioFunciones struct {
	Funciones map[string]*Funcion
}

// NewDirectorioFunciones es un constructor
func NewDirectorioFunciones() *DirectorioFunciones {
	return &DirectorioFunciones{Funciones: make(map[string]*Funcion)}
}

// Agregar añade una función al directorio de funciones. Añade los parametros
// a la tabla de variables
func (d *DirectorioFunciones) Agregar(nombre string, tipo types.Tipo, parametros []*Variable) (*Funcion, error) {
    if _, existe := d.Funciones[nombre]; existe {
        return nil, fmt.Errorf("función '%s' ya declarada", nombre)
    }
    f := &Funcion{
        Nombre:      nombre,
        TipoRetorno: tipo,
        Parametros:  parametros,
        Variables:   NewTablaVariables(),
    }
    for _, v := range parametros {
        if err := f.Variables.Agregar(v); err != nil {
          return nil, fmt.Errorf("en función '%s': %w", nombre, err)
      }
    }
    d.Funciones[nombre] = f
    return f, nil
}

// Buscar regresa si existe una función y el objeto correspondiente
// en base a un nombre
func (d *DirectorioFunciones) Buscar(nombre string) (*Funcion, bool) {
    f, existe := d.Funciones[nombre]
    return f, existe
}

