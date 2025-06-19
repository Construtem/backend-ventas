package controllers

import (
	"errors"
	"strconv"
	"time"

	"backend-ventas/api/dtos"
	"backend-ventas/api/mappers"
	"backend-ventas/api/models"

	"gorm.io/gorm"
)

// GET de todas las cotizaciones, con filtros por: Fechas, ClienteID, VendedorID, UbicacionID,
// Estado, AprobadorPorID, Fecha Aprobacion, ClienteNombre, ClienteEmail, VendedorNombre, UbicacionNombre.
// Ademas tiene paginacion con valor por defecto LIMIT=10 y PAGE=1.
func GetCotizaciones(db *gorm.DB, queryParams map[string]string) ([]dtos.CotizacionListaResponse, int64, error) {
	var cotizacionesModel []models.Cotizacion
	var totalRecords int64

	query := db.Model(&models.Cotizacion{})

	// --- Aplicar Filtros ---
	// Filtro por rangos de Fechas
	if fechaInicioStr, ok := queryParams["fecha_inicio"]; ok && fechaInicioStr != "" {
		if fechaInicio, err := time.Parse("2006-01-02", fechaInicioStr); err == nil {
			query = query.Where("fecha >= ?", fechaInicio)
		}
	}
	if fechaFinStr, ok := queryParams["fecha_fin"]; ok && fechaFinStr != "" {
		if fechaFin, err := time.Parse("2006-01-02", fechaFinStr); err == nil {
			query = query.Where("fecha <= ?", fechaFin.Add(24*time.Hour-time.Second))
		}
	}

	// Filtros por IDs directos
	if clienteIDStr, ok := queryParams["cliente_id"]; ok && clienteIDStr != "" {
		if clienteID, err := strconv.ParseUint(clienteIDStr, 10, 64); err == nil {
			query = query.Where("cliente_id = ?", clienteID)
		}
	}
	if vendedorIDStr, ok := queryParams["vendedor_id"]; ok && vendedorIDStr != "" {
		if vendedorID, err := strconv.ParseUint(vendedorIDStr, 10, 64); err == nil {
			query = query.Where("vendedor_id = ?", vendedorID)
		}
	}
	if ubicacionIDStr, ok := queryParams["ubicacion_id"]; ok && ubicacionIDStr != "" {
		if ubicacionID, err := strconv.ParseUint(ubicacionIDStr, 10, 64); err == nil {
			query = query.Where("ubicacion_id = ?", ubicacionID)
		}
	}
	if estado, ok := queryParams["estado"]; ok && estado != "" {
		query = query.Where("estado = ?", estado)
	}
	if aprobadaPorIDStr, ok := queryParams["aprobada_por_id"]; ok && aprobadaPorIDStr != "" {
		if aprobadaPorID, err := strconv.ParseUint(aprobadaPorIDStr, 10, 64); err == nil {
			query = query.Where("aprobada_por_id = ?", aprobadaPorID)
		}
	}
	// Filtro por rango de Fecha Aprobacion
	if fechaAprobacionInicioStr, ok := queryParams["fecha_aprobacion_inicio"]; ok && fechaAprobacionInicioStr != "" {
		if fechaAprobacionInicio, err := time.Parse("2006-01-02", fechaAprobacionInicioStr); err == nil {
			query = query.Where("fecha_aprobacion >= ?", fechaAprobacionInicio)
		}
	}
	if fechaAprobacionFinStr, ok := queryParams["fecha_aprobacion_fin"]; ok && fechaAprobacionFinStr != "" {
		if fechaAprobacionFin, err := time.Parse("2006-01-02", fechaAprobacionFinStr); err == nil {
			query = query.Where("fecha_aprobacion <= ?", fechaAprobacionFin.Add(24*time.Hour-time.Second))
		}
	}

	// Filtro por Cliente (Nombre o Email) - Requiere JOINs
	if clienteNombre, ok := queryParams["cliente_nombre"]; ok && clienteNombre != "" {
		query = query.Joins("JOIN usuarios AS c ON cotizaciones.cliente_id = c.id").
			Where("c.nombre ILIKE ?", "%"+clienteNombre+"%")
	}
	if clienteEmail, ok := queryParams["cliente_email"]; ok && clienteEmail != "" {
		query = query.Joins("JOIN usuarios AS c ON cotizaciones.cliente_id = c.id").
			Where("c.email ILIKE ?", "%"+clienteEmail+"%")
	}

	// Filtro por Nombre de Vendedor - Requiere JOIN
	if vendedorNombre, ok := queryParams["vendedor_nombre"]; ok && vendedorNombre != "" {
		query = query.Joins("JOIN usuarios AS v ON cotizaciones.vendedor_id = v.id").
			Where("v.nombre ILIKE ?", "%"+vendedorNombre+"%")
	}

	// Filtro por Nombre de Ubicación - Requiere JOIN
	if ubicacionNombre, ok := queryParams["ubicacion_nombre"]; ok && ubicacionNombre != "" {
		query = query.Joins("JOIN ubicaciones AS u ON cotizaciones.ubicacion_id = u.id").
			Where("u.nombre ILIKE ?", "%"+ubicacionNombre+"%")
	}

	// Conteo Total para Paginacion
	countQuery := query
	if err := countQuery.Count(&totalRecords).Error; err != nil {
		return nil, 0, err
	}

	// Paginación
	limit := 10
	page := 1

	if limitStr, ok := queryParams["limit"]; ok && limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}
	if pageStr, ok := queryParams["page"]; ok && pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	offset := (page - 1) * limit
	query = query.Limit(limit).Offset(offset)

	// --- Precargar Relaciones (Preload) ---
	query = query.Preload("Cliente").
		Preload("Vendedor").
		Preload("Ubicacion").
		Preload("AprobadaPor").
		Preload("DetalleCotizacion.Producto")

	if err := query.Find(&cotizacionesModel).Error; err != nil {
		return nil, 0, err
	}

	// Mapear los modelos de GORM a DTOs para la respuesta (formato simplificado)
	cotizacionesResponse := make([]dtos.CotizacionListaResponse, len(cotizacionesModel))
	for i, cot := range cotizacionesModel {
		cotizacionesResponse[i] = mappers.MapCotizacionToListaDTO(cot)
	}

	return cotizacionesResponse, totalRecords, nil
}

// GetCotizacionByID obtiene una cotización por su ID, incluyendo sus relaciones, y la mapea a un DTO.
func GetCotizacionByID(db *gorm.DB, id uint) (dtos.CotizacionResponse, error) {
	var cotizacionModel models.Cotizacion

	err := db.Preload("Cliente").
		Preload("Vendedor").
		Preload("Ubicacion").
		Preload("AprobadaPor").
		Preload("DetalleCotizacion.Producto").
		First(&cotizacionModel, id).
		Error

	if err != nil {
		return dtos.CotizacionResponse{}, err
	}

	cotizacionResponse := mappers.MapCotizacionToDTO(cotizacionModel)
	return cotizacionResponse, nil
}

// CreateCotizacion crea una nueva cotización en la base de datos.
func CreateCotizacion(db *gorm.DB, cotizacion *models.Cotizacion) error {
	tx := db.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	// Validar que no haya productos repetidos en el detalle
	productoSet := make(map[uint]bool)
	for _, d := range cotizacion.DetalleCotizacion {
		if productoSet[d.ProductoID] {
			tx.Rollback()
			return errors.New("no se puede repetir el mismo producto en el detalle de la cotización")
		}
		productoSet[d.ProductoID] = true
	}

	// SET fecha por defecto si no se proporciona
	if cotizacion.Fecha.IsZero() {
		cotizacion.Fecha = time.Now()
	}
	// SET estado por defecto si no se proporciona
	if cotizacion.Estado == "" {
		cotizacion.Estado = "Pendiente"
	}

	// Valida y maneja campos de aprobacion condicionalmente
	if cotizacion.Estado == "Aprobada" {
		if cotizacion.AprobadaPorID == nil || cotizacion.FechaAprobacion == nil || cotizacion.FechaAprobacion.IsZero() {
			tx.Rollback()
			return errors.New("los campos 'aprobadaPorID' y 'fechaAprobacion' son requeridos si el estado es 'Aprobada'")
		}
	} else {
		cotizacion.AprobadaPorID = nil
		cotizacion.FechaAprobacion = nil
	}

	// Deja que GORM inserte la cotización y sus detalles en cascada
	if err := tx.Create(&cotizacion).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		return err
	}

	return nil
}

// UpdateEstadoCotizacion actualiza el estado de una cotización y la mapea a un DTO para la respuesta.
func UpdateEstadoCotizacion(db *gorm.DB, id uint, nuevoEstado string, usuarioID *uint) (dtos.CotizacionResponse, error) {
	var cotizacion models.Cotizacion

	if err := db.First(&cotizacion, id).Error; err != nil {
		return dtos.CotizacionResponse{}, err
	}

	switch nuevoEstado {
	case "Pendiente", "Aprobada", "Rechazada", "Vencida":
	default:
		return dtos.CotizacionResponse{}, errors.New("estado proporcionado invalido")
	}

	cotizacion.Estado = nuevoEstado

	if nuevoEstado == "Aprobada" {
		if usuarioID != nil {
			cotizacion.AprobadaPorID = usuarioID
			now := time.Now()
			cotizacion.FechaAprobacion = &now
		} else {
			return dtos.CotizacionResponse{}, errors.New("es necesario el id del usuario para aprobar la cotizacion")
		}
	} else {
		cotizacion.AprobadaPorID = nil
		cotizacion.FechaAprobacion = nil
	}

	if err := db.Save(&cotizacion).Error; err != nil {
		return dtos.CotizacionResponse{}, err
	}

	var updatedCotizacionModel models.Cotizacion
	err := db.Preload("Cliente").
		Preload("Vendedor").
		Preload("Ubicacion").
		Preload("AprobadaPor").
		First(&updatedCotizacionModel, id).Error

	if err != nil {
		return dtos.CotizacionResponse{}, errors.New("cotización actualizada, pero error al recargar para la respuesta: " + err.Error())
	}

	return mappers.MapCotizacionToDTO(updatedCotizacionModel), nil
}

// UpdateCotizacionDetalle actualiza los detalles de una cotización.
// Utiliza la clave compuesta (CotizacionID, ProductoID) para identificar y modificar/eliminar/añadir detalles.
func UpdateCotizacionDetalle(db *gorm.DB, cotizacionID uint, nuevoDetalle []models.DetalleCotizacion) (dtos.CotizacionResponse, error) {
	var cotizacion models.Cotizacion

	tx := db.Begin()
	if tx.Error != nil {
		return dtos.CotizacionResponse{}, tx.Error
	}

	if err := tx.Preload("DetalleCotizacion").First(&cotizacion, cotizacionID).Error; err != nil {
		tx.Rollback()
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return dtos.CotizacionResponse{}, errors.New("cotización no encontrada")
		}
		return dtos.CotizacionResponse{}, err
	}

	// Crear un mapa de los detalles actuales por ProductoID para facilitar la comparación
	// (Dentro de una misma cotización, ProductoID debe ser único para cada detalle)
	actualDetalleMap := make(map[uint]models.DetalleCotizacion)
	for _, d := range cotizacion.DetalleCotizacion {
		actualDetalleMap[d.ProductoID] = d
	}

	// Iterar sobre los nuevos detalles: actualizar o agregar
	for _, nuevoDetalleItem := range nuevoDetalle {
		if actual, ok := actualDetalleMap[nuevoDetalleItem.ProductoID]; ok {
			// Detalle existente: verificar si necesita actualización
			if actual.Cantidad != nuevoDetalleItem.Cantidad || actual.PrecioUnitario != nuevoDetalleItem.PrecioUnitario {
				actual.Cantidad = nuevoDetalleItem.Cantidad
				actual.PrecioUnitario = nuevoDetalleItem.PrecioUnitario

				// GORM usará CotizacionID y ProductoID para actualizar la fila correcta
				if err := tx.Save(&actual).Error; err != nil {
					tx.Rollback()
					return dtos.CotizacionResponse{}, err
				}
			}
			// Eliminar del mapa de detalles existentes, ya que ha sido procesado (actualizado).
			delete(actualDetalleMap, nuevoDetalleItem.ProductoID)
		} else { // Agrega el Detalle nuevo
			nuevoDetalleItem.CotizacionID = cotizacionID // Asignar el ID de la cotización padre
			// GORM creará una nueva fila usando CotizacionID y ProductoID como clave compuesta
			if err := tx.Create(&nuevoDetalleItem).Error; err != nil {
				tx.Rollback()
				return dtos.CotizacionResponse{}, err
			}
		}
	}

	// Procesar los detalles restantes en actualDetalleMap (eliminar).
	// Cualquier detalle que quede en actualDetalleMap no estaba en la nueva lista, por lo que debe ser eliminado.
	for _, detailToDelete := range actualDetalleMap {
		// GORM eliminará la fila usando CotizacionID y ProductoID
		if err := tx.Delete(&detailToDelete).Error; err != nil {
			tx.Rollback()
			return dtos.CotizacionResponse{}, err
		}
	}

	// Recargar la cotización con los detalles actualizados para la respuesta.
	var updatedCotizacionModel models.Cotizacion
	if err := tx.Preload("Cliente").
		Preload("Vendedor").
		Preload("Ubicacion").
		Preload("AprobadaPor").
		Preload("DetalleCotizacion.Producto").
		First(&updatedCotizacionModel, cotizacionID).Error; err != nil {
		tx.Rollback()
		return dtos.CotizacionResponse{}, err
	}

	if err := tx.Commit().Error; err != nil {
		return dtos.CotizacionResponse{}, err
	}

	return mappers.MapCotizacionToDTO(updatedCotizacionModel), nil
}
