package routes

import (
	"backend-ventas/api/database"
	"backend-ventas/api/handlers"

	"github.com/gin-gonic/gin"
)

// RegisterUbicacionRoutes añade los endpoints de sucursales al router
func SucursalesRoutes(r *gin.Engine) {
	sucursales := r.Group("/api/sucursales")
	{
		sucursales.GET("", handlers.ObtenerSucursales(database.DB))
		sucursales.GET("/:id", handlers.ObtenerSucursal(database.DB))
	}
}
