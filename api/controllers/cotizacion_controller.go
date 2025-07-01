package controllers

import (
	"backend-ventas/api/database"
	"backend-ventas/api/models"
)

// Métodos para cotizaciones SIMPLIFICADAS
func ListarCotizacionesSimples() ([]models.Cotizacion, error) {
	var cotizaciones []models.Cotizacion
	err := database.DB.Preload("Cliente").Preload("Usuario").Preload("Items.Producto").Find(&cotizaciones).Error
	return cotizaciones, err
}

func ObtenerCotizacionSimplePorID(id uint) (models.Cotizacion, error) {
	var cotizacion models.Cotizacion
	err := database.DB.Preload("Cliente").Preload("Usuario").Preload("Items.Producto").First(&cotizacion, id).Error
	return cotizacion, err
}

func ObtenerItemsSimplesCotizacion(id uint) ([]models.CotizacionItem, error) {
	var items []models.CotizacionItem
	err := database.DB.Preload("Producto").Preload("Sucursal").Where("cotizacion_id = ?", id).Find(&items).Error
	return items, err
}

// Métodos para cotizaciones COMPLETAS
func ListarCotizacionesCompletas() ([]models.Cotizacion, error) {
	var cotizaciones []models.Cotizacion
	err := database.DB.Preload("Cliente").Preload("Usuario").Preload("Items.Producto").Preload("Items.Sucursal").Find(&cotizaciones).Error
	return cotizaciones, err
}

func ObtenerCotizacionCompletaPorID(id uint) (models.Cotizacion, error) {
	var cotizacion models.Cotizacion
	err := database.DB.Preload("Cliente").Preload("Usuario").Preload("Items.Producto").Preload("Items.Sucursal").First(&cotizacion, id).Error
	return cotizacion, err
}

func ObtenerItemsCompletosCotizacion(id uint) ([]models.CotizacionItem, error) {
	var items []models.CotizacionItem
	err := database.DB.Preload("Producto").Preload("Sucursal").Where("cotizacion_id = ?", id).Find(&items).Error
	return items, err
}
