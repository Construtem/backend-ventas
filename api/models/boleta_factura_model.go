package models

import (
	"time"
)

type BoletaFactura struct {
	ID           uint      `gorm:"primaryKey"`
	PedidoID     uint      `gorm:"column:pedido_id;not null;unique"`  // One-to-one con Pedido
	Pedido       Pedido    `gorm:"foreignKey:PedidoID;references:ID"` // Belongs To Pedido
	Tipo         string    `gorm:"column:tipo;type:VARCHAR(10);not null;check:tipo IN ('Boleta', 'Factura')"`
	FechaEmision time.Time `gorm:"column:fecha_emision;not null;default:NOW()"`
	Total        float64   `gorm:"column:total;type:NUMERIC(12,2);not null"` // Considerar usar 'decimal' si la precisión es crítica
	RutCliente   *string   `gorm:"column:rut_cliente;type:VARCHAR(20)"`      // Puede ser nulo
}

func (BoletaFactura) TableName() string {
	return "boletas_facturas"
}
