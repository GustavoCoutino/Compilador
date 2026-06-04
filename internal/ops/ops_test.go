package ops

import (
	"testing"
)

func TestSimbolo(t *testing.T){
	tests := []struct {
		name      string
		op int
		expected  string
	}{
		{"mas", MAS, "+"},
		{"menos", MENOS, "-"},
		{"por", POR, "*"},
		{"entre", ENTRE, "/"},
		{"mayor", MAYOR, ">"},
		{"menor", MENOR, "<"},
		{"igual", IGUAL, "=="},
		{"diferente", DIFERENTE, "!="},
		{"asignavar", ASIGNAVAR, "="},
		{"imprime", IMPRIME, "escribe"},
		{"retorno", RETORNO, "retorno"},
		{"si", SI, "si"},
		{"sino", SINO, "sino"},
		{"mientras", MIENTRAS, "mientras"},
		{"goto", GOTO, "goto"},
		{"gotof", GOTOF, "gotof"},
		{"gotot", GOTOT, "gotot"},
		{"endfunc", ENDFUNC, "endfunc"},
		{"era", ERA, "era"},
		{"param", PARAM, "param"},
		{"gosub", GOSUB, "gosub"},
		{"inicio", INICIO, "inicio"},
		{"fin", FIN, "fin"},
		{"desconocido", 9999, "?"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Simbolo(tt.op)
			if got != tt.expected {
				t.Errorf("Deberia regresar %s", tt.expected)
			}
		})
	}
}