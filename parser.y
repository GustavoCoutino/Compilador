%{
package main

import (
    "fmt"
    "os"
    "strconv"
    "strings"
    "gustavocoutino.compilador/internal/types"
    "gustavocoutino.compilador/internal/semantics"
    "gustavocoutino.compilador/internal/quadruples"
    "gustavocoutino.compilador/internal/ops"
    "gustavocoutino.compilador/internal/virtual_machine"
)
%}

%union {
    entero  int
    flotante float64
    texto   string
    tipo types.Tipo
}

%token <texto> ID LITERAL
%token <entero> CTE_ENT
%token <flotante> CTE_FLOT

%type <tipo> Tipo TipoRetorno

%token PROGRAMA INICIO FIN VARS ENTERO FLOTANTE NULA RETORNO
%token SI SINO MIENTRAS HAZ ESCRIBE

%token MAS MENOS POR ENTRE
%token MAYOR MENOR IGUAL DIFERENTE
%token ASIGNAVAR

%token COMA PCOMA DOSPUNTOS
%token LPARENTESIS RPARENTESIS
%token LLLAVE RLLAVE
%token LCORCHETE RCORCHETE

%%

Programa : PROGRAMA ID PCOMA {
    quadruples.CrearCuadruploGotoInicio()
    semantics.IniciarPrograma($2)
} VarsOpt FuncsOpt INICIO {
    quadruples.ActualizarSalto()
} Cuerpo FIN {
    quadruples.CrearCuadruploFin()
    semantics.ImprimirTablaConstantes()
    semantics.ImprimirDirectorioFunciones()
    quadruples.ImprimirCuadruplos()
} ;
VarsOpt: Vars | /* vacío */ ;
FuncsOpt: FuncsOpt Funcs | /* vacío */ ;
Vars: VARS LLLAVE DeclaracionLista RLLAVE ;
DeclaracionLista: Declaracion | DeclaracionLista Declaracion ;
Declaracion: IdsLista DOSPUNTOS Tipo PCOMA {
    semantics.DeclararIdsActuales($3)
} ;
IdsLista: ID IdsListaExtension {
    semantics.AgregarIdActual($1)
};
IdsListaExtension: IdsListaExtension COMA ID {
    semantics.AgregarIdActual($3)
}| /* vacío */ ;
Tipo: ENTERO { $$ = types.TipoEntero } | FLOTANTE { $$ = types.TipoFlotante } ;
Cuerpo: LLLAVE EstatutosLista RLLAVE ;
EstatutosLista: EstatutosLista Estatuto | /* vacío */ ;
Funcs: TipoRetorno ID {
    semantics.DeclararVariable($2, $1)
} LPARENTESIS ParametrosOpt RPARENTESIS {
    semantics.IniciarFuncionConParametros($2, $1)
    semantics.AsignarCuadruploInicio(quadruples.ContadorActual())
} LLLAVE VarsOpt EstatutosLista RLLAVE PCOMA {
    semantics.TerminarFuncion()
    quadruples.GenerarCuadruploAcabarFunc()
} ;
TipoRetorno: Tipo | NULA {
    $$ = types.TipoNula
} ;
ParametrosOpt: ParametrosLista | /* vacío */ ;
ParametrosLista: ID DOSPUNTOS Tipo {
    semantics.AgregarParametroActual($1, $3)
} ParametrosListaExtension;
ParametrosListaExtension: ParametrosListaExtension COMA ID DOSPUNTOS Tipo {
    semantics.AgregarParametroActual($3, $5)
} | /* vacío */ ;
Estatuto: Asigna | Condicion | Ciclo | Llamada PCOMA | Imprime | EstatutoBloque | Retorno ;
Retorno: RETORNO Expresion {
    quadruples.GenerarCuadruploRetorno()
} PCOMA ;
EstatutoBloque: LCORCHETE EstatutosBloqueLista RCORCHETE ; 
EstatutosBloqueLista: EstatutosBloqueLista Estatuto | Estatuto ;
Asigna: ID ASIGNAVAR Expresion {
    quadruples.GenerarCuadruploAsigna($1)
} PCOMA ;
Expresion: Exp | Exp Operadores Exp {
    quadruples.GenerarCuadruplo()
} ;
Operadores: MAYOR {
    quadruples.EmpujarOperador(ops.MAYOR)
} | MENOR {
    quadruples.EmpujarOperador(ops.MENOR)
} | IGUAL {
    quadruples.EmpujarOperador(ops.IGUAL)
} | DIFERENTE {
    quadruples.EmpujarOperador(ops.DIFERENTE)
} ;
Exp: Exp MAS {
    quadruples.EmpujarOperador(ops.MAS)
} Termino {
    quadruples.GenerarCuadruplo()
} | Exp MENOS {
    quadruples.EmpujarOperador(ops.MENOS)
} Termino {
    quadruples.GenerarCuadruplo()
} | Termino ;
Termino: Termino POR {
    quadruples.EmpujarOperador(ops.POR)
} Factor {
    quadruples.GenerarCuadruplo()
} | Termino ENTRE {
    quadruples.EmpujarOperador(ops.ENTRE)
} Factor {
    quadruples.GenerarCuadruplo()
} | Factor ;
Factor: LPARENTESIS Expresion RPARENTESIS | FactorOpt | MAS FactorOpt | MENOS FactorOpt {
    quadruples.GenerarOperandoMenos()
} | Llamada ;
FactorOpt: ID {
    if variable, existe := semantics.BuscarVariable($1); existe {
        quadruples.EmpujarOperando(variable.Direccion, variable.Tipo)
    } else {
        semantics.ErrorSemantico(fmt.Sprintf("variable '%s' no declarada", $1))
    }
} | CTE ;
CTE: CTE_ENT {
    direccion, tipo := semantics.ProcesarConstante(strconv.Itoa($1), types.TipoEntero)
    quadruples.EmpujarOperando(direccion, tipo)
} | CTE_FLOT {
    direccion, tipo := semantics.ProcesarConstante(strconv.FormatFloat($1, 'g', -1, 64), types.TipoFlotante)
    quadruples.EmpujarOperando(direccion, tipo)
} ;
Condicion: SI LPARENTESIS Expresion RPARENTESIS {
    quadruples.EmpujarSalto(ops.GOTOF)
} Cuerpo SinoOpt PCOMA ;
SinoOpt: SINO {
    quadruples.ActualizarSino()
} Cuerpo {
    quadruples.ActualizarSalto()
} | /* vacío */ {
    quadruples.ActualizarSalto()
} ;
Ciclo: MIENTRAS {
    quadruples.GuardarMientrasUbicacion()
} LPARENTESIS Expresion RPARENTESIS {
    quadruples.CrearCuadruploMientrasGotof()
} HAZ Cuerpo PCOMA {
    quadruples.ActualizarMientras()
} ;
Imprime: ESCRIBE LPARENTESIS ImprimeLista RPARENTESIS PCOMA ;
ImprimeLista: ImprimeEl | ImprimeLista COMA ImprimeEl ;
ImprimeEl: Expresion {
    quadruples.GenerarCuadruploEscribe()
} | LITERAL {
    dir, _ := semantics.ProcesarConstante($1, types.TipoLiteral)
    quadruples.EmpujarOperando(dir, types.TipoLiteral)
    quadruples.GenerarCuadruploEscribe()
} ;
Llamada: ID {
    quadruples.GuardarNombreFuncionActual($1)
    quadruples.GenerarCuadruploEra()
} LPARENTESIS ArgumentosOpt RPARENTESIS {
    quadruples.GenerarCuadruploGosub()
    quadruples.GenerarCuadruploResultadoLlamada()
} ;
ArgumentosOpt: ArgumentosLista | /* vacío */ ;
ArgumentosLista: Expresion {
    quadruples.GenerarCuadruploParametro()
} | ArgumentosLista COMA Expresion {
    quadruples.GenerarCuadruploParametro()
} ;

%%

// Lexer es una representacion de un analizador lexico.
// Lexer contiene el indice del caracter actual de la entrda de caracteres,
// el indice del caracter leido, la entrada de caracteres, el caracter
// actual, la linea del programa, y la columna correspondiente a un caracter en la entrada
type Lexer struct{
    input string
    position int
    readPosition int
    ch byte
    line int
    column int
}



// Lista de palabras reservadas
var palabrasReservadas = map[string]int{
    "programa": PROGRAMA,
    "inicio": INICIO,
    "fin": FIN,
    "vars": VARS,
    "entero": ENTERO,
    "flotante": FLOTANTE,
    "nula":     NULA,
    "si":       SI,
    "sino":     SINO,
    "mientras": MIENTRAS,
    "haz":      HAZ,
    "escribe":  ESCRIBE,
    "retorno":   RETORNO,
}

// Lex es invocada por el parser para realizar el
// analisis lexico de la entrada
func (l *Lexer) Lex(lval *yySymType) int {
    tok := l.Next(lval)
    semantics.LineaActual = l.line
    return tok
}

// Next analiza el siguiente token en la entrada de
// caracteres usando las expresiones regulares y la lista
// de simbolos terminales de la gramatica
func (l *Lexer) Next(lval *yySymType) int {
    l.skipWhitespace()
    if l.ch == 0 {
        return 0  
    }
    switch l.ch {
        default:
            if isLetter(l.ch){
                token := l.readIdentifier()
                if tipo, esReservada := palabrasReservadas[token]; esReservada{
                    return tipo
                }
                lval.texto = token
                return ID
            } else if isDigit(l.ch) {
                return l.readNumber(lval)
            } else {
                l.Error(fmt.Sprintf("Carácter no reconocido: %c", l.ch))
            }
        case '+':
            l.readChar()
            return MAS
        case '-':
            l.readChar()
            return MENOS
        case '*':
            l.readChar()
            return POR
        case '/':
            l.readChar()
            return ENTRE
        case '>':
            l.readChar()
            return MAYOR
        case '<':
            l.readChar()
            return MENOR
        case ',':
            l.readChar()
            return COMA
        case ';':
            l.readChar()
            return PCOMA
        case ':':
            l.readChar()
            return DOSPUNTOS
        case '(':
            l.readChar()
            return LPARENTESIS
        case ')':
            l.readChar()
            return RPARENTESIS
        case '{':
            l.readChar()
            return LLLAVE
        case '}':
            l.readChar()
            return RLLAVE
        case '[':
            l.readChar()
            return LCORCHETE
        case ']':
            l.readChar()
            return RCORCHETE
        case '=':
            l.readChar()
            if l.ch == '=' {
                l.readChar()
                return IGUAL 
            } else {
                return ASIGNAVAR
            }
        case '!':
            l.readChar()
            if l.ch == '=' {
                l.readChar()
                return DIFERENTE 
            } else {
                l.Error(fmt.Sprintf("Carácter no reconocido: %c", l.ch))
            }
        case '"':
            l.readChar()
            position := l.position
            for isLetter(l.ch) || isDigit(l.ch) || isWhitespace(l.ch) || isOther(l.ch) {
                l.readChar()
            }
            if l.ch == '"' {
                lval.texto = l.input[position:l.position]
                l.readChar()
                return LITERAL
            }
            l.Error(fmt.Sprint("Formato de literal incorrecto"))
    }
    return 0
}

// readIdentifier lee el identificador y regresa
// dicho id
func (l *Lexer) readIdentifier() string {
    position := l.position
    for isLetter(l.ch) || isDigit(l.ch) || l.ch == '_' || l.ch == '-'  {
        l.readChar()
        
    }
    return l.input[position:l.position]
}

// readNumber lee un numero de acuerdo a las expresiones
// regulares definidas y regresa este numero. Este numero puede ser un entero
// o un flotante, y este flotante puede estar escrito en
// notacion cientifica
func (l *Lexer) readNumber(lval *yySymType) int {
    position := l.position
    for isDigit(l.ch) {
        l.readChar()
    }

    if l.ch == '.' && isDigit(l.peekChar()){
        l.readChar()
        for isDigit(l.ch) {
            l.readChar()
        }
        if l.ch == 'e' {
            if l.peekChar() == '+' || l.peekChar() == '-' {
                l.readChar()
                l.readChar()
                if !isDigit(l.ch){
                    l.Error("Notación científica inválida: se esperaban dígitos después del exponente")
                }
                for isDigit(l.ch) {
                    l.readChar()
                }
            } else {
                l.Error("Notación científica inválida: se esperaban un carácter '+' o '-' despues del exponente")
            }
        } 
        valor, _ := strconv.ParseFloat(l.input[position:l.position], 64)
        lval.flotante = valor
        return CTE_FLOT
    }
    valor, _ := strconv.Atoi(l.input[position:l.position])
    lval.entero = valor
    return CTE_ENT
}

// readChar lee el siguiente caracter.
// Actualiza la posicion actual de la entrada de
// caracteres, y la posicion del caracter actual leido
func (l *Lexer) readChar(){
    if l.readPosition >= len(l.input){
        l.ch = 0
    } else {
        l.ch = l.input[l.readPosition]
    }
    l.position = l.readPosition
    l.readPosition += 1
    if l.ch == '\n' {
        l.line++
        l.column = 0
    } else {
        l.column++
    }
}

// peekChar ve el siguiente caracter sin
// consumirlo, para distinguir tokens de dos
// caracteres como '==' o '!=', o detectar
// el punto decimal de un flotante
func (l *Lexer) peekChar() byte {
    if l.readPosition >= len(l.input){
        return 0
    } else {
        return l.input[l.readPosition]
    }
}

// isLetter revisa si el caracter se encuentra entre la
// a y z (no distingue entre mayusculas y minisculas)
func isLetter(ch byte) bool {
    return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z'
}  

// isDigit revisa si el caracter es un numero
func isDigit(ch byte) bool {
    return ch >= '0' && ch <= '9'
}

// isWhitespace revisa si el caracter es un espacio en blanco
func isWhitespace(ch byte) bool {
    return ch == ' ' || ch == '\t' || ch == '\n' || ch == '\r'
}

// isOther revisa si el caracter es un símbolo permitido dentro de un literal
// (puntuación o signos comunes distintos de la comilla doble)
func isOther(ch byte) bool {
    switch ch {
    case '-', '_', '=', '?', '.', ',', ';', ':', '!', '@', '#', '$',
        '%', '^', '&', '*', '(', ')', '+', '/', '<', '>',
        '[', ']', '{', '}', '|', '\\', '~', '\'', '`':
        return true
    }
    return false
}

// skipWhitespace se salta los espacios en blanco y avanza el lexer
func (l *Lexer) skipWhitespace() {
    for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
        l.readChar()
    }
}

// Error es la implementacion del método de error
// de la interfaz de Lexer para imprimir un error
func (l *Lexer) Error(s string) {
    fmt.Fprintf(os.Stderr, "Error de sintaxis en línea %d, columna %d: %s\n", l.line, l.column, s)
    lineas := strings.Split(l.input, "\n")
    if l.line-1 < len(lineas) {
        fmt.Fprintf(os.Stderr, "  %s\n", lineas[l.line-1])
        fmt.Fprintf(os.Stderr, "  %s^\n", strings.Repeat(" ", l.column-1))
    }
}

func main() {
     if len(os.Args) < 2 {
        fmt.Println("Porfavor provee el nombre de un archivo")
        return
    }
    fileName := os.Args[1]
    data, err := os.ReadFile(fileName)
    if err != nil {
        fmt.Println("Error al leer el archivo:", err)
        return
    }
    lexer := &Lexer{input: string(data), position: 0, readPosition: 0, line: 1, column: 0}
    lexer.readChar()
    ok := yyParse(lexer)
    if ok == 0 {
        if semantics.HasError() {
            fmt.Println("El análisis semántico tiene errores")
            return
        }
        fmt.Println("Compilación exitosa")
        vm := virtualmachine.NewVM(quadruples.GetFilaCuadruplos())
        vm.Ejecutar()
        vm.ImprimirMapaMemoriaGlobal()
    } else if ok == 1 {
        fmt.Println("El análisis léxico contiene errores")
    } else if ok == 2 {
        fmt.Println("Agotamiento de memoria")
    }
}