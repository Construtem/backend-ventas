package models

import (
	"time"
)

type DetallePedido struct {
	ID             uint      `gorm:"primaryKey"`
	PedidoID       uint      `gorm:"column:pedido_id;primaryKey"`
	ProductoID     uint      `gorm:"column:producto_id;primaryKey"`
	Cantidad       int       `gorm:"column:cantidad;not null"`
	PrecioUnitario float64   `gorm:"column:precio_unitario;type:NUMERIC(10,2);not null"`
	CreatedAt      time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt      time.Time `gorm:"column:updated_at;autoUpdateTime"`
	Pedido         Pedido    `gorm:"foreignKey:PedidoID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Producto       Producto  `gorm:"foreignKey:ProductoID;references:ID"`
}

func (DetallePedido) TableName() string {
	return "detalle_pedido"
}
