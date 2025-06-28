package models

type Rol struct {
	ID     uint   `gorm:"primaryKey"`
	Nombre string `gorm:"column:nombre;type:VARCHAR(50);not null;unique"`
}

func (Rol) TableName() string {
	return "roles"
}
