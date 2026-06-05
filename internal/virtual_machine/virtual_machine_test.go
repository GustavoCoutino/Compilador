package virtualmachine

import (
	"testing"

	"gustavocoutino.compilador/internal/data_structures/queue"
	"gustavocoutino.compilador/internal/data_structures/stack"
	"gustavocoutino.compilador/internal/memory"
	"gustavocoutino.compilador/internal/ops"
	"gustavocoutino.compilador/internal/quadruples"
	"gustavocoutino.compilador/internal/types"
)

func nuevaVMVacia() *VM {
	vm := &VM{
		fila:   queue.New[quadruples.Quadruple](),
		global: make(map[int]any),
		pilaAR: stack.New[map[int]any](),
		pilaIP: stack.New[int](),
	}
	vm.pilaAR.Push(make(map[int]any))
	return vm
}

func TestSegmentoDe(t *testing.T) {
	tests := []struct {
		name      string
		direccion int
		expected  memory.Segmento
	}{
		{"global entero", 0, memory.Global},
		{"global flotante", 2000, memory.Global},
		{"local entero", 6000, memory.Local},
		{"temporal entero", 12000, memory.Temporal},
		{"constante entero", 18000, memory.Constante},
		{"constante literal", 22000, memory.Constante},
	}

	vm := nuevaVMVacia()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := vm.segmentoDe(tt.direccion)
			if got != tt.expected {
				t.Errorf("segmentoDe(%d) = %d; Esperado %d", tt.direccion, got, tt.expected)
			}
		})
	}
}

func TestParsearValor(t *testing.T) {
	tests := []struct {
		name     string
		nombre   string
		tipo     types.Tipo
		expected any
	}{
		{"entero", "42", types.TipoEntero, 42},
		{"flotante", "3.14", types.TipoFlotante, 3.14},
		{"literal", "hola", types.TipoLiteral, "hola"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parsearValor(tt.nombre, tt.tipo)
			if got != tt.expected {
				t.Errorf("parsearValor(%q, %s) = %v; Esperado %v", tt.nombre, tt.tipo, got, tt.expected)
			}
		})
	}
}

func TestWriteYLeerGlobal(t *testing.T) {
	vm := nuevaVMVacia()

	vm.write(0, 42)

	if got := vm.global[0]; got != 42 {
		t.Errorf("vm.global[0] = %v; Esperado 42", got)
	}
	if got := vm.leer(0); got != 42 {
		t.Errorf("vm.leer(0) = %v; Esperado 42", got)
	}
}

func TestWriteYLeerLocal(t *testing.T) {
	vm := nuevaVMVacia()
	dirLocal := 6000

	vm.write(dirLocal, 99)

	ar, _ := vm.pilaAR.Top()
	if got := ar[dirLocal]; got != 99 {
		t.Errorf("ar[%d] = %v; Esperado 99", dirLocal, got)
	}
	if got := vm.leer(dirLocal); got != 99 {
		t.Errorf("vm.leer(%d) = %v; Esperado 99", dirLocal, got)
	}
}

func TestLeerNumero(t *testing.T) {
	vm := nuevaVMVacia()
	vm.global[0] = 7

	got := vm.leerNumero(0)

	if got != 7 {
		t.Errorf("leerNumero(0) = %v; Esperado 7", got)
	}
}

func TestEjecutarMas(t *testing.T) {
	vm := nuevaVMVacia()
	vm.global[18000] = 5
	vm.global[18001] = 3
	vm.fila.Push(quadruples.Quadruple{Operador: ops.MAS, Izquierda: 18000, Derecha: 18001, Resultado: 12000})
	vm.fila.Push(quadruples.Quadruple{Operador: ops.FIN, Izquierda: -1, Derecha: -1, Resultado: -1})

	vm.Ejecutar()

	ar, _ := vm.pilaAR.Top()
	if got := ar[12000]; got != 8 {
		t.Errorf("resultado en temporal = %v; Esperado 8", got)
	}
}

func TestEjecutarAsignacion(t *testing.T) {
	vm := nuevaVMVacia()
	vm.global[18000] = 42
	vm.fila.Push(quadruples.Quadruple{Operador: ops.ASIGNAVAR, Izquierda: 18000, Derecha: -1, Resultado: 0})
	vm.fila.Push(quadruples.Quadruple{Operador: ops.FIN, Izquierda: -1, Derecha: -1, Resultado: -1})

	vm.Ejecutar()

	if got := vm.global[0]; got != 42 {
		t.Errorf("vm.global[0] = %v; Esperado 42", got)
	}
}

func TestEjecutarGoto(t *testing.T) {
	vm := nuevaVMVacia()
	vm.global[18000] = 1
	vm.global[18001] = 2
	vm.fila.Push(quadruples.Quadruple{Operador: ops.GOTO, Izquierda: -1, Derecha: -1, Resultado: 2})
	vm.fila.Push(quadruples.Quadruple{Operador: ops.ASIGNAVAR, Izquierda: 18000, Derecha: -1, Resultado: 0})
	vm.fila.Push(quadruples.Quadruple{Operador: ops.ASIGNAVAR, Izquierda: 18001, Derecha: -1, Resultado: 0})
	vm.fila.Push(quadruples.Quadruple{Operador: ops.FIN, Izquierda: -1, Derecha: -1, Resultado: -1})

	vm.Ejecutar()

	if got := vm.global[0]; got != 2 {
		t.Errorf("vm.global[0] = %v; Esperado 2 (debe brincarse el primer ASIGNAVAR)", got)
	}
}

func TestEjecutarGotof(t *testing.T) {
	tests := []struct {
		name      string
		condicion int
		expected  int
	}{
		{"condicion verdadera ejecuta cuerpo", 1, 100},
		{"condicion falsa brinca cuerpo", 0, 200},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vm := nuevaVMVacia()
			vm.global[18000] = tt.condicion
			vm.global[18001] = 100
			vm.global[18002] = 200
			vm.fila.Push(quadruples.Quadruple{Operador: ops.GOTOF, Izquierda: 18000, Derecha: -1, Resultado: 3})
			vm.fila.Push(quadruples.Quadruple{Operador: ops.ASIGNAVAR, Izquierda: 18001, Derecha: -1, Resultado: 0})
			vm.fila.Push(quadruples.Quadruple{Operador: ops.GOTO, Izquierda: -1, Derecha: -1, Resultado: 4})
			vm.fila.Push(quadruples.Quadruple{Operador: ops.ASIGNAVAR, Izquierda: 18002, Derecha: -1, Resultado: 0})
			vm.fila.Push(quadruples.Quadruple{Operador: ops.FIN, Izquierda: -1, Derecha: -1, Resultado: -1})

			vm.Ejecutar()

			if got := vm.global[0]; got != tt.expected {
				t.Errorf("vm.global[0] = %v; Esperado %d", got, tt.expected)
			}
		})
	}
}

func TestEjecutarOperacionesAritmeticas(t *testing.T) {
	tests := []struct {
		name     string
		operador int
		expected int
	}{
		{"suma", ops.MAS, 13},
		{"resta", ops.MENOS, 7},
		{"multiplicacion", ops.POR, 30},
		{"division", ops.ENTRE, 3},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vm := nuevaVMVacia()
			vm.global[18000] = 10
			vm.global[18001] = 3
			vm.fila.Push(quadruples.Quadruple{Operador: tt.operador, Izquierda: 18000, Derecha: 18001, Resultado: 12000})
			vm.fila.Push(quadruples.Quadruple{Operador: ops.FIN, Izquierda: -1, Derecha: -1, Resultado: -1})

			vm.Ejecutar()

			ar, _ := vm.pilaAR.Top()
			if got := ar[12000]; got != tt.expected {
				t.Errorf("resultado = %v; Esperado %d", got, tt.expected)
			}
		})
	}
}

func TestEjecutarRelacionales(t *testing.T) {
	tests := []struct {
		name     string
		operador int
		izq      int
		der      int
		expected int
	}{
		{"mayor verdadero", ops.MAYOR, 10, 3, 1},
		{"mayor falso", ops.MAYOR, 3, 10, 0},
		{"menor verdadero", ops.MENOR, 3, 10, 1},
		{"igual verdadero", ops.IGUAL, 5, 5, 1},
		{"igual falso", ops.IGUAL, 5, 6, 0},
		{"diferente verdadero", ops.DIFERENTE, 5, 6, 1},
		{"diferente falso", ops.DIFERENTE, 5, 5, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vm := nuevaVMVacia()
			vm.global[18000] = tt.izq
			vm.global[18001] = tt.der
			vm.fila.Push(quadruples.Quadruple{Operador: tt.operador, Izquierda: 18000, Derecha: 18001, Resultado: 12000})
			vm.fila.Push(quadruples.Quadruple{Operador: ops.FIN, Izquierda: -1, Derecha: -1, Resultado: -1})

			vm.Ejecutar()

			ar, _ := vm.pilaAR.Top()
			if got := ar[12000]; got != tt.expected {
				t.Errorf("resultado = %v; Esperado %d", got, tt.expected)
			}
		})
	}
}
