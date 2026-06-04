# Compilador de lenguaje Patito (Baby Duck)

Patito es un lenguaje imperativo y procedural didáctico para enseñar los fundamentos de desarrollo de compiladores. Soporta declaracion de variables de tipo entero y flotante, funciones con retornos de tipo y nulas, operaciones de impresion de expresiones y cadenas de caracteres, llamada de funciones, ciclos, condicionales, y recursión.

Al ser un ejercicio académico, este compilador no tiene las mismas fases que un compilador promedio: carece de árbol de sintaxis abstracta, y el proceso lexico, sintáctico, y semántico ocurren en una sola pasada del archivo de entrada

## Lexer

Lee carácter por carácter y asocia cadenas compuestas de estas unidades en tokens. Los tokens identificados en esta fase del compilador involucran operadores (+, -, \*, /), palabras reservadas (mientras, si, sino), e identificadores. Se utiliza la herramienta `goyacc` para generar un parser en Go, lo cual requiere de una implementacion de la interfaz Lexer y un método de errores, los cuales se llaman para cada token por dicho parser. Las expresiones regulares correspondientes al lenguaje se implementan programáticamente por pura cuestión de aprendizaje.

## Parser

Se utiliza `goyacc` para generar un parser **LALR(1)** a partir de la gramática definida en [`parser.y`](parser.y). Cada regla de producción puede llevar acciones semánticas (escritas en Go) que se ejecutan al momento de la reducción de esa regla, durante la misma pasada que el análisis sintáctico.

Como el compilador carece de árbol de sintaxis abstracta, no existe una fase posterior que recorra un árbol — las acciones semánticas hacen directamente, en el momento del parseo:

- inserción y consulta de variables, parámetros, constantes y funciones en sus tablas correspondientes;
- verificación de tipos contra el cubo semántico;
- asignación de direcciones virtuales a variables, constantes y temporales;
- generación de los cuádruplos, que se encolan en la fila de salida.

Las pilas auxiliares de operandos, operadores y saltos constituyen el "estado" que reemplaza al árbol durante el parsing. Guardan información parcial que se consume al reducir las reglas correspondientes.

## Análisis semántico

Generar una fila de cuadruplos que será utilizada por la máquina virtual, y realiza la validación semántica necesaria para los cuádruplos (verificar cantidad correcta y tipo de parámetros, asignación correcta de variables, comparación válida de expresiones) Debido a la falta de AST, los puntos neurálgicos de las acciones semánticas se hacen simultaneamente al proceso de parsing. La verificación de tipos de asignación y de expresione se realiza con un cubo semántico. Los cuadruplos generados en cada regla se insertan a una fila. A la par, en este proceso se les asigna una dirección en memoria a las variables y temporales (locales y globales) que se encuentren en el análisis.

# Manejo de memoria

Se encarga de asignar una dirección de memoria a los resultados de los cuadruplos generados en el análisis semántico, asi como a las variables y constantes identificadas. Se utiliza la siguiente formula para decidir el segmento de memoria que se le asigna a un resultado:

`dirección = (segmento * len(tiposDir) + idxTipo) * tamanoBloque`

Este asignador también se encarga de liberar memoria una vez que una función local termine su compilación.

Las divisiones de los segmentos son las siguientes (con `tamanoBloque = 2000` y los tres tipos direccionables `entero`, `flotante`, `cadena`):

| Segmento      | entero        | flotante      | cadena        |
| ------------- | ------------- | ------------- | ------------- |
| **Global**    | 0 – 1999      | 2000 – 3999   | 4000 – 5999   |
| **Local**     | 6000 – 7999   | 8000 – 9999   | 10000 – 11999 |
| **Temporal**  | 12000 – 13999 | 14000 – 15999 | 16000 – 17999 |
| **Constante** | 18000 – 19999 | 20000 – 21999 | 22000 – 23999 |

A partir de una dirección, el segmento al que pertenece se puede recuperar invirtiendo la fórmula: `segmento = dirección / tamanoBloque / len(tiposDir)`. La máquina virtual usa exactamente eso para decidir si una dirección vive en la memoria global o en el registro de activación actual.

# Máquina virtual

Última fase del compilador en donde se recibe la fila de cuádruplos generada en la pasada
anterior y la interpreta secuencialmente, ejecutando la acción que corresponde
al operador de cada cuádruplo (sumar, comparar, brincar, llamar, imprimir, etc.).

La VM mantiene cuatro estructuras vivas durante toda la ejecución:

- **`ip`** — apuntador al cuádruplo actual. Avanza secuencialmente excepto cuando
  un `GOTO`/`GOTOF`/`GOTOT`/`GOSUB`/`ENDFUNC`/`RETORNO` lo redirige.
- **`global`** — mapa de variables globales, constantes (cargadas al arranque desde la tabla de
  constantes del compilador) y los globales de retorno de cada función.
- **`pilaAR`** — pila de registros de activación. Cada elemento es un mapa con los locales y temporales de UNA llamada activa.
  Se empuja uno nuevo al iniciar una llamada y se saca al terminarla.
- **`pilaIP`** — pila de direcciones de retorno. `GOSUB` empuja el siguiente cuadruplo a una llamada de función,
  `ENDFUNC`/`RETORNO` lo sacan para volver al cuádruplo siguiente a la llamada.
- **`pendiente`** - pila de AR pendiente. Previene que se pierde el AR del parámetro que se empuja a la pila de AR.

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

Los programas de prueba viven en [`program_tests/`](program_tests/) y están organizados en dos grandes grupos:

- **`success/`** — programas que deben compilar y ejecutarse sin errores. Cada subcarpeta cubre un aspecto del lenguaje:
  - `mientras/` — ciclos.
  - `si/`, `sino/` — condicionales simples y con `sino`.
  - `escribe/` — impresión de literales, expresiones y multi-argumento.
  - `operadores/` — aritméticos, relacionales y precedencia con paréntesis.
  - `funciones/` — funciones con y sin valor de retorno, llamadas como estatuto y dentro de expresiones.
- **`failure/`** — programas semánticamente inválidos que el compilador debe rechazar:
  - `asignacion/` — variables no declaradas, tipos incompatibles.
  - `parametros/` — número incorrecto de argumentos, tipos que no coinciden.
- **`flujo_completo.patito`** — programa integrador que ejercita todas las reglas de la gramática en un solo archivo (vars con ambos tipos, los tres tipos de retorno, `si`/`sino`/`mientras`, `escribe`, llamadas como estatuto y como factor, bloques `[...]`, paréntesis y unarios).

### Correr tests unitatios

```
go test ./...                              # todos los paquetes
go test ./internal/semantics/              # solo un paquete
go test ./internal/semantics/ -run Buscar  # un subset por nombre
go test ./... -v                           # con detalle de cada test
```
