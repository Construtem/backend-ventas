package dtos

import "time"

// HistorialCotizacionResponse representa un registro del historial de una cotizacion.
type HistorialCotizacionResponse struct {
	ID       uint                      `json:"id"`
	Fecha    time.Time                 `json:"fecha"`
	Accion   string                    `json:"accion"`
	Usuario  *UsuarioAprobadorResponse `json:"usuario,omitempty"`
	Detalles *string                   `json:"detalles,omitempty"`
}
