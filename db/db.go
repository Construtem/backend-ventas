package db

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Conectar() {
	// Datos de conexión
	joho := godotenv.Load()
	if joho != nil {
		log.Fatal("Error cargando .env:", joho)
	}

	dbHost := os.Getenv("DB_HOST")
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbPort := os.Getenv("DB_PORT")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		dbHost, dbUser, dbPass, dbName, dbPort)

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("\n\n\t\t>>>>> ERROR AL CONECTAR LA BASE DE DATOS:\n\n", err)
	}

	fmt.Print("\n\n\t\t>>>>> BASE DE DATOS CONECTADA \n\n")
}