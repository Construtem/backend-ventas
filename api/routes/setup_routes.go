package routes

import (
	"backend-ventas/api/handlers" // Asegúrate de que este path sea correcto
	"net/http"                    // Necesitas importar net/http para las constantes de estado HTTP

	"github.com/gin-gonic/gin"
	"gorm.io/gorm" // Necesitas importar gorm para el tipo *gorm.DB
)

// SetupRoutes ahora recibe la conexión a la base de datos 'db' como parámetro.
func SetupRoutes(router *gin.Engine, db *gorm.DB) {
	// Puedes añadir middleware global aquí si es necesario
	// router.Use(gin.Logger())
	// router.Use(gin.Recovery())

	// Crea un grupo base para las rutas de la API, por ejemplo, /api/
	// Esto es opcional, pero buena práctica para versionar tu API.
	api := router.Group("/api") // Se recomienda agrupar tus rutas de API bajo un prefijo
	{
		// Rutas de Cotizacion (ahora inyectando 'db')
		api.GET("/cotizaciones", handlers.GetCotizacionesHandler(db))
		api.POST("/cotizaciones", handlers.CreateCotizacionHandler(db))
		api.GET("/cotizaciones/:id", handlers.GetCotizacionByIDHandler(db))
		api.PATCH("/cotizaciones/:id/estado", handlers.UpdateEstadoCotizacionHandler(db))
		api.PATCH("/cotizaciones/:id/detalles", handlers.UpdateCotizacionDetalleHandler(db))
		api.GET("/cotizaciones/:id/historial", handlers.GetHistorialCotizacionHandler(db))

		// !!! IMPORTANTE: Estas rutas necesitan ser refactorizadas en sus handlers
		// para aceptar 'db' y retornar gin.HandlerFunc, de forma consistente.
		// Por ahora, si 'GetUsuarios' y 'GetClientes' siguen usando 'database.DB' global,
		// estas líneas provocarán inconsistencias con el patrón de inyección.
		// **Propongo cómo deberían quedar abajo.**

		// Si GetUsuarios y GetClientes SÍ acceden a una variable global database.DB,
		// entonces esta sintaxis es correcta para la ruta:
		// api.GET("/users", handlers.GetUsuarios) // Si GetUsuarios tiene func(c *gin.Context)
		// api.GET("/users/:id", handlers.GetUsuarioByID)
		// api.GET("/clientes", handlers.GetClientes)

		// SIN EMBARGO, para mantener la CONSISTENCIA y el patrón de "Gin puro",
		// si quieres que estos handlers también reciban 'db', DEBES refactorizarlos así:

		// Propuesta de rutas refactorizadas para Usuarios y Clientes:
		// api.GET("/users", handlers.GetUsuariosHandler(db)) // Asumiendo que GetUsuariosHandler existe y toma 'db'
		// api.GET("/users/:id", handlers.GetUsuarioByIDHandler(db)) // Asumiendo que GetUsuarioByIDHandler existe y toma 'db'
		// api.GET("/clientes", handlers.GetClientesHandler(db)) // Asumiendo que GetClientesHandler existe y toma 'db'
	}

	// Rutas de salud o raíz opcionales (fuera del grupo /api si no son parte de la API principal)
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Welcome to Backend Ventas API"}) // Usar http.StatusOK
	})
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"}) // Usar http.StatusOK
	})
}
