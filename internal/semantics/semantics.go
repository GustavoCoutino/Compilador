package semantics

import (
	"fmt"
	"os"
	"sort"
	"strings"

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

func GetConstantes() map[string]*symbols.Variable {
	return tablaConstantes
}

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

func ObtenerFuncionDireccion () int {
    if funcion, existe := tablaVariablesGlobal.Buscar(funcionActual.Nombre); !existe {
		return -1
	} else {
        return funcion.Direccion
    }
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
    if funcionActual == nil {
        return
    }
    tempEnteros := memory.Contador(memory.Temporal, types.TipoEntero)
    tempFlotantes := memory.Contador(memory.Temporal, types.TipoFlotante)
    funcionActual.Recursos = len(funcionActual.Variables.Variables) + tempEnteros + tempFlotantes
    memory.LiberarMemoria()
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

func imprimirTablaVariables(tabla *symbols.TablaVariables, sangria string) {
	vars := make([]*symbols.Variable, 0, len(tabla.Variables))
	for _, v := range tabla.Variables {
		vars = append(vars, v)
	}
	sort.Slice(vars, func(i, j int) bool { return vars[i].Direccion < vars[j].Direccion })

	fmt.Printf("%s%-14s %-10s %-10s\n", sangria, "Nombre", "Tipo", "Dirección")
	for _, v := range vars {
		fmt.Printf("%s%-14s %-10s %-10d\n", sangria, v.Nombre, v.Tipo, v.Direccion)
	}
}

func ImprimirTablaVariablesGlobal() {
	fmt.Println("Tabla de variables global")
	imprimirTablaVariables(tablaVariablesGlobal, "")
	fmt.Println()
}

func ImprimirDirectorioFunciones() {
	funcs := make([]*symbols.Funcion, 0, len(directorioFunciones.Funciones))
	for _, f := range directorioFunciones.Funciones {
		funcs = append(funcs, f)
	}
	sort.Slice(funcs, func(i, j int) bool { return funcs[i].DirInicio < funcs[j].DirInicio })

	fmt.Println("Directorio de funciones")
	fmt.Printf("%-14s %-10s %-8s %-9s %s\n", "Función", "Retorno", "Inicio", "Recursos", "Parámetros")
	for _, f := range funcs {
		partes := make([]string, len(f.Parametros))
		for i, p := range f.Parametros {
			partes[i] = fmt.Sprintf("%s:%s", p.Nombre, p.Tipo)
		}
		fmt.Printf("%-14s %-10s %-8d %-9d %s\n", f.Nombre, f.TipoRetorno, f.DirInicio, f.Recursos, strings.Join(partes, ", "))
	}
	fmt.Println()

	for _, f := range funcs {
		if len(f.Variables.Variables) == 0 {
			continue
		}
		fmt.Printf("  %s:\n", f.Nombre)
		imprimirTablaVariables(f.Variables, "  ")
		fmt.Println()
	}
}

func ErrorSemantico(msg string) {
	errorContador++
    fmt.Fprintf(os.Stderr, "Error semántico: %s\n", msg)
}

func HasError() bool {
	return errorContador > 0
}