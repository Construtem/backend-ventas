package controllers

import (
	"backend-ventas/api/database"
	"backend-ventas/api/dtos"
	"backend-ventas/api/mappers"
	"backend-ventas/api/models"
	"fmt"
	"gorm.io/gorm"
	"log"
	"math"
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
func CrearCotizacion(rutCliente, userID, tipoDespacho string, descripcion *string, costoEnvio float64, total *float64) (models.Cotizacion, error) {
	cotizacion := models.Cotizacion{
		RutCliente:   rutCliente,
		UserID:       userID,
		TipoDespacho: tipoDespacho,
		Total:        total,
		CostoEnvio:   costoEnvio,
		Descripcion:  descripcion,
		Estado:       "pendiente",
		EstadoPago:   "pendiente",
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
		EstadoPagado:  "pendiente",
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
	log.Printf("se está buscando la cotizacion %d", id)
	err := database.DB.
		Preload("Cliente.Direcciones").
		Preload("Usuario").
		Preload("Items.Producto").
		Preload("Items.Sucursal").
		First(&cot, id).Error
	if err != nil {
		return nil, err
	}

	/*-----------------------------------------------------------
	  1) Traer descuentos de stock_sucursal para los ítems
	-----------------------------------------------------------*/
	// Mapear claves sku+sucursal_id → % descuento
	// Mapear claves sku+sucursal_id → % descuento
	type row struct{ Descuento float64 } // ← float64
	descs := make(map[string]int)        // ← guardamos int 0-100

	for _, it := range cot.Items {
		var r row
		key := fmt.Sprintf("%s|%d", it.ProductoID, it.SucursalID)

		if err := database.DB.Raw(`
        SELECT COALESCE(descuento,0) AS descuento
          FROM stock_sucursal
         WHERE sku = ? AND sucursal_id = ?`,
			it.ProductoID, it.SucursalID).
			Scan(&r).Error; err != nil {
			return nil, err
		}

		pct := int(math.Round(r.Descuento)) // 10.00 → 10
		if pct < 0 {
			pct = 0
		}
		descs[key] = pct
	}

	/*-----------------------------------------------------------
	  2) Construir DTO y totales
	-----------------------------------------------------------*/
	var subtotal float64
	var descuentoTotal float64
	itemsDTO := make([]dtos.CheckoutItemDTO, 0, len(cot.Items))

	for _, it := range cot.Items {
		if it.Producto == nil || it.Sucursal == nil {
			continue
		}
		key := fmt.Sprintf("%s|%d", it.ProductoID, it.SucursalID)
		pctDesc := descs[key] // 0-100
		precio := it.Producto.Precio
		sub := float64(it.Cantidad) * precio     // sin descuento
		ahorro := sub * float64(pctDesc) / 100.0 // monto rebajado

		subtotal += sub
		descuentoTotal += ahorro

		itemsDTO = append(itemsDTO, dtos.CheckoutItemDTO{
			SKU:        it.Producto.SKU,
			Nombre:     it.Producto.Nombre,
			Cantidad:   it.Cantidad,
			PrecioUnit: precio,
			Subtotal:   sub,
			Descuento:  pctDesc, // porcentaje entero
			Sucursal:   it.Sucursal.Nombre,
		})
	}

	subtotalConDesc := subtotal - descuentoTotal
	base := subtotalConDesc + cot.CostoEnvio
	iva := base * 0.19
	total := base + iva

	// primera dirección
	dirDTO := dtos.DireccionCliente{}
	if len(cot.Cliente.Direcciones) > 0 {
		dirDTO = mappers.DirClienteToDTO(&cot.Cliente.Direcciones[0])
	}

	resp := &dtos.CheckoutCotizacionResponse{
		ID:           cot.ID,
		FechaCrea:    cot.FechaCrea.Format("2006-01-02 15:04:05"),
		Estado:       cot.Estado,
		TipoDespacho: cot.TipoDespacho,
		EstadoPago:   cot.EstadoPago,

		Cliente:   mappers.ClienteToDTO(cot.Cliente),
		Usuario:   mappers.UsuarioToDTO(cot.Usuario),
		Direccion: dirDTO,

		Items:          itemsDTO,
		CostoEnvio:     cot.CostoEnvio, // ← valor del despacho
		SubtotalNeto:   subtotal,
		DescuentoTotal: descuentoTotal,
		IVA:            iva,
		Total:          total,
	}

	if cot.PreviewCotizacion != nil {
		resp.PreviewID = &cot.PreviewCotizacion.ID
	}
	return resp, nil
}

// ListarCotizacionesPorCliente devuelve TODAS las cotizaciones de un RUT
func ListarCotizacionesPorCliente(rut string) ([]models.Cotizacion, error) {
	var cots []models.Cotizacion
	err := database.DB.
		Preload("Cliente").
		Preload("Usuario").
		Preload("Items.Producto").
		Preload("Items.Sucursal").
		Where("rut_cliente = ?", rut).
		Order("fecha_crea DESC").
		Find(&cots).Error // ←  Find (no First)
	return cots, err
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

func ObtenerDespachoDestinoCotizacion(id int) (*models.DirCliente, error) {
	var despacho models.Despacho
	err := database.DB.Preload("DestinoObj").Where("cotizacion_id = ?", id).First(&despacho).Error

	return despacho.DestinoObj, err
}
func EliminarCotizacionPorID(id string) error {
	return database.DB.Transaction(func(tx *gorm.DB) error {
		var cot models.Cotizacion
		if err := tx.First(&cot, id).Error; err != nil {
			return err
		}

		if cot.Estado != "rechazada" {
			return fmt.Errorf("solo se pueden eliminar cotizaciones rechazadas")
		}

		// Eliminar ítems asociados
		if err := tx.Where("cotizacion_id = ?", id).Delete(&models.CotizacionItem{}).Error; err != nil {
			return err
		}

		// Eliminar la cotización
		if err := tx.Delete(&cot).Error; err != nil {
			return err
		}

		return nil
	})
}
