package handlers

import (
	// "crypto/rand"
	// "encoding/hex"
	"log"

	"net/http"
	"strconv"

	"backend-ventas/api/controllers"
	"backend-ventas/api/dtos"
	"backend-ventas/api/mappers"
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

		var result []dtos.CotizacionSimplificadaResponse
		for _, cot := range cotizaciones {
			var totalItems int
			var totalPrecio float64
			for _, item := range cot.Items {
				if item.Producto != nil {
					totalItems += item.Cantidad
					totalPrecio += float64(item.Cantidad) * item.Producto.Precio
				}
			}
			totalPrecio += cot.CostoEnvio

			res := dtos.CotizacionSimplificadaResponse{
				ID:           cot.ID,
				FechaCrea:    cot.FechaCrea,
				Estado:       cot.Estado,
				CostoEnvio:   cot.CostoEnvio,
				UserID:       cot.UserID,
				TipoDespacho: cot.TipoDespacho,
				EstadoPago:   cot.EstadoPago,
				TotalItems:   totalItems,
				TotalPrecio:  totalPrecio,
			}

			if cot.Usuario != nil {
				res.Nombre = cot.Usuario.Nombre
			}
			if cl := cot.Cliente; cl != nil {
				res.Cliente.Rut = cl.Rut
				res.Cliente.Nombre = cl.Nombre
				if cl.Telefono != nil {
					res.Cliente.Telefono = *cl.Telefono
				}
				if cl.Email != nil {
					res.Cliente.Email = *cl.Email
				}
				if cl.RazonSocial != nil {
					res.Cliente.RazonSocial = *cl.RazonSocial
				}
			}
			for _, item := range cot.Items {
				if item.Producto != nil {
					res.Items = append(res.Items, struct {
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
			result = append(result, res)
		}
		c.JSON(http.StatusOK, result)
	}
}

// GET /cotizaciones/:id - Obtener detalle simplificado

func ObtenerCotizacionSimplificada(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		cotID, err := strconv.Atoi(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "ID de cotización inválido"})
			return
		}
		cot, err := controllers.ObtenerCotizacionSimplePorID(cotID)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Cotización no encontrada"})
			return
		}

		var totalItems int
		var totalPrecio float64
		for _, item := range cot.Items {
			if item.Producto != nil {
				totalItems += item.Cantidad
				totalPrecio += float64(item.Cantidad) * item.Producto.Precio
			}
		}
		totalPrecio += cot.CostoEnvio
		response := dtos.CotizacionSimplificadaResponse{
			ID:           cot.ID,
			FechaCrea:    cot.FechaCrea,
			Estado:       cot.Estado,
			CostoEnvio:   cot.CostoEnvio,
			UserID:       cot.UserID,
			Nombre:       "",
			TipoDespacho: cot.TipoDespacho,
			EstadoPago:   cot.EstadoPago,
			TotalItems:   totalItems,
			TotalPrecio:  totalPrecio,
		}
		if cot.Usuario != nil {
			response.Nombre = cot.Usuario.Nombre
		}
		if cl := cot.Cliente; cl != nil {
			response.Cliente.Rut = cl.Rut
			response.Cliente.Nombre = cl.Nombre
			if cl.Telefono != nil {
				response.Cliente.Telefono = *cl.Telefono
			}
			if cl.Email != nil {
				response.Cliente.Email = *cl.Email
			}
			if cl.RazonSocial != nil {
				response.Cliente.RazonSocial = *cl.RazonSocial
			}
		}
		for _, item := range cot.Items {
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

		var resp []dtos.CotizacionResponse
		for _, cot := range cotizaciones {
			var totalItems int
			var totalPrecio float64
			for _, item := range cot.Items {
				totalItems += item.Cantidad
				if item.Producto != nil {
					totalPrecio += float64(item.Cantidad) * item.Producto.Precio
				} else {
					// Opcional: loguea o maneja el error, por ejemplo:
					log.Printf("Producto nil en item: %+v", item)
				}
			}
			totalPrecio += cot.CostoEnvio

			cr := dtos.CotizacionResponse{
				ID:           cot.ID,
				FechaCrea:    cot.FechaCrea,
				Estado:       cot.Estado,
				CostoEnvio:   cot.CostoEnvio,
				RutCliente:   cot.RutCliente,
				UserID:       cot.UserID,
				TipoDespacho: cot.TipoDespacho,
				Total:        cot.Total,
				Descripcion:  cot.Descripcion,
				EstadoPago:   cot.EstadoPago,
				Cliente:      cot.Cliente,
				Usuario:      cot.Usuario,
				TotalItems:   totalItems,
				TotalPrecio:  totalPrecio,
			}
			for _, it := range cot.Items {
				safeAppendItemResponse(&cr.Items, it)
			}
			resp = append(resp, cr)
		}
		c.JSON(http.StatusOK, resp)
	}
}

// GET /cotizaciones/:rut/historial
func ObtenerCotizacionesPorClienteID(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		rut := c.Param("id")
		if rut == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Rut del cliente inválido"})
			return
		}
		cots, err := controllers.ObtenerCotizacionesPorClienteID(rut)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Cotizaciones no encontradas"})
			return
		}

		resp := make([]dtos.CotizacionResponse, 0)
		for _, cot := range cots {
			var totalItems int
			var totalPrecio float64
			for _, it := range cot.Items {
				addItemTotals(it, &totalItems, &totalPrecio)
			}
			totalPrecio += cot.CostoEnvio

			cr := dtos.CotizacionResponse{
				ID:           cot.ID,
				FechaCrea:    cot.FechaCrea,
				Estado:       cot.Estado,
				CostoEnvio:   cot.CostoEnvio,
				RutCliente:   cot.RutCliente,
				UserID:       cot.UserID,
				TipoDespacho: cot.TipoDespacho,
				Total:        cot.Total,
				Descripcion:  cot.Descripcion,
				EstadoPago:   cot.EstadoPago,
				Cliente:      cot.Cliente,
				Usuario:      cot.Usuario,
				TotalItems:   totalItems,
				TotalPrecio:  totalPrecio,
				Items:        make([]dtos.CotizacionItemResponse, 0),
			}
			for _, it := range cot.Items {
				safeAppendItemResponse(&cr.Items, it)
			}
			resp = append(resp, cr)
		}
		c.JSON(http.StatusOK, resp)
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
		for _, item := range cot.Items {
			if item.Producto != nil {
				totalItems += item.Cantidad
				totalPrecio += float64(item.Cantidad) * item.Producto.Precio
			}
		}
		totalPrecio += cot.CostoEnvio
		response := dtos.CotizacionResponse{
			ID:           cot.ID,
			FechaCrea:    cot.FechaCrea,
			Estado:       cot.Estado,
			CostoEnvio:   cot.CostoEnvio,
			RutCliente:   cot.RutCliente,
			UserID:       cot.UserID,
			TipoDespacho: cot.TipoDespacho,
			Descripcion:  cot.Descripcion,
			EstadoPago:   cot.EstadoPago,
			Total:        cot.Total,
			Cliente:      cot.Cliente,
			Usuario:      cot.Usuario,
			TotalItems:   totalItems,
			TotalPrecio:  totalPrecio,
		}
		for _, it := range cot.Items {
			safeAppendItemResponse(&response.Items, it)
		}
		c.JSON(http.StatusOK, response)
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
				SKU:      it.Producto.SKU,
				Nombre:   it.Producto.Nombre,
				Cantidad: it.Cantidad,
				Sucursal: it.Sucursal.Nombre,
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
		if err := db.Where("rut = ?", req.RutCliente).First(&cliente).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Cliente no encontrado"})
			return
		}
		var usuario models.Usuario
		if err := db.Where("email = ?", req.UserID).First(&usuario).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Usuario no encontrado"})
			return
		}

		cot, err := controllers.CrearCotizacion(req.RutCliente, req.UserID, req.TipoDespacho, req.Descripcion, req.CostoEnvio, req.Total)
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
		// Verificar que el producto existe
		var producto models.Producto
		if err := db.Where("sku = ?", req.ProductoID).First(&producto).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Producto no encontrado"})
			return
		}
		if err := db.First(&models.Sucursal{}, req.SucursalID).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Sucursal no encontrada"})
			return
		}
		// Verificar si ya existe el item
		var existingItem models.CotizacionItem
		err = db.Where("cotizacion_id = ? AND sku = ? AND sucursal_id = ?", cotID, req.ProductoID, req.SucursalID).First(&existingItem).Error
		if err == nil {
			// El item ya existe, actualizar cantidad
			existingItem.Cantidad += req.Cantidad
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
			EstadoPago:    prev.EstadoPagado,
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

// GET /api/cotizaciones/checkout/:id
func ObtenerCotizacionCheckout(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
			return
		}
		dto, err := controllers.ObtenerCotizacionCheckout(id)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Cotización no encontrada"})
			return
		}
		c.JSON(http.StatusOK, dto)
	}
}

func ObtenerHistorialCotizaciones(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		rut := c.Param("rut") // ← ahora el parámetro se llama :rut
		if rut == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "RUT inválido"})
			return
		}

		cots, err := controllers.ListarCotizacionesPorCliente(rut)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener cotizaciones"})
			return
		}

		// --- mapear a DTO (reaprovechamos la lógica que ya tenías) ---
		resp := make([]dtos.CotizacionResponse, 0, len(cots))
		for _, cot := range cots {
			var totalItems int
			var totalPrecio float64
			for _, it := range cot.Items {
				if it.Producto != nil {
					totalItems += it.Cantidad
					totalPrecio += float64(it.Cantidad) * it.Producto.Precio
				}
			}
			totalPrecio += cot.CostoEnvio

			// Buscar el despacho de la cotización asociada
			direccionObj, err := controllers.ObtenerDespachoDestinoCotizacion(cot.ID)
			var direccionDTO *dtos.DireccionCliente
			if err == nil && direccionObj != nil {
				dir := mappers.DirClienteToDTO(direccionObj)
				direccionDTO = &dir
			}

			cr := dtos.CotizacionResponse{
				ID:           cot.ID,
				FechaCrea:    cot.FechaCrea,
				Estado:       cot.Estado,
				EstadoPago:   cot.EstadoPago,
				CostoEnvio:   cot.CostoEnvio,
				RutCliente:   cot.RutCliente,
				UserID:       cot.UserID,
				TipoDespacho: cot.TipoDespacho,
				Direccion:    direccionDTO,
				Total:        cot.Total,
				Descripcion:  cot.Descripcion,
				Cliente:      cot.Cliente,
				Usuario:      cot.Usuario,
				TotalItems:   totalItems,
				TotalPrecio:  totalPrecio,
			}
			for _, it := range cot.Items {
				safeAppendItemResponse(&cr.Items, it)
			}
			resp = append(resp, cr)
		}
		c.JSON(http.StatusOK, resp)
	}
}

func ActualizarEstadoPagoCotizacion(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		cotID, err := strconv.Atoi(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "ID de cotización inválido"})
			return
		}
		var req dtos.UpdateEstadoPagoCotizacionRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Datos ingresados inválidos", "details": err.Error()})
			return
		}

		_, err = controllers.ActualizarEstadoPagoCotizacion(cotID, req.EstadoPago)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Error al actualizar el Estado del pago"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"mensaje": "Estado de pago actualizado correctamente"})
	}
}

func ObtenerEstadoPagoCotizacion(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		cotID, err := strconv.Atoi(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "ID de cotización inválido"})
			return
		}
		var res dtos.EstadoPagoCotizacionResponse
		cotizacion, err := controllers.ObtenerEstadoPagoCotizacionPorID(cotID)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Cotización no encontrada"})
			return
		}
		res.ID = cotizacion.ID
		res.EstadoPago = cotizacion.EstadoPago

		c.JSON(http.StatusOK, res)
	}
}
func EliminarCotizacion(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		if id == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "ID de cotización inválido"})
			return
		}

		err := controllers.EliminarCotizacionPorID(id)
		if err != nil {
			// Detectar si el error es por estado no rechazado
			if err.Error() == "solo se pueden eliminar cotizaciones rechazadas" {
				c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al eliminar cotización"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"mensaje": "Cotización eliminada correctamente"})
	}
}

func TestObtenerDespachoDestinoCotizacion(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var dirCliente *models.DirCliente
		dirCliente, err := controllers.ObtenerDespachoDestinoCotizacion(4)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Direccion no encontrada"})
		}
		c.JSON(http.StatusOK, dirCliente)
	}
}
