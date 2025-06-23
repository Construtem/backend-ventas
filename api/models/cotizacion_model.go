package models

import "time"

type Cotizacion struct {
	ID                uint                `gorm:"primaryKey"`
	Fecha             time.Time           `gorm:"type:date;not null" json:"fecha,omitempty"`
	ClienteID         uint                `gorm:"column:cliente_id;not null" json:"cliente_id"`
	Cliente           Cliente             `gorm:"foreignKey:ClienteID;references:ID"`
	VendedorID        uint                `gorm:"column:vendedor_id;not null" json:"vendedor_id"`
	Vendedor          Usuario             `gorm:"foreignKey:VendedorID;references:ID"`
	UbicacionID       uint                `gorm:"column:ubicacion_id;not null" json:"ubicacion_id"`
	Ubicacion         Ubicacion           `gorm:"foreignKey:UbicacionID;references:ID"`
	Estado            string              `gorm:"column:estado;type:VARCHAR(50);not null" json:"estado,omitempty"`
	AprobadaPorID     *uint               `gorm:"column:aprobada_por" json:"aprobada_por,omitempty"`
	AprobadaPor       *Usuario            `gorm:"foreignKey:AprobadaPorID;references:ID"`
	FechaAprobacion   *time.Time          `gorm:"column:fecha_aprobacion;type:date" json:"fecha_aprobacion,omitempty"`
	DetalleCotizacion []DetalleCotizacion `gorm:"foreignKey:CotizacionID;references:ID" json:"detalle_cotizacion"`
}

func (Cotizacion) TableName() string {
	return "cotizaciones"
}
