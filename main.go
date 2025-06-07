package main

import (
	"backend-ventas/api/database" // Paquete de conexión a la base de datos
	"backend-ventas/api/routes"   // Paquete de rutas
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"            // Importa la librería
	_ "github.com/joho/godotenv/autoload" // Carga automática del archivo .env
)

func main() {
	// Cargar variables de entorno desde el archivo .env
	if err := godotenv.Load(); err != nil {
		fmt.Println("Error al cargar el archivo .env")
		return
	}

	// Conecta a la base de datos e inicializa GORM
	database.InitDB() // Llama a la función de inicialización de la DB

	fmt.Println("Corriendo en localhost:8080 !!!")

	router := gin.Default()
	routes.SetupRoutes(router)

	router.Run(":8080")
}
