package database

import (
	"backend-ventas/api/models"

	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/joho/godotenv/autoload" // Carga automática del archivo .env
	"gorm.io/driver/postgres"             // Driver de PostgreSQL
	"gorm.io/gorm"                        // ORM de GORM
)

var DB *gorm.DB

func InitDB() {
	dbHost := os.Getenv("DB_HOST")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbPort := os.Getenv("DB_PORT")
	dbSSLMode := os.Getenv("DB_SSLMODE") // 'disable' para dev, 'require'/'verify-full' para prod

	if dbHost == "" || dbUser == "" || dbPassword == "" || dbName == "" || dbPort == "" {
		log.Fatal("Error: Una o más variables de entorno de la base de datos no están configuradas.")
	}

	// Cadena de conexión (DSN)
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		dbHost, dbUser, dbPassword, dbName, dbPort, dbSSLMode)

	var err error
	// Intentar abrir la conexión a la base de datos con reintentos
	const maxRetries = 10 // Número máximo de intentos
	for i := range maxRetries {
		DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err == nil {
			// Si la conexión fue exitosa, salimos del bucle
			break
		}
		// Si falló, logueamos el error y esperamos antes de reintentar
		log.Printf("Intento %d/%d: Fallo al conectar a la base de datos: %v. Reintentando en 2 segundos...", i+1, maxRetries, err)
		time.Sleep(2 * time.Second) // Espera 2 segundos antes del siguiente intento
	}

	// Después de todos los reintentos, si `err` todavía no es `nil`, significa que no se pudo conectar.
	if err != nil {
		log.Fatalf("Fallo CRÍTICO: No se pudo conectar a la base de datos después de %d intentos: %v", maxRetries, err)
	}
	// Si llegamos aquí, la conexión fue exitosa

	log.Println("Conexión a la base de datos establecida exitosamente.")

	migrationErr := DB.AutoMigrate(
		&models.Usuario{},
		&models.Cliente{},
		&models.HistorialCotizacion{},
	)

	if migrationErr != nil {
		log.Fatalf("Fallo la automigración de la base de datos: %v", err)
	}
}

// CloseDB cierra la conexión a la base de datos subyacente de GORM
func CloseDB() {
	if DB != nil {
		sqlDB, err := DB.DB()
		if err != nil {
			log.Printf("Error obteniendo *sql.DB de GORM: %v", err)
		} else {
			sqlDB.Close()
			log.Println("Conexión a PostgreSQL (GORM) cerrada.")
		}
	}
}
