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
	default:
		return "?"
	}
}