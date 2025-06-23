package handlers

import (
	"backend-ventas/api/database"
	"backend-ventas/api/models"
	"strconv"

	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetClientes(c *gin.Context) {
	var clientes []models.Cliente

	// --- Lógica de Paginación ---
	// Leer parámetros de query para paginación: /productos?page=1&limit=10
	pageStr := c.DefaultQuery("page", "1") // Gin's DefaultQuery para valor por defecto
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	limitStr := c.DefaultQuery("limit", "10") // Valor por defecto 10
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		limit = 10
	}
	offset := (page - 1) * limit

	query := database.DB.Model(&models.Cliente{})

	//--- Lógica de Filtrado (Ejemplos) ---
	//productos?nombre=JaneDoe&activo=true&categoria_id=1
	nombre := c.Query("nombre")
	if nombre != "" {
		// Búsqueda parcial por nombre
		query = query.Where("nombre ILIKE ?", "%"+nombre+"%")
	}

	activoStr := c.Query("activo")
	if activoStr != "" {
		esActivo, err := strconv.ParseBool(activoStr)
		if err == nil {
			query = query.Where("activo = ?", esActivo)
		}
	}

	// Ejecutar la consulta Find con paginación y filtros
	result := query.Offset(offset).Limit(limit).Find(&clientes)

	if result.Error != nil {
		log.Printf("Error al obtener clientes: %v", result.Error)
		// Gin usa c.JSON y c.AbortWithStatusJSON para detener la ejecución y enviar respuesta
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Error interno del servidor al obtener clientes"})
		return
	}

	if result.RowsAffected == 0 { // GORM devuelve result.RowsAffected = 0 y result.Error = nil si no encuentra registros
		c.JSON(http.StatusOK, []models.Cliente{}) // Devuelve un array vacío si no hay resultados
		return
	}

	c.JSON(http.StatusOK, clientes)
}
