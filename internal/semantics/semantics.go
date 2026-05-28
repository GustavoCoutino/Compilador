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

func AsignarCuadruploInicio(inicio int) {
    if funcionActual != nil {
        funcionActual.DirInicio = inicio
    }       
}

func DeclararVariable(nombre string, tipo types.Tipo) {
    var tabla *symbols.TablaVariables
    segmento := memory.Global
    if funcionActual != nil {
        tabla = funcionActual.Variables
        segmento = memory.Local
    } else {
        tabla = tablaVariablesGlobal
    }
    variable := &symbols.Variable{Nombre: nombre, Tipo: tipo}
    if err := tabla.Agregar(variable); err != nil {
        ErrorSemantico(err.Error())
        return
    }
    direccion, err := memory.Asignar(segmento, tipo)
    if err != nil {
        ErrorSemantico(err.Error())
        return
    }
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
        parametro.Direccion = direccion
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
    funcionActual.Recursos = len(funcionActual.Variables.Variables)
    memory.New().LiberarMemoria()
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

func BuscarFuncion(nombre string) (*symbols.Funcion, bool) {
    return directorioFunciones.Buscar(nombre)
}

func IniciarPrograma(nombre string) {
    variable := &symbols.Variable{Nombre: nombre, Tipo: types.TipoNula}
	if err := tablaVariablesGlobal.Agregar(variable); err != nil {
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
        fmt.Println("Recursos")
        fmt.Println(v.Recursos)
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