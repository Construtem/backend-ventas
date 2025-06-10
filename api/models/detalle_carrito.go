package models // El paquete de modelos

// DetalleCarrito representa la tabla Detalle_Carrito
type DetalleCarrito struct {
	IDDetalleCarrito  uint    `gorm:"primaryKey;column:id_detalle_carrito" json:"id_detalle_carrito"`
	IDCarrito         uint    `json:"id_carrito"`
	IDProducto        uint    `json:"id_producto"` // Asumiendo que `Productos` tiene `id_producto` de tipo `uint`
	Cantidad          int     `json:"cantidad"`
	PrecioUnitario    float64 `json:"precio_unitario"`
	DescuentoAplicado float64 `json:"descuento_aplicado"`
}
