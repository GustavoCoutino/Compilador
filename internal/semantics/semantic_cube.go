package semantics

import (
	"gustavocoutino.compilador/internal/ops"
	"gustavocoutino.compilador/internal/types"
)

type claveCubo struct {
	izquierda types.Tipo
	derecha types.Tipo
	operador int
}

var cuboSemantico = map[claveCubo]types.Tipo {
	// Izquierda: entero, derecha: entero
	{izquierda: types.TipoEntero, derecha: types.TipoEntero, operador: ops.MAS}: types.TipoEntero,
	{izquierda: types.TipoEntero, derecha: types.TipoEntero, operador: ops.MENOS}: types.TipoEntero,
	{izquierda: types.TipoEntero, derecha: types.TipoEntero, operador: ops.ENTRE}: types.TipoEntero,
	{izquierda: types.TipoEntero, derecha: types.TipoEntero, operador: ops.POR}: types.TipoEntero,
	{izquierda: types.TipoEntero, derecha: types.TipoEntero, operador: ops.MAYOR}: types.TipoEntero,
	{izquierda: types.TipoEntero, derecha: types.TipoEntero, operador: ops.MENOR}: types.TipoEntero,
	{izquierda: types.TipoEntero, derecha: types.TipoEntero, operador: ops.IGUAL}: types.TipoEntero,
	{izquierda: types.TipoEntero, derecha: types.TipoEntero, operador: ops.DIFERENTE}: types.TipoEntero,
	{izquierda: types.TipoEntero, derecha: types.TipoEntero, operador: ops.ASIGNAVAR}: types.TipoEntero,
	// Izquierda: entero, derecha: flotante
	{izquierda: types.TipoEntero, derecha: types.TipoFlotante, operador: ops.MAS}: types.TipoFlotante,
	{izquierda: types.TipoEntero, derecha: types.TipoFlotante, operador: ops.MENOS}: types.TipoFlotante,
	{izquierda: types.TipoEntero, derecha: types.TipoFlotante, operador: ops.ENTRE}: types.TipoFlotante,
	{izquierda: types.TipoEntero, derecha: types.TipoFlotante, operador: ops.POR}: types.TipoFlotante,
	{izquierda: types.TipoEntero, derecha: types.TipoFlotante, operador: ops.MAYOR}: types.TipoEntero,
	{izquierda: types.TipoEntero, derecha: types.TipoFlotante, operador: ops.MENOR}: types.TipoEntero,
	{izquierda: types.TipoEntero, derecha: types.TipoFlotante, operador: ops.IGUAL}: types.TipoEntero,
	{izquierda: types.TipoEntero, derecha: types.TipoFlotante, operador: ops.DIFERENTE}: types.TipoEntero,
	{izquierda: types.TipoEntero, derecha: types.TipoFlotante, operador: ops.ASIGNAVAR}: types.TipoError,
	// Izquierda: flotante, derecha: entero
	{izquierda: types.TipoFlotante, derecha: types.TipoEntero, operador: ops.MAS}: types.TipoFlotante,
	{izquierda: types.TipoFlotante, derecha: types.TipoEntero, operador: ops.MENOS}: types.TipoFlotante,
	{izquierda: types.TipoFlotante, derecha: types.TipoEntero, operador: ops.ENTRE}: types.TipoFlotante,
	{izquierda: types.TipoFlotante, derecha: types.TipoEntero, operador: ops.POR}: types.TipoFlotante,
	{izquierda: types.TipoFlotante, derecha: types.TipoEntero, operador: ops.MAYOR}: types.TipoEntero,
	{izquierda: types.TipoFlotante, derecha: types.TipoEntero, operador: ops.MENOR}: types.TipoEntero,
	{izquierda: types.TipoFlotante, derecha: types.TipoEntero, operador: ops.IGUAL}: types.TipoEntero,
	{izquierda: types.TipoFlotante, derecha: types.TipoEntero, operador: ops.DIFERENTE}: types.TipoEntero,
	{izquierda: types.TipoFlotante, derecha: types.TipoEntero, operador: ops.ASIGNAVAR}: types.TipoFlotante,
	// Izquierda: flotante, derecha: flotante
	{izquierda: types.TipoFlotante, derecha: types.TipoFlotante, operador: ops.MAS}: types.TipoFlotante,
	{izquierda: types.TipoFlotante, derecha: types.TipoFlotante, operador: ops.MENOS}: types.TipoFlotante,
	{izquierda: types.TipoFlotante, derecha: types.TipoFlotante, operador: ops.ENTRE}: types.TipoFlotante,
	{izquierda: types.TipoFlotante, derecha: types.TipoFlotante, operador: ops.POR}: types.TipoFlotante,
	{izquierda: types.TipoFlotante, derecha: types.TipoFlotante, operador: ops.MAYOR}: types.TipoEntero,
	{izquierda: types.TipoFlotante, derecha: types.TipoFlotante, operador: ops.MENOR}: types.TipoEntero,
	{izquierda: types.TipoFlotante, derecha: types.TipoFlotante, operador: ops.IGUAL}: types.TipoEntero,
	{izquierda: types.TipoFlotante, derecha: types.TipoFlotante, operador: ops.DIFERENTE}: types.TipoEntero,
	{izquierda: types.TipoFlotante, derecha: types.TipoFlotante, operador: ops.ASIGNAVAR}: types.TipoFlotante,

}

func ValidarSemantica(izq, der types.Tipo, op int) types.Tipo {
	if tipo, ok := cuboSemantico[claveCubo{izq, der, op}]; ok  {
		return tipo
	}
	return types.TipoError
}