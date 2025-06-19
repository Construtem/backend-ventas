package models

type Ubicacion struct {
	ID        uint    `gorm:"primaryKey"`
	Nombre    string  `gorm:"column:nombre;type:TEXT;not null;unique"`
	Tipo      string  `gorm:"column:tipo;type:VARCHAR(20);not null;check:tipo IN ('Tienda','Bodega')"` // GORM no crea la restricción CHECK automáticamente, se debe manejar en la migración SQL o a nivel de aplicación.
	Direccion *string `gorm:"column:direccion;type:TEXT"`                                              // Puede ser nulo
}

func (Ubicacion) TableName() string {
	return "ubicaciones"
}
