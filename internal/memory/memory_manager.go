package memory

import (
	"fmt"
	"sort"

	"gustavocoutino.compilador/internal/types"
)

type Segmento int

const (
	Global Segmento = iota
	Local
	Temporal
	Constante
)

const tamanoBloque = 2000 


var TiposDir = []types.Tipo{types.TipoEntero, types.TipoFlotante, types.TipoLiteral}

// base regresa la direccion de memoria asignada a una variable/temporal/constante
func base(s Segmento, t types.Tipo) int {
	idxTipo := 0
	for i, tt := range TiposDir {
		if tt == t {
			idxTipo = i
		}
	}
	return (int(s)*len(TiposDir) + idxTipo) * tamanoBloque
}


// contadores cuenta cuantas direcciones han sido asignadas por segmento
type MemoryManager struct {
	contadores map[int]int 
}

// New es el constructor del manejador de memoria
func New() *MemoryManager {
	return &MemoryManager{contadores: make(map[int]int)}
}

// Asignar asigna una direccion virtual a una variable/temporal/constante.
// Incrementa el contador de direcciones asignadas y regresa la direccion
func (m *MemoryManager) Asignar(s Segmento, t types.Tipo) (int, error) {
	b := base(s, t)
	if m.contadores[b] >= tamanoBloque {
		return 0, fmt.Errorf("memoria agotada en segmento %d tipo %s", s, t)
	}
	dir := b + m.contadores[b]
	m.contadores[b]++
	return dir, nil
}

// LiberarMemoria libera memoria de un segmento. Se llama
// despues de generar los cuadruplos de una funcion local
func (m *MemoryManager) LiberarMemoria(){
	for _, segmento := range []Segmento{Local, Temporal}{
		for _, tipo := range TiposDir {
            delete(m.contadores, base(segmento, tipo))
        }
	}	
}

// Contador regresa la cantidad de direcciones asignadas a un segmento
func (m *MemoryManager) Contador(s Segmento, t types.Tipo) int {
    return m.contadores[base(s, t)]
}

// El estado global del asignador de memoria
var gestor = New()

// Mapa para la impresion de direcciones de memoria
var nombres = map[int]string{}

// Asignar asigna una direccion de memoria
func Asignar(s Segmento, t types.Tipo) (int, error) {
	return gestor.Asignar(s, t)
}

// LiberarMemoria llama al liberador de memoria
func LiberarMemoria() {
	gestor.LiberarMemoria()
}

// Contador regresa el numero de direcciones asignadas a un segmento
func Contador(s Segmento, t types.Tipo) int {
    return gestor.Contador(s, t)
}

// RegistrarNombre le asigna un nombre a una direccion para motivos de impresion
func RegistrarNombre(direccion int, etiqueta string) {
	nombres[direccion] = etiqueta
}

// segmentoNombre regresa el nombre del segmento para imprimir
// la tabla de direcciones
func segmentoNombre(direccion int) string {
	switch Segmento(direccion/tamanoBloque/len(TiposDir)) {
	case Global:
		return "global"
	case Local:
		return "local"
	case Temporal:
		return "temporal"
	case Constante:
		return "constante"
	default:
		return "?"
	}
}

func ImprimirMemoria() {
	dirs := make([]int, 0, len(nombres))
	for d := range nombres {
		dirs = append(dirs, d)
	}
	sort.Ints(dirs)

	fmt.Println("Asignador de memoria")
	fmt.Printf("%-10s %-10s %-14s\n", "Dirección", "Segmento", "Nombre")
	for _, d := range dirs {
		fmt.Printf("%-10d %-10s %-14s\n", d, segmentoNombre(d), nombres[d])
	}
	fmt.Println()
}