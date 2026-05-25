package quadruples

import (
	"fmt"

	"gustavocoutino.compilador/internal/data_structures/queue"
	"gustavocoutino.compilador/internal/data_structures/stack"
	"gustavocoutino.compilador/internal/memory"
	"gustavocoutino.compilador/internal/ops"
	"gustavocoutino.compilador/internal/semantics"
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
)


func operandoStr(direccion int) string {
	if direccion < 0 {
		return "_"
	}
	if etiqueta, ok := memory.NombreDe(direccion); ok {
		return etiqueta
	}
	return fmt.Sprintf("%d", direccion)
}

func EmpujarOperando(direccion int, tipo types.Tipo) {
	pilaOperandos.Push(direccion)
	pilaOperandosType.Push(tipo)
}

func EmpujarOperador(operador int){
	pilaOperadores.Push(operador)
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

func GenerarEscribeCuadruplo(){
	operando, _ := pilaOperandos.Pop()
	pilaOperandosType.Pop()
	filaCuadruplos.Push(Quadruple{ops.IMPRIME, -1, -1, operando})
}

func GenerarRetornoCuadruplo(){
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

func GenerarAsignaCuadruplo(nombre string) {
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
		fmt.Printf("%-4d %-6s %-8s %-8s %-8s\n",
			i,
			ops.Simbolo(q.operador),
			operandoStr(q.izquierda),
			operandoStr(q.derecha),
			operandoStr(q.resultado))
	}
}
