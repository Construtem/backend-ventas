package routes

import (
	"backend-ventas/api/database"
	"backend-ventas/api/handlers"
	"backend-ventas/api/middleware"

	"github.com/gin-gonic/gin"
)

func CotizacionRoutes(r *gin.Engine) {
	// Grupo de rutas para cotizaciones con prefijo /api
	cotizaciones := r.Group("/api/cotizaciones")
	{
		cotizaciones.GET("", handlers.ObtenerCotizacionesSimplificadas(database.DB))
		cotizaciones.POST("", handlers.CrearCotizacion(database.DB))
		cotizaciones.GET(":id", handlers.ObtenerCotizacionSimplificada(database.DB))
		cotizaciones.PUT(":id", handlers.ActualizarCotizacion(database.DB))
		cotizaciones.PATCH(":id/estado", middleware.AuthRoles("gerente"), handlers.ActualizarEstadoCotizacion(database.DB))
		cotizaciones.GET(":id/items", handlers.ObtenerItemsCotizacion(database.DB)) // devuelve null
		cotizaciones.POST(":id/items", handlers.AgregarItemCotizacion(database.DB))
		cotizaciones.POST(":id/preview", handlers.CrearPreviewCotizacion(database.DB))
		cotizaciones.GET("completas", handlers.ObtenerCotizaciones(database.DB))
		cotizaciones.GET("completa/:id", handlers.ObtenerCotizacionPorID(database.DB))
		cotizaciones.GET(":id/historial", handlers.ObtenerCotizacionesPorClienteID(database.DB))
		cotizaciones.GET(":id/items/simple", handlers.ObtenerItemsSimplesCotizacion(database.DB)) //devuelve null
		cotizaciones.DELETE(":id/items/:producto_id/:sucursal_id", handlers.EliminarItemCotizacion(database.DB))
		cotizaciones.GET("/checkout/:id", handlers.ObtenerCotizacionCheckout(database.DB))
		cotizaciones.GET("/status/:id", handlers.ObtenerEstadoPagoCotizacion(database.DB))
		cotizaciones.PATCH("/status/:id", handlers.ActualizarEstadoPagoCotizacion(database.DB))
		cotizaciones.GET("/cliente/:rut/historial", handlers.ObtenerHistorialCotizaciones(database.DB)) //Devuelve todas las cotizaciones de un cliente por su rut
		cotizaciones.GET("/despacho", handlers.TestObtenerDespachoDestinoCotizacion(database.DB))
		cotizaciones.DELETE("/:id", handlers.EliminarCotizacion(database.DB))
	}
}
