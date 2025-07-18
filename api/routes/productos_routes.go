package routes

import (
	"backend-ventas/api/handlers"
	"github.com/gin-gonic/gin"
)

func ProductosRoutes(r *gin.Engine) {
	prod := r.Group("/api/productos")
	{
		// Ruta para listar productos con inventario: EJEMPLO: /api/productos/inventario?sucursal_id={id}&page=1&limit=10
		prod.GET("/inventario", handlers.ObtenerProductosInventario())
	}
}
