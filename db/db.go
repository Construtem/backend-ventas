package db

import (
	"fmt"
	"log"

	"backend-ventas/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Conectar() {
	// Datos de conexión
	host := "aws-0-sa-east-1.pooler.supabase.com"
	user := "postgres.ksnafogaaqwnoajvtivq"
	password := "x1qtjvPeWkC5FMAK"
	dbname := "postgres"
	port := "5432"

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		host, user, password, dbname, port)

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Error al conectar a la base de datos:", err)
	}

	// Migrar la estructura del modelo (opcional, crea la tabla si no existe)
	err = DB.AutoMigrate(&models.Cliente{})
	if err != nil {
		log.Fatal("Error al migrar tabla cliente:", err)
	}

	fmt.Println("Base de datos conectada con GORM")
}
