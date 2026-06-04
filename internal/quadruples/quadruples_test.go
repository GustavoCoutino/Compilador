package quadruples

import (
	"testing"

	"gustavocoutino.compilador/internal/data_structures/queue"
	"gustavocoutino.compilador/internal/data_structures/stack"
	"gustavocoutino.compilador/internal/ops"
	"gustavocoutino.compilador/internal/types"
)

func resetState() {
	pilaOperadores = stack.New[int]()
	pilaOperandos = stack.New[int]()
	pilaOperandosType = stack.New[types.Tipo]()
	filaCuadruplos = queue.New[Quadruple]()
	pilaDeSaltos = stack.New[int]()
	contadorParametro = 0
	funcionLlamada = nil
}

func TestContadorActual(t *testing.T) {
	resetState()

	if got := ContadorActual(); got != 0 {
		t.Errorf("ContadorActual inicial = %d; Esperado 0", got)
	}

	filaCuadruplos.Push(Quadruple{ops.MAS, 0, 1, 2})
	filaCuadruplos.Push(Quadruple{ops.MENOS, 0, 1, 2})

	if got := ContadorActual(); got != 2 {
		t.Errorf("ContadorActual tras 2 pushes = %d; Esperado 2", got)
	}
}

func TestGetFilaCuadruplos(t *testing.T) {
	resetState()

	fila := GetFilaCuadruplos()

	if fila == nil {
		t.Fatal("GetFilaCuadruplos regresó nil")
	}
	if fila != filaCuadruplos {
		t.Error("GetFilaCuadruplos no regresa la fila global")
	}
}

func TestEmpujarOperando(t *testing.T) {
	resetState()

	EmpujarOperando(100, types.TipoEntero)

	if got, _ := pilaOperandos.Top(); got != 100 {
		t.Errorf("tope pilaOperandos = %d; Esperado 100", got)
	}
	if got, _ := pilaOperandosType.Top(); got != types.TipoEntero {
		t.Errorf("tope pilaOperandosType = %s; Esperado %s", got, types.TipoEntero)
	}
}

func TestEmpujarOperador(t *testing.T) {
	resetState()

	EmpujarOperador(ops.MAS)

	if got, _ := pilaOperadores.Top(); got != ops.MAS {
		t.Errorf("tope pilaOperadores = %d; Esperado %d", got, ops.MAS)
	}
}

func TestCrearCuadruploFin(t *testing.T) {
	resetState()

	CrearCuadruploFin()

	if got := ContadorActual(); got != 1 {
		t.Fatalf("ContadorActual = %d; Esperado 1", got)
	}
	q, _ := filaCuadruplos.Back()
	if q.Operador != ops.FIN {
		t.Errorf("Operador = %d; Esperado %d (FIN)", q.Operador, ops.FIN)
	}
}

func TestCrearCuadruploGotoInicio(t *testing.T) {
	resetState()

	CrearCuadruploGotoInicio()

	q, _ := filaCuadruplos.Back()
	if q.Operador != ops.GOTO {
		t.Errorf("Operador = %d; Esperado %d (GOTO)", q.Operador, ops.GOTO)
	}
	if q.Resultado != -1 {
		t.Errorf("Resultado = %d; Esperado -1 (pendiente)", q.Resultado)
	}
	indice, _ := pilaDeSaltos.Top()
	if indice != 0 {
		t.Errorf("tope pilaDeSaltos = %d; Esperado 0", indice)
	}
}

func TestGuardarMientrasUbicacion(t *testing.T) {
	resetState()
	filaCuadruplos.Push(Quadruple{ops.MAS, 0, 1, 2})

	GuardarMientrasUbicacion()

	if got, _ := pilaDeSaltos.Top(); got != 1 {
		t.Errorf("tope pilaDeSaltos = %d; Esperado 1", got)
	}
}

func TestActualizarSalto(t *testing.T) {
	resetState()
	filaCuadruplos.Push(Quadruple{ops.GOTOF, 0, -1, -1})
	pilaDeSaltos.Push(0)
	filaCuadruplos.Push(Quadruple{ops.MAS, 1, 2, 3})

	ActualizarSalto()

	q := filaCuadruplos.Find(0)
	if q.Resultado != 2 {
		t.Errorf("Resultado parchado = %d; Esperado 2", q.Resultado)
	}
}

func TestEmpujarSalto(t *testing.T) {
	resetState()
	filaCuadruplos.Push(Quadruple{ops.MAYOR, 10, 20, 4000})

	EmpujarSalto(ops.GOTOF)

	q, _ := filaCuadruplos.Back()
	if q.Operador != ops.GOTOF {
		t.Errorf("Operador = %d; Esperado %d (GOTOF)", q.Operador, ops.GOTOF)
	}
	if q.Izquierda != 4000 {
		t.Errorf("Izquierda = %d; Esperado 4000 (resultado del cuadruplo anterior)", q.Izquierda)
	}
	indice, _ := pilaDeSaltos.Top()
	if indice != 1 {
		t.Errorf("tope pilaDeSaltos = %d; Esperado 1 (índice del GOTOF)", indice)
	}
}

func TestActualizarMientras(t *testing.T) {
	resetState()
	pilaDeSaltos.Push(0)
	filaCuadruplos.Push(Quadruple{ops.MAYOR, 10, 20, 4000})
	filaCuadruplos.Push(Quadruple{ops.GOTOF, 4000, -1, -1})
	pilaDeSaltos.Push(1)
	filaCuadruplos.Push(Quadruple{ops.IMPRIME, -1, -1, 4000})

	ActualizarMientras()

	gotoQuad, _ := filaCuadruplos.Back()
	if gotoQuad.Operador != ops.GOTO {
		t.Errorf("último cuadruplo = %d; Esperado %d (GOTO)", gotoQuad.Operador, ops.GOTO)
	}
	if gotoQuad.Resultado != 0 {
		t.Errorf("destino del GOTO = %d; Esperado 0 (inicio del ciclo)", gotoQuad.Resultado)
	}
	gotofQuad := filaCuadruplos.Find(1)
	if gotofQuad.Resultado != 4 {
		t.Errorf("destino del GOTOF parchado = %d; Esperado 4 (después del GOTO)", gotofQuad.Resultado)
	}
}

func TestCrearCuadruploMientrasGotof(t *testing.T) {
	resetState()
	filaCuadruplos.Push(Quadruple{ops.MAYOR, 10, 20, 4000})

	CrearCuadruploMientrasGotof()

	q, _ := filaCuadruplos.Back()
	if q.Operador != ops.GOTOF {
		t.Errorf("Operador = %d; Esperado %d (GOTOF)", q.Operador, ops.GOTOF)
	}
	if q.Izquierda != 4000 {
		t.Errorf("Izquierda = %d; Esperado 4000", q.Izquierda)
	}
	indice, _ := pilaDeSaltos.Top()
	if indice != 1 {
		t.Errorf("tope pilaDeSaltos = %d; Esperado 1", indice)
	}
}

func TestGenerarCuadruploAcabarFunc(t *testing.T) {
	resetState()

	GenerarCuadruploAcabarFunc()

	q, _ := filaCuadruplos.Back()
	if q.Operador != ops.ENDFUNC {
		t.Errorf("Operador = %d; Esperado %d (ENDFUNC)", q.Operador, ops.ENDFUNC)
	}
}
