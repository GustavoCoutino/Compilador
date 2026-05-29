package types

type Tipo int

const (
	TipoEntero Tipo = iota
	TipoFlotante
	TipoError
	TipoNula
	TipoLiteral
)

func (t Tipo) String() string {
	switch t {
		case TipoEntero:
			return "entero"
		case TipoFlotante:
			return "flotante"
		case TipoError:
			return "error"
		case TipoNula:
			return "nula"
		default:
			return "error"
	}
}