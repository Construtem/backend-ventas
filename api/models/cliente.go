package models

import (
	"time"
)

// Cliente representa un cliente en el sistema
type Cliente struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Nombre      string    `gorm:"not null" json:"nombre"`
	Telefono    string    `json:"telefono"`
	Email       string    `json:"email"`
	RazonSocial string    `gorm:"column:razon_social" json:"razon_social"`
	Rut         string    `json:"rut"`
	TipoID      uint      `gorm:"column:tipo_id;not null" json:"tipo_id"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	// Relaciones
	TipoCliente  *TipoCliente `gorm:"foreignKey:TipoID" json:"tipo_cliente,omitempty"`
	Direcciones  []DirCliente `gorm:"foreignKey:ClienteID" json:"direcciones,omitempty"`
	Cotizaciones []Cotizacion `gorm:"foreignKey:ClienteID" json:"cotizaciones,omitempty"`
}

// TableName especifica el nombre de la tabla en la base de datos
func (Cliente) TableName() string {
	return "clientes"
}

// TipoCliente representa el tipo de cliente
type TipoCliente struct {
	ID     uint   `gorm:"primaryKey" json:"id"`
	Nombre string `gorm:"not null" json:"nombre"`
}

func (TipoCliente) TableName() string {
	return "tipo_cliente"
}

// DirCliente representa las direcciones de un cliente
type DirCliente struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	ClienteID uint      `gorm:"column:cliente_id;not null" json:"cliente_id"`
	Nombre    string    `gorm:"not null" json:"nombre"`
	Direccion string    `gorm:"not null" json:"direccion"`
	Comuna    string    `gorm:"not null" json:"comuna"`
	Ciudad    string    `gorm:"not null" json:"ciudad"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	// Relaciones
	Cliente *Cliente `gorm:"foreignKey:ClienteID" json:"cliente,omitempty"`
}

func (DirCliente) TableName() string {
	return "dir_cliente"
}

// BeforeCreate hook para validaciones antes de crear
func (c *Cliente) BeforeCreate() error {
	// Aquí puedes agregar validaciones adicionales si es necesario
	return nil
}

// BeforeUpdate hook para validaciones antes de actualizar
func (c *Cliente) BeforeUpdate() error {
	// Aquí puedes agregar validaciones adicionales si es necesario
	return nil
}
