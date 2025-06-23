package handlers

import (
	"backend-ventas/api/database"
	"backend-ventas/api/models"
	"fmt"

	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// GetProductos maneja la ruta GET /productos
// Permite filtrar y paginar productos.
func GetUsuarios(c *gin.Context) { // En Gin, el contexto es *gin.Context
	var usuarios []models.Usuario

	// --- Lógica de Paginación ---
	// Leer parámetros de query para paginación: /productos?page=1&limit=10
	// pageStr := c.DefaultQuery("page", "1") // Gin's DefaultQuery para valor por defecto
	// page, err := strconv.Atoi(pageStr)
	// if err != nil || page < 1 {
	// 	page = 1
	// }

	// limitStr := c.DefaultQuery("limit", "10") // Valor por defecto 10
	// limit, err := strconv.Atoi(limitStr)
	// if err != nil || limit < 1 {
	// 	limit = 10
	// }
	// offset := (page - 1) * limit

	query := database.DB.Model(&models.Usuario{}) // Crea una consulta base

	//--- Lógica de Filtrado (Ejemplos) ---
	//productos?nombre=JaneDoe&activo=true&categoria_id=1
	nombre := c.Query("nombre") // Gin's Query para obtener un parámetro de URL
	if nombre != "" {
		// Búsqueda parcial por nombre
		query = query.Where("nombre ILIKE ?", "%"+nombre+"%")
	}

	// activoStr := c.Query("activo")
	// if activoStr != "" {
	// 	esActivo, err := strconv.ParseBool(activoStr)
	// 	if err == nil {
	// 		query = query.Where("activo = ?", esActivo)
	// 	}
	// }

	// categoriaIDStr := c.Query("categoria_id")
	// if categoriaIDStr != "" {
	// 	id, err := strconv.Atoi(categoriaIDStr)
	// 	if err == nil {
	// 		query = query.Where("categoria_id = ?", id)
	// 	}
	// }

	// Ejecutar la consulta Find con paginación y filtros
	//result := query.Offset(offset).Limit(limit).Find(&productos)
	result := query.Unscoped().Find(&usuarios)

	if result.Error != nil {
		log.Printf("Error al obtener productos: %v", result.Error)
		// Gin usa c.JSON y c.AbortWithStatusJSON para detener la ejecución y enviar respuesta
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Error interno del servidor al obtener productos"})
		return
	}

	if result.RowsAffected == 0 { // GORM devuelve result.RowsAffected = 0 y result.Error = nil si no encuentra registros
		c.JSON(http.StatusOK, []models.Usuario{}) // Devuelve un array vacío si no hay resultados
		return
	}

	c.JSON(http.StatusOK, usuarios)
}

func GetUsuarioByID(c *gin.Context) {
	id := c.Param("id") // Obtener el ID de la URL

	if id == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "ID de usuario inválido"})
		return
	}

	var usuario models.Usuario
	// Usamos First para esperar un solo registro y manejar gorm.ErrRecordNotFound
	result := database.DB.Unscoped().First(&usuario, id)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": fmt.Sprintf("Usuario con ID %s no encontrado", id)})
			return
		}
		log.Printf("Error al obtener usuario por ID %s: %v", id, result.Error)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Error interno del servidor al obtener el usuario"})
		return
	}

	c.JSON(http.StatusOK, usuario)
}
