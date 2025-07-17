package controllers

import (
	"backend-ventas/api/database"
	"backend-ventas/api/models"
)

// ListarSucursales devuelve todas las ubicaciones cuyo tipo_id = 2 (sucursales)
func ListarSucursales() ([]models.Sucursales, error) {
	var sucursales []models.Sucursales
	err := database.DB.
		Where("tipo_id = ?", 2).
		Find(&sucursales).Error
	return sucursales, err
}

// ObtenerSucursalPorID devuelve una sucursal por su id, asegurándose de que sea tipo_id = 2
func ObtenerSucursalPorID(id int) (models.Sucursales, error) {
	var sucursal models.Sucursales
	err := database.DB.
		Where("id = ? AND tipo_id = ?", id, 2).
		First(&sucursal).Error
	return sucursal, err
}
