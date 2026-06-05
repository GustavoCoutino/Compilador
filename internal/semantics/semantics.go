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
	directorioFunciones = symbols.NewDirectorioFunciones() // directorio de funciones 
	tablaVariablesGlobal = symbols.NewTablaVariables() // tabla de variables global
	funcionActual *symbols.Funcion // guarda la funcion actual leida
	listaIdsActual []string // lista de identificadors declarados juntos
	listaParametrosActual []*symbols.Variable // lista de parametros de la funcion actual
	errorContador int // contador de errores semánticos
	tablaConstantes = map[string]*symbols.Variable{} // tabla de constantes
	LineaActual int // linea actual de compilacion para rastrear error
)

// GetConstantes es un getter de la tabla de constantes
func GetConstantes() map[string]*symbols.Variable {
	return tablaConstantes
}

// AsignarCuadruploInicio le asigna a la funcion actual el cuadruplo
// de inicio
func AsignarCuadruploInicio(inicio int) {
    if funcionActual != nil {
        funcionActual.DirInicio = inicio
    } 
}

// DeclararVariable declara una variable en su tabla respectiva de variables.
// Se le asigna a la variable declarada un segmento de memoria basado en su alcance
// y tipo.
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

// ProcesarConstante registra una variable en la tabla de constantes y le asigna una direccion en memoria.
// Regresa la direccion y tipo de la constante
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

// BuscarVariable busca una variable por nombre en un alcance local o 
// global. Regresa la variable encontrada y estado (se encontro o no)
func BuscarVariable(nombre string) (*symbols.Variable, bool) {
	if funcionActual != nil {
		if variable, ok := funcionActual.Variables.Buscar(nombre); ok {
			return variable, ok
		}
	}
	return tablaVariablesGlobal.Buscar(nombre)
}

// ObtenerFuncionDireccion regresa la direccion asignada en memoria de una funcion
func ObtenerFuncionDireccion () int {
    if funcion, existe := tablaVariablesGlobal.Buscar(funcionActual.Nombre); !existe {
		return -1
	} else {
        return funcion.Direccion
    }
} 

// IniciarFuncion comienza el proceso de agregar una funcion al directorio de funciones. 
// Le asigna a los parametros de la función una direccion virtual
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

// TipoRetornoActual regresa el tipo de retorno de 
// la función que se está leyendo
func TipoRetornoActual() (types.Tipo, bool) {
	if funcionActual == nil {
		return types.TipoNula, false
	}
	return funcionActual.TipoRetorno, true
}

// TerminarFuncion reinicia los estados del análisis semántico,
// libera memoria de la función local, y calcula los recursos.
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

// AgregarIdActual agrega un identificador a la lista de ids
func AgregarIdActual(id string) {
	listaIdsActual = append(listaIdsActual, id)
}

// DeclararIdsActuales toma los elementos de listaIdsActual
// y los declara en la tabla de variables
func DeclararIdsActuales(tipo types.Tipo) {
	for _, nombre := range listaIdsActual {
		DeclararVariable(nombre, tipo)
	}
	listaIdsActual = nil
}

// AgregarParametroActual agrega un parametro a la lista de parámetros de una función
func AgregarParametroActual(nombre string, tipo types.Tipo) {
	listaParametrosActual = append(listaParametrosActual, &symbols.Variable{Nombre: nombre, Tipo: tipo})
}

// IniciarFuncionConParametros llama a IniciarFuncion para
func IniciarFuncionConParametros(nombre string, tipoRetorno types.Tipo) {
	IniciarFuncion(nombre, tipoRetorno, listaParametrosActual)
	listaParametrosActual = nil
}

// BuscarFuncion regresa una funcion del directorio de funciones
func BuscarFuncion(nombre string) (*symbols.Funcion, bool) {
    return directorioFunciones.Buscar(nombre)
}

// IniciarPrograma guarda el nombre del programa en la tabla de variables
// global y le asigna direccion virtual
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

func ImprimirTablaConstantes() {
	constantes := make([]*symbols.Variable, 0, len(tablaConstantes))
	for _, v := range tablaConstantes {
		constantes = append(constantes, v)
	}
	sort.Slice(constantes, func(i, j int) bool { return constantes[i].Direccion < constantes[j].Direccion })

	fmt.Println("Tabla de constantes")
	fmt.Printf("%-14s %-10s\n", "Valor", "Dirección")
	for _, c := range constantes {
		fmt.Printf("%-14s %-10d\n", c.Nombre, c.Direccion)
	}
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

// ErrorSemantico incrementa el contador de errores semánticos
// e imprime un mensaje de error
func ErrorSemantico(msg string) {
	errorContador++
	if LineaActual > 0 {
        fmt.Fprintf(os.Stderr, "Error semántico [línea %d]: %s\n", LineaActual, msg)
    } else {
        fmt.Fprintf(os.Stderr, "Error semántico: %s\n", msg)
    }
}

// HasError regresa si el análisis semántico tiene errores
func HasError() bool {
	return errorContador > 0
}