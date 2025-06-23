package models

type Cliente struct {
	ID        uint    `gorm:"primaryKey"`
	Nombre    string  `gorm:"column:nombre;type:TEXT;not null"`
	Telefono  *string `gorm:"column:telefono;type:VARCHAR(20);"`
	Email     *string `gorm:"column:email;type:TEXT"`
	Direccion *string `gorm:"column:direccion;type:TEXT"`
}

func (Cliente) TableName() string {
	return "clientes"
}
