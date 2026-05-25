package types

import (
	"testing"
)

const (
	entero = "entero"
	flotante = "flotante"
	nula = "nula"
	errorString = "error"
)

func TestString(t *testing.T){
	tests := []struct {
		name string
		tipo Tipo
		expectedType string
	}{
		{
			name:      "regresar tipo constante",
			tipo:      TipoEntero,
			expectedType: entero,
		},
		{
			name:      "regresar tipo flotante",
			tipo:      TipoFlotante,
			expectedType: flotante,
		},
		{
			name:      "regresar tipo nula",
			tipo:      TipoNula,
			expectedType: nula,
		},
		{
			name:      "regresar tipo error",
			tipo:      TipoError,
			expectedType: errorString,
		},
		
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.tipo.String() != tt.expectedType {
				t.Errorf("tipo incorrecto: %q != %q", tt.tipo.String(), tt.expectedType)
			}
		})
	}
}