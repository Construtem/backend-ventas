package routes

import (
	"backend-ventas/api/handlers"

	"github.com/gin-gonic/gin"
)

// VentaRoutes registra las rutas relacionadas con las ventas usando Gin.
func VentaRoutes(rg *gin.RouterGroup) {
	// Grupo de rutas para /ventas
	ventas := rg.Group("/ventas")
	{
		ventas.GET("/", handlers.GetVentas)    // GET /ventas
		ventas.POST("/", handlers.CreateVenta) // POST /ventas

		ventas.GET("/:id", handlers.GetVentaByID)   // GET /ventas/:id
		ventas.PUT("/:id", handlers.UpdateVenta)    // PUT /ventas/:id
		ventas.DELETE("/:id", handlers.DeleteVenta) // DELETE /ventas/:id
	}
}
