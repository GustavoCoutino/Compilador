package symbols

import (
	"fmt"

	"gustavocoutino.compilador/internal/types"
)

type Variable struct {
	Nombre string // nombre de la variable
	Tipo types.Tipo // tipo de la variable
	Direccion int /// direccion virtual
}

// TablaVariables es una estructura que tiene un 
// atributo: Variables, el cual es un mapa con llave string (nombre de la variable)
// y valor Variable (objeto de la variable)
type TablaVariables struct {
	Variables map[string]*Variable
}

// NewTablaVariables es un constructor
func NewTablaVariables() *TablaVariables {
	return &TablaVariables{Variables: make(map[string]*Variable)} 
}

// Agregar añade una variable a la tabla de variables global o local.
func (t *TablaVariables) Agregar(v *Variable) error {
	if _, existe := t.Variables[v.Nombre]; existe {
		return fmt.Errorf("variable '%s' ya declarada", v.Nombre)
	}
	t.Variables[v.Nombre] = v
	return nil
}

// Buscar regresa si existe una variable y el objeto correspondiente
// en base a un nombre
func (t *TablaVariables) Buscar(nombre string) (*Variable, bool) {
    v, existe := t.Variables[nombre]
    return v, existe
}