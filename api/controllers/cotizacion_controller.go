package controllers

import (
	"backend-ventas/api/database"
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

func ObtenerItemsCompletosCotizacion(id int) ([]models.CotizacionItem, error) {
	var items []models.CotizacionItem
	err := database.DB.Preload("Producto").Preload("Sucursal").Where("cotizacion_id = ?", id).Find(&items).Error
	return items, err
}

// Métodos para crear cotizaciones
func CrearCotizacion(rutCliente, userID, tipoDespacho string, costoEnvio float64) (models.Cotizacion, error) {
	cotizacion := models.Cotizacion{
		RutCliente:   rutCliente,
		UserID:       userID,
		TipoDespacho: tipoDespacho,
		CostoEnvio:   costoEnvio,
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
func ActualizarCotizacion(id int, estado *string, costoEnvio *float64, tipoDespacho *string, total *float64) (models.Cotizacion, error) {
	var cotizacion models.Cotizacion
	err := database.DB.First(&cotizacion, id).Error
	if err != nil {
		return cotizacion, err
	}

	if estado != nil {
		cotizacion.Estado = *estado
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

	err = database.DB.Save(&cotizacion).Error
	return cotizacion, err
}

// Métodos para eliminar items de cotización
func EliminarItemCotizacion(cotizacionID int, productoID string, sucursalID int) error {
	return database.DB.Where("cotizacion_id = ? AND producto_id = ? AND sucursal_id = ?",
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
