package symbols

import (
	"fmt"

	"gustavocoutino.compilador/internal/types"
)

type Variable struct {
	Nombre string
	Tipo types.Tipo
	Direccion int
}

type TablaVariables struct {
	Variables map[string]*Variable
}

func NewTablaVariables() *TablaVariables {
	return &TablaVariables{Variables: make(map[string]*Variable)} 
}

func (t *TablaVariables) Agregar(v *Variable) error {
	if _, existe := t.Variables[v.Nombre]; existe {
		return fmt.Errorf("variable '%s' ya declarada", v.Nombre)
	}
	t.Variables[v.Nombre] = v
	return nil
}

func (t *TablaVariables) Buscar(nombre string) (*Variable, bool) {
    v, existe := t.Variables[nombre]
    return v, existe
}