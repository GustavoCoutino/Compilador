package memory

import (
	"fmt"

	"gustavocoutino.compilador/internal/types"
)

type Segmento int

const (
	Global Segmento = iota
	Local
	Temporal
	Constante
)

const tamanoBloque = 1000 

var tiposDir = []types.Tipo{types.TipoEntero, types.TipoFlotante}

func base(s Segmento, t types.Tipo) int {
	idxTipo := 0
	for i, tt := range tiposDir {
		if tt == t {
			idxTipo = i
		}
	}
	return (int(s)*len(tiposDir) + idxTipo) * tamanoBloque
}

type MemoryManager struct {
	contadores map[int]int 
}

func New() *MemoryManager {
	return &MemoryManager{contadores: make(map[int]int)}
}

func (m *MemoryManager) Asignar(s Segmento, t types.Tipo) (int, error) {
	b := base(s, t)
	if m.contadores[b] >= tamanoBloque {
		return 0, fmt.Errorf("memoria agotada en segmento %d tipo %s", s, t)
	}
	dir := b + m.contadores[b]
	m.contadores[b]++
	return dir, nil
}

func (m *MemoryManager) LiberarMemoria(){
	for _, segmento := range []Segmento{Local, Temporal}{
		for _, tipo := range tiposDir {
            delete(m.contadores, base(segmento, tipo))
        }
	}	
}

// Las siguientes lineas de codigo solo existe para parte 
// de las pruebas de imprimir la fila de cuadruplos, 
// no tienen ningun proposito practico
var gestor = New()

var nombres = map[int]string{}

func Asignar(s Segmento, t types.Tipo) (int, error) {
	return gestor.Asignar(s, t)
}

func RegistrarNombre(direccion int, etiqueta string) {
	nombres[direccion] = etiqueta
}

func NombreDe(direccion int) (string, bool) {
	n, ok := nombres[direccion]
	return n, ok
}