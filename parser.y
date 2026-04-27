%{
package main

import (
    "fmt"
    "os"
)
%}

%union {
    entero  int
    flotante float64
    texto   string
}

%token <texto> ID LITERAL
%token <entero> CTE_ENT
%token <flotante> CTE_FLOT

%token PROGRAMA INICIO FIN VARS ENTERO FLOTANTE NULA
%token SI SINO MIENTRAS HAZ ESCRIBE

%token MAS MENOS POR ENTRE
%token MAYOR MENOR IGUAL DIFERENTE
%token ASIGNAVAR

%token COMA PCOMA DOSPUNTOS
%token LPARENTESIS RPARENTESIS
%token LLLAVE RLLAVE
%token LCORCHETE RCORCHETE

%%

Programa : PROGRAMA ID PCOMA VarsOpt FuncsOpt INICIO Cuerpo FIN ;
VarsOpt: Vars | /* vacío */ ;
FuncsOpt: FuncsOpt Funcs | /* vacío */ ;
Vars: VARS Declaracion DeclaracionLista ; 
DeclaracionLista: DeclaracionLista Declaracion | /* vacío */ ;
Declaracion: IdsLista DOSPUNTOS Tipo PCOMA ;
IdsLista: ID IdsListaExtension ;
IdsListaExtension: IdsListaExtension COMA ID | /* vacío */ ;
Tipo: ENTERO | FLOTANTE ;
Cuerpo: LLLAVE EstatutosLista RLLAVE ;
EstatutosLista: EstatutosLista Estatuto | /* vacío */ ;
Funcs: TipoRetorno ID LPARENTESIS ParametrosOpt RPARENTESIS LLLAVE VarsOpt Cuerpo RLLAVE PCOMA ;
TipoRetorno: Tipo | NULA ;
ParametrosOpt: ParametrosLista | /* vacío */ ;
ParametrosLista: ID DOSPUNTOS Tipo ParametrosListaExtension;
ParametrosListaExtension: ParametrosListaExtension COMA ID DOSPUNTOS Tipo | /* vacío */ ;
Estatuto: Asigna | Condicion | Ciclo | Llamada PCOMA | Imprime | EstatutoBloque ;
EstatutoBloque: LCORCHETE EstatutosBloqueLista RCORCHETE ; 
EstatutosBloqueLista: Estatuto EstatutosBloqueLista | Estatuto ;
Asigna: ID ASIGNAVAR Expresion PCOMA ;
Expresion: Exp | Exp Operadores Exp ;
Operadores: MAYOR | MENOR | IGUAL | DIFERENTE ;
Exp: Exp MAS Termino | Exp MENOS Termino | Termino ;
Termino: Termino POR Factor | Termino ENTRE Factor | Factor ;
Factor: LPARENTESIS Expresion RPARENTESIS | MAS FactorOpt | MENOS FactorOpt | Llamada ;
FactorOpt: ID | CTE ;
CTE: CTE_ENT | CTE_FLOT ;
Condicion: SI LPARENTESIS Expresion RPARENTESIS Cuerpo SinoOpt PCOMA ;
SinoOpt: SINO Cuerpo | /* vacío */ ;
Ciclo: MIENTRAS LPARENTESIS Expresion RPARENTESIS HAZ Cuerpo PCOMA ;
Imprime: ESCRIBE LPARENTESIS ImprimeLista RPARENTESIS PCOMA ;
ImprimeLista: ImprimeEl | ImprimeLista COMA ImprimeEl ;
ImprimeEl: Expresion | LITERAL ;
Llamada: ID LPARENTESIS ArgumentosOpt RPARENTESIS ;
ArgumentosOpt: ArgumentosLista | /* vacío */ ;
ArgumentosLista: Expresion | ArgumentosLista COMA Expresion ;

%%

type Lexer struct{}

func (l *Lexer) Lex(lval *yySymType) int {
    return 0
}

func (l *Lexer) Error(s string) {
    fmt.Fprintln(os.Stderr, "Error de sintaxis:", s)
}

func main() {
    yyParse(&Lexer{})
}