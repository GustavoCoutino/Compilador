package virtualmachine

import (
	"fmt"
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
	pendiente map[int]interface{} //
}

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

func (vm *VM) Ejecutar(){
	ip := 0
	for {
		q := vm.fila.Find(ip)
		switch q.Operador {
		case ops.GOTO:
			ip = q.Resultado
			continue
		case ops.GOTOF:
			if vm.leerEntero(q.Izquierda) != 1 {
				ip = q.Resultado
				continue
			}
		case ops.GOTOT:
			if vm.leerEntero(q.Izquierda) != 0 {
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
			continue
		case ops.ASIGNAVAR:
			vm.write(q.Resultado, vm.leer(q.Izquierda))
		case ops.IGUAL:
			var v int
			if vm.leerEntero(q.Izquierda) == vm.leerEntero(q.Derecha) { 
				v = 1 
			}
			vm.write(q.Resultado, v)
		case ops.DIFERENTE:
			var v int
			if vm.leerEntero(q.Izquierda) != vm.leerEntero(q.Derecha) { 
				v = 1 
			}
			vm.write(q.Resultado, v)
		case ops.MAYOR:
			var v int
			if vm.leerEntero(q.Izquierda) > vm.leerEntero(q.Derecha) { 
				v = 1 
			}
			vm.write(q.Resultado, v)
		case ops.MENOR:
			var v int
			if vm.leerEntero(q.Izquierda) < vm.leerEntero(q.Derecha) { 
				v = 1 
			}
			vm.write(q.Resultado, v)
		case ops.MAS:
			vm.write(q.Resultado, vm.leerEntero(q.Izquierda) + vm.leerEntero(q.Derecha))
		case ops.MENOS:
			vm.write(q.Resultado, vm.leerEntero(q.Izquierda) - vm.leerEntero(q.Derecha))
		case ops.POR:
			vm.write(q.Resultado, vm.leerEntero(q.Izquierda) * vm.leerEntero(q.Derecha))
		case ops.ENTRE:
			vm.write(q.Resultado, vm.leerEntero(q.Izquierda) / vm.leerEntero(q.Derecha))
		case ops.IMPRIME:
			fmt.Println(vm.leer(q.Resultado))
		case ops.RETORNO:
			vm.write(q.Resultado, vm.leer(q.Izquierda))
			vm.pilaAR.Pop()
			retorno, _ := vm.pilaIP.Pop()
			ip = retorno
			continue
		case ops.ENDFUNC:
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

func (vm *VM) leerEntero(dir int) int {
	return vm.leer(dir).(int)
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

func (vm *VM) guardarSiguienteCuadruploDeLlamada(i int){
	vm.pilaIP.Push(i)
}

func (vm *VM) segmentoDe(dir int) memory.Segmento {
	return memory.Segmento(dir / 2000 / len(memory.TiposDir))
}

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

func parsearValor(nombre string, tipo types.Tipo) any {
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