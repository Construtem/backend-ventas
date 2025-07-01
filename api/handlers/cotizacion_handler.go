package handlers

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"
	"strconv"
	"time"

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
				totalItems += item.Cantidad
				totalPrecio += float64(item.Cantidad) * item.Producto.Precio
			}
			totalPrecio += cotizacion.CostoEnvio
			simplificada := dtos.CotizacionSimplificadaResponse{
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
				simplificada.Cliente.Telefono = cotizacion.Cliente.Telefono
				simplificada.Cliente.Email = cotizacion.Cliente.Email
				simplificada.Cliente.Rut = cotizacion.Cliente.Rut
				simplificada.Cliente.RazonSocial = cotizacion.Cliente.RazonSocial
			}
			for _, item := range cotizacion.Items {
				simplificada.Items = append(simplificada.Items, struct {
					SKU      uint   `json:"sku"`
					Nombre   string `json:"nombre"`
					Cantidad int    `json:"cantidad"`
				}{
					SKU:      item.Producto.SKU,
					Nombre:   item.Producto.Nombre,
					Cantidad: item.Cantidad,
				})
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
		cotizacionID, err := strconv.ParseUint(id, 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "ID de cotización inválido"})
			return
		}
		cotizacion, err := controllers.ObtenerCotizacionSimplePorID(uint(cotizacionID))
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Cotización no encontrada"})
			return
		}
		var totalItems int
		var totalPrecio float64
		for _, item := range cotizacion.Items {
			totalItems += item.Cantidad
			totalPrecio += float64(item.Cantidad) * item.Producto.Precio
		}
		totalPrecio += cotizacion.CostoEnvio
		response := dtos.CotizacionSimplificadaResponse{
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
			response.Cliente.Telefono = cotizacion.Cliente.Telefono
			response.Cliente.Email = cotizacion.Cliente.Email
			response.Cliente.Rut = cotizacion.Cliente.Rut
			response.Cliente.RazonSocial = cotizacion.Cliente.RazonSocial
		}
		for _, item := range cotizacion.Items {
			response.Items = append(response.Items, struct {
				SKU      uint   `json:"sku"`
				Nombre   string `json:"nombre"`
				Cantidad int    `json:"cantidad"`
			}{
				SKU:      item.Producto.SKU,
				Nombre:   item.Producto.Nombre,
				Cantidad: item.Cantidad,
			})
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
				totalPrecio += float64(item.Cantidad) * item.Producto.Precio
			}
			totalPrecio += cotizacion.CostoEnvio
			cotizacionResponse := dtos.CotizacionResponse{
				ID:           cotizacion.ID,
				FechaCrea:    cotizacion.FechaCrea,
				Estado:       cotizacion.Estado,
				CostoEnvio:   cotizacion.CostoEnvio,
				ClienteID:    cotizacion.ClienteID,
				UserID:       cotizacion.UserID,
				TipoDespacho: cotizacion.TipoDespacho,
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
		cotizacionID, err := strconv.ParseUint(id, 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "ID de cotización inválido"})
			return
		}
		cotizacion, err := controllers.ObtenerCotizacionCompletaPorID(uint(cotizacionID))
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Cotización no encontrada"})
			return
		}
		var totalItems int
		var totalPrecio float64
		for _, item := range cotizacion.Items {
			totalItems += item.Cantidad
			totalPrecio += float64(item.Cantidad) * item.Producto.Precio
		}
		totalPrecio += cotizacion.CostoEnvio
		response := dtos.CotizacionResponse{
			ID:           cotizacion.ID,
			FechaCrea:    cotizacion.FechaCrea,
			Estado:       cotizacion.Estado,
			CostoEnvio:   cotizacion.CostoEnvio,
			ClienteID:    cotizacion.ClienteID,
			UserID:       cotizacion.UserID,
			TipoDespacho: cotizacion.TipoDespacho,
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

// GET /cotizaciones/:id/items/simple - Obtener items simples
func ObtenerItemsSimplesCotizacion(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		cotizacionID, err := strconv.ParseUint(id, 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "ID de cotización inválido"})
			return
		}
		items, err := controllers.ObtenerItemsSimplesCotizacion(uint(cotizacionID))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener items de la cotización"})
			return
		}
		type ItemSimple struct {
			SKU         uint   `json:"sku"`
			Nombre      string `json:"nombre"`
			Descripcion string `json:"descripcion"`
			Marca       string `json:"marca"`
			Cantidad    int    `json:"cantidad"`
			Sucursal    string `json:"sucursal"`
		}
		var itemsResponse []ItemSimple
		for _, item := range items {
			itemsResponse = append(itemsResponse, ItemSimple{
				SKU:         item.Producto.SKU,
				Nombre:      item.Producto.Nombre,
				Descripcion: item.Producto.Descripcion,
				Marca:       item.Producto.Marca,
				Cantidad:    item.Cantidad,
				Sucursal:    item.Sucursal.Nombre,
			})
		}
		c.JSON(http.StatusOK, gin.H{
			"cotizacion_id": cotizacionID,
			"items":         itemsResponse,
			"total_items":   len(itemsResponse),
		})
	}
}

// GET /cotizaciones/:id/items - Obtener items completos
func ObtenerItemsCotizacion(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		cotizacionID, err := strconv.ParseUint(id, 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "ID de cotización inválido"})
			return
		}
		items, err := controllers.ObtenerItemsCompletosCotizacion(uint(cotizacionID))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener items de la cotización"})
			return
		}
		var itemsResponse []dtos.CotizacionItemResponse
		for _, item := range items {
			itemsResponse = append(itemsResponse, dtos.CotizacionItemResponse{
				CotizacionID: item.CotizacionID,
				ProductoID:   item.ProductoID,
				SucursalID:   item.SucursalID,
				Cantidad:     item.Cantidad,
				Producto:     item.Producto,
				Sucursal:     item.Sucursal,
			})
		}
		c.JSON(http.StatusOK, gin.H{
			"cotizacion_id": cotizacionID,
			"items":         itemsResponse,
			"total_items":   len(itemsResponse),
		})
	}
}

// POST /cotizaciones - Crear una cotización
func CrearCotizacion(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var request dtos.CreateCotizacionRequest
		if err := c.BindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "JSON inválido"})
			fmt.Print("\n\t\t<<<< JSON INVALIDO AL CREAR COTIZACION >>>>\n")
			return
		}

		// Validar que el cliente existe
		var cliente models.Cliente
		if err := db.First(&cliente, request.ClienteID).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Cliente no encontrado"})
			fmt.Printf("\n\t\t<<<< CLIENTE ID %d NO ENCONTRADO >>>>\n", request.ClienteID)
			return
		}

		// Validar que el usuario existe
		var usuario models.Usuario
		if err := db.Where("email = ?", request.UserID).First(&usuario).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Usuario no encontrado"})
			fmt.Printf("\n\t\t<<<< USUARIO %s NO ENCONTRADO >>>>\n", request.UserID)
			return
		}

		// Crear la cotización
		cotizacion := models.Cotizacion{
			ClienteID:    request.ClienteID,
			UserID:       request.UserID,
			TipoDespacho: request.TipoDespacho,
			CostoEnvio:   request.CostoEnvio,
			Estado:       "PENDIENTE", // Estado por defecto
		}

		if err := db.Create(&cotizacion).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al crear cotización"})
			fmt.Print("\n\t\t<<<< ERROR AL CREAR COTIZACION >>>>\n")
			return
		}

		response := dtos.CreateCotizacionResponse{
			ID:        cotizacion.ID,
			FechaCrea: cotizacion.FechaCrea,
			Estado:    cotizacion.Estado,
			Mensaje:   "Cotización creada exitosamente",
		}

		c.JSON(http.StatusCreated, response)
		fmt.Printf("\n\t\t<<<< COTIZACION ID %d CREADA CON EXITO >>>>\n", cotizacion.ID)
	}
}

// POST /cotizaciones/[id]/items - Agregar un producto a una cotización
func AgregarItemCotizacion(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		cotizacionID, err := strconv.ParseUint(id, 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "ID de cotización inválido"})
			return
		}

		var request dtos.AddItemRequest
		if err := c.BindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "JSON inválido"})
			fmt.Print("\n\t\t<<<< JSON INVALIDO AL AGREGAR ITEM >>>>\n")
			return
		}

		// Verificar que la cotización existe
		var cotizacion models.Cotizacion
		if err := db.First(&cotizacion, cotizacionID).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Cotización no encontrada"})
			fmt.Printf("\n\t\t<<<< COTIZACION ID %s NO ENCONTRADA >>>>\n", id)
			return
		}

		// Verificar que el producto existe
		var producto models.Producto
		if err := db.First(&producto, request.ProductoID).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Producto no encontrado"})
			fmt.Printf("\n\t\t<<<< PRODUCTO ID %d NO ENCONTRADO >>>>\n", request.ProductoID)
			return
		}

		// Verificar que la sucursal existe
		var sucursal models.Sucursal
		if err := db.First(&sucursal, request.SucursalID).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Sucursal no encontrada"})
			fmt.Printf("\n\t\t<<<< SUCURSAL ID %d NO ENCONTRADA >>>>\n", request.SucursalID)
			return
		}

		// Verificar stock disponible
		var stock models.StockSucursal
		if err := db.Where("sucursal_id = ? AND producto_id = ?", request.SucursalID, request.ProductoID).
			First(&stock).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Producto no disponible en esta sucursal"})
			fmt.Printf("\n\t\t<<<< PRODUCTO NO DISPONIBLE EN SUCURSAL >>>>\n")
			return
		}

		if stock.Cantidad < request.Cantidad {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Stock insuficiente"})
			fmt.Printf("\n\t\t<<<< STOCK INSUFICIENTE >>>>\n")
			return
		}

		// Verificar si el item ya existe en la cotización
		var existingItem models.CotizacionItem
		if err := db.Where("cotizacion_id = ? AND producto_id = ? AND sucursal_id = ?",
			cotizacionID, request.ProductoID, request.SucursalID).First(&existingItem).Error; err == nil {
			// El item ya existe, actualizar cantidad
			newCantidad := existingItem.Cantidad + request.Cantidad
			if newCantidad > stock.Cantidad {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Stock insuficiente para la cantidad total"})
				return
			}

			if err := db.Model(&existingItem).Update("cantidad", newCantidad).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al actualizar cantidad"})
				return
			}

			response := dtos.AddItemResponse{
				CotizacionID: uint(cotizacionID),
				ProductoID:   request.ProductoID,
				SucursalID:   request.SucursalID,
				Cantidad:     newCantidad,
				Mensaje:      "Cantidad actualizada exitosamente",
			}

			c.JSON(http.StatusOK, response)
			fmt.Printf("\n\t\t<<<< CANTIDAD ACTUALIZADA EN COTIZACION ID %s >>>>\n", id)
			return
		}

		// Crear nuevo item
		item := models.CotizacionItem{
			CotizacionID: uint(cotizacionID),
			ProductoID:   request.ProductoID,
			SucursalID:   request.SucursalID,
			Cantidad:     request.Cantidad,
		}

		if err := db.Create(&item).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al agregar item a la cotización"})
			fmt.Print("\n\t\t<<<< ERROR AL AGREGAR ITEM >>>>\n")
			return
		}

		response := dtos.AddItemResponse{
			CotizacionID: uint(cotizacionID),
			ProductoID:   request.ProductoID,
			SucursalID:   request.SucursalID,
			Cantidad:     request.Cantidad,
			Mensaje:      "Item agregado exitosamente",
		}

		c.JSON(http.StatusCreated, response)
		fmt.Printf("\n\t\t<<<< ITEM AGREGADO A COTIZACION ID %s >>>>\n", id)
	}
}

// PUT /cotizaciones/[id] - Editar una cotización
func EditarCotizacion(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		cotizacionID, err := strconv.ParseUint(id, 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "ID de cotización inválido"})
			return
		}

		var request dtos.UpdateCotizacionRequest
		if err := c.BindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "JSON inválido"})
			fmt.Print("\n\t\t<<<< JSON INVALIDO AL EDITAR COTIZACION >>>>\n")
			return
		}

		// Buscar la cotización
		var cotizacion models.Cotizacion
		if err := db.First(&cotizacion, cotizacionID).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Cotización no encontrada"})
			fmt.Printf("\n\t\t<<<< COTIZACION ID %s NO ENCONTRADA >>>>\n", id)
			return
		}

		// Actualizar campos si fueron proporcionados
		updates := make(map[string]interface{})
		if request.Estado != nil {
			updates["estado"] = *request.Estado
		}
		if request.CostoEnvio != nil {
			updates["costo_envio"] = *request.CostoEnvio
		}
		if request.TipoDespacho != nil {
			updates["tipo_despacho"] = *request.TipoDespacho
		}

		if len(updates) == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "No se proporcionaron campos para actualizar"})
			return
		}

		// Actualizar la cotización
		if err := db.Model(&cotizacion).Updates(updates).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al actualizar cotización"})
			fmt.Print("\n\t\t<<<< ERROR AL ACTUALIZAR COTIZACION >>>>\n")
			return
		}

		// Obtener la cotización actualizada con relaciones
		if err := db.Preload("Cliente").Preload("Usuario").Preload("Items.Producto").Preload("Items.Sucursal").
			First(&cotizacion, cotizacionID).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener cotización actualizada"})
			return
		}

		// Calcular totales
		var totalItems int
		var totalPrecio float64
		for _, item := range cotizacion.Items {
			totalItems += item.Cantidad
			totalPrecio += float64(item.Cantidad) * item.Producto.Precio
		}
		totalPrecio += cotizacion.CostoEnvio

		// Construir respuesta
		response := dtos.CotizacionResponse{
			ID:           cotizacion.ID,
			FechaCrea:    cotizacion.FechaCrea,
			Estado:       cotizacion.Estado,
			CostoEnvio:   cotizacion.CostoEnvio,
			ClienteID:    cotizacion.ClienteID,
			UserID:       cotizacion.UserID,
			TipoDespacho: cotizacion.TipoDespacho,
			Cliente:      cotizacion.Cliente,
			Usuario:      cotizacion.Usuario,
			TotalItems:   totalItems,
			TotalPrecio:  totalPrecio,
		}

		// Convertir items a respuesta
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
		fmt.Printf("\n\t\t<<<< COTIZACION ID %s ACTUALIZADA >>>>\n", id)
	}
}

// POST /cotizaciones/[id]/preview - Crear preview de cotización
func CrearPreviewCotizacion(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		cotizacionID, err := strconv.ParseUint(id, 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "ID de cotización inválido"})
			return
		}

		// Verificar que la cotización existe
		var cotizacion models.Cotizacion
		if err := db.Preload("Items.Producto").Preload("Items.Sucursal").
			First(&cotizacion, cotizacionID).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Cotización no encontrada"})
			fmt.Printf("\n\t\t<<<< COTIZACION ID %s NO ENCONTRADA >>>>\n", id)
			return
		}

		// Verificar que la cotización tiene items
		if len(cotizacion.Items) == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "La cotización no tiene items"})
			fmt.Printf("\n\t\t<<<< COTIZACION SIN ITEMS >>>>\n")
			return
		}

		// Generar token de acceso único
		tokenBytes := make([]byte, 16)
		rand.Read(tokenBytes)
		tokenAcceso := hex.EncodeToString(tokenBytes)

		// Calcular fecha de expiración (24 horas desde ahora)
		fechaExpiracion := time.Now().Add(24 * time.Hour)

		// Calcular totales
		var totalPrecio float64
		var totalConDescuento float64

		for _, item := range cotizacion.Items {
			// Obtener descuento de la sucursal
			var stock models.StockSucursal
			if err := db.Where("sucursal_id = ? AND producto_id = ?", item.SucursalID, item.ProductoID).
				First(&stock).Error; err == nil {
				// Aplicar descuento si existe
				precioConDescuento := item.Producto.Precio * (1 - stock.Descuento/100)
				totalConDescuento += float64(item.Cantidad) * precioConDescuento
			} else {
				totalConDescuento += float64(item.Cantidad) * item.Producto.Precio
			}
			totalPrecio += float64(item.Cantidad) * item.Producto.Precio
		}

		// Agregar costo de envío
		totalPrecio += cotizacion.CostoEnvio
		totalConDescuento += cotizacion.CostoEnvio
		totalFinal := totalConDescuento

		// Obtener status de pago "PENDIENTE" (ID 1)
		statusPagoID := uint(1)

		// Crear preview
		preview := models.PreviewCotizacion{
			CotizacionID:    uint(cotizacionID),
			TokenAcceso:     tokenAcceso,
			FechaExpiracion: fechaExpiracion,
			StatusPagoID:    statusPagoID,
		}

		if err := db.Create(&preview).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al crear preview de cotización"})
			fmt.Print("\n\t\t<<<< ERROR AL CREAR PREVIEW >>>>\n")
			return
		}

		response := dtos.PreviewCotizacionResponse{
			ID:                preview.ID,
			CotizacionID:      preview.CotizacionID,
			TokenAcceso:       preview.TokenAcceso,
			FechaExpiracion:   preview.FechaExpiracion,
			StatusPagoID:      preview.StatusPagoID,
			TotalPrecio:       totalPrecio,
			TotalConDescuento: totalConDescuento,
			TotalFinal:        totalFinal,
			Mensaje:           "Preview de cotización creado exitosamente",
		}

		c.JSON(http.StatusCreated, response)
		fmt.Printf("\n\t\t<<<< PREVIEW DE COTIZACION ID %s CREADO >>>>\n", id)
	}
}

// DELETE /cotizaciones/:cotizacion_id/items/:item_id - Eliminar o disminuir cantidad de un item de una cotización
func EliminarItemCotizacion(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		cotizacionID, err := strconv.ParseUint(c.Param("cotizacion_id"), 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "ID de cotización inválido"})
			return
		}
		itemID, err := strconv.ParseUint(c.Param("item_id"), 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "ID de item inválido"})
			return
		}

		// Buscar el item
		var item models.CotizacionItem
		if err := db.First(&item, itemID).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Item no encontrado"})
			return
		}
		if item.CotizacionID != uint(cotizacionID) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "El item no pertenece a la cotización indicada"})
			return
		}

		// Leer cantidad opcional
		var body struct {
			Cantidad *int `json:"cantidad"`
		}
		_ = c.ShouldBindJSON(&body)

		if body.Cantidad == nil || *body.Cantidad >= item.Cantidad {
			// Eliminar el item completo
			if err := db.Delete(&item).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al eliminar el item"})
				return
			}
			c.JSON(http.StatusOK, gin.H{"mensaje": "Item eliminado de la cotización"})
			return
		}

		// Disminuir cantidad
		nuevaCantidad := item.Cantidad - *body.Cantidad
		if err := db.Model(&item).Update("cantidad", nuevaCantidad).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al actualizar la cantidad del item"})
			return
		}
		item.Cantidad = nuevaCantidad
		c.JSON(http.StatusOK, gin.H{
			"mensaje": "Cantidad actualizada",
			"item":    item,
		})
	}
}
