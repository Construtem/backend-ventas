package models

type Categoria struct {
	ID     uint   `gorm:"primaryKey"`
	Nombre string `gorm:"column:nombre;type:TEXT;not null;unique"`
}

func (Categoria) TableName() string {
	return "categorias"
}
