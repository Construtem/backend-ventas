package dtos

import (
	"backend-ventas/api/models"
	"time"
)

// DTOs para crear cotización
type CreateCotizacionRequest struct {
	ClienteID    uint    `json:"cliente_id" binding:"required"`
	UserID       string  `json:"user_id" binding:"required"`
	TipoDespacho string  `json:"tipo_despacho"`
	CostoEnvio   float64 `json:"costo_envio"`
}

type CreateCotizacionResponse struct {
	ID        uint      `json:"id"`
	FechaCrea time.Time `json:"fecha_crea"`
	Estado    string    `json:"estado"`
	Mensaje   string    `json:"mensaje"`
}

// DTOs para agregar item a cotización
type AddItemRequest struct {
	ProductoID uint `json:"producto_id" binding:"required"`
	SucursalID uint `json:"sucursal_id" binding:"required"`
	Cantidad   int  `json:"cantidad" binding:"required,min=1"`
}

type AddItemResponse struct {
	CotizacionID uint   `json:"cotizacion_id"`
	ProductoID   uint   `json:"producto_id"`
	SucursalID   uint   `json:"sucursal_id"`
	Cantidad     int    `json:"cantidad"`
	Mensaje      string `json:"mensaje"`
}

// DTOs para actualizar cotización
type UpdateCotizacionRequest struct {
	Estado       *string  `json:"estado"`
	CostoEnvio   *float64 `json:"costo_envio"`
	TipoDespacho *string  `json:"tipo_despacho"`
}

// DTOs para respuestas de items
type CotizacionItemResponse struct {
	CotizacionID uint             `json:"cotizacion_id"`
	ProductoID   uint             `json:"producto_id"`
	SucursalID   uint             `json:"sucursal_id"`
	Cantidad     int              `json:"cantidad"`
	Producto     *models.Producto `json:"producto,omitempty"`
	Sucursal     *models.Sucursal `json:"sucursal,omitempty"`
}

type CotizacionResponse struct {
	ID           uint                     `json:"id"`
	FechaCrea    time.Time                `json:"fecha_crea"`
	Estado       string                   `json:"estado"`
	CostoEnvio   float64                  `json:"costo_envio"`
	ClienteID    uint                     `json:"cliente_id"`
	UserID       string                   `json:"user_id"`
	TipoDespacho string                   `json:"tipo_despacho"`
	Cliente      *models.Cliente          `json:"cliente,omitempty"`
	Usuario      *models.Usuario          `json:"usuario,omitempty"`
	Items        []CotizacionItemResponse `json:"items"`
	TotalItems   int                      `json:"total_items"`
	TotalPrecio  float64                  `json:"total_precio"`
}

// DTOs para listar cotizaciones
type CotizacionListResponse struct {
	ID          uint            `json:"id"`
	FechaCrea   time.Time       `json:"fecha_crea"`
	Estado      string          `json:"estado"`
	ClienteID   uint            `json:"cliente_id"`
	Cliente     *models.Cliente `json:"cliente,omitempty"`
	TotalItems  int             `json:"total_items"`
	TotalPrecio float64         `json:"total_precio"`
}

// DTOs para preview de cotización
type PreviewCotizacionRequest struct {
	TokenAcceso     string    `json:"token_acceso"`
	FechaExpiracion time.Time `json:"fecha_expiracion"`
}

type PreviewCotizacionResponse struct {
	ID                uint      `json:"id"`
	CotizacionID      uint      `json:"cotizacion_id"`
	TokenAcceso       string    `json:"token_acceso"`
	FechaExpiracion   time.Time `json:"fecha_expiracion"`
	StatusPagoID      uint      `json:"status_pago_id"`
	TotalPrecio       float64   `json:"total_precio"`
	TotalConDescuento float64   `json:"total_con_descuento"`
	TotalFinal        float64   `json:"total_final"`
	Mensaje           string    `json:"mensaje"`
}

// DTO simplificado para cotización con datos específicos
type CotizacionSimplificadaResponse struct {
	// Datos de la cotización
	FechaCrea    time.Time `json:"fecha_crea"`
	Estado       string    `json:"estado"`
	CostoEnvio   float64   `json:"costo_envio"`
	UserID       string    `json:"user_id"`
	Nombre       string    `json:"nombre"` // Nombre del usuario
	TipoDespacho string    `json:"tipo_despacho"`

	// Datos del cliente
	Cliente struct {
		Nombre      string `json:"nombre"`
		Telefono    string `json:"telefono"`
		Email       string `json:"email"`
		Rut         string `json:"rut"`
		RazonSocial string `json:"razon_social"`
	} `json:"cliente"`

	// Items simplificados
	Items []struct {
		SKU      uint   `json:"sku"`
		Nombre   string `json:"nombre"`
		Cantidad int    `json:"cantidad"`
	} `json:"items"`

	// Totales
	TotalItems  int     `json:"total_items"`
	TotalPrecio float64 `json:"total_precio"`
}
