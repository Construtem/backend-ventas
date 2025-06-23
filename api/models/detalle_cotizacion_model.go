package models

type DetalleCotizacion struct {
	CotizacionID   uint       `gorm:"primaryKey;column:cotizacion_id"`
	Cotizacion     Cotizacion `gorm:"foreignKey:CotizacionID;references:ID"`
	ProductoID     uint       `gorm:"primaryKey;column:producto_id" json:"producto_id"`
	Producto       Producto   `gorm:"foreignKey:ProductoID;references:ID"`
	Cantidad       int        `gorm:"column:cantidad;not null" json:"cantidad"`
	PrecioUnitario float64    `gorm:"column:precio_unitario;type:decimal(10,2);not null" json:"precio_unitario"`
}

func (DetalleCotizacion) TableName() string {
	return "detalle_cotizacion"
}
