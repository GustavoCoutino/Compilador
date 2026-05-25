# Compilador Patito

Compilador para el lenguaje procedural **Patito**, escrito en Go por Gustavo Coutiño Ocampo - A01412203.

Actualmente el compilador realiza:

- **Análisis léxico**: convierte la entrada en una secuencia de tokens
  (palabras reservadas, identificadores, constantes enteras y flotantes,
  literales de cadena, operadores y signos de puntuación). Implementado
  manualmente como una máquina de estados sobre bytes.
- **Análisis sintáctico**: valida que la secuencia de tokens respete la
  gramática del lenguaje Patito. Implementado con `goyacc` (LALR(1))
  a partir de la gramática definida en `parser.y`.
- **Análisis semántico**: conecta la definición de variables y funciones
  con sus usos, guardando sus datos y alcance correcto. Implementando
  con un mapa de directories y un mapa de variables por alcance.
  Si el programa de entrada es válido, el compilador imprime
  `El análisis léxico fue exitoso`. Si no, imprime un mensaje de
  `Error de sintaxis` describiendo la falla e imprimiendo el error.

## Requisitos

- Go 1.21+
- `goyacc` (`go install golang.org/x/tools/cmd/goyacc@latest`)
- `make`

## Cómo ejecutar

El `makefile` contiene los siguientes comandos:

```sh
make tidy                              # ordena y formatea dependencias
make build                             # genera parser, compila y corre programa1.patito
make build PROGRAM=programa4.patito    # corre el archivo que se indique
make erase                             # borra parser.go, y.output y el binario
```

El nombre del archivo de entrada se pasa como variable `PROGRAM`.
Si se omite, se usa `programa1.patito` por defecto.

## Especificación del lenguaje

### Tokens (expresiones regulares)

El analizador léxico ([`parser.y`](parser.y)) reconoce los siguientes tokens.
Las expresiones regulares describen el comportamiento de la máquina de estados
implementada manualmente.

| Token                 | Expresión regular              | Ejemplo              |
| --------------------- | ------------------------------ | -------------------- |
| `ID` (identificador)  | `[A-Za-z][A-Za-z0-9_-]*`       | `contador`, `base_1` |
| `CTE_ENT` (entero)    | `[0-9]+`                       | `42`                 |
| `CTE_FLOT` (flotante) | `[0-9]+\.[0-9]+(e[+-][0-9]+)?` | `3.14`, `1.0e-3`     |
| `LITERAL` (cadena)    | `"[^"]*"`                      | `"hola adios"`       |

> Un flotante requiere punto decimal; la notación científica usa `e` minúscula
> con signo explícito (`e+` / `e-`).

**Operadores y signos**

| Categoría    | Tokens                                                                                      |
| ------------ | ------------------------------------------------------------------------------------------- |
| Aritméticos  | `+` &nbsp; `-` &nbsp; `*` &nbsp; `/`                                                        |
| Relacionales | `>` &nbsp; `<` &nbsp; `==` &nbsp; `!=`                                                      |
| Asignación   | `=`                                                                                         |
| Puntuación   | `,` &nbsp; `;` &nbsp; `:` &nbsp; `(` &nbsp; `)` &nbsp; `{` &nbsp; `}` &nbsp; `[` &nbsp; `]` |

**Palabras reservadas**

`programa` &nbsp; `inicio` &nbsp; `fin` &nbsp; `vars` &nbsp; `entero` &nbsp;
`flotante` &nbsp; `nula` &nbsp; `si` &nbsp; `sino` &nbsp; `mientras` &nbsp;
`haz` &nbsp; `escribe` &nbsp; `retorno`

### Gramática

Reglas de producción (notación BNF; `ε` es la cadena vacía; los terminales
`id`, `cte_ent`, `cte_flot` y `literal` corresponden a los tokens de arriba):

```text
Programa          → "programa" id ";" VarsOpt FuncsOpt "inicio" Cuerpo "fin"

VarsOpt           → Vars | ε
Vars              → "vars" "{" DeclaracionLista "}"
DeclaracionLista  → Declaracion | DeclaracionLista Declaracion
Declaracion       → IdsLista ":" Tipo ";"
IdsLista          → id IdsListaExtension
IdsListaExtension → IdsListaExtension "," id | ε
Tipo              → "entero" | "flotante"

FuncsOpt          → FuncsOpt Funcs | ε
Funcs             → TipoRetorno id "(" ParametrosOpt ")" "{" VarsOpt EstatutosLista "}" ";"
TipoRetorno       → Tipo | "nula"
ParametrosOpt     → ParametrosLista | ε
ParametrosLista   → id ":" Tipo ParametrosListaExtension
ParametrosListaExtension → ParametrosListaExtension "," id ":" Tipo | ε

Cuerpo            → "{" EstatutosLista "}"
EstatutosLista    → EstatutosLista Estatuto | ε
Estatuto          → Asigna | Condicion | Ciclo | Llamada ";" | Imprime
                  | EstatutoBloque | Retorno
EstatutoBloque    → "[" EstatutosBloqueLista "]"
EstatutosBloqueLista → EstatutosBloqueLista Estatuto | Estatuto

Asigna            → id "=" Expresion ";"
Retorno           → "retorno" Expresion ";"
Condicion         → "si" "(" Expresion ")" Cuerpo SinoOpt ";"
SinoOpt           → "sino" Cuerpo | ε
Ciclo             → "mientras" "(" Expresion ")" "haz" Cuerpo ";"
Imprime           → "escribe" "(" ImprimeLista ")" ";"
ImprimeLista      → ImprimeEl | ImprimeLista "," ImprimeEl
ImprimeEl         → Expresion | literal

Expresion         → Exp | Exp Operadores Exp
Operadores        → ">" | "<" | "==" | "!="
Exp               → Exp "+" Termino | Exp "-" Termino | Termino
Termino           → Termino "*" Factor | Termino "/" Factor | Factor
Factor            → "(" Expresion ")" | FactorOpt | "+" FactorOpt | "-" FactorOpt | Llamada
FactorOpt         → id | CTE
CTE               → cte_ent | cte_flot

Llamada           → id "(" ArgumentosOpt ")"
ArgumentosOpt     → ArgumentosLista | ε
ArgumentosLista   → Expresion | ArgumentosLista "," Expresion
```

## Programas de prueba

| Archivo            | Tipo       | Descripción                                                               |
| ------------------ | ---------- | ------------------------------------------------------------------------- |
| `programa1.patito` | Correcto   | Programa válido con `vars`, función, `si/sino`, `mientras` y `escribe`    |
| `programa2.patito` | Incorrecto | Falta el encabezado `programa <id>;` y notación científica inválida       |
| `programa3.patito` | Incorrecto | Falta `;` tras encabezado y tipo faltante en `vars`                       |
| `programa4.patito` | Correcto   | Cubre `nula`, `==`, `!=`, `-`, `e-3`, `escribe` multi-arg, bloque `[...]` |
| `programa5.patito` | Correcto   | Programa mínimo válido (sin `vars`, sin funciones, cuerpo vacío)          |
| `programa6.patito` | Correcto   | Literales con varios símbolos                                             |
| `programa7.patito` | Incorrecto | Faltan `;`, operando faltante                                             |
| `programa8.patito` | Incorrecto | Verifica límite de unario en gramática (`-id` válido, `-(expr)` inválido) |

## To Do

- [ ] **Código intermedio** — generación a partir de las
      acciones de gramática.
- [ ] **Máquina virtual** — intérprete del código intermedio generado.
