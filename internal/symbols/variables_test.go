package symbols

import (
	"testing"

	"gustavocoutino.compilador/internal/types"
)

func TestNewTablaVariables(t *testing.T) {
	tabla := NewTablaVariables()

	if tabla == nil {
		t.Fatal("NewTablaVariables devolvió nil")
	}
	if tabla.Variables == nil {
		t.Error("Variables = nil; Esperado mapa inicializado")
	}
	if len(tabla.Variables) != 0 {
		t.Errorf("len(Variables) = %d; Esperado 0", len(tabla.Variables))
	}
}

func TestTablaVariablesAgregar(t *testing.T) {
	tests := []struct {
		name        string
		preexisting []struct {
			nombre string
			tipo   types.Tipo
		}
		nombre    string
		tipo      types.Tipo
		expectErr bool
	}{
		{
			name:      "agregar nueva variable",
			nombre:    "x",
			tipo:      types.TipoConstante,
			expectErr: false,
		},
		{
			name: "agregar variable duplicada",
			preexisting: []struct {
				nombre string
				tipo   types.Tipo
			}{{"x", types.TipoConstante}},
			nombre:    "x",
			tipo:      types.TipoFlotante,
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tabla := NewTablaVariables()
			for _, pre := range tt.preexisting {
				tabla.Agregar(pre.nombre, pre.tipo)
			}

			err := tabla.Agregar(tt.nombre, tt.tipo)

			if tt.expectErr {
				if err == nil {
					t.Fatal("err = nil; Esperado error")
				}
				return
			}
			if err != nil {
				t.Fatalf("err = %v; Esperado nil", err)
			}
			v, ok := tabla.Variables[tt.nombre]
			if !ok {
				t.Fatalf("variable %q no agregada", tt.nombre)
			}
			if v.Nombre != tt.nombre {
				t.Errorf("Nombre = %q; Esperado %q", v.Nombre, tt.nombre)
			}
			if v.Tipo != tt.tipo {
				t.Errorf("Tipo = %s; Esperado %s", v.Tipo, tt.tipo)
			}
		})
	}
}

func TestTablaVariablesAgregarNoSobrescribe(t *testing.T) {
	tabla := NewTablaVariables()
	tabla.Agregar("x", types.TipoConstante)

	tabla.Agregar("x", types.TipoFlotante)

	if got := tabla.Variables["x"].Tipo; got != types.TipoConstante {
		t.Errorf("Tipo = %s; Esperado %s (duplicado no debe sobrescribir)", got, types.TipoConstante)
	}
}

func TestTablaVariablesBuscar(t *testing.T) {
	tests := []struct {
		name         string
		seed         map[string]types.Tipo
		nombre       string
		expectFound  bool
		expectedTipo types.Tipo
	}{
		{
			name:         "variable existente",
			seed:         map[string]types.Tipo{"x": types.TipoConstante},
			nombre:       "x",
			expectFound:  true,
			expectedTipo: types.TipoConstante,
		},
		{
			name:        "variable inexistente",
			seed:        map[string]types.Tipo{"x": types.TipoConstante},
			nombre:      "y",
			expectFound: false,
		},
		{
			name:        "tabla vacía",
			nombre:      "x",
			expectFound: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tabla := NewTablaVariables()
			for nombre, tipo := range tt.seed {
				tabla.Agregar(nombre, tipo)
			}

			v, ok := tabla.Buscar(tt.nombre)

			if ok != tt.expectFound {
				t.Fatalf("ok = %v; Esperado %v", ok, tt.expectFound)
			}
			if !tt.expectFound {
				if v != nil {
					t.Errorf("variable = %+v; Esperado nil", v)
				}
				return
			}
			if v.Tipo != tt.expectedTipo {
				t.Errorf("Tipo = %s; Esperado %s", v.Tipo, tt.expectedTipo)
			}
		})
	}
}
