package quadruples

import (
	"fmt"

	"gustavocoutino.compilador/internal/data_structures/queue"
	"gustavocoutino.compilador/internal/data_structures/stack"
	"gustavocoutino.compilador/internal/memory"
	"gustavocoutino.compilador/internal/ops"
	"gustavocoutino.compilador/internal/semantics"
	"gustavocoutino.compilador/internal/symbols"
	"gustavocoutino.compilador/internal/types"
)

type Quadruple struct {
	Operador int
	Izquierda int
	Derecha int
	Resultado int

}

var (
	pilaOperadores = stack.New[int]() // pila de operadores
	pilaOperandos = stack.New[int]() // pila de operandos
	pilaOperandosType = stack.New[types.Tipo]() // pila de los tipos de los operandos
	filaCuadruplos = queue.New[Quadruple]() // pila de cuadruplos
	pilaDeSaltos = stack.New[int]() // pila de saltos
	contadorParametro = 0 // contador de parametros para operación PARAM
	funcionLlamada *symbols.Funcion // función llamada más reciente
)

// ContadorActual regresa la longitud de la fila de cuadruplos
func ContadorActual() int {
    return filaCuadruplos.Len()
}

// GetFilaCuadruplos es un getter que regresa la fila de cuadruplos
func GetFilaCuadruplos() *queue.Queue[Quadruple] {
	return filaCuadruplos
}

// CrearCuadruploGotoInicio genera el cuadruplo para hacer el salto a main
func CrearCuadruploGotoInicio(){
	pilaDeSaltos.Push(ContadorActual())
	filaCuadruplos.Push(Quadruple{ops.GOTO, -1, -1, -1})
}

// CrearCuadruploFin genera el cuadruplo de fin
func CrearCuadruploFin(){
	filaCuadruplos.Push(Quadruple{ops.FIN, -1, -1, -1})
}

// EmpujarOperando empuja un operando y su tipo a las pilas
// respectivas
func EmpujarOperando(direccion int, tipo types.Tipo) {
	pilaOperandos.Push(direccion)
	pilaOperandosType.Push(tipo)
}

// EmpujarOperador empuja un operador a la pila de operadores
func EmpujarOperador(operador int){
	pilaOperadores.Push(operador)
}

// GuardarMientrasUbicacion guarda el salto de la operación
// mientras en la pila de saltos (que corresponde el siguiente cuadruplo)
func GuardarMientrasUbicacion(){
	pilaDeSaltos.Push(ContadorActual())
}

// CrearCuadruploMientrasGotof genera el cuadruplo de la operación
// gotof sin el resultado del salto. Guarda el cuadruplo actual en la
// pila de saltos
func CrearCuadruploMientrasGotof(){
	temporal, _ := filaCuadruplos.Back()
	filaCuadruplos.Push(Quadruple{ops.GOTOF, temporal.Resultado, -1, -1})
	pilaDeSaltos.Push(ContadorActual()-1)
}

// ActualizarSino actualiza el resultado del cuadruplo de salto de sino. Actualiza el cuadruplo gotof del si,
// y crea un cuadruplo goto para saltar el sino
func ActualizarSino(){
    falso, _ := pilaDeSaltos.Pop()                        
    filaCuadruplos.Push(Quadruple{ops.GOTO, -1, -1, -1})  
    pilaDeSaltos.Push(ContadorActual() - 1)            
    c := filaCuadruplos.Find(falso)
    c.Resultado = ContadorActual()                        
}

// ActualizarMientras actualiza el cuadruplo gotof
// generado y crea un cuadruplo goto para indicar 
// el salto en el ultimo cuadruplo del while
func ActualizarMientras(){
	gotofWhile, _ := pilaDeSaltos.Pop()
	comienzoWhile, _ := pilaDeSaltos.Pop()
	filaCuadruplos.Push(Quadruple{ops.GOTO, -1, -1, comienzoWhile})
	c := filaCuadruplos.Find(gotofWhile)
	c.Resultado = ContadorActual()
}

// EmpujarSalto crea el cuadruplo de salto. Obtiene el temporal/variable más reciente
// y lo empuja la pila de cuaruplos. El cuadruplo actual se inserta en la pila de saltos
func EmpujarSalto(operador int){
	temporal, _ := filaCuadruplos.Back()
	filaCuadruplos.Push(Quadruple{operador, temporal.Resultado, -1, -1})
	i := ContadorActual() - 1
	pilaDeSaltos.Push(i)
}

// ActualizarSalto actualiza el resultado del salto pendiente con el siguiente
// cuadruplo
func ActualizarSalto(){
	i, _ := pilaDeSaltos.Pop()
	q := filaCuadruplos.Find(i)
	q.Resultado = ContadorActual()
}

// GuardarNombreFuncionActual busca si la funcion existe y la guarda en
// funcionLlamada
func GuardarNombreFuncionActual(nombre string) {
	funcion, existe := semantics.BuscarFuncion(nombre)
	if !existe {
		semantics.ErrorSemantico("Funcion no ha sido declarada")
		funcionLlamada = nil
		return
	}
	funcionLlamada = funcion
}

// GenerarCuadruploAcabarFunc genera el cuadruplo del fin del programa
func GenerarCuadruploAcabarFunc(){
	filaCuadruplos.Push(Quadruple{ops.ENDFUNC, -1, -1, -1})
}

// GenerarCuadruplo es un generador de cuadruplos generico. Saca los operandos izquierdo y derecho 
// (junto a su tipo), el operador, y valida el tipo del resultado con el cubo semantico. Al resultado se le 
// asigna un segmento de memoria. Las cuatro partes se empujan a la fila de cuadruplos, y el resultado se empuja
// a la pila de operandos.
func GenerarCuadruplo(){
	operandoDerecho, _ := pilaOperandos.Pop()
	tipoDerecho, _ := pilaOperandosType.Pop()
	operandoIzquierdo, _ := pilaOperandos.Pop()
	tipoIzquierdo, _ := pilaOperandosType.Pop()
	operador, _ := pilaOperadores.Pop()
	tipoResultado := semantics.ValidarSemantica(tipoIzquierdo, tipoDerecho, operador)
	if tipoResultado == types.TipoError {
		semantics.ErrorSemantico("tipos incompatibles en la operación")
		return
	}
	resultado, err := memory.Asignar(memory.Temporal, tipoResultado)
	if err != nil {
		semantics.ErrorSemantico(err.Error())
		return
	}
	memory.RegistrarNombre(resultado, fmt.Sprintf("t%d", resultado))
	filaCuadruplos.Push(Quadruple{operador, operandoIzquierdo, operandoDerecho, resultado})
	pilaOperandos.Push(resultado)
	pilaOperandosType.Push(tipoResultado)
}

// GenerarCuadruploParametro genera el cuadruplo de parametro al llamar una funcion. Valida
// que el tipo del argumento sea correcto, y que la cantidad de parametros no haya excedido lo
// definido. Se crea un cuadruplo para el parametro y se aumenta el contador
func GenerarCuadruploParametro(){
	if funcionLlamada != nil {
		argumento, _ := pilaOperandos.Pop()
		tipoArgumento, _ := pilaOperandosType.Pop()
		if contadorParametro >= len(funcionLlamada.Parametros) {
			semantics.ErrorSemantico("demasiados argumentos en la llamada")
			return
		}
		parametro := funcionLlamada.Parametros[contadorParametro]
		if tipoArgumento != parametro.Tipo {
			semantics.ErrorSemantico(fmt.Sprintf("tipo del argumento %d no coincide", contadorParametro+1))
			return
		}
		filaCuadruplos.Push(Quadruple{ops.PARAM, argumento, -1, parametro.Direccion})
		contadorParametro++
	}
}

// GenerarCuadruploEra genera el cuadruplo ERA para indicar activacion de memoria
func GenerarCuadruploEra(){
	if funcionLlamada != nil {
		contadorParametro = 0
		filaCuadruplos.Push(Quadruple{ops.ERA, -1, -1, funcionLlamada.DirInicio})
	}
}

// GenerarCuadruploGosub genera el cuadruplo gosub. Verifica la cantidad de parametros pasados
func GenerarCuadruploGosub(){
	if funcionLlamada != nil {
		if contadorParametro != len(funcionLlamada.Parametros) {
			semantics.ErrorSemantico(fmt.Sprintf(
				"faltan argumentos en la llamada a '%s': se esperaban %d, se recibieron %d",
				funcionLlamada.Nombre, len(funcionLlamada.Parametros), contadorParametro))
			return
		}
		filaCuadruplos.Push(Quadruple{ops.GOSUB, -1, -1, funcionLlamada.DirInicio})
	}
}

// GenerarCuadruploEscribe genera el cuadruplo para la operacion de escribir
func GenerarCuadruploEscribe(){
    operando, _ := pilaOperandos.Pop()
    pilaOperandosType.Pop()
    filaCuadruplos.Push(Quadruple{ops.IMPRIME, -1, -1, operando})
}

// GenerarCuadruploRetorno genera el cuadruplo de retorno de una funcion.
// Verifica que el tipo del valor de retorno sea igual al tipo de retorno de
// la funcion, que el retorno se encuentro dentro de una función. Obtiene la direccion 
// de memoria asociada con la función y la usa en el cuadruplo de retorno
func GenerarCuadruploRetorno(){
	operando, _ := pilaOperandos.Pop()
	tipoValor, _ := pilaOperandosType.Pop()
	tipoRetorno, dentroDeFuncion := semantics.TipoRetornoActual()
	if !dentroDeFuncion {
		semantics.ErrorSemantico("retorno fuera de una función")
		return
	}
	if tipoRetorno == types.TipoNula {
		semantics.ErrorSemantico("una función nula no puede retornar un valor")
		return
	}
	if tipoValor != tipoRetorno {
		semantics.ErrorSemantico("el tipo del retorno no coincide con el de la función")
		return
	}
	dirFuncMem := semantics.ObtenerFuncionDireccion()
	filaCuadruplos.Push(Quadruple{ops.RETORNO, operando, -1, dirFuncMem})
}

// GenerarCuadruploAsigna genera el cuadruplo de asignacion de variable. Verifica
// que la variable exista, y checa con el cubo semantico que la asignacion sea correcta. 
// Empuja el cuadruplo a la fila de cuadruplos
func GenerarCuadruploAsigna(nombre string) {
	operando, _ := pilaOperandos.Pop()
	tipoValor, _ := pilaOperandosType.Pop()

	variable, existe := semantics.BuscarVariable(nombre)
	if !existe {
		semantics.ErrorSemantico(fmt.Sprintf("variable '%s' no declarada", nombre))
		return
	}
	if semantics.ValidarSemantica(variable.Tipo, tipoValor, ops.ASIGNAVAR) == types.TipoError {
		semantics.ErrorSemantico(fmt.Sprintf("tipos incompatibles en la asignación a '%s'", nombre))
		return
	}
	filaCuadruplos.Push(Quadruple{ops.ASIGNAVAR, operando, -1, variable.Direccion})
}

// GenerarCuadruploResultadoLlamada genera el cuadruplo del resultado de llamada a una función.
// Busca la variable asociada con el nombre de la función y le asigna una direccion en memoria al
// resultado temporal. Se empuja el temporal a la pila de operandos junto a su tipo.
func GenerarCuadruploResultadoLlamada(){
	if funcionLlamada == nil {
		return
	}
	if funcionLlamada.TipoRetorno == types.TipoNula {
		return
	}
	global, _ := semantics.BuscarVariable(funcionLlamada.Nombre)
	temporal, err := memory.Asignar(memory.Temporal, global.Tipo)
	if err != nil {
        semantics.ErrorSemantico(err.Error())
        return
    }
    memory.RegistrarNombre(temporal, fmt.Sprintf("t%d", temporal))
    filaCuadruplos.Push(Quadruple{ops.ASIGNAVAR, global.Direccion, -1, temporal})
    pilaOperandos.Push(temporal)
    pilaOperandosType.Push(global.Tipo)
}

// GenerarOperandoMenos genera el cuadruplo para las constantes negativas.
// Obtenemos la direccion de la constante cero para simular la expresion
// 0 - [constante]
func GenerarOperandoMenos() {
    operando, _ := pilaOperandos.Pop()
    tipo, _ := pilaOperandosType.Pop()

    ceroDir, _ := semantics.ProcesarConstante("0", tipo)
    temp, err := memory.Asignar(memory.Temporal, tipo)
    if err != nil {
        semantics.ErrorSemantico(err.Error())
        return
    }
    memory.RegistrarNombre(temp, fmt.Sprintf("t%d", temp))
    filaCuadruplos.Push(Quadruple{ops.MENOS, ceroDir, operando, temp})
    pilaOperandos.Push(temp)
    pilaOperandosType.Push(tipo)
}

func ImprimirCuadruplos() {
	fmt.Printf("%-4s %-6s %-8s %-8s %-8s\n", "#", "op", "izq", "der", "res")
	for i, q := range filaCuadruplos.Items {
		fmt.Printf("%-4d %-6s %-8d %-8d %-8d\n",
			i,
			ops.Simbolo(q.Operador),
			q.Izquierda,
			q.Derecha,
			q.Resultado)
	}
}
