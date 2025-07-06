package main

import (
	"backend-ventas/api/database"
	apiRoutes "backend-ventas/api/routes"
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	// Conectar a la base de datos
	database.InitDB()

	// Configurar todas las rutas usando SetupRoutes
	apiRoutes.SetupRoutes(router)

	// Obtener puerto desde variables de entorno o usar 8080 por defecto
	puerto := os.Getenv("PORT")
	if puerto == "" {
		puerto = "8080"
	}

	router.Run(":" + puerto) // listen and serve on 0.0.0.0:puerto
	fmt.Printf("\n\n\t\t>>>>> Servidor corriendo en localhost:%s !!!\n\n", puerto)
}
