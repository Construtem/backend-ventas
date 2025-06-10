package handlers

import (
	"fmt"
	"strconv"
)

// parseIDFromParam convierte un string de parámetro de Gin a uint.
// Esto es un helper si quieres consolidar la lógica de parsing de IDs de Gin.
func parseIDFromParam(idStr string) (uint, error) {
	if idStr == "" {
		return 0, fmt.Errorf("el ID de parámetro está vacío")
	}
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return 0, fmt.Errorf("el ID '%s' no es un número válido: %w", idStr, err)
	}
	return uint(id), nil
}
