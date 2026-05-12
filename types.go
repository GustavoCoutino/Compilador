package main

type Tipo int

const (
	TipoConstante Tipo = iota
	TipoFlotante
	TipoError
	TipoNula
)

func (t Tipo) String() string {
	switch t {
		case TipoConstante:
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