package models // El paquete de modelos

import "time"

// Venta representa la tabla Ventas
type Venta struct {
	IDVenta   uint      `gorm:"primaryKey;column:id_venta" json:"id_venta"`
	IDCarrito uint      `json:"id_carrito"`
	Fecha     time.Time `json:"fecha"`
	Total     float64   `json:"total"`
	Estado    string    `json:"estado"`
	// Puedes definir la relación con CarritoCompras si el modelo CarritoCompras ya existe en models.
	// GORM infiere la relación por `IDCarrito` si la convención de nombres es seguida.
	// Carrito CarritoCompras `gorm:"foreignKey:IDCarrito"`
}

func (Venta) TableName() string {
	return "ventas"
}
