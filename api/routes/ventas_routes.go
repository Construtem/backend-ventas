package routes

import (
	"backend-ventas/api/handlers"
	"github.com/gin-gonic/gin"
)

func VentasRoutes(r *gin.Engine) {
	v := r.Group("/api/ventas")
	{
		v.GET("/realizadas", handlers.ObtenerResumenVentas())
	}
}
