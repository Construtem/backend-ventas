package routes

import (
	"backend-ventas/api/auth"     // Asegúrate de importar tu paquete de autenticación
	"backend-ventas/api/database" // <--- NUEVO: Necesitas la DB para obtener el rol real del usuario
	"backend-ventas/api/models"   // <--- NUEVO: Necesitas los modelos para consultar el usuario
	"backend-ventas/api/utils"
	"net/http" // Para http.StatusOK
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm" // Para manejar errores de GORM
)

// SetupRoutes configura todas las rutas de tu aplicación.
func SetupRoutes(router *gin.Engine) {
	// Rutas de autenticación (login, registro, etc.)
	SetupAuthRoutes(router)

	// Ruta pública de ejemplo (no requiere autenticación)
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "pong"})
	})

	// --- Rutas protegidas por rol, utilizan el middleware AutorizarRol ---

	// Grupo de rutas para administradores
	admin := router.Group("/admin")
	admin.Use(auth.AutorizarRol("administrador")) // Aplica el middleware a todas las rutas de este grupo
	{
		admin.GET("/dashboard", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "Bienvenido al panel de administrador!"})
		})
		// admin.POST("/usuarios", CreateUserHandler)
		// admin.GET("/reportes", GetAdminReports)
	}

	// Grupo de rutas para vendedores
	vendedor := router.Group("/vendedor")
	vendedor.Use(auth.AutorizarRol("vendedor"))
	{
		vendedor.GET("/productos", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "Lista de productos para vendedores."})
		})
		// vendedor.POST("/ventas", CreateSaleHandler)
	}

	// Grupo de rutas para clientes
	cliente := router.Group("/cliente")
	cliente.Use(auth.AutorizarRol("cliente"))
	{
		cliente.GET("/perfil", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "Acceso al perfil del cliente."})
		})
		// cliente.GET("/mis-pedidos", GetMyOrdersHandler)
	}

	router.GET("/test", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "prueba de conexión exitosa",
		})
	})

	// --- TEMPORAL: Endpoint para generar tokens de prueba (¡QUITAR EN PRODUCCIÓN!) ---
	router.GET("/generate-test-token/:id", func(c *gin.Context) {
		userIDStr := c.Param("id")
		userID, err := strconv.ParseUint(userIDStr, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "ID de usuario inválido."})
			return
		}

		// **MODIFICACIÓN CLAVE AQUÍ:**
		// Necesitamos obtener el rol del usuario desde la base de datos para pasárselo a GenerateJWT.
		var usuario models.Usuario
		if err := database.DB.Preload("Rol").First(&usuario, uint(userID)).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				c.JSON(http.StatusNotFound, gin.H{"error": "Usuario no encontrado para generar token."})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al buscar usuario para generar token: " + err.Error()})
			return
		}

		if usuario.Rol == nil || usuario.Rol.Nombre == "" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "El usuario no tiene un rol asignado o el rol no tiene nombre."})
			return
		}

		// Genera el token, pasando el ID del usuario Y su nombre de rol
		tokenString, err := utils.GenerateJWT(usuario.ID, usuario.Rol.Nombre) // <-- MODIFICADO
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al generar token de prueba: " + err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"token": tokenString})
	})
}
