package handlers

import (
	"backend-ventas/api/database"
	"backend-ventas/api/models"
	"fmt"      // Para fmt.Errorf
	"net/http" // Para http.Status...
	"strings"  // Para strings.Contains
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// GetVentas maneja la obtención de todas las ventas
func GetVentas(c *gin.Context) {
	var ventas []models.Venta
	if result := database.DB.Find(&ventas); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener ventas: " + result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, ventas)
}

// GetVentaByID maneja la obtención de una venta por su ID
func GetVentaByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := parseIDFromParam(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de venta inválido: " + err.Error()})
		return
	}

	var venta models.Venta
	result := database.DB.First(&venta, id)
	if result.Error == gorm.ErrRecordNotFound {
		c.JSON(http.StatusNotFound, gin.H{"error": "Venta no encontrada"})
		return
	}
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener venta: " + result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, venta)
}

// CreateVenta maneja la creación de una nueva venta a partir de un carrito
func CreateVenta(c *gin.Context) {
	var reqBody struct {
		IDCarrito uint `json:"id_carrito"`
	}

	if err := c.ShouldBindJSON(&reqBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "JSON de solicitud inválido: " + err.Error()})
		return
	}

	// Iniciar una transacción GORM para asegurar la atomicidad
	var newVenta models.Venta // Declarar la variable aquí para que sea accesible fuera del closure
	err := database.DB.Transaction(func(tx *gorm.DB) error {
		var carrito models.CarritoCompras
		result := tx.Preload("Detalles").First(&carrito, reqBody.IDCarrito)
		if result.Error == gorm.ErrRecordNotFound {
			return gorm.ErrRecordNotFound
		}
		if result.Error != nil {
			return result.Error
		}

		if carrito.Estado != "ABIERTO" {
			return fmt.Errorf("el carrito no está en un estado válido (%s) para generar una venta", carrito.Estado)
		}

		if len(carrito.Detalles) == 0 {
			return fmt.Errorf("el carrito está vacío, no se puede generar una venta")
		}

		var totalVenta float64 = 0.0
		for _, item := range carrito.Detalles {
			totalVenta += float64(item.Cantidad) * (item.PrecioUnitario - item.DescuentoAplicado)
		}

		newVenta = models.Venta{ // Asignar a la variable declarada fuera del closure
			IDCarrito: carrito.IDCarrito,
			Fecha:     time.Now(),
			Total:     totalVenta,
			Estado:    "PENDIENTE_PAGO",
		}
		if result := tx.Create(&newVenta); result.Error != nil {
			return result.Error
		}

		if result := tx.Model(&carrito).Update("Estado", "COMPLETADO"); result.Error != nil {
			return result.Error
		}

		return nil
	})

	if err == gorm.ErrRecordNotFound {
		c.JSON(http.StatusNotFound, gin.H{"error": "Carrito no encontrado"})
		return
	}
	if err != nil && (strings.Contains(err.Error(), "estado válido") || strings.Contains(err.Error(), "carrito está vacío")) {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al crear venta: " + err.Error()})
		return
	}

	// newVenta ya contendrá el ID asignado por GORM después de la transacción
	c.JSON(http.StatusCreated, gin.H{
		"message":     "Venta creada exitosamente",
		"id_venta":    newVenta.IDVenta,
		"total_venta": newVenta.Total,
	})
}

// UpdateVenta maneja la actualización de una venta existente
func UpdateVenta(c *gin.Context) {
	idStr := c.Param("id")
	id, err := parseIDFromParam(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de venta inválido: " + err.Error()})
		return
	}

	var updatedVentaData models.Venta
	if err := c.ShouldBindJSON(&updatedVentaData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "JSON de solicitud inválido: " + err.Error()})
		return
	}

	var existingVenta models.Venta
	result := database.DB.First(&existingVenta, id)
	if result.Error == gorm.ErrRecordNotFound {
		c.JSON(http.StatusNotFound, gin.H{"error": "Venta no encontrada"})
		return
	}
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al buscar venta: " + result.Error.Error()})
		return
	}

	updateMap := make(map[string]interface{})
	// Solo actualizar si el valor proporcionado no es el valor cero predeterminado para su tipo
	if updatedVentaData.Total != 0 {
		updateMap["total"] = updatedVentaData.Total
	}
	if updatedVentaData.Estado != "" {
		updateMap["estado"] = updatedVentaData.Estado
	}

	if len(updateMap) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No hay campos válidos para actualizar"})
		return
	}

	result = database.DB.Model(&existingVenta).Updates(updateMap)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al actualizar venta: " + result.Error.Error()})
		return
	}
	if result.RowsAffected == 0 {
		c.JSON(http.StatusOK, gin.H{"message": "No se realizó ninguna actualización (los datos ya son iguales o la venta no fue encontrada)"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Venta actualizada exitosamente"})
}

// DeleteVenta maneja la "eliminación suave" de una venta.
// Aquí, en lugar de eliminar el registro de la base de datos, actualizamos su estado a "Eliminada" o similar.
// Se debe implementear una columna para "eliminación suave".
func DeleteVenta(c *gin.Context) {
	idStr := c.Param("id")
	id, err := parseIDFromParam(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de venta inválido: " + err.Error()})
		return
	}

	// 1. Encontrar la venta primero para asegurarnos de que exista
	var venta models.Venta
	findResult := database.DB.First(&venta, id)
	if findResult.Error == gorm.ErrRecordNotFound {
		c.JSON(http.StatusNotFound, gin.H{"error": "Venta no encontrada."})
		return
	}
	if findResult.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al buscar la venta: " + findResult.Error.Error()})
		return
	}

	// 2. Realizar la "eliminación suave" actualizando el estado
	// Asignamos el nuevo estado para "eliminar" la venta lógicamente.
	venta.Estado = "Eliminada"

	// Guardar los cambios en la base de datos
	updateResult := database.DB.Save(&venta)
	if updateResult.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al marcar la venta como eliminada: " + updateResult.Error.Error()})
		return
	}

	// 3. Enviar una respuesta de éxito (200 OK o 204 No Content)
	c.JSON(http.StatusOK, gin.H{"message": "Venta eliminada exitosamente."})
}
