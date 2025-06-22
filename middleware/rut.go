package middleware

import (
	"strconv"
	"strings"
)

func VerificacionRut(rut string) bool {
	rut = strings.ToUpper(strings.ReplaceAll(strings.ReplaceAll(rut, ".", ""), "-", ""))
	if len(rut) < 2 {
		return false
	}

	num, dv := rut[:len(rut)-1], rut[len(rut)-1:]

	numInt, err := strconv.Atoi(num)
	if err != nil {
		return false
	}

	var suma, factor = 0, 2
	for numInt > 0 {
		d := numInt % 10
		suma += d * factor
		factor++
		if factor > 7 {
			factor = 2
		}
		numInt /= 10
	}

	dvEsperado := 11 - (suma % 11)
	var dvCalculado string
	if dvEsperado == 11 {
		dvCalculado = "0"
	} else if dvEsperado == 10 {
		dvCalculado = "K"
	} else {
		dvCalculado = strconv.Itoa(dvEsperado)
	}

	return dv == dvCalculado
}
