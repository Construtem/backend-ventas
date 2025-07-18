package main

import (
	"backend-ventas/api/database"
	//"backend-ventas/api/routes"
	"backend-ventas/handlers"
	"backend-ventas/services"

	apiRoutes "backend-ventas/api/routes"
	"fmt"
	"os"

	//_ "github.com/joho/godotenv/autoload" // Carga automática del archivo .env
	"github.com/gin-contrib/cors" //Para configurar CORS
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	// Conectar a la base de datos
	database.InitDB()

	services.InitFirebase() // Inicializar Firebase para autenticación (auth/verify)

	// Configurando CORS

	router.Use(cors.New(cors.Config{
		AllowOrigins: []string{ // Lista de URLs permitidas para CORS
			"*",
			// os.Getenv("FRONT_VENTAS_URL"),      // URL del frontend de ventas
			// os.Getenv("FRONT_INVENTARIO_URL"),  // URL del frontend de inventario
			// os.Getenv("FRONT_FACTURACION_URL"), // URL del frontend de facturación
			// os.Getenv("BACK_FACTURACION_URL"),  // URL del backend de facturación
		},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	router.POST("/auth/verify", handlers.VerifyToken)

	// Configurar todas las rutas usando SetupRoutes
	apiRoutes.SetupRoutes(router)
	apiRoutes.ClienteRoutes(router)

	// Obtener puerto desde variables de entorno o usar 8080 por defecto
	puerto := os.Getenv("PORT")
	if puerto == "" {
		puerto = "8080"
	}

	router.Run(":" + puerto) // listen and serve on 0.0.0.0:puerto
	fmt.Printf("\n\n\t\t>>>>> Servidor corriendo en localhost:%s !!!\n\n", puerto)
}
