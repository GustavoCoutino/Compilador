package symbols

import (
	"testing"

	"gustavocoutino.compilador/internal/types"
)

func TestNewDirectorioFunciones(t *testing.T) {
	dir := NewDirectorioFunciones()

	if dir == nil {
		t.Fatal("NewDirectorioFunciones devolvió nil")
	}
	if dir.Funciones == nil {
		t.Error("Funciones = nil; Esperado mapa inicializado")
	}
	if len(dir.Funciones) != 0 {
		t.Errorf("len(Funciones) = %d; Esperado 0", len(dir.Funciones))
	}
}

func TestDirectorioFuncionesAgregar(t *testing.T) {
	tests := []struct {
		name        string
		preexisting []string
		nombre      string
		tipoRetorno types.Tipo
		parametros  []*Variable
		expectErr   bool
	}{
		{
			name:        "agregar funcion sin parametros",
			nombre:      "main",
			tipoRetorno: types.TipoNula,
			parametros:  nil,
		},
		{
			name:        "agregar funcion con parametros",
			nombre:      "suma",
			tipoRetorno: types.TipoConstante,
			parametros: []*Variable{
				{Nombre: "a", Tipo: types.TipoConstante},
				{Nombre: "b", Tipo: types.TipoConstante},
			},
		},
		{
			name:        "agregar funcion duplicada",
			preexisting: []string{"f"},
			nombre:      "f",
			tipoRetorno: types.TipoConstante,
			expectErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := NewDirectorioFunciones()
			for _, nombre := range tt.preexisting {
				dir.Agregar(nombre, types.TipoNula, nil)
			}

			f, err := dir.Agregar(tt.nombre, tt.tipoRetorno, tt.parametros)

			if tt.expectErr {
				if err == nil {
					t.Fatal("err = nil; Esperado error")
				}
				if f != nil {
					t.Errorf("funcion = %+v; Esperado nil cuando hay error", f)
				}
				return
			}
			if err != nil {
				t.Fatalf("err = %v; Esperado nil", err)
			}
			if f == nil {
				t.Fatal("funcion = nil; Esperado no nil")
			}
			if f.Nombre != tt.nombre {
				t.Errorf("Nombre = %q; Esperado %q", f.Nombre, tt.nombre)
			}
			if f.TipoRetorno != tt.tipoRetorno {
				t.Errorf("TipoRetorno = %s; Esperado %s", f.TipoRetorno, tt.tipoRetorno)
			}
			if len(f.Parametros) != len(tt.parametros) {
				t.Errorf("len(Parametros) = %d; Esperado %d", len(f.Parametros), len(tt.parametros))
			}
			if f.Variables == nil {
				t.Error("Variables = nil; Esperado tabla inicializada")
			}
			if dir.Funciones[tt.nombre] != f {
				t.Errorf("directorio no apunta a la funcion devuelta")
			}
		})
	}
}

func TestDirectorioFuncionesAgregarNoSobrescribe(t *testing.T) {
	dir := NewDirectorioFunciones()
	dir.Agregar("f", types.TipoNula, nil)
	original := dir.Funciones["f"]

	dir.Agregar("f", types.TipoConstante, nil)

	if dir.Funciones["f"] != original {
		t.Error("funcion fue sobrescrita en duplicado")
	}
	if dir.Funciones["f"].TipoRetorno != types.TipoNula {
		t.Errorf("TipoRetorno = %s; Esperado %s", dir.Funciones["f"].TipoRetorno, types.TipoNula)
	}
}

func TestDirectorioFuncionesBuscar(t *testing.T) {
	tests := []struct {
		name        string
		seed        []string
		nombre      string
		expectFound bool
	}{
		{
			name:        "funcion existente",
			seed:        []string{"f", "g"},
			nombre:      "f",
			expectFound: true,
		},
		{
			name:        "funcion inexistente",
			seed:        []string{"f"},
			nombre:      "g",
			expectFound: false,
		},
		{
			name:        "directorio vacío",
			nombre:      "f",
			expectFound: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := NewDirectorioFunciones()
			for _, nombre := range tt.seed {
				dir.Agregar(nombre, types.TipoNula, nil)
			}

			f, ok := dir.Buscar(tt.nombre)

			if ok != tt.expectFound {
				t.Fatalf("ok = %v; Esperado %v", ok, tt.expectFound)
			}
			if !tt.expectFound {
				if f != nil {
					t.Errorf("funcion = %+v; Esperado nil", f)
				}
				return
			}
			if f.Nombre != tt.nombre {
				t.Errorf("Nombre = %q; Esperado %q", f.Nombre, tt.nombre)
			}
		})
	}
}
