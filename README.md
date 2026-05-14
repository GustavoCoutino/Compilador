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

- [ ] **Validaciones semánticas** — verificación de tipos, declaración previa
      de variables y funciones, llamadas, etc.
- [ ] **Código intermedio** — generación a partir de las
      acciones de gramática.
- [ ] **Máquina virtual** — intérprete del código intermedio generado.
