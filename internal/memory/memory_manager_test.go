package memory

import (
	"testing"

	"gustavocoutino.compilador/internal/types"
)

func TestNewManager(t *testing.T) {
	m := New()

	if m == nil {
		t.Fatal("New devolvió nil")
	}
	if m.contadores == nil {
		t.Error("contadores = nil; Esperado mapa inicializado")
	}
	if len(m.contadores) != 0 {
		t.Errorf("len(contadores) = %d; Esperado 0", len(m.contadores))
	}
}

func TestMemoryManagerAsignar(t *testing.T) {
	tests := []struct {
		name string
		calls []struct {
			segmento Segmento
			tipo     types.Tipo
		}
		expectedDirs []int
	}{
		{
			name: "secuencial mismo segmento y tipo",
			calls: []struct {
				segmento Segmento
				tipo     types.Tipo
			}{
				{Global, types.TipoEntero},
				{Global, types.TipoEntero},
				{Global, types.TipoEntero},
			},
			expectedDirs: []int{0, 1, 2},
		},
		{
			name: "tipos distintos del mismo segmento",
			calls: []struct {
				segmento Segmento
				tipo     types.Tipo
			}{
				{Global, types.TipoEntero},
				{Global, types.TipoFlotante},
				{Global, types.TipoLiteral},
			},
			expectedDirs: []int{0, 2000, 4000},
		},
		{
			name: "segmentos distintos",
			calls: []struct {
				segmento Segmento
				tipo     types.Tipo
			}{
				{Global, types.TipoEntero},
				{Local, types.TipoEntero},
				{Temporal, types.TipoEntero},
				{Constante, types.TipoEntero},
			},
			expectedDirs: []int{0, 6000, 12000, 18000},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := New()
			for i, c := range tt.calls {
				got, err := m.Asignar(c.segmento, c.tipo)
				if err != nil {
					t.Fatalf("Asignar[%d] err = %v", i, err)
				}
				if got != tt.expectedDirs[i] {
					t.Errorf("Asignar[%d] = %d; Esperado %d", i, got, tt.expectedDirs[i])
				}
			}
		})
	}
}

func TestMemoryManagerAsignarDesborde(t *testing.T) {
	m := New()
	for i := range tamanoBloque {
		if _, err := m.Asignar(Global, types.TipoEntero); err != nil {
			t.Fatalf("Asignar[%d] err inesperado: %v", i, err)
		}
	}

	_, err := m.Asignar(Global, types.TipoEntero)

	if err == nil {
		t.Error("err = nil; Esperado error de memoria agotada")
	}
}

func TestMemoryManagerLiberarMemoria(t *testing.T) {
	m := New()
	m.Asignar(Global, types.TipoEntero)
	m.Asignar(Local, types.TipoEntero)
	m.Asignar(Temporal, types.TipoEntero)
	m.Asignar(Constante, types.TipoEntero)

	m.LiberarMemoria()

	dirLocal, _ := m.Asignar(Local, types.TipoEntero)
	if dirLocal != base(Local, types.TipoEntero) {
		t.Errorf("Local no reseteado: dir = %d; Esperado %d", dirLocal, base(Local, types.TipoEntero))
	}
	dirTemporal, _ := m.Asignar(Temporal, types.TipoEntero)
	if dirTemporal != base(Temporal, types.TipoEntero) {
		t.Errorf("Temporal no reseteado: dir = %d; Esperado %d", dirTemporal, base(Temporal, types.TipoEntero))
	}
	dirGlobal, _ := m.Asignar(Global, types.TipoEntero)
	if dirGlobal == base(Global, types.TipoEntero) {
		t.Error("Global reseteado incorrectamente")
	}
	dirConstante, _ := m.Asignar(Constante, types.TipoEntero)
	if dirConstante == base(Constante, types.TipoEntero) {
		t.Error("Constante reseteado incorrectamente")
	}
}

func TestMemoryManagerContador(t *testing.T) {
	tests := []struct {
		name      string
		calls     []types.Tipo
		consulta  types.Tipo
		expected  int
	}{
		{
			name:     "sin asignaciones",
			calls:    []types.Tipo{},
			consulta: types.TipoEntero,
			expected: 0,
		},
		{
			name:     "varias del mismo tipo",
			calls:    []types.Tipo{types.TipoEntero, types.TipoEntero, types.TipoEntero},
			consulta: types.TipoEntero,
			expected: 3,
		},
		{
			name:     "tipos mezclados",
			calls:    []types.Tipo{types.TipoEntero, types.TipoFlotante, types.TipoEntero},
			consulta: types.TipoFlotante,
			expected: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := New()
			for _, tipo := range tt.calls {
				m.Asignar(Local, tipo)
			}

			got := m.Contador(Local, tt.consulta)

			if got != tt.expected {
				t.Errorf("Contador = %d; Esperado %d", got, tt.expected)
			}
		})
	}
}

func TestPackageAsignar(t *testing.T) {
	gestor = New()

	dir1, err := Asignar(Global, types.TipoEntero)
	if err != nil {
		t.Fatalf("Asignar err = %v", err)
	}
	dir2, _ := Asignar(Global, types.TipoEntero)

	if dir2 != dir1+1 {
		t.Errorf("Asignar consecutivo no incrementa: %d después de %d", dir2, dir1)
	}
}

func TestPackageLiberarMemoria(t *testing.T) {
	gestor = New()
	Asignar(Local, types.TipoEntero)
	Asignar(Local, types.TipoEntero)

	LiberarMemoria()

	dir, _ := Asignar(Local, types.TipoEntero)
	if dir != base(Local, types.TipoEntero) {
		t.Errorf("dir = %d; Esperado %d (segmento reseteado)", dir, base(Local, types.TipoEntero))
	}
}

func TestPackageContador(t *testing.T) {
	gestor = New()
	Asignar(Temporal, types.TipoEntero)
	Asignar(Temporal, types.TipoEntero)

	got := Contador(Temporal, types.TipoEntero)

	if got != 2 {
		t.Errorf("Contador = %d; Esperado 2", got)
	}
}

func TestRegistrarNombre(t *testing.T) {
	nombres = map[int]string{}

	RegistrarNombre(42, "miVariable")

	if got, ok := nombres[42]; !ok || got != "miVariable" {
		t.Errorf("nombres[42] = %q, %v; Esperado %q, true", got, ok, "miVariable")
	}
}

func TestSegmentoNombre(t *testing.T) {
	tests := []struct {
		name      string
		direccion int
		expected  string
	}{
		{"global", 0, "global"},
		{"local", 6000, "local"},
		{"temporal", 12000, "temporal"},
		{"constante", 18000, "constante"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := segmentoNombre(tt.direccion)
			if got != tt.expected {
				t.Errorf("segmentoNombre(%d) = %q; Esperado %q", tt.direccion, got, tt.expected)
			}
		})
	}
}
