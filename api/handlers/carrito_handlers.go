package handlers

import (
	"backend-ventas/api/database"
	"backend-ventas/api/models"
	"net/http" // Se usa para constantes como http.StatusBadRequest

	// Necesario para strconv.ParseUint
	"time"

	"github.com/gin-gonic/gin" // Importa Gin
	"gorm.io/gorm"
)

// CreateCarrito maneja la creación de un nuevo carrito de compras
func CreateCarrito(c *gin.Context) {
	var carrito models.CarritoCompras
	// c.ShouldBindJSON intenta mapear el cuerpo de la solicitud JSON al struct carrito.
	if err := c.ShouldBindJSON(&carrito); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "JSON de solicitud inválido: " + err.Error()})
		return
	}

	if carrito.Fecha.IsZero() {
		carrito.Fecha = time.Now()
	}
	if carrito.Estado == "" {
		carrito.Estado = "ABIERTO"
	}

	result := database.DB.Create(&carrito)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al crear carrito: " + result.Error.Error()})
		return
	}

	c.JSON(http.StatusCreated, carrito) // Gin maneja automáticamente el Content-Type: application/json
}

// GetCarritoByID maneja la obtención de un carrito por su ID, incluyendo sus detalles
func GetCarritoByID(c *gin.Context) {
	idStr := c.Param("id")             // Obtiene el parámetro 'id' de la URL de Gin
	id, err := parseIDFromParam(idStr) // O usa strconv.ParseUint(idStr, 10, 32) directamente aquí
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de carrito inválido: " + err.Error()})
		return
	}

	var carrito models.CarritoCompras
	result := database.DB.Preload("Detalles").First(&carrito, id)
	if result.Error == gorm.ErrRecordNotFound {
		c.JSON(http.StatusNotFound, gin.H{"error": "Carrito no encontrado"})
		return
	}
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener carrito: " + result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, carrito)
}

// AddItemToCarrito maneja la adición de un producto a un carrito específico
func AddItemToCarrito(c *gin.Context) {
	idCarritoStr := c.Param("id") // Obtiene el ID del carrito de la URL
	idCarrito, err := parseIDFromParam(idCarritoStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de carrito inválido en la ruta: " + err.Error()})
		return
	}

	var item models.DetalleCarrito
	if err := c.ShouldBindJSON(&item); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "JSON de solicitud inválido para el ítem: " + err.Error()})
		return
	}

	item.IDCarrito = idCarrito // Asegura que el ID del carrito del ítem coincida con el de la URL

	var carrito models.CarritoCompras
	res := database.DB.First(&carrito, idCarrito)
	if res.Error == gorm.ErrRecordNotFound {
		c.JSON(http.StatusNotFound, gin.H{"error": "Carrito no encontrado"})
		return
	}
	if res.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al verificar carrito: " + res.Error.Error()})
		return
	}
	if carrito.Estado != "ABIERTO" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "El carrito no está abierto para añadir ítems"})
		return
	}

	var existingItem models.DetalleCarrito
	result := database.DB.Where("id_carrito = ? AND id_producto = ?", item.IDCarrito, item.IDProducto).First(&existingItem)

	if result.Error == gorm.ErrRecordNotFound {
		if createResult := database.DB.Create(&item); createResult.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al añadir producto al carrito: " + createResult.Error.Error()})
			return
		}
		c.JSON(http.StatusCreated, item)
	} else if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al verificar producto en carrito: " + result.Error.Error()})
		return
	} else {
		// Actualizar cantidad si el producto ya existe en el carrito
		newQuantity := existingItem.Cantidad + item.Cantidad
		updateResult := database.DB.Model(&existingItem).Update("Cantidad", newQuantity)
		if updateResult.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al actualizar cantidad del producto en el carrito: " + updateResult.Error.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"message":            "Cantidad de producto actualizada en el carrito",
			"id_detalle_carrito": existingItem.IDDetalleCarrito,
			"nueva_cantidad":     newQuantity,
		})
	}
}

// UpdateItemInCarrito maneja la actualización de un ítem específico en el carrito
func UpdateItemInCarrito(c *gin.Context) {
	idCarritoStr := c.Param("id")
	idDetalleCarritoStr := c.Param("item_id") // Gin usa los nombres que definiste en la ruta (:item_id)

	idCarrito, err := parseIDFromParam(idCarritoStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de carrito inválido en la ruta: " + err.Error()})
		return
	}
	idDetalleCarrito, err := parseIDFromParam(idDetalleCarritoStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de ítem de carrito inválido en la ruta: " + err.Error()})
		return
	}

	var updatedItem models.DetalleCarrito
	if err := c.ShouldBindJSON(&updatedItem); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "JSON de solicitud inválido: " + err.Error()})
		return
	}

	var existingItem models.DetalleCarrito
	result := database.DB.Where("id_carrito = ?", idCarrito).First(&existingItem, idDetalleCarrito)
	if result.Error == gorm.ErrRecordNotFound {
		c.JSON(http.StatusNotFound, gin.H{"error": "Ítem de carrito no encontrado en este carrito"})
		return
	}
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al buscar ítem de carrito: " + result.Error.Error()})
		return
	}

	updateData := make(map[string]interface{})
	if updatedItem.Cantidad != 0 {
		updateData["cantidad"] = updatedItem.Cantidad
	}
	if updatedItem.PrecioUnitario != 0 {
		updateData["precio_unitario"] = updatedItem.PrecioUnitario
	}
	if updatedItem.DescuentoAplicado != 0 {
		updateData["descuento_aplicado"] = updatedItem.DescuentoAplicado
	}

	if len(updateData) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No hay campos válidos para actualizar"})
		return
	}

	result = database.DB.Model(&existingItem).Updates(updateData)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al actualizar ítem en carrito: " + result.Error.Error()})
		return
	}
	if result.RowsAffected == 0 {
		c.JSON(http.StatusOK, gin.H{"message": "No se realizó ninguna actualización (los datos ya son iguales o el ítem no fue encontrado)"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Ítem de carrito actualizado exitosamente"})
}

// DeleteItemFromCarrito maneja la eliminación de un ítem específico del carrito
func DeleteItemFromCarrito(c *gin.Context) {
	idCarritoStr := c.Param("id")
	idDetalleCarritoStr := c.Param("item_id")

	idCarrito, err := parseIDFromParam(idCarritoStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de carrito inválido en la ruta: " + err.Error()})
		return
	}
	idDetalleCarrito, err := parseIDFromParam(idDetalleCarritoStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de ítem de carrito inválido en la ruta: " + err.Error()})
		return
	}

	result := database.DB.Where("id_carrito = ?", idCarrito).Delete(&models.DetalleCarrito{}, idDetalleCarrito)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al eliminar ítem del carrito: " + result.Error.Error()})
		return
	}
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Ítem de carrito no encontrado en este carrito"})
		return
	}

	c.Status(http.StatusNoContent) // Para 204 No Content, Gin usa c.Status()
}
