package models

import (
	"time"
)

type Pedido struct {
	ID              uint            `gorm:"primaryKey"`
	Fecha           time.Time       `gorm:"column:fecha;not null;default:NOW()"`
	ClienteID       uint            `gorm:"column:cliente_id;not null"`
	Cliente         Cliente         `gorm:"foreignKey:ClienteID"`
	VendedorID      uint            `gorm:"column:vendedor_id;not null"`
	Vendedor        Usuario         `gorm:"foreignKey:VendedorID"`
	CotizacionID    *uint           `gorm:"column:cotizacion_id"`    // Puede ser nulo
	Cotizacion      *Cotizacion     `gorm:"foreignKey:CotizacionID"` // Belongs To Cotizacion (Puede ser nulo)
	UbicacionID     uint            `gorm:"column:ubicacion_id;not null"`
	Ubicacion       Ubicacion       `gorm:"foreignKey:UbicacionID"`
	Estado          string          `gorm:"column:estado;type:VARCHAR(20);not null;default:'Pendiente';check:estado IN ('Pendiente','Despachado','Cancelado','Completado')"`
	FechaDespacho   *time.Time      `gorm:"column:fecha_despacho"`      // Usar *time.Time para campos nulos
	DespachadoPorID *uint           `gorm:"column:despachado_por"`      // Usar *uint para campos nulos
	DespachadoPor   *Usuario        `gorm:"foreignKey:DespachadoPorID"` // Belongs To Usuario (Puede ser nulo)
	DetallePedido   []DetallePedido `gorm:"foreignKey:PedidoID"`        // Has Many DetallePedido
	BoletaFactura   *BoletaFactura  `gorm:"foreignKey:PedidoID"`        // Has One BoletaFactura
}

func (Pedido) TableName() string {
	return "pedidos"
}
