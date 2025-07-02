package models

import "time"

// HistorialCotizacion representa los cambios relevantes de una cotizacion.
type HistorialCotizacion struct {
	ID           uint       `gorm:"primaryKey"`
	CotizacionID uint       `gorm:"column:cotizacion_id;not null"`
	Cotizacion   Cotizacion `gorm:"foreignKey:CotizacionID"`
	Accion       string     `gorm:"column:accion;type:VARCHAR(100);not null"`
	UsuarioID    *uint      `gorm:"column:usuario_id"`
	Usuario      *Usuario   `gorm:"foreignKey:UsuarioID"`
	Fecha        time.Time  `gorm:"column:fecha;autoCreateTime"`
	Detalles     *string    `gorm:"column:detalles;type:TEXT"`
}

func (HistorialCotizacion) TableName() string {
	return "historial_cotizaciones"
}
