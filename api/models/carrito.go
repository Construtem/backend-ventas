package models // El paquete de modelos

import "time"

// CarritoCompras representa la tabla Carrito_Compras
type CarritoCompras struct {
	IDCarrito uint      `gorm:"primaryKey;column:id_carrito" json:"id_carrito"`
	IDUsuario uint      `json:"id_usuario"`
	IDCliente uint      `json:"id_cliente"`
	Fecha     time.Time `json:"fecha"`
	Estado    string    `json:"estado"`
	// Relación: Un Carrito tiene muchos Detalles de Carrito
	Detalles []DetalleCarrito `gorm:"foreignKey:IDCarrito" json:"detalles,omitempty"`
	// Si tienes modelos Usuario y Cliente, puedes añadir las relaciones aquí:
	// Usuario models.Usuario `gorm:"foreignKey:IDUsuario"`
	// Cliente models.Cliente `gorm:"foreignKey:IDCliente"`
}
