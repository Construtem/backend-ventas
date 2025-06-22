package models

type Cliente struct {
	ID           uint   `gorm:"primaryKey" json:"id"`
	Nombre       string `gorm:"column:nombre_cliente" json:"nombre"`
	Rut          string `json:"rut"`
	Correo       string `gorm:"unique" json:"correo"`
	Telefono     string `json:"telefono"`
	RegionComuna string `json:"region_comuna"`
}

// Tabla se llama "cliente" explícitamente (opcional)
func (Cliente) TableName() string {
	return "cliente"
}
