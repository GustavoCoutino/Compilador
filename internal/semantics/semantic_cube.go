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
	{izquierda: types.TipoConstante, derecha: types.TipoConstante, operador: ops.MAS}: types.TipoConstante,
	{izquierda: types.TipoConstante, derecha: types.TipoConstante, operador: ops.MENOS}: types.TipoConstante,
	{izquierda: types.TipoConstante, derecha: types.TipoConstante, operador: ops.ENTRE}: types.TipoConstante,
	{izquierda: types.TipoConstante, derecha: types.TipoConstante, operador: ops.POR}: types.TipoConstante,
	{izquierda: types.TipoConstante, derecha: types.TipoConstante, operador: ops.MAYOR}: types.TipoConstante,
	{izquierda: types.TipoConstante, derecha: types.TipoConstante, operador: ops.MENOR}: types.TipoConstante,
	{izquierda: types.TipoConstante, derecha: types.TipoConstante, operador: ops.IGUAL}: types.TipoConstante,
	{izquierda: types.TipoConstante, derecha: types.TipoConstante, operador: ops.DIFERENTE}: types.TipoConstante,
	{izquierda: types.TipoConstante, derecha: types.TipoConstante, operador: ops.ASIGNAVAR}: types.TipoConstante,
	// Izquierda: entero, derecha: flotante
	{izquierda: types.TipoConstante, derecha: types.TipoFlotante, operador: ops.MAS}: types.TipoFlotante,
	{izquierda: types.TipoConstante, derecha: types.TipoFlotante, operador: ops.MENOS}: types.TipoFlotante,
	{izquierda: types.TipoConstante, derecha: types.TipoFlotante, operador: ops.ENTRE}: types.TipoFlotante,
	{izquierda: types.TipoConstante, derecha: types.TipoFlotante, operador: ops.POR}: types.TipoFlotante,
	{izquierda: types.TipoConstante, derecha: types.TipoFlotante, operador: ops.MAYOR}: types.TipoConstante,
	{izquierda: types.TipoConstante, derecha: types.TipoFlotante, operador: ops.MENOR}: types.TipoConstante,
	{izquierda: types.TipoConstante, derecha: types.TipoFlotante, operador: ops.IGUAL}: types.TipoConstante,
	{izquierda: types.TipoConstante, derecha: types.TipoFlotante, operador: ops.DIFERENTE}: types.TipoConstante,
	{izquierda: types.TipoConstante, derecha: types.TipoFlotante, operador: ops.ASIGNAVAR}: types.TipoError,
	// Izquierda: flotante, derecha: entero
	{izquierda: types.TipoConstante, derecha: types.TipoFlotante, operador: ops.MAS}: types.TipoFlotante,
	{izquierda: types.TipoConstante, derecha: types.TipoFlotante, operador: ops.MENOS}: types.TipoFlotante,
	{izquierda: types.TipoConstante, derecha: types.TipoFlotante, operador: ops.ENTRE}: types.TipoFlotante,
	{izquierda: types.TipoConstante, derecha: types.TipoFlotante, operador: ops.POR}: types.TipoFlotante,
	{izquierda: types.TipoConstante, derecha: types.TipoFlotante, operador: ops.MAYOR}: types.TipoConstante,
	{izquierda: types.TipoConstante, derecha: types.TipoFlotante, operador: ops.MENOR}: types.TipoConstante,
	{izquierda: types.TipoConstante, derecha: types.TipoFlotante, operador: ops.IGUAL}: types.TipoConstante,
	{izquierda: types.TipoConstante, derecha: types.TipoFlotante, operador: ops.DIFERENTE}: types.TipoConstante,
	{izquierda: types.TipoConstante, derecha: types.TipoFlotante, operador: ops.ASIGNAVAR}: types.TipoFlotante,
	// Izquierda: flotante, derecha: flotante
	{izquierda: types.TipoConstante, derecha: types.TipoFlotante, operador: ops.MAS}: types.TipoFlotante,
	{izquierda: types.TipoConstante, derecha: types.TipoFlotante, operador: ops.MENOS}: types.TipoFlotante,
	{izquierda: types.TipoConstante, derecha: types.TipoFlotante, operador: ops.ENTRE}: types.TipoFlotante,
	{izquierda: types.TipoConstante, derecha: types.TipoFlotante, operador: ops.POR}: types.TipoFlotante,
	{izquierda: types.TipoConstante, derecha: types.TipoFlotante, operador: ops.MAYOR}: types.TipoConstante,
	{izquierda: types.TipoConstante, derecha: types.TipoFlotante, operador: ops.MENOR}: types.TipoConstante,
	{izquierda: types.TipoConstante, derecha: types.TipoFlotante, operador: ops.IGUAL}: types.TipoConstante,
	{izquierda: types.TipoConstante, derecha: types.TipoFlotante, operador: ops.DIFERENTE}: types.TipoConstante,
	{izquierda: types.TipoConstante, derecha: types.TipoFlotante, operador: ops.ASIGNAVAR}: types.TipoFlotante,

}

func validarSemantica(izq, der types.Tipo, op int) types.Tipo {
	if tipo, ok := cuboSemantico[claveCubo{izq, der, op}]; ok  {
		return tipo
	}
	return types.TipoError
}