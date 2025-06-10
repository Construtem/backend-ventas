package routes

import (
	"github.com/gin-gonic/gin"
)

// SetupRoutes configura todas las rutas de la API usando el router de Gin.
func SetupRoutes(router *gin.Engine) {
	// Puedes añadir middleware global aquí si es necesario
	// router.Use(gin.Logger())
	// router.Use(gin.Recovery())

	// Crea un grupo base para las rutas de la API, por ejemplo, /api/
	// Esto es opcional, pero buena práctica para versionar tu API.
	// Si no quieres /api/, simplemente usa router.Group("/")
	api := router.Group("/api/")
	{
		// Llama a las funciones que registran las rutas específicas
		CarritoRoutes(api) // Pasa el grupo de rutas a la función de carritos
		VentaRoutes(api)   // Pasa el grupo de rutas a la función de ventas

		// Puedes añadir otras rutas aquí, ej:
		// UserRoutes(api)
		// ProductRoutes(api)
	}

	// Rutas de salud o raíz opcionales
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Welcome to Backend Ventas API"})
	})
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})
}
