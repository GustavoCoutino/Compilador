package semantics

import (
	"os"
	"testing"

	"gustavocoutino.compilador/internal/symbols"
	"gustavocoutino.compilador/internal/types"
)

func silenciarStderr(t *testing.T) {
	t.Helper()
	original := os.Stderr
	devnull, err := os.OpenFile(os.DevNull, os.O_WRONLY, 0)
	if err != nil {
		t.Fatalf("no se pudo abrir %s: %v", os.DevNull, err)
	}
	os.Stderr = devnull
	t.Cleanup(func() {
		os.Stderr = original
		devnull.Close()
	})
}

func TestDeclararVariable(t *testing.T) {
	tests := []struct {
		name         string
		nombre       string
		tipo         types.Tipo
		enFuncion    bool
		expectedName string
		expectedTipo types.Tipo
	}{
		{"declarar variable global entera", "x", types.TipoConstante, false, "x", types.TipoConstante},
		{"declarar variable global flotante", "y", types.TipoFlotante, false, "y", types.TipoFlotante},
		{"declarar variable local flotante", "z", types.TipoFlotante, true, "z", types.TipoFlotante},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tablaVariablesGlobal = symbols.NewTablaVariables()
			funcionActual = nil
			if tt.enFuncion {
				funcionActual = &symbols.Funcion{Variables: symbols.NewTablaVariables()}
			}

			DeclararVariable(tt.nombre, tt.tipo)

			tabla := tablaVariablesGlobal
			if funcionActual != nil {
				tabla = funcionActual.Variables
			}
			v, ok := tabla.Variables[tt.nombre]
			if !ok {
				t.Fatalf("variable %q no fue declarada", tt.nombre)
			}
			if v.Nombre != tt.expectedName {
				t.Errorf("Nombre = %q; Esperado %q", v.Nombre, tt.expectedName)
			}
			if v.Tipo != tt.expectedTipo {
				t.Errorf("Tipo = %s; Esperado %s", v.Tipo, tt.expectedTipo)
			}
		})
	}
}

func TestBuscarVariable( t *testing.T){
	tests := []struct {
		name         string
		nombre       string
		tipo 	     types.Tipo
		enFuncion    bool
		expectedName string
	}{
		{"buscar variable global", "x", types.TipoConstante, false, "x"},
		{"buscar variable local", "z", types.TipoFlotante, true, "z"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tablaVariablesGlobal = symbols.NewTablaVariables()
			funcionActual = nil
			if tt.enFuncion {
				funcionActual = &symbols.Funcion{Variables: symbols.NewTablaVariables()}
			}
			DeclararVariable(tt.nombre, tt.tipo)
			variable, ok := BuscarVariable(tt.nombre)
			if !ok {
				t.Fatalf("variable %q no encontrada", tt.nombre)
			}
			if variable.Nombre != tt.expectedName {
				t.Errorf("Nombre = %q; Esperado %q", variable.Nombre, tt.expectedName)
			}
		})
	}
}

func TestIniciarFuncion(t *testing.T) {
	tests := []struct {
		name        string
		nombre      string
		tipoRetorno types.Tipo
		parametros  []*symbols.Variable
	}{
		{"funcion sin parametros", "main", types.TipoNula, nil},
		{"funcion con parametros", "suma", types.TipoConstante, []*symbols.Variable{
			{Nombre: "a", Tipo: types.TipoConstante},
			{Nombre: "b", Tipo: types.TipoConstante},
		}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			directorioFunciones = symbols.NewDirectorioFunciones()
			funcionActual = nil

			IniciarFuncion(tt.nombre, tt.tipoRetorno, tt.parametros)

			if funcionActual == nil {
				t.Fatal("funcionActual no fue establecida")
			}
			if funcionActual.Nombre != tt.nombre {
				t.Errorf("Nombre = %q; Esperado %q", funcionActual.Nombre, tt.nombre)
			}
			if funcionActual.TipoRetorno != tt.tipoRetorno {
				t.Errorf("TipoRetorno = %s; Esperado %s", funcionActual.TipoRetorno, tt.tipoRetorno)
			}
			if len(funcionActual.Parametros) != len(tt.parametros) {
				t.Errorf("len(Parametros) = %d; Esperado %d", len(funcionActual.Parametros), len(tt.parametros))
			}
			if _, ok := directorioFunciones.Funciones[tt.nombre]; !ok {
				t.Errorf("funcion %q no agregada al directorio", tt.nombre)
			}
		})
	}
}

func TestTerminarFuncion(t *testing.T) {
	funcionActual = &symbols.Funcion{Nombre: "f", Variables: symbols.NewTablaVariables()}

	TerminarFuncion()

	if funcionActual != nil {
		t.Errorf("funcionActual = %+v; Esperado nil", funcionActual)
	}
}

func TestAgregarIdActual(t *testing.T) {
	listaIdsActual = nil

	AgregarIdActual("a")
	AgregarIdActual("b")
	AgregarIdActual("c")

	expected := []string{"a", "b", "c"}
	if len(listaIdsActual) != len(expected) {
		t.Fatalf("len(listaIdsActual) = %d; Esperado %d", len(listaIdsActual), len(expected))
	}
	for i, id := range expected {
		if listaIdsActual[i] != id {
			t.Errorf("listaIdsActual[%d] = %q; Esperado %q", i, listaIdsActual[i], id)
		}
	}
}

func TestDeclararIdsActuales(t *testing.T) {
	tests := []struct {
		name      string
		ids       []string
		tipo      types.Tipo
		enFuncion bool
	}{
		{"declarar ids globales", []string{"a", "b", "c"}, types.TipoConstante, false},
		{"declarar ids locales", []string{"x", "y"}, types.TipoFlotante, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tablaVariablesGlobal = symbols.NewTablaVariables()
			funcionActual = nil
			if tt.enFuncion {
				funcionActual = &symbols.Funcion{Variables: symbols.NewTablaVariables()}
			}
			listaIdsActual = append([]string(nil), tt.ids...)

			DeclararIdsActuales(tt.tipo)

			tabla := tablaVariablesGlobal
			if funcionActual != nil {
				tabla = funcionActual.Variables
			}
			for _, id := range tt.ids {
				v, ok := tabla.Variables[id]
				if !ok {
					t.Errorf("variable %q no declarada", id)
					continue
				}
				if v.Tipo != tt.tipo {
					t.Errorf("Tipo de %q = %s; Esperado %s", id, v.Tipo, tt.tipo)
				}
			}
			if listaIdsActual != nil {
				t.Errorf("listaIdsActual = %v; Esperado nil", listaIdsActual)
			}
		})
	}
}

func TestAgregarParametroActual(t *testing.T) {
	listaParametrosActual = nil

	AgregarParametroActual("a", types.TipoConstante)
	AgregarParametroActual("b", types.TipoFlotante)

	if len(listaParametrosActual) != 2 {
		t.Fatalf("len(listaParametrosActual) = %d; Esperado 2", len(listaParametrosActual))
	}
	if listaParametrosActual[0].Nombre != "a" || listaParametrosActual[0].Tipo != types.TipoConstante {
		t.Errorf("listaParametrosActual[0] = %+v; Esperado {a, TipoConstante}", listaParametrosActual[0])
	}
	if listaParametrosActual[1].Nombre != "b" || listaParametrosActual[1].Tipo != types.TipoFlotante {
		t.Errorf("listaParametrosActual[1] = %+v; Esperado {b, TipoFlotante}", listaParametrosActual[1])
	}
}

func TestIniciarFuncionConParametros(t *testing.T) {
	directorioFunciones = symbols.NewDirectorioFunciones()
	funcionActual = nil
	listaParametrosActual = []*symbols.Variable{
		{Nombre: "a", Tipo: types.TipoConstante},
		{Nombre: "b", Tipo: types.TipoFlotante},
	}

	IniciarFuncionConParametros("suma", types.TipoConstante)

	if funcionActual == nil {
		t.Fatal("funcionActual no fue establecida")
	}
	if funcionActual.Nombre != "suma" {
		t.Errorf("Nombre = %q; Esperado %q", funcionActual.Nombre, "suma")
	}
	if len(funcionActual.Parametros) != 2 {
		t.Errorf("len(Parametros) = %d; Esperado 2", len(funcionActual.Parametros))
	}
	if listaParametrosActual != nil {
		t.Errorf("listaParametrosActual = %v; Esperado nil", listaParametrosActual)
	}
}

func TestDeclararVariableRedeclaracion(t *testing.T) {
	silenciarStderr(t)
	tablaVariablesGlobal = symbols.NewTablaVariables()
	funcionActual = nil
	errorContador = 0

	DeclararVariable("x", types.TipoConstante)
	DeclararVariable("x", types.TipoFlotante)

	v, ok := tablaVariablesGlobal.Variables["x"]
	if !ok {
		t.Fatal("variable x no existe")
	}
	if v.Tipo != types.TipoConstante {
		t.Errorf("Tipo = %s; Esperado %s (redeclaración no debe sobrescribir)", v.Tipo, types.TipoConstante)
	}
	if errorContador != 1 {
		t.Errorf("errorContador = %d; Esperado 1 tras redeclaración", errorContador)
	}
}

func TestBuscarVariableEdgeCases(t *testing.T) {
	tests := []struct {
		name         string
		setup        func()
		nombre       string
		expectFound  bool
		expectedTipo types.Tipo
	}{
		{
			name: "variable inexistente",
			setup: func() {
				tablaVariablesGlobal = symbols.NewTablaVariables()
				funcionActual = nil
			},
			nombre:      "inexistente",
			expectFound: false,
		},
		{
			name: "global encontrada desde scope local",
			setup: func() {
				tablaVariablesGlobal = symbols.NewTablaVariables()
				tablaVariablesGlobal.Agregar("x", types.TipoConstante)
				funcionActual = &symbols.Funcion{Variables: symbols.NewTablaVariables()}
			},
			nombre:       "x",
			expectFound:  true,
			expectedTipo: types.TipoConstante,
		},
		{
			name: "local opaca a global",
			setup: func() {
				tablaVariablesGlobal = symbols.NewTablaVariables()
				tablaVariablesGlobal.Agregar("x", types.TipoConstante)
				funcionActual = &symbols.Funcion{Variables: symbols.NewTablaVariables()}
				funcionActual.Variables.Agregar("x", types.TipoFlotante)
			},
			nombre:       "x",
			expectFound:  true,
			expectedTipo: types.TipoFlotante,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			v, ok := BuscarVariable(tt.nombre)
			if ok != tt.expectFound {
				t.Fatalf("ok = %v; Esperado %v", ok, tt.expectFound)
			}
			if tt.expectFound && v.Tipo != tt.expectedTipo {
				t.Errorf("Tipo = %s; Esperado %s", v.Tipo, tt.expectedTipo)
			}
		})
	}
}

func TestIniciarFuncionDuplicada(t *testing.T) {
	silenciarStderr(t)
	directorioFunciones = symbols.NewDirectorioFunciones()
	funcionActual = nil
	errorContador = 0

	IniciarFuncion("f", types.TipoNula, nil)
	primera := funcionActual

	IniciarFuncion("f", types.TipoConstante, nil)

	if funcionActual == primera {
		t.Errorf("funcionActual no fue reseteada tras declaración duplicada")
	}
	if funcionActual != nil {
		t.Errorf("funcionActual = %+v; Esperado nil tras error", funcionActual)
	}
	if got := directorioFunciones.Funciones["f"].TipoRetorno; got != types.TipoNula {
		t.Errorf("TipoRetorno = %s; Esperado %s (redeclaración no debe sobrescribir)", got, types.TipoNula)
	}
	if errorContador != 1 {
		t.Errorf("errorContador = %d; Esperado 1 tras función duplicada", errorContador)
	}
}

func TestTerminarFuncionIdempotente(t *testing.T) {
	funcionActual = nil

	TerminarFuncion()

	if funcionActual != nil {
		t.Errorf("funcionActual = %+v; Esperado nil", funcionActual)
	}
}

func TestDeclararIdsActualesVacio(t *testing.T) {
	tablaVariablesGlobal = symbols.NewTablaVariables()
	funcionActual = nil
	listaIdsActual = nil

	DeclararIdsActuales(types.TipoConstante)

	if len(tablaVariablesGlobal.Variables) != 0 {
		t.Errorf("se declararon %d variables sin ids en la lista", len(tablaVariablesGlobal.Variables))
	}
}

func TestIniciarFuncionConParametrosLimpiaListaSiempre(t *testing.T) {
	silenciarStderr(t)
	directorioFunciones = symbols.NewDirectorioFunciones()
	funcionActual = nil
	errorContador = 0
	directorioFunciones.Agregar("f", types.TipoNula, nil)

	listaParametrosActual = []*symbols.Variable{
		{Nombre: "a", Tipo: types.TipoConstante},
	}

	IniciarFuncionConParametros("f", types.TipoNula)

	if listaParametrosActual != nil {
		t.Errorf("listaParametrosActual = %v; Esperado nil (debe limpiarse aún si IniciarFuncion falla)", listaParametrosActual)
	}
	if errorContador != 1 {
		t.Errorf("errorContador = %d; Esperado 1 tras función duplicada", errorContador)
	}
}

func TestIniciarPrograma(t *testing.T) {
	tests := []struct {
		name              string
		setup             func()
		programaNombre    string
		silenciar         bool
		expectInVars      bool
		expectInFunciones bool
		expectErrores     int
	}{
		{
			name: "programa nuevo se registra en ambas tablas",
			setup: func() {
				tablaVariablesGlobal = symbols.NewTablaVariables()
				directorioFunciones = symbols.NewDirectorioFunciones()
				funcionActual = nil
				errorContador = 0
			},
			programaNombre:    "miPrograma",
			expectInVars:      true,
			expectInFunciones: true,
			expectErrores:     0,
		},
		{
			name: "colisiona con variable global preexistente",
			setup: func() {
				tablaVariablesGlobal = symbols.NewTablaVariables()
				directorioFunciones = symbols.NewDirectorioFunciones()
				funcionActual = nil
				errorContador = 0
				tablaVariablesGlobal.Agregar("foo", types.TipoConstante)
			},
			programaNombre:    "foo",
			silenciar:         true,
			expectInVars:      true,
			expectInFunciones: true,
			expectErrores:     1,
		},
		{
			name: "colisiona con función preexistente",
			setup: func() {
				tablaVariablesGlobal = symbols.NewTablaVariables()
				directorioFunciones = symbols.NewDirectorioFunciones()
				funcionActual = nil
				errorContador = 0
				directorioFunciones.Agregar("bar", types.TipoNula, nil)
			},
			programaNombre:    "bar",
			silenciar:         true,
			expectInVars:      true,
			expectInFunciones: true,
			expectErrores:     1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.silenciar {
				silenciarStderr(t)
			}
			tt.setup()

			IniciarPrograma(tt.programaNombre)

			if _, ok := tablaVariablesGlobal.Variables[tt.programaNombre]; ok != tt.expectInVars {
				t.Errorf("presente en tablaVariablesGlobal = %v; Esperado %v", ok, tt.expectInVars)
			}
			if _, ok := directorioFunciones.Funciones[tt.programaNombre]; ok != tt.expectInFunciones {
				t.Errorf("presente en directorioFunciones = %v; Esperado %v", ok, tt.expectInFunciones)
			}
			if funcionActual != nil {
				t.Errorf("funcionActual = %+v; Esperado nil tras IniciarPrograma", funcionActual)
			}
			if errorContador != tt.expectErrores {
				t.Errorf("errorContador = %d; Esperado %d", errorContador, tt.expectErrores)
			}
		})
	}
}

func TestIniciarProgramaNoEntraEnScope(t *testing.T) {
	tablaVariablesGlobal = symbols.NewTablaVariables()
	directorioFunciones = symbols.NewDirectorioFunciones()
	funcionActual = nil
	errorContador = 0

	IniciarPrograma("mainProg")
	DeclararVariable("x", types.TipoConstante)

	if _, ok := tablaVariablesGlobal.Variables["x"]; !ok {
		t.Error("variable x debió declararse en tablaVariablesGlobal después de IniciarPrograma")
	}
	if funcionActual != nil {
		t.Errorf("funcionActual = %+v; Esperado nil tras IniciarPrograma", funcionActual)
	}
}

func TestHasError(t *testing.T) {
	silenciarStderr(t)
	errorContador = 0

	if HasError() {
		t.Error("HasError() = true; Esperado false sin errores")
	}

	ErrorSemantico("error de prueba")

	if !HasError() {
		t.Error("HasError() = false; Esperado true tras ErrorSemantico")
	}
}

func TestErrorSemanticoContador(t *testing.T) {
	silenciarStderr(t)
	errorContador = 0

	ErrorSemantico("primero")
	ErrorSemantico("segundo")
	ErrorSemantico("tercero")

	if errorContador != 3 {
		t.Errorf("errorContador = %d; Esperado 3", errorContador)
	}
}