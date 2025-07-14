package handlers

import (
	"crypto/rand"
	"encoding/hex"
	"log"
	"net/http"
	"strconv"

	"backend-ventas/api/controllers"
	"backend-ventas/api/dtos"
	"backend-ventas/api/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// GET /cotizaciones - Obtener todas las cotizaciones (simplificadas)
func ObtenerCotizacionesSimplificadas(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		cotizaciones, err := controllers.ListarCotizacionesSimples()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener cotizaciones"})
			return
		}
		// Calcular totales y armar respuesta simplificada
		var result []dtos.CotizacionSimplificadaResponse
		for _, cotizacion := range cotizaciones {
			var totalItems int
			var totalPrecio float64
			for _, item := range cotizacion.Items {
				if item.Producto != nil {
					totalItems += item.Cantidad
					totalPrecio += float64(item.Cantidad) * item.Producto.Precio
				}
			}
			totalPrecio += cotizacion.CostoEnvio
			// Al armar la respuesta simplificada, incluir el ID
			simplificada := dtos.CotizacionSimplificadaResponse{
				ID:           cotizacion.ID,
				FechaCrea:    cotizacion.FechaCrea,
				Estado:       cotizacion.Estado,
				CostoEnvio:   cotizacion.CostoEnvio,
				UserID:       cotizacion.UserID,
				Nombre:       "",
				TipoDespacho: cotizacion.TipoDespacho,
				TotalItems:   totalItems,
				TotalPrecio:  totalPrecio,
			}
			if cotizacion.Usuario != nil {
				simplificada.Nombre = cotizacion.Usuario.Nombre
			}
			if cotizacion.Cliente != nil {
				simplificada.Cliente.Nombre = cotizacion.Cliente.Nombre
				if cotizacion.Cliente.Telefono != nil {
					simplificada.Cliente.Telefono = *cotizacion.Cliente.Telefono
				}
				if cotizacion.Cliente.Email != nil {
					simplificada.Cliente.Email = *cotizacion.Cliente.Email
				}
				simplificada.Cliente.Rut = cotizacion.Cliente.Rut
				if cotizacion.Cliente.RazonSocial != nil {
					simplificada.Cliente.RazonSocial = *cotizacion.Cliente.RazonSocial
				}
			}
			for _, item := range cotizacion.Items {
				if item.Producto != nil {
					simplificada.Items = append(simplificada.Items, struct {
						SKU      string `json:"sku"`
						Nombre   string `json:"nombre"`
						Cantidad int    `json:"cantidad"`
					}{
						SKU:      item.Producto.SKU,
						Nombre:   item.Producto.Nombre,
						Cantidad: item.Cantidad,
					})
				}
			}
			result = append(result, simplificada)
		}
		c.JSON(http.StatusOK, result)
	}
}

// GET /cotizaciones/:id - Obtener detalle simplificado
func ObtenerCotizacionSimplificada(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		cotizacionID, err := strconv.Atoi(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "ID de cotización inválido"})
			return
		}
		cotizacion, err := controllers.ObtenerCotizacionSimplePorID(cotizacionID)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Cotización no encontrada"})
			return
		}
		var totalItems int
		var totalPrecio float64
		for _, item := range cotizacion.Items {
			if item.Producto != nil {
				totalItems += item.Cantidad
				totalPrecio += float64(item.Cantidad) * item.Producto.Precio
			}
		}
		totalPrecio += cotizacion.CostoEnvio
		response := dtos.CotizacionSimplificadaResponse{
			ID:           cotizacion.ID,
			FechaCrea:    cotizacion.FechaCrea,
			Estado:       cotizacion.Estado,
			CostoEnvio:   cotizacion.CostoEnvio,
			UserID:       cotizacion.UserID,
			Nombre:       "",
			TipoDespacho: cotizacion.TipoDespacho,
			TotalItems:   totalItems,
			TotalPrecio:  totalPrecio,
		}
		if cotizacion.Usuario != nil {
			response.Nombre = cotizacion.Usuario.Nombre
		}
		if cotizacion.Cliente != nil {
			response.Cliente.Nombre = cotizacion.Cliente.Nombre
			if cotizacion.Cliente.Telefono != nil {
				response.Cliente.Telefono = *cotizacion.Cliente.Telefono
			}
			if cotizacion.Cliente.Email != nil {
				response.Cliente.Email = *cotizacion.Cliente.Email
			}
			response.Cliente.Rut = cotizacion.Cliente.Rut
			if cotizacion.Cliente.RazonSocial != nil {
				response.Cliente.RazonSocial = *cotizacion.Cliente.RazonSocial
			}
		}
		for _, item := range cotizacion.Items {
			if item.Producto != nil {
				response.Items = append(response.Items, struct {
					SKU      string `json:"sku"`
					Nombre   string `json:"nombre"`
					Cantidad int    `json:"cantidad"`
				}{
					SKU:      item.Producto.SKU,
					Nombre:   item.Producto.Nombre,
					Cantidad: item.Cantidad,
				})
			}
		}
		c.JSON(http.StatusOK, response)
	}
}

// GET /cotizaciones/completas - Obtener todas las cotizaciones completas
func ObtenerCotizaciones(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		cotizaciones, err := controllers.ListarCotizacionesCompletas()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener cotizaciones"})
			return
		}
		var response []dtos.CotizacionResponse
		for _, cotizacion := range cotizaciones {
			var totalItems int
			var totalPrecio float64
			for _, item := range cotizacion.Items {
				totalItems += item.Cantidad
				if item.Producto != nil {
					totalPrecio += float64(item.Cantidad) * item.Producto.Precio
				} else {
				// Opcional: loguea o maneja el error, por ejemplo:
					log.Printf("Producto nil en item: %+v", item)
				}
			}
			totalPrecio += cotizacion.CostoEnvio
			cotizacionResponse := dtos.CotizacionResponse{
				ID:           cotizacion.ID,
				FechaCrea:    cotizacion.FechaCrea,
				Estado:       cotizacion.Estado,
				CostoEnvio:   cotizacion.CostoEnvio,
				RutCliente:   cotizacion.RutCliente,
				UserID:       cotizacion.UserID,
				TipoDespacho: cotizacion.TipoDespacho,
				Total:        cotizacion.Total,
				Cliente:      cotizacion.Cliente,
				TotalItems:   totalItems,
				TotalPrecio:  totalPrecio,
			}
			for _, item := range cotizacion.Items {
				cotizacionResponse.Items = append(cotizacionResponse.Items, dtos.CotizacionItemResponse{
					CotizacionID: item.CotizacionID,
					ProductoID:   item.ProductoID,
					SucursalID:   item.SucursalID,
					Cantidad:     item.Cantidad,
					Producto:     item.Producto,
					Sucursal:     item.Sucursal,
				})
			}
			response = append(response, cotizacionResponse)
		}
		c.JSON(http.StatusOK, response)
	}
}

// GET /cotizaciones/completa/:id - Obtener detalle completo
func ObtenerCotizacionPorID(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		cotizacionID, err := strconv.Atoi(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "ID de cotización inválido"})
			return
		}
		cotizacion, err := controllers.ObtenerCotizacionCompletaPorID(cotizacionID)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Cotización no encontrada"})
			return
		}
		var totalItems int
		var totalPrecio float64
		for _, item := range cotizacion.Items {
			if item.Producto != nil {
				totalItems += item.Cantidad
				totalPrecio += float64(item.Cantidad) * item.Producto.Precio
			}
		}
		totalPrecio += cotizacion.CostoEnvio
		response := dtos.CotizacionResponse{
			ID:           cotizacion.ID,
			FechaCrea:    cotizacion.FechaCrea,
			Estado:       cotizacion.Estado,
			CostoEnvio:   cotizacion.CostoEnvio,
			RutCliente:   cotizacion.RutCliente,
			UserID:       cotizacion.UserID,
			TipoDespacho: cotizacion.TipoDespacho,
			Total:        cotizacion.Total,
			Cliente:      cotizacion.Cliente,
			Usuario:      cotizacion.Usuario,
			TotalItems:   totalItems,
			TotalPrecio:  totalPrecio,
		}
		for _, item := range cotizacion.Items {
			response.Items = append(response.Items, dtos.CotizacionItemResponse{
				CotizacionID: item.CotizacionID,
				ProductoID:   item.ProductoID,
				SucursalID:   item.SucursalID,
				Cantidad:     item.Cantidad,
				Producto:     item.Producto,
				Sucursal:     item.Sucursal,
			})
		}
		c.JSON(http.StatusOK, response)
	}
}

// GET /cotizaciones/:id/items/simples - Obtener items simples de una cotización
func ObtenerItemsSimplesCotizacion(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		cotizacionID, err := strconv.Atoi(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "ID de cotización inválido"})
			return
		}
		items, err := controllers.ObtenerItemsSimplesCotizacion(cotizacionID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener items"})
			return
		}
		type ItemSimple struct {
			SKU         string `json:"sku"`
			Nombre      string `json:"nombre"`
			Descripcion string `json:"descripcion"`
			Marca       string `json:"marca"`
			Cantidad    int    `json:"cantidad"`
			Sucursal    string `json:"sucursal"`
		}
		var result []ItemSimple
		for _, item := range items {
			itemSimple := ItemSimple{
				SKU:         item.Producto.SKU,
				Nombre:      item.Producto.Nombre,
				Descripcion: item.Producto.Descripcion,
				Marca:       item.Producto.Marca,
				Cantidad:    item.Cantidad,
				Sucursal:    item.Sucursal.Nombre,
			}
			result = append(result, itemSimple)
		}
		c.JSON(http.StatusOK, result)
	}
}

// GET /cotizaciones/:id/items - Obtener items completos de una cotización
func ObtenerItemsCotizacion(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		cotizacionID, err := strconv.Atoi(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "ID de cotización inválido"})
			return
		}
		items, err := controllers.ObtenerItemsCompletosCotizacion(cotizacionID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener items"})
			return
		}
		var response []dtos.CotizacionItemResponse
		for _, item := range items {
			response = append(response, dtos.CotizacionItemResponse{
				CotizacionID: item.CotizacionID,
				ProductoID:   item.ProductoID,
				SucursalID:   item.SucursalID,
				Cantidad:     item.Cantidad,
				Producto:     item.Producto,
				Sucursal:     item.Sucursal,
			})
		}
		c.JSON(http.StatusOK, response)
	}
}

// POST /cotizaciones - Crear nueva cotización
func CrearCotizacion(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var request dtos.CreateCotizacionRequest
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Datos de entrada inválidos", "details": err.Error()})
			return
		}

		// Verificar que el cliente existe
		var cliente models.Cliente
		if err := db.Where("rut = ?", request.RutCliente).First(&cliente).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Cliente no encontrado"})
			return
		}

		// Verificar que el usuario existe
		var usuario models.Usuario
		if err := db.Where("email = ?", request.UserID).First(&usuario).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Usuario no encontrado"})
			return
		}

		cotizacion, err := controllers.CrearCotizacion(request.RutCliente, request.UserID, request.TipoDespacho, request.CostoEnvio)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al crear cotización"})
			return
		}

		response := dtos.CreateCotizacionResponse{
			ID:        cotizacion.ID,
			FechaCrea: cotizacion.FechaCrea,
			Estado:    cotizacion.Estado,
			Mensaje:   "Cotización creada exitosamente",
		}
		c.JSON(http.StatusCreated, response)
	}
}

// POST /cotizaciones/:id/items - Agregar item a cotización
func AgregarItemCotizacion(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		cotizacionID, err := strconv.Atoi(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "ID de cotización inválido"})
			return
		}
		var request dtos.AddItemRequest
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Datos de entrada inválidos", "details": err.Error()})
			return
		}
		// Verificar que la cotización existe
		var cotizacion models.Cotizacion
		if err := db.First(&cotizacion, cotizacionID).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Cotización no encontrada"})
			return
		}
		// Verificar que el producto existe
		var producto models.Producto
		if err := db.Where("sku = ?", request.ProductoID).First(&producto).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Producto no encontrado"})
			return
		}
		// Verificar que la sucursal existe
		var sucursal models.Sucursal
		if err := db.First(&sucursal, request.SucursalID).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Sucursal no encontrada"})
			return
		}
		// Verificar si ya existe el item
		var existingItem models.CotizacionItem
		err = db.Where("cotizacion_id = ? AND sku = ? AND sucursal_id = ?", cotizacionID, request.ProductoID, request.SucursalID).First(&existingItem).Error
		if err == nil {
			// El item ya existe, actualizar cantidad
			existingItem.Cantidad += request.Cantidad
			if err := db.Save(&existingItem).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al actualizar item"})
				return
			}
			response := dtos.AddItemResponse{
				CotizacionID: existingItem.CotizacionID,
				ProductoID:   existingItem.ProductoID,
				SucursalID:   existingItem.SucursalID,
				Cantidad:     existingItem.Cantidad,
				Mensaje:      "Cantidad actualizada exitosamente",
			}
			c.JSON(http.StatusOK, response)
			return
		}
		// Crear nuevo item usando el controlador
		item, err := controllers.AgregarItemCotizacion(cotizacionID, request.ProductoID, request.SucursalID, request.Cantidad)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al agregar item"})
			return
		}
		response := dtos.AddItemResponse{
			CotizacionID: item.CotizacionID,
			ProductoID:   item.ProductoID,
			SucursalID:   item.SucursalID,
			Cantidad:     item.Cantidad,
			Mensaje:      "Item agregado exitosamente",
		}
		c.JSON(http.StatusCreated, response)
	}
}

// PUT /cotizaciones/:id - Editar cotización
func EditarCotizacion(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		cotizacionID, err := strconv.Atoi(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "ID de cotización inválido"})
			return
		}
		var request dtos.UpdateCotizacionRequest
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Datos de entrada inválidos", "details": err.Error()})
			return
		}
		// Actualizar cotización usando el controlador
		_, err = controllers.ActualizarCotizacion(cotizacionID, request.Estado, request.CostoEnvio, request.TipoDespacho, request.Total)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al actualizar cotización"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"mensaje": "Cotización actualizada exitosamente"})
	}
}

// POST /cotizaciones/:id/preview - Crear preview de cotización
func CrearPreviewCotizacion(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		cotizacionID, err := strconv.Atoi(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "ID de cotización inválido"})
			return
		}
		var request dtos.PreviewCotizacionRequest
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Datos de entrada inválidos", "details": err.Error()})
			return
		}
		// Verificar que la cotización existe
		var cotizacion models.Cotizacion
		if err := db.Preload("Items.Producto").Preload("Items.Sucursal").First(&cotizacion, cotizacionID).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Cotización no encontrada"})
			return
		}
		// Calcular totales
		var subtotal float64
		for _, item := range cotizacion.Items {
			subtotal += float64(item.Cantidad) * item.Producto.Precio
		}
		tax := subtotal * 0.19 // IVA 19%
		total := subtotal + tax + cotizacion.CostoEnvio
		preview, err := controllers.CrearPreviewCotizacion(&cotizacionID, request.IssuedAt, subtotal, tax, total)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al crear preview"})
			return
		}
		response := dtos.PreviewCotizacionResponse{
			ID:            preview.ID,
			CotizacionID:  preview.CotizacionID,
			IssuedAt:      preview.IssuedAt,
			Subtotal:      preview.Subtotal,
			Tax:           preview.Tax,
			Total:         preview.Total,
			PaymentStatus: preview.PaymentStatus,
			StatusPagado:  preview.StatusPagado,
			CreatedAt:     preview.CreatedAt,
			UpdatedAt:     preview.UpdatedAt,
			Mensaje:       "Preview creado exitosamente",
		}
		c.JSON(http.StatusCreated, response)
	}
}

// DELETE /cotizaciones/:id/items/:producto_id/:sucursal_id - Eliminar item de cotización
func EliminarItemCotizacion(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		productoID := c.Param("producto_id")
		sucursalID := c.Param("sucursal_id")
		cotizacionID, err := strconv.Atoi(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "ID de cotización inválido"})
			return
		}
		sucursalIDInt, err := strconv.Atoi(sucursalID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "ID de sucursal inválido"})
			return
		}
		// Eliminar el item usando el controlador
		err = controllers.EliminarItemCotizacion(cotizacionID, productoID, sucursalIDInt)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al eliminar item"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"mensaje": "Item eliminado exitosamente"})
	}
}

// Función auxiliar para generar token
func generateToken() string {
	bytes := make([]byte, 16)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}
