package models

import "gorm.io/gorm"

type Usuario struct {
	gorm.Model             // Incluye ID, CreatedAt, UpdatedAt, DeletedAt
	ID          uint       `gorm:"primaryKey"` // ID del usuario, clave primaria
	Nombre      string     `gorm:"not null"`
	Email       string     `gorm:"unique;not null"`
	Contrasena  string     `gorm:"not null"` // Hash bcrypt
	RolID       uint       `gorm:"not null"` // Relación con tabla roles
	Rol         *Rol       `gorm:"foreignKey:RolID"`
	UbicacionID *uint      // Puede ser NULL (tienda asignada)
	Ubicacion   *Ubicacion `gorm:"foreignKey:UbicacionID"` // Relación con tabla ubicaciones
	Activo      bool       `gorm:"default:true"`           // Indica si el usuario está activo
}

type Rol struct {
	gorm.Model        // Incluye ID, CreatedAt, UpdatedAt, DeletedAt
	ID         uint   `gorm:"primaryKey"`
	Nombre     string `gorm:"unique;not null"`
}

func (Rol) TableName() string {
	return "roles" // <--- Esto le dice a GORM que use la tabla "roles"
}

type Ubicacion struct {
	gorm.Model        // Incluye ID, CreatedAt, UpdatedAt, DeletedAt
	ID         uint   `gorm:"primaryKey"`
	Nombre     string `gorm:"unique;not null"`
	Tipo       string `gorm:"not null"` // Ej: "Tienda", "Almacen", "Oficina"
	Direccion  string `gorm:"not null"`
}

func (Ubicacion) TableName() string {
	return "ubicaciones" // <--- Esto le dice a GORM que use la tabla "ubicaciones"
}
