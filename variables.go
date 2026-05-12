package main

import "fmt"

type Variable struct {
	Nombre string
	Tipo Tipo
	Direccion int
}

type TablaVariables struct {
	Variables map[string]*Variable
}

func NewTablaVariables() *TablaVariables {
	return &TablaVariables{Variables: make(map[string]*Variable)} 
}

func (t *TablaVariables) Agregar(nombre string, tipo Tipo) error {
	if _, existe := t.Variables[nombre]; existe {
		return fmt.Errorf("variable '%s' ya declarada", nombre)
	}
	t.Variables[nombre] = &Variable{Nombre: nombre, Tipo: tipo}
	return nil
}

func (t *TablaVariables) Buscar(nombre string) (*Variable, bool) {
    v, existe := t.Variables[nombre]
    return v, existe
}