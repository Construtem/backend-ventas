package routes

import (
	"backend-ventas/api/handlers"

	"github.com/gin-gonic/gin"
)

// CarritoRoutes registra las rutas relacionadas con los carritos de compras usando Gin.
func CarritoRoutes(rg *gin.RouterGroup) {
	// Grupo de rutas para /carritos
	carritos := rg.Group("/carritos")
	{
		carritos.POST("/", handlers.CreateCarrito)    // POST /carritos
		carritos.GET("/:id", handlers.GetCarritoByID) // GET /carritos/:id

		// Rutas para ítems dentro de un carrito
		items := carritos.Group("/:id/items")
		{
			items.POST("/", handlers.AddItemToCarrito)                // POST /carritos/:id/items
			items.PUT("/:item_id", handlers.UpdateItemInCarrito)      // PUT /carritos/:id/items/:item_id
			items.DELETE("/:item_id", handlers.DeleteItemFromCarrito) // DELETE /carritos/:id/items/:item_id
		}
	}
}
