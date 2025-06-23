package main

import (
	"fmt"
	"log"
	"os"

	"backend-ventas/api/routes" // Paquete de rutas

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// Verifica variables de entorno para la conexión a la DB
	dbHost := os.Getenv("DB_HOST")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbPort := os.Getenv("DB_PORT")

	if dbHost == "" || dbUser == "" || dbPassword == "" || dbName == "" || dbPort == "" {
		log.Fatal("Error: Una o más variables de entorno de la base de datos (DB_HOST, DB_USER, DB_PASSWORD, DB_NAME, DB_PORT) no están configuradas")
	}

	// Cadena de conexión DSN para PostgreSQL
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=America/Santiago", dbHost, dbUser, dbPassword, dbName, dbPort)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{}) // Conecta a la base de datos e inicializa GORM
	if err != nil {
		log.Fatalf("Error al conectar a la base de datos: %v", err)
	}
	log.Println("Conexión a la base de datos exitosa.")

	router := gin.Default()

	// --- Configura todas las rutas ---
	routes.SetupRoutes(router, db)

	port := os.Getenv("API_PORT")
	if port == "" {
		port = "8080" // Puerto por defecto
	}

	log.Printf("Servidor Gin iniciado en el puerto :%s", port)
	router.Run(":" + port)
}
