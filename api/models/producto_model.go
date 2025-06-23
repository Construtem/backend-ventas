package models

type Producto struct {
	ID          uint       `gorm:"primaryKey"`
	Codigo      *string    `gorm:"column:codigo;type:VARCHAR(50);unique"`
	Nombre      string     `gorm:"column:nombre;type:TEXT;not null"`
	Descripcion *string    `gorm:"column:descripcion;type:TEXT"`
	CategoriaID *uint      `gorm:"column:categoria_id"`                             // Puede ser nulo
	Categoria   *Categoria `gorm:"foreignKey:CategoriaID"`                          // Belongs To Categoria (Puede ser nulo)
	ProveedorID *uint      `gorm:"column:proveedor_id"`                             // Puede ser nulo
	Proveedor   *Proveedor `gorm:"foreignKey:ProveedorID"`                          // Belongs To Proveedor (Puede ser nulo)
	PrecioCosto float64    `gorm:"column:precio_costo;type:NUMERIC(10,2);not null"` // Considerar usar 'decimal' si la precisión es crítica
	PrecioVenta float64    `gorm:"column:precio_venta;type:NUMERIC(10,2);not null"` // Considerar usar 'decimal' si la precisión es crítica
	Activo      bool       `gorm:"column:activo;type:BOOLEAN;not null;default:true"`
}

func (Producto) TableName() string {
	return "productos"
}
