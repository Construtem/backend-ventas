package controllers

import (
	"backend-ventas/api/database"
	"backend-ventas/api/models"
)

func ContarVentasRealizadas() (int64, error) {
	var total int64
	err := database.DB.
		Model(&models.Cotizacion{}).
		Where("estado_pago = ?", "pagado").
		Count(&total).Error
	return total, err
}
