package handlers

import (
	"net/http"
	"strconv"

	"backend-ventas/api/controllers"
	"backend-ventas/api/dtos"
	"backend-ventas/api/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func addItemTotals(it models.CotizacionItem, totalItems *int, totalPrecio *float64) {
	if it.Producto == nil {
		return
	}
	*totalItems += it.Cantidad
	*totalPrecio += float64(it.Cantidad) * it.Producto.Precio
}

func safeAppendItemResponse(dst *[]dtos.CotizacionItemResponse, it models.CotizacionItem) {
	if it.Producto == nil || it.Sucursal == nil {
		return
	}
	*dst = append(*dst, dtos.CotizacionItemResponse{
		CotizacionID: it.CotizacionID,
		ProductoID:   it.ProductoID,
		SucursalID:   it.SucursalID,
		Cantidad:     it.Cantidad,
		Producto:     it.Producto,
		Sucursal:     it.Sucursal,
	})
}

// GET /cotizaciones - Obtener todas las cotizaciones simplificadas
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
			totalItems += item.Cantidad
			totalPrecio += float64(item.Cantidad) * item.Producto.Precio
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
				RutCliente:   cotizacion.RutCliente,
				UserID:       cotizacion.UserID,
				TipoDespacho: cotizacion.TipoDespacho,
				Total:        cotizacion.Total,
				Descripcion:  cotizacion.Descripcion,
				Cliente:      cotizacion.Cliente,
				Usuario:      cotizacion.Usuario,
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

// GET /cotizaciones/:rut/historial
func ObtenerCotizacionesPorClienteID(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		rut := c.Param("id")
		clienteID := rut
		if clienteID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Rut del cliente inválido"})
		}
		cortizaciones, err := controllers.ObtenerCotizacionesPorClienteID(clienteID)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Cotizaciones no encontradas"})
		}
		var response []dtos.CotizacionResponse
		for _, cotizacion := range cortizaciones {
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
				RutCliente:   cotizacion.RutCliente,
				UserID:       cotizacion.UserID,
				TipoDespacho: cotizacion.TipoDespacho,
				Total:        cotizacion.Total,
				Descripcion:  cotizacion.Descripcion,
				Cliente:      cotizacion.Cliente,
				Usuario:      cotizacion.Usuario,
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
		cotID, err := strconv.Atoi(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "ID de cotización inválido"})
			return
		}

		cot, err := controllers.ObtenerCotizacionCompletaPorID(cotID)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Cotización no encontrada"})
			return
		}

		var totalItems int
		var totalPrecio float64
		for _, it := range cot.Items {
			addItemTotals(it, &totalItems, &totalPrecio)
		}
		totalPrecio += cot.CostoEnvio

		resp := dtos.CotizacionResponse{
			ID:           cot.ID,
			FechaCrea:    cot.FechaCrea,
			Estado:       cot.Estado,
			CostoEnvio:   cot.CostoEnvio,
			RutCliente:   cot.RutCliente,
			UserID:       cot.UserID,
			TipoDespacho: cot.TipoDespacho,
			Total:        cot.Total,
			Descripcion:  cot.Descripcion,
			Cliente:      cot.Cliente,
			Usuario:      cot.Usuario,
			TotalItems:   totalItems,
			TotalPrecio:  totalPrecio,
		}
		for _, it := range cot.Items {
			safeAppendItemResponse(&resp.Items, it)
		}
		c.JSON(http.StatusOK, resp)
	}
}

// GET /cotizaciones/:id/items/simples - Obtener items simples de una cotización

func ObtenerItemsSimplesCotizacion(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		cotID, err := strconv.Atoi(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "ID de cotización inválido"})
			return
		}

		items, err := controllers.ObtenerItemsSimplesCotizacion(cotID)
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
		for _, it := range items {
			if it.Producto == nil || it.Sucursal == nil {
				continue
			}
			result = append(result, ItemSimple{
				SKU:         it.Producto.SKU,
				Nombre:      it.Producto.Nombre,
				Descripcion: it.Producto.Descripcion,
				Marca:       it.Producto.Marca,
				Cantidad:    it.Cantidad,
				Sucursal:    it.Sucursal.Nombre,
			})
		}
		c.JSON(http.StatusOK, result)
	}
}

// GET /cotizaciones/:id/items - Obtener items completos de una cotización

func ObtenerItemsCotizacion(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		cotID, err := strconv.Atoi(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "ID de cotización inválido"})
			return
		}

		items, err := controllers.ObtenerItemsCompletosCotizacion(cotID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener items"})
			return
		}

		var resp []dtos.CotizacionItemResponse
		for _, it := range items {
			safeAppendItemResponse(&resp, it)
		}
		c.JSON(http.StatusOK, resp)
	}
}

// POST /cotizaciones - Crear nueva cotización
func CrearCotizacion(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req dtos.CreateCotizacionRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Datos de entrada inválidos", "details": err.Error()})
			return
		}

		// Verificaciones básicas (cliente / usuario)
		var cliente models.Cliente
		if err := db.First(&cliente, req.RutCliente).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Cliente no encontrado"})
			return
		}
		var usuario models.Usuario
		if err := db.First(&usuario, req.UserID).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Usuario no encontrado"})
			return
		}

		cot, err := controllers.CrearCotizacion(req.RutCliente, req.UserID, req.TipoDespacho, req.Descripcion, req.CostoEnvio)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al crear cotización"})
			return
		}

		c.JSON(http.StatusCreated, dtos.CreateCotizacionResponse{
			ID:        cot.ID,
			FechaCrea: cot.FechaCrea,
			Estado:    cot.Estado,
			Mensaje:   "Cotización creada exitosamente",
		})
	}
}

// POST /cotizaciones/:id/items - Agregar item a cotización
func AgregarItemCotizacion(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		cotID, err := strconv.Atoi(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "ID de cotización inválido"})
			return
		}

		var req dtos.AddItemRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Datos de entrada inválidos", "details": err.Error()})
			return
		}

		// Verificar existencia de cotización / producto / sucursal
		if err := db.First(&models.Cotizacion{}, cotID).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Cotización no encontrada"})
			return
		}
		if err := db.First(&models.Producto{}, req.ProductoID).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Producto no encontrado"})
			return
		}
		if err := db.First(&models.Sucursal{}, req.SucursalID).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Sucursal no encontrada"})
			return
		}

		// Si el ítem existe se actualiza cantidad; si no, se crea
		item, err := controllers.AgregarItemCotizacion(cotID, req.ProductoID, req.SucursalID, req.Cantidad)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al agregar item"})
			return
		}

		c.JSON(http.StatusCreated, dtos.AddItemResponse{
			CotizacionID: item.CotizacionID,
			ProductoID:   item.ProductoID,
			SucursalID:   item.SucursalID,
			Cantidad:     item.Cantidad,
			Mensaje:      "Item agregado/actualizado exitosamente",
		})
	}
}

// PUT /cotizaciones/:id - Editar cotización
func ActualizarCotizacion(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		cotID, err := strconv.Atoi(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "ID de cotización inválido"})
			return
		}

		var req dtos.UpdateCotizacionRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Datos ingresados inválidos", "details": err.Error()})
			return
		}

		if _, err = controllers.ActualizarCotizacion(cotID, req.CostoEnvio, req.TipoDespacho, req.Total, req.Descripcion); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al actualizar cotización"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"mensaje": "Cotización actualizada exitosamente"})
	}
}

// PATCH /cotizaciones/:id/estado
func ActualizarEstadoCotizacion(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		cotID, err := strconv.Atoi(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "ID de cotización inválido"})
			return
		}

		var req dtos.UpdateEstadoCotizacionRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Datos ingresados inválidos", "details": err.Error()})
			return
		}

		if _, err = controllers.ActualizarEstadoCotizacion(cotID, req.Estado); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al actualizar el estado"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"mensaje": "Estado actualizado correctamente"})
	}
}

// POST /cotizaciones/:id/preview - Crear preview de cotización
func CrearPreviewCotizacion(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		cotID, err := strconv.Atoi(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "ID de cotización inválido"})
			return
		}

		var req dtos.PreviewCotizacionRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Datos de entrada inválidos", "details": err.Error()})
			return
		}

		// Cotización + relaciones
		var cot models.Cotizacion
		if err := db.Preload("Items.Producto").
			Preload("Items.Sucursal").
			First(&cot, cotID).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Cotización no encontrada"})
			return
		}

		var subtotal float64
		for _, it := range cot.Items {
			if it.Producto != nil {
				subtotal += float64(it.Cantidad) * it.Producto.Precio
			}
		}
		tax := subtotal * 0.19 // IVA 19 %
		total := subtotal + tax + cot.CostoEnvio

		prev, err := controllers.CrearPreviewCotizacion(&cotID, req.IssuedAt, subtotal, tax, total)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al crear preview"})
			return
		}

		c.JSON(http.StatusCreated, dtos.PreviewCotizacionResponse{
			ID:            prev.ID,
			CotizacionID:  prev.CotizacionID,
			IssuedAt:      prev.IssuedAt,
			Subtotal:      prev.Subtotal,
			Tax:           prev.Tax,
			Total:         prev.Total,
			PaymentStatus: prev.PaymentStatus,
			StatusPagado:  prev.StatusPagado,
			CreatedAt:     prev.CreatedAt,
			UpdatedAt:     prev.UpdatedAt,
			Mensaje:       "Preview creado exitosamente",
		})
	}
}

// DELETE /cotizaciones/:id/items/:producto_id/:sucursal_id - Eliminar item de cotización
func EliminarItemCotizacion(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		prodID := c.Param("producto_id")
		sucID := c.Param("sucursal_id")

		cotID, err := strconv.Atoi(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "ID de cotización inválido"})
			return
		}
		sucIDInt, err := strconv.Atoi(sucID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "ID de sucursal inválido"})
			return
		}

		if err := controllers.EliminarItemCotizacion(cotID, prodID, sucIDInt); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al eliminar item"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"mensaje": "Item eliminado exitosamente"})
	}
}
