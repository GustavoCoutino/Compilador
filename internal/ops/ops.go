package ops

const (
	MAS = iota
	MENOS
	POR
	ENTRE
	MAYOR
	MENOR
	IGUAL
	DIFERENTE
	ASIGNAVAR
	IMPRIME
	RETORNO
	SI
	SINO
	MIENTRAS
	GOTO
	GOTOF
	GOTOT
	ENDFUNC
	ERA
	PARAM
	GOSUB
	FIN
)

func Simbolo(op int) string {
	switch op {
	case MAS:
		return "+"
	case MENOS:
		return "-"
	case POR:
		return "*"
	case ENTRE:
		return "/"
	case MAYOR:
		return ">"
	case MENOR:
		return "<"
	case IGUAL:
		return "=="
	case DIFERENTE:
		return "!="
	case ASIGNAVAR:
		return "="
	case IMPRIME:
		return "escribe"
	case RETORNO:
		return "retorno"
	case SI:	
		return "si"
	case SINO:
		return "sino"
	case MIENTRAS:
		return "mientras"
	case GOTO:
		return "goto"
	case GOTOF:
		return "gotof"
	case GOTOT:
		return "gotot"
	case ENDFUNC:
		return "endfunc"
	case ERA:
		return "era"
	case PARAM:
		return "param"
	case GOSUB:
		return "gosub"
	default:
		return "?"
	}
}