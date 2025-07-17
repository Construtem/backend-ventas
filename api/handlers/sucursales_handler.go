package handlers

import (
	"net/http"
	"strconv"

	"backend-ventas/api/controllers"
	"backend-ventas/api/dtos"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// ObtenerSucursales construye el handler para GET /sucursales
func ObtenerSucursales(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		sucursales, err := controllers.ListarSucursales()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener sucursales"})
			return
		}

		var resp []dtos.UbicacionResponse
		for _, s := range sucursales {
			resp = append(resp, dtos.UbicacionResponse{
				ID:        s.ID,
				Nombre:    s.Nombre,
				Telefono:  s.Telefono,
				Direccion: s.Direccion,
				Comuna:    s.Comuna,
				Ciudad:    s.Ciudad,
			})
		}

		c.JSON(http.StatusOK, resp)
	}
}

// ObtenerSucursal construye el handler para GET /sucursales/:id
func ObtenerSucursal(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		idParam := c.Param("id")
		id, err := strconv.Atoi(idParam)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "ID de sucursal inválido"})
			return
		}

		sucursal, err := controllers.ObtenerSucursalPorID(id)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Sucursal no encontrada"})
			return
		}

		resp := dtos.UbicacionResponse{
			ID:        sucursal.ID,
			Nombre:    sucursal.Nombre,
			Telefono:  sucursal.Telefono,
			Direccion: sucursal.Direccion,
			Comuna:    sucursal.Comuna,
			Ciudad:    sucursal.Ciudad,
		}

		c.JSON(http.StatusOK, resp)
	}
}
