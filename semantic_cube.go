package main

type claveCubo struct {
	izquierda Tipo
	derecha Tipo
	operador int
}

var cuboSemantico = map[claveCubo]Tipo {
	// Izquierda: entero, derecha: entero
	{izquierda: TipoConstante, derecha:TipoConstante, operador: MAS}: TipoConstante,
	{izquierda: TipoConstante, derecha:TipoConstante, operador: MENOS}: TipoConstante,
	{izquierda: TipoConstante, derecha:TipoConstante, operador: ENTRE}: TipoConstante,
	{izquierda: TipoConstante, derecha:TipoConstante, operador: POR}: TipoConstante,
	{izquierda: TipoConstante, derecha:TipoConstante, operador: MAYOR}: TipoConstante,
	{izquierda: TipoConstante, derecha:TipoConstante, operador: MENOR}: TipoConstante,
	{izquierda: TipoConstante, derecha:TipoConstante, operador: IGUAL}: TipoConstante,
	{izquierda: TipoConstante, derecha:TipoConstante, operador: DIFERENTE}: TipoConstante,
	{izquierda: TipoConstante, derecha:TipoConstante, operador: ASIGNAVAR}: TipoConstante,
	// Izquierda: entero, derecha: flotante
	{izquierda: TipoConstante, derecha:TipoFlotante, operador: MAS}: TipoFlotante,
	{izquierda: TipoConstante, derecha:TipoFlotante, operador: MENOS}: TipoFlotante,
	{izquierda: TipoConstante, derecha:TipoFlotante, operador: ENTRE}: TipoFlotante,
	{izquierda: TipoConstante, derecha:TipoFlotante, operador: POR}: TipoFlotante,
	{izquierda: TipoConstante, derecha:TipoFlotante, operador: MAYOR}: TipoConstante,
	{izquierda: TipoConstante, derecha:TipoFlotante, operador: MENOR}: TipoConstante,
	{izquierda: TipoConstante, derecha:TipoFlotante, operador: IGUAL}: TipoConstante,
	{izquierda: TipoConstante, derecha:TipoFlotante, operador: DIFERENTE}: TipoConstante,
	{izquierda: TipoConstante, derecha:TipoFlotante, operador: ASIGNAVAR}: TipoError,
	// Izquierda: flotante, derecha: entero
	{izquierda: TipoConstante, derecha:TipoFlotante, operador: MAS}: TipoFlotante,
	{izquierda: TipoConstante, derecha:TipoFlotante, operador: MENOS}: TipoFlotante,
	{izquierda: TipoConstante, derecha:TipoFlotante, operador: ENTRE}: TipoFlotante,
	{izquierda: TipoConstante, derecha:TipoFlotante, operador: POR}: TipoFlotante,
	{izquierda: TipoConstante, derecha:TipoFlotante, operador: MAYOR}: TipoConstante,
	{izquierda: TipoConstante, derecha:TipoFlotante, operador: MENOR}: TipoConstante,
	{izquierda: TipoConstante, derecha:TipoFlotante, operador: IGUAL}: TipoConstante,
	{izquierda: TipoConstante, derecha:TipoFlotante, operador: DIFERENTE}: TipoConstante,
	{izquierda: TipoConstante, derecha:TipoFlotante, operador: ASIGNAVAR}: TipoFlotante,
	// Izquierda: flotante, derecha: flotante
	{izquierda: TipoConstante, derecha:TipoFlotante, operador: MAS}: TipoFlotante,
	{izquierda: TipoConstante, derecha:TipoFlotante, operador: MENOS}: TipoFlotante,
	{izquierda: TipoConstante, derecha:TipoFlotante, operador: ENTRE}: TipoFlotante,
	{izquierda: TipoConstante, derecha:TipoFlotante, operador: POR}: TipoFlotante,
	{izquierda: TipoConstante, derecha:TipoFlotante, operador: MAYOR}: TipoConstante,
	{izquierda: TipoConstante, derecha:TipoFlotante, operador: MENOR}: TipoConstante,
	{izquierda: TipoConstante, derecha:TipoFlotante, operador: IGUAL}: TipoConstante,
	{izquierda: TipoConstante, derecha:TipoFlotante, operador: DIFERENTE}: TipoConstante,
	{izquierda: TipoConstante, derecha:TipoFlotante, operador: ASIGNAVAR}: TipoFlotante,

}

func validarSemantica(izq, der Tipo, op int) Tipo {
	if tipo, ok := cuboSemantico[claveCubo{izq, der, op}]; ok  {
		return tipo
	}
	return TipoError
}