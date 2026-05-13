package main

import (
	"fmt"
	"os"
)

var (
	directorioFunciones = NewDirectorioFunciones()
	tablaVariablesGlobal = NewTablaVariables()
	funcionActual *Funcion
	listaIdsActual []string
	listaParametrosActual []*Variable
)

func DeclararVariable(nombre string, tipo Tipo) {
    var tabla *TablaVariables
    if funcionActual != nil {
        tabla = funcionActual.Variables
    } else {
        tabla = tablaVariablesGlobal
    }
    if err := tabla.Agregar(nombre, tipo); err != nil {
        ErrorSemantico(err.Error())
    }
}

func BuscarVariable(nombre string) (*Variable, bool) {
	if funcionActual != nil {
		if variable, ok := funcionActual.Variables.Buscar(nombre); ok {
			return variable, ok
		}
	}
	return tablaVariablesGlobal.Buscar(nombre)
}

func IniciarFuncion(nombre string, tipoRetorno Tipo, parametros []*Variable) {
	f, err := directorioFunciones.Agregar(nombre, tipoRetorno, parametros)
    if err != nil {
        ErrorSemantico(err.Error())
        return
    }
    funcionActual = f
}

func TerminarFuncion() {
	funcionActual = nil
}

func ImprimirTablaVariablesGlobal() {
	for k, v := range tablaVariablesGlobal.Variables {
		fmt.Println(k)
		fmt.Println(v.Tipo)
		fmt.Print("\n")
	}
}

func ImprimirDirectorioFunciones() {
	for _, v := range directorioFunciones.Funciones {
		fmt.Println(v.Nombre)
		fmt.Println(v.TipoRetorno)
		fmt.Println("Parametros:")
		for _, vv := range v.Parametros {
			fmt.Println(vv.Nombre)
			fmt.Println(vv.Tipo)
		}
		fmt.Println("Variables")
		for _, vv := range v.Variables.Variables {
			fmt.Println(vv.Nombre)
			fmt.Println(vv.Tipo)
		}
		fmt.Print("\n")
	}
}

func ErrorSemantico(msg string) {
    fmt.Fprintf(os.Stderr, "Error semántico: %s\n", msg)
}