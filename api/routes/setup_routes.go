package routes

import (
	// Asegúrate de que este path sea correcto
	"net/http" // Necesitas importar net/http para las constantes de estado HTTP

	"github.com/gin-gonic/gin"
	// Necesitas importar gorm para el tipo *gorm.DB
)

// SetupRoutes configura todas las rutas de la API
func SetupRoutes(router *gin.Engine) {
	// Middleware global
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// Rutas de salud y raíz
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Backend Ventas API",
			"version": "1.0.0",
			"status":  "running",
			"endpoints": gin.H{
				"cotizaciones": "/api/cotizaciones",
				"sucursales":   "/api/sucursales",
				"productos":    "/api/productos",
				"ventas":       "/api/ventas",
				"health":       "/health",
			},
		})
	})

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"message": "API is healthy",
		})
	})

	// Configurar rutas de cotizaciones
	CotizacionRoutes(router)
	SucursalesRoutes(router)
	ProductosRoutes(router)
	VentasRoutes(router)

	// Aquí agregar más grupos de rutas en el futuro:
	// - ProductoRoutes(router)
	// - UsuarioRoutes(router)
	// - SucursalRoutes(router)
	// - etc.
}
