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
	operador int
	izquierda int
	derecha int
	resultado int

}

var (
	pilaOperadores = stack.New[int]()
	pilaOperandos = stack.New[int]()
	pilaOperandosType = stack.New[types.Tipo]()
	filaCuadruplos = queue.New[Quadruple]()
	pilaDeSaltos = stack.New[int]()
	contadorParametro = 0
	funcionLlamada *symbols.Funcion
)

func ContadorActual() int {
    return filaCuadruplos.Len()
}

func CrearCuadruploGotoInicio(){
	pilaDeSaltos.Push(ContadorActual())
	filaCuadruplos.Push(Quadruple{ops.GOTO, -1, -1, -1})
}

func CrearCuadruploFin(){
	filaCuadruplos.Push(Quadruple{ops.FIN, -1, -1, -1})
}

func EmpujarOperando(direccion int, tipo types.Tipo) {
	pilaOperandos.Push(direccion)
	pilaOperandosType.Push(tipo)
}

func EmpujarOperador(operador int){
	pilaOperadores.Push(operador)
}

func GuardarMientrasUbicacion(){
	pilaDeSaltos.Push(ContadorActual())
}

func CrearCuadruploMientrasGotof(){
	temporal, _ := filaCuadruplos.Back()
	filaCuadruplos.Push(Quadruple{ops.GOTOF, temporal.resultado, -1, -1})
	pilaDeSaltos.Push(ContadorActual()-1)
}

func ActualizarSino(){
    falso, _ := pilaDeSaltos.Pop()                        
    filaCuadruplos.Push(Quadruple{ops.GOTO, -1, -1, -1})  
    pilaDeSaltos.Push(ContadorActual() - 1)            
    c := filaCuadruplos.Find(falso)
    c.resultado = ContadorActual()                        
}

func ActualizarMientras(){
	gotofWhile, _ := pilaDeSaltos.Pop()
	comienzoWhile, _ := pilaDeSaltos.Pop()
	filaCuadruplos.Push(Quadruple{ops.GOTO, -1, -1, comienzoWhile})
	c := filaCuadruplos.Find(gotofWhile)
	c.resultado = ContadorActual()
}

func EmpujarSalto(operador int){
	temporal, _ := filaCuadruplos.Back()
	filaCuadruplos.Push(Quadruple{operador, temporal.resultado, -1, -1})
	i := ContadorActual() - 1
	pilaDeSaltos.Push(i)
}

func ActualizarSalto(){
	i, _ := pilaDeSaltos.Pop()
	q := filaCuadruplos.Find(i)
	q.resultado = ContadorActual()
}

func GuardarNombreFuncionActual(nombre string) {
	funcion, existe := semantics.BuscarFuncion(nombre)
	if !existe {
		semantics.ErrorSemantico("Funcion no ha sido declarada")
		return
	}
	funcionLlamada = funcion
}


func GenerarCuadruploAcabarFunc(){
	filaCuadruplos.Push(Quadruple{ops.ENDFUNC, -1, -1, -1})
}

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

func GenerarCuadruploEra(){
	if funcionLlamada != nil {
		contadorParametro = 0
		filaCuadruplos.Push(Quadruple{ops.ERA, -1, -1, funcionLlamada.DirInicio})
	}
}

func GenerarCuadruploGosub(){
	if funcionLlamada != nil {
		filaCuadruplos.Push(Quadruple{ops.GOSUB, -1, -1, funcionLlamada.DirInicio})
	}
}

func GenerarCuadruploEscribe(){
    operando, _ := pilaOperandos.Pop()
    pilaOperandosType.Pop()
    filaCuadruplos.Push(Quadruple{ops.IMPRIME, -1, -1, operando})
}

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
	filaCuadruplos.Push(Quadruple{ops.RETORNO, -1, -1, operando})
}

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

func ImprimirCuadruplos() {
	fmt.Printf("%-4s %-6s %-8s %-8s %-8s\n", "#", "op", "izq", "der", "res")
	for i, q := range filaCuadruplos.Items {
		fmt.Printf("%-4d %-6s %-8d %-8d %-8d\n",
			i,
			ops.Simbolo(q.operador),
			q.izquierda,
			q.derecha,
			q.resultado)
	}
}
