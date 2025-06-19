package models

type Usuario struct {
	ID          uint       `gorm:"primaryKey"`
	Nombre      string     `gorm:"column:nombre;type:TEXT;not null"`
	Email       string     `gorm:"column:email;type:TEXT;not null;unique"`
	Contrasena  string     `gorm:"column:contrasena;type:TEXT;not null"`
	RolID       uint       `gorm:"column:rol_id;not null"`
	Rol         Rol        `gorm:"foreignKey:RolID"`       // Belongs To Rol
	UbicacionID *uint      `gorm:"column:ubicacion_id"`    // Usar puntero para campos nulos
	Ubicacion   *Ubicacion `gorm:"foreignKey:UbicacionID"` // Belongs To Ubicacion (Puede ser nulo)
}

func (Usuario) TableName() string {
	return "usuarios"
}
