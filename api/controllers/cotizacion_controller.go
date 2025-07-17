package controllers

import (
	"backend-ventas/api/database"
	"backend-ventas/api/dtos"
	"backend-ventas/api/mappers"
	"backend-ventas/api/models"
	"time"
)

// Métodos para cotizaciones SIMPLIFICADAS
func ListarCotizacionesSimples() ([]models.Cotizacion, error) {
	var cotizaciones []models.Cotizacion
	err := database.DB.Preload("Cliente").Preload("Usuario").Preload("Items.Producto").Find(&cotizaciones).Error
	return cotizaciones, err
}

func ObtenerCotizacionSimplePorID(id int) (models.Cotizacion, error) {
	var cotizacion models.Cotizacion
	err := database.DB.Preload("Cliente").Preload("Usuario").Preload("Items.Producto").First(&cotizacion, id).Error
	return cotizacion, err
}

func ObtenerItemsSimplesCotizacion(id int) ([]models.CotizacionItem, error) {
	var items []models.CotizacionItem
	// Optimización: usar índices y limitar campos si es necesario
	err := database.DB.Preload("Producto").Preload("Sucursal").Where("cotizacion_id = ?", id).Find(&items).Error
	return items, err
}

// Métodos para cotizaciones COMPLETAS
func ListarCotizacionesCompletas() ([]models.Cotizacion, error) {
	var cotizaciones []models.Cotizacion
	err := database.DB.Preload("Cliente").Preload("Usuario").Preload("Items.Producto").Preload("Items.Sucursal").Find(&cotizaciones).Error
	return cotizaciones, err
}

func ObtenerCotizacionCompletaPorID(id int) (models.Cotizacion, error) {
	var cotizacion models.Cotizacion
	err := database.DB.Preload("Cliente").Preload("Usuario").Preload("Items.Producto").Preload("Items.Sucursal").First(&cotizacion, id).Error
	return cotizacion, err
}

func ObtenerCotizacionesPorClienteID(id string) ([]models.Cotizacion, error) {
	var cotizaciones []models.Cotizacion
	err := database.DB.Preload("Cliente").Preload("Usuario").Preload("Items.Producto").Preload("Items.Sucursal").First(&cotizaciones, id).Error
	return cotizaciones, err
}

func ObtenerItemsCompletosCotizacion(id int) ([]models.CotizacionItem, error) {
	var items []models.CotizacionItem
	err := database.DB.Preload("Producto").Preload("Sucursal").Where("cotizacion_id = ?", id).Find(&items).Error
	return items, err
}

// Métodos para crear cotizaciones
func CrearCotizacion(rutCliente, userID, tipoDespacho string, descripcion *string, costoEnvio float64) (models.Cotizacion, error) {
	cotizacion := models.Cotizacion{
		RutCliente:   rutCliente,
		UserID:       userID,
		TipoDespacho: tipoDespacho,
		CostoEnvio:   costoEnvio,
		Descripcion:  descripcion,
		Estado:       "pendiente",
		FechaCrea:    time.Now(),
	}

	err := database.DB.Create(&cotizacion).Error
	return cotizacion, err
}

// Métodos para agregar items a cotización
func AgregarItemCotizacion(cotizacionID int, productoID string, sucursalID int, cantidad int) (models.CotizacionItem, error) {
	item := models.CotizacionItem{
		CotizacionID: cotizacionID,
		ProductoID:   productoID,
		SucursalID:   sucursalID,
		Cantidad:     cantidad,
	}

	err := database.DB.Create(&item).Error
	return item, err
}

// Métodos para actualizar cotizaciones
func ActualizarCotizacion(id int, costoEnvio *float64, tipoDespacho *string, total *float64, descripcion *string) (models.Cotizacion, error) {
	var cotizacion models.Cotizacion
	err := database.DB.First(&cotizacion, id).Error
	if err != nil {
		return cotizacion, err
	}

	if costoEnvio != nil {
		cotizacion.CostoEnvio = *costoEnvio
	}
	if tipoDespacho != nil {
		cotizacion.TipoDespacho = *tipoDespacho
	}
	if total != nil {
		cotizacion.Total = total
	}
	if descripcion != nil {
		cotizacion.Descripcion = descripcion
	}

	err = database.DB.Save(&cotizacion).Error
	return cotizacion, err
}

// Método para actualizar estado de una cotizacione
func ActualizarEstadoCotizacion(id int, estado string) (models.Cotizacion, error) {
	var cotizacion models.Cotizacion
	err := database.DB.First(&cotizacion, id).Error
	if err != nil {
		return cotizacion, err
	}
	if estado != "" {
		cotizacion.Estado = estado
	}

	err = database.DB.Save(&cotizacion).Error
	return cotizacion, err
}

// Métodos para eliminar items de cotización
func EliminarItemCotizacion(cotizacionID int, productoID string, sucursalID int) error {
	return database.DB.Where("cotizacion_id = ? AND sku = ? AND sucursal_id = ?",
		cotizacionID, productoID, sucursalID).Delete(&models.CotizacionItem{}).Error
}

// Métodos para preview de cotización
func CrearPreviewCotizacion(cotizacionID *int, issuedAt time.Time, subtotal, tax, total float64) (models.PreviewCotizacion, error) {
	var preview models.PreviewCotizacion
	// Buscar si ya existe un preview para esta cotización
	err := database.DB.Where("cotizacion_id = ?", cotizacionID).First(&preview).Error
	if err == nil {
		// Ya existe, actualizarlo
		preview.IssuedAt = issuedAt
		preview.Subtotal = subtotal
		preview.Tax = tax
		preview.Total = total
		preview.UpdatedAt = time.Now()
		err = database.DB.Save(&preview).Error
		return preview, err
	}
	// Si no existe, crearlo
	preview = models.PreviewCotizacion{
		CotizacionID:  cotizacionID,
		IssuedAt:      issuedAt,
		Subtotal:      subtotal,
		Tax:           tax,
		Total:         total,
		PaymentStatus: "pending",
		StatusPagado:  false,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}
	err = database.DB.Create(&preview).Error
	return preview, err
}

// ObtenerCotizacionCheckout devuelve la cotización con cliente, dirección,
// usuario, ítems precargados y totales calculados.
func ObtenerCotizacionCheckout(id int) (*dtos.CheckoutCotizacionResponse, error) {
	var cot models.Cotizacion
	err := database.DB.
		Preload("Cliente").
		Preload("Cliente.Direcciones").
		Preload("Usuario").
		Preload("Items"). // ← NUEVO
		Preload("Items.Producto").
		Preload("Items.Sucursal").
		Preload("PreviewCotizacion").
		First(&cot, id).Error
	if err != nil {
		return nil, err
	}

	// -------- construir DTO --------
	var subtotal float64
	itemsDTO := make([]dtos.CheckoutItemDTO, 0) // ← evita null

	for _, it := range cot.Items {
		// mientras depuras, comenta el filtro o registra por qué es nil
		if it.Producto == nil || it.Sucursal == nil {
			// log.Printf("ítem descartado: %+v", it)
			continue
		}

		sub := float64(it.Cantidad) * it.Producto.Precio
		subtotal += sub

		itemsDTO = append(itemsDTO, dtos.CheckoutItemDTO{
			SKU:        it.Producto.SKU,
			Nombre:     it.Producto.Nombre,
			Cantidad:   it.Cantidad,
			PrecioUnit: it.Producto.Precio,
			Subtotal:   sub,
			Sucursal:   it.Sucursal.Nombre,
		})
	}

	iva := subtotal * 0.19
	total := subtotal + iva + cot.CostoEnvio

	// primera dirección si existe
	dirDTO := dtos.DireccionCliente{}
	if len(cot.Cliente.Direcciones) > 0 {
		dirDTO = mappers.DirClienteToDTO(&cot.Cliente.Direcciones[0])
	}

	resp := &dtos.CheckoutCotizacionResponse{
		ID:           cot.ID,
		FechaCrea:    cot.FechaCrea.Format("2006-01-02 15:04:05"),
		Estado:       cot.Estado,
		CostoEnvio:   cot.CostoEnvio,
		TipoDespacho: cot.TipoDespacho,

		Cliente:   mappers.ClienteToDTO(cot.Cliente),
		Usuario:   mappers.UsuarioToDTO(cot.Usuario),
		Direccion: dirDTO,

		Items:        itemsDTO,
		SubtotalNeto: subtotal,
		IVA:          iva,
		Total:        total,
	}

	if cot.PreviewCotizacion != nil {
		resp.PreviewID = &cot.PreviewCotizacion.ID
	}
	return resp, nil
}

func ActualizarEstadoPagoCotizacion(id int, estadoPago string) (models.Cotizacion, error) {
	var cotizacion models.Cotizacion
	err := database.DB.First(&cotizacion, id).Error

	if err != nil {
		return cotizacion, err
	}

	if estadoPago != "" {
		cotizacion.EstadoPago = estadoPago
	}

	err = database.DB.Save(&cotizacion).Error
	return cotizacion, err

}

func ObtenerEstadoPagoCotizacionPorID(id int) (models.Cotizacion, error) {
	var cotizacion models.Cotizacion
	err := database.DB.First(&cotizacion, id).Error

	return cotizacion, err
}
