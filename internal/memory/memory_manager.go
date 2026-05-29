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


var tiposDir = []types.Tipo{types.TipoEntero, types.TipoFlotante, types.TipoLiteral}

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

func (m *MemoryManager) Contador(s Segmento, t types.Tipo) int {
    return m.contadores[base(s, t)]
}

var gestor = New()

var nombres = map[int]string{}

func Asignar(s Segmento, t types.Tipo) (int, error) {
	return gestor.Asignar(s, t)
}

func LiberarMemoria() {
	gestor.LiberarMemoria()
}

func Contador(s Segmento, t types.Tipo) int {
    return gestor.Contador(s, t)
}

func RegistrarNombre(direccion int, etiqueta string) {
	nombres[direccion] = etiqueta
}

func NombreDe(direccion int) (string, bool) {
	n, ok := nombres[direccion]
	return n, ok
}

func segmentoNombre(direccion int) string {
	switch Segmento(direccion/tamanoBloque/len(tiposDir)) {
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