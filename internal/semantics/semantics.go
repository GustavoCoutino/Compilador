package semantics

import (
	"fmt"
	"os"

	"gustavocoutino.compilador/internal/memory"
	"gustavocoutino.compilador/internal/symbols"
	"gustavocoutino.compilador/internal/types"
)

var (
	directorioFunciones = symbols.NewDirectorioFunciones()
	tablaVariablesGlobal = symbols.NewTablaVariables()
	funcionActual *symbols.Funcion
	listaIdsActual []string
	listaParametrosActual []*symbols.Variable
	errorContador int
	tablaConstantes = map[string]*symbols.Variable{}
)

func DeclararVariable(nombre string, tipo types.Tipo) {
    var tabla *symbols.TablaVariables
    segmento := memory.Global
    if funcionActual != nil {
        tabla = funcionActual.Variables
        segmento = memory.Local
    } else {
        tabla = tablaVariablesGlobal
    }
    if err := tabla.Agregar(nombre, tipo); err != nil {
        ErrorSemantico(err.Error())
        return
    }
    direccion, err := memory.Asignar(segmento, tipo)
    if err != nil {
        ErrorSemantico(err.Error())
        return
    }
    variable, _ := tabla.Buscar(nombre)
    variable.Direccion = direccion
    memory.RegistrarNombre(direccion, nombre)
}

func ProcesarConstante(literal string, tipo types.Tipo) (int, types.Tipo) {
    if variable, ok := tablaConstantes[literal]; ok {
        return variable.Direccion, variable.Tipo
    }
    direccion, err := memory.Asignar(memory.Constante, tipo)
    if err != nil {
        ErrorSemantico(err.Error())
        return -1, types.TipoError
    }
    tablaConstantes[literal] = &symbols.Variable{Nombre: literal, Tipo: tipo, Direccion: direccion}
    memory.RegistrarNombre(direccion, literal)
    return direccion, tipo
}

func BuscarVariable(nombre string) (*symbols.Variable, bool) {
	if funcionActual != nil {
		if variable, ok := funcionActual.Variables.Buscar(nombre); ok {
			return variable, ok
		}
	}
	return tablaVariablesGlobal.Buscar(nombre)
}

func IniciarFuncion(nombre string, tipoRetorno types.Tipo, parametros []*symbols.Variable) {
	f, err := directorioFunciones.Agregar(nombre, tipoRetorno, parametros)
    if err != nil {
        ErrorSemantico(err.Error())
		funcionActual = nil
        return
    }
    funcionActual = f
    for _, parametro := range f.Parametros {
        direccion, err := memory.Asignar(memory.Local, parametro.Tipo)
        if err != nil {
            ErrorSemantico(err.Error())
            continue
        }
        if variable, ok := f.Variables.Buscar(parametro.Nombre); ok {
            variable.Direccion = direccion
        }
        memory.RegistrarNombre(direccion, parametro.Nombre)
    }
}

func TipoRetornoActual() (types.Tipo, bool) {
	if funcionActual == nil {
		return types.TipoNula, false
	}
	return funcionActual.TipoRetorno, true
}

func TerminarFuncion() {
	funcionActual = nil
}

func AgregarIdActual(id string) {
	listaIdsActual = append(listaIdsActual, id)
}

func DeclararIdsActuales(tipo types.Tipo) {
	for _, nombre := range listaIdsActual {
		DeclararVariable(nombre, tipo)
	}
	listaIdsActual = nil
}

func AgregarParametroActual(nombre string, tipo types.Tipo) {
	listaParametrosActual = append(listaParametrosActual, &symbols.Variable{Nombre: nombre, Tipo: tipo})
}

func IniciarFuncionConParametros(nombre string, tipoRetorno types.Tipo) {
	IniciarFuncion(nombre, tipoRetorno, listaParametrosActual)
	listaParametrosActual = nil
}

func IniciarPrograma(nombre string) {
	if err := tablaVariablesGlobal.Agregar(nombre, types.TipoNula); err != nil {
        ErrorSemantico(err.Error())
	}
    if _, err := directorioFunciones.Agregar(nombre, types.TipoNula, nil); err != nil {
        ErrorSemantico(err.Error())
    }
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
		fmt.Println("Variables")
		for _, vv := range v.Variables.Variables {
			fmt.Println(vv.Nombre)
			fmt.Println(vv.Tipo)
		}
		fmt.Print("\n")
	}
}

func ErrorSemantico(msg string) {
	errorContador++
    fmt.Fprintf(os.Stderr, "Error semántico: %s\n", msg)
}

func HasError() bool {
	return errorContador > 0
}