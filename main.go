package main

import (
	"backend-ventas/api/database" // Paquete de conexión a la base de datos
	"backend-ventas/api/routes"   // Paquete de rutas (ahora con funciones para Gin)
	"fmt"
	"os" // Para manejar variables de entorno

	"github.com/gin-gonic/gin"
)

func main() {
	if os.Getenv("DB_HOST") == "" {
		fmt.Println("Advertencia: Las variables de entorno pueden no haberse cargado correctamente.")
		// Añadir log.Fatal("Error: Las variables de entorno no se han cargado.")
	}

	// Conecta a la base de datos e inicializa GORM
	database.InitDB()        // Llama a la función de inicialización de la DB
	defer database.CloseDB() // Asegúrate de cerrar la conexión de la DB

	// Inicializa el router de Gin
	router := gin.Default()

	// Configura todas las rutas usando la función SetupRoutes del paquete routes
	routes.SetupRoutes(router)

	// Obtener el puerto de las variables de entorno, con un valor por defecto
	port := os.Getenv("API_PORT")
	if port == "" {
		port = "8080" // Puerto por defecto
	}

	router.Run(":" + port) // Inicia el servidor Gin
}
