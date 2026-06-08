package virtualmachine

import (
	"fmt"
	"sort"
	"strconv"

	"gustavocoutino.compilador/internal/data_structures/queue"
	"gustavocoutino.compilador/internal/data_structures/stack"
	"gustavocoutino.compilador/internal/memory"
	"gustavocoutino.compilador/internal/ops"
	"gustavocoutino.compilador/internal/quadruples"
	"gustavocoutino.compilador/internal/semantics"
	"gustavocoutino.compilador/internal/types"
)

type VM struct {
	fila      *queue.Queue[quadruples.Quadruple] // fila de cuadruplos
    global    map[int]interface{} // mapa de memoria
    pilaAR    *stack.Stack[map[int]interface{}] // pila de registros de activacion
    pilaIP    *stack.Stack[int] // pila para regresar despues de gosub
	pendiente map[int]interface{} // pila para AR en construccion
}

var resultados []interface{}

// NewVM crea la instancia de la maquina virtual. Le asigna la fila de cuadruplos,
// un mapa de memoria vacio, la pila de registros de activacion, y la pila del cuadruplo pendiente
// despues de llamar a gosub. Instancia el mapa de memoria con las constantes globales de la tabla de constantes
func NewVM(fila *queue.Queue[quadruples.Quadruple]) *VM {
    vm := &VM{
        fila:   fila,
        global: make(map[int]interface{}),
        pilaAR: stack.New[map[int]interface{}](),
        pilaIP: stack.New[int](),
    }
	for _, c := range semantics.GetConstantes() {   
        vm.global[c.Direccion] = parsearValor(c.Nombre, c.Tipo)
    }
    vm.pilaAR.Push(make(map[int]interface{}))  
    return vm
}

// Ejecutar inicia la maquina virtual. Itera sobre todos los cuadruplos
// y al encontrase con un operando, realiza una accion especifica (salto, impresion
// asignacion, actualizar una pilia, acabar ejecución)
func (vm *VM) Ejecutar(){
	ip := 0
	for {
		q := vm.fila.Find(ip)
		switch q.Operador {
		case ops.GOTO:
			ip = q.Resultado
			continue
		case ops.GOTOF:
			if vm.leerNumero(q.Izquierda) != 1 {
				ip = q.Resultado
				continue
			}
		case ops.GOTOT:
			if vm.leerNumero(q.Izquierda) != 0 {
				ip = q.Resultado
				continue
			}
		case ops.ERA:
			vm.pendiente = make(map[int]interface{})
		case ops.PARAM:
			vm.pendiente[q.Resultado] = vm.leer(q.Izquierda)
		case ops.GOSUB:
			vm.pilaAR.Push(vm.pendiente)
			vm.pendiente = nil
			vm.pilaIP.Push(ip+1)
			ip = q.Resultado
			fmt.Println("Llamada de funcion")
			vm.ImprimirMapaMemoriaAR()
			continue
		case ops.ASIGNAVAR:
			vm.write(q.Resultado, vm.leer(q.Izquierda))
		case ops.IGUAL:
			var v int
			if vm.leerNumero(q.Izquierda) == vm.leerNumero(q.Derecha) { 
				v = 1 
			}
			vm.write(q.Resultado, v)
		case ops.DIFERENTE:
			var v int
			if vm.leerNumero(q.Izquierda) != vm.leerNumero(q.Derecha) { 
				v = 1 
			}
			vm.write(q.Resultado, v)
		case ops.MAYOR:
			var v int
			if vm.leerNumero(q.Izquierda) > vm.leerNumero(q.Derecha) { 
				v = 1 
			}
			vm.write(q.Resultado, v)
		case ops.MENOR:
			var v int
			if vm.leerNumero(q.Izquierda) < vm.leerNumero(q.Derecha) { 
				v = 1 
			}
			vm.write(q.Resultado, v)
		case ops.MAS:
			vm.write(q.Resultado, vm.aritmetica(q.Izquierda, q.Derecha, '+'))
		case ops.MENOS:
			vm.write(q.Resultado, vm.aritmetica(q.Izquierda, q.Derecha, '-'))
		case ops.POR:
			vm.write(q.Resultado, vm.aritmetica(q.Izquierda, q.Derecha, '*'))
		case ops.ENTRE:
			vm.write(q.Resultado, vm.aritmetica(q.Izquierda, q.Derecha, '/'))
		case ops.IMPRIME:
			resultados = append(resultados, vm.leer(q.Resultado))
		case ops.RETORNO:
			vm.write(q.Resultado, vm.leer(q.Izquierda))
			fmt.Println("Retorno de funcion")
			vm.ImprimirMapaMemoriaAR()
			vm.ImprimirMapaMemoriaGlobal()
			vm.pilaAR.Pop()
			retorno, _ := vm.pilaIP.Pop()
			ip = retorno
			continue
		case ops.ENDFUNC:
			fmt.Println("Final de funcion")
			vm.ImprimirMapaMemoriaAR()
			vm.pilaAR.Pop()
			retorno, _ := vm.pilaIP.Pop()
			ip = retorno
			continue
		case ops.FIN:
			return
		
		}
		ip++
	}
}

// leerNumero lee una dirección como float64, sirva para comparar enteros o flotantes
func (vm *VM) leerNumero(dir int) float64 {
	return toFloat(vm.leer(dir))
}

// toFloat convierte un valor entero o flotante a float64
func toFloat(v interface{}) float64 {
	switch n := v.(type) {
	case int:
		return float64(n)
	case float64:
		return n
	}
	return 0
}

// aritmetica opera dos direcciones; entero si ambos son enteros, flotante si alguno es flotante
func (vm *VM) aritmetica(izq, der int, op byte) interface{} {
	a, b := vm.leer(izq), vm.leer(der)
	if ai, ok := a.(int); ok {
		if bi, ok := b.(int); ok {
			switch op {
			case '+':
				return ai + bi
			case '-':
				return ai - bi
			case '*':
				return ai * bi
			case '/':
				return ai / bi
			}
		}
	}
	x, y := toFloat(a), toFloat(b)
	switch op {
	case '+':
		return x + y
	case '-':
		return x - y
	case '*':
		return x * y
	case '/':
		return x / y
	}
	return nil
}

func (vm *VM) write(dir int, valor interface{}){
	switch vm.segmentoDe(dir){
	case memory.Global, memory.Constante:
		vm.global[dir] = valor
	case memory.Local, memory.Temporal:
		ar, _ := vm.pilaAR.Top()
        ar[dir] = valor
	}
}

// segmentoDe regresa el segmento correspondeinte a una direccion 
func (vm *VM) segmentoDe(dir int) memory.Segmento {
	return memory.Segmento(dir / 2000 / len(memory.TiposDir))
}

// leer regresa el valor de una dirección virtual
func (vm *VM) leer(dir int) interface{} {
	switch vm.segmentoDe(dir){
	case memory.Global, memory.Constante:
		return vm.global[dir]
	case memory.Local, memory.Temporal:
		ar, _ := vm.pilaAR.Top()
        return ar[dir]
	}
	return nil
}

// parsearValor convierte los valores de la tabla de constantes en enteros y
// tipos (ya que estan guardados como string)
func parsearValor(nombre string, tipo types.Tipo) interface {} {
	switch tipo {
	case types.TipoEntero:
		r, _ := strconv.Atoi(nombre)
		return r
	case types.TipoFlotante:
		r, _ := strconv.ParseFloat(nombre, 64)
		return r
	}
	return nombre
}

func (vm *VM) ImprimirMapaMemoriaGlobal(){
	fmt.Println("Memoria global:")
    fmt.Printf("  %-10s %s\n", "Dirección", "Valor")
    dirs := make([]int, 0, len(vm.global))
    for d := range vm.global {
        dirs = append(dirs, d)
    }
    sort.Ints(dirs)
    for _, d := range dirs {
        fmt.Printf("  %-10d %v\n", d, vm.global[d])
    }
}

func (vm *VM) ImprimirMapaMemoriaAR(){
	fmt.Printf("\nPila de registros de activación")
    for i, ar := range vm.pilaAR.Items {
        fmt.Printf("  AR[%d]:\n", i)
        arDirs := make([]int, 0, len(ar))
        for d := range ar {
            arDirs = append(arDirs, d)
        }
        sort.Ints(arDirs)
        for _, d := range arDirs {
            fmt.Printf("    %-10d %v\n", d, ar[d])
        }
    }
    fmt.Println()
}

func (vm *VM) ImprimirResultadosVM(){
	fmt.Println("Resultados de ejecución de programa")
	for _, v := range resultados {
		fmt.Println(v)
	}
}