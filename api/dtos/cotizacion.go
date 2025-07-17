package dtos

import (
	"backend-ventas/api/models"
	"time"
)

// DTOs para crear cotización
type CreateCotizacionRequest struct {
	RutCliente   string  `json:"rut_cliente" binding:"required"`
	UserID       string  `json:"user_id" binding:"required"`
	TipoDespacho string  `json:"tipo_despacho"`
	CostoEnvio   float64 `json:"costo_envio"`
	Descripcion  *string `json:"descripcion"`
}

type CreateCotizacionResponse struct {
	ID        int       `json:"id"`
	FechaCrea time.Time `json:"fecha_crea"`
	Estado    string    `json:"estado"`
	Mensaje   string    `json:"mensaje"`
}

// DTOs para agregar item a cotización
type AddItemRequest struct {
	ProductoID string `json:"producto_id" binding:"required"`
	SucursalID int    `json:"sucursal_id" binding:"required"`
	Cantidad   int    `json:"cantidad" binding:"required,min=1"`
}

type AddItemResponse struct {
	CotizacionID int    `json:"cotizacion_id"`
	ProductoID   string `json:"producto_id"`
	SucursalID   int    `json:"sucursal_id"`
	Cantidad     int    `json:"cantidad"`
	Mensaje      string `json:"mensaje"`
}

// DTOs para actualizar cotización
type UpdateCotizacionRequest struct {
	CostoEnvio   *float64 `json:"costo_envio"`
	TipoDespacho *string  `json:"tipo_despacho"`
	Total        *float64 `json:"total"`
	Descripcion  *string  `json:"descripcion"`
}

// DTOs para actualizar estado cotizacion
type UpdateEstadoCotizacionRequest struct {
	Estado string `json:"estado"`
}

// DTOs para respuestas de items
type CotizacionItemResponse struct {
	CotizacionID int              `json:"cotizacion_id"`
	ProductoID   string           `json:"producto_id"`
	SucursalID   int              `json:"sucursal_id"`
	Cantidad     int              `json:"cantidad"`
	Producto     *models.Producto `json:"producto,omitempty"`
	Sucursal     *models.Sucursal `json:"sucursal,omitempty"`
}

type CotizacionResponse struct {
	ID           int                      `json:"id"`
	FechaCrea    time.Time                `json:"fecha_crea"`
	Estado       string                   `json:"estado"`
	CostoEnvio   float64                  `json:"costo_envio"`
	RutCliente   string                   `json:"rut_cliente"`
	UserID       string                   `json:"user_id"`
	TipoDespacho string                   `json:"tipo_despacho"`
	EstadoPago   string                   `json:"estado_pago"`
	Total        *float64                 `json:"total"`
	Descripcion  *string                  `json:"descripcion"`
	Cliente      *models.Cliente          `json:"cliente,omitempty"`
	Usuario      *models.Usuario          `json:"usuario,omitempty"`
	Items        []CotizacionItemResponse `json:"items"`
	TotalItems   int                      `json:"total_items"`
	TotalPrecio  float64                  `json:"total_precio"`
}

// DTOs para listar cotizaciones
type CotizacionListResponse struct {
	ID          int             `json:"id"`
	FechaCrea   time.Time       `json:"fecha_crea"`
	Estado      string          `json:"estado"`
	RutCliente  string          `json:"rut_cliente"`
	Descripcion *string         `json:"descripcion"`
	Cliente     *models.Cliente `json:"cliente,omitempty"`
	TotalItems  int             `json:"total_items"`
	TotalPrecio float64         `json:"total_precio"`
	EstadoPago  string          `json:"estado_pago"`
}

// DTOs para preview de cotización
type PreviewCotizacionRequest struct {
	IssuedAt time.Time `json:"issued_at"`
	Subtotal float64   `json:"subtotal"`
	Tax      float64   `json:"tax"`
	Total    float64   `json:"total"`
}

type PreviewCotizacionResponse struct {
	ID                        int       `json:"id"`
	CotizacionID              *int      `json:"cotizacion_id"`
	IssuedAt                  time.Time `json:"issued_at"`
	Subtotal                  float64   `json:"subtotal"`
	Tax                       float64   `json:"tax"`
	Total                     float64   `json:"total"`
	PaymentStatus             string    `json:"payment_status"`
	EstadoPago                string    `json:"estado_pago"`
	SuccessfulPaymentIntentID *int      `json:"successful_payment_intent_id"`
	CreatedAt                 time.Time `json:"created_at"`
	UpdatedAt                 time.Time `json:"updated_at"`
	Mensaje                   string    `json:"mensaje"`
}

// DTO simplificado para cotización con datos específicos
type CotizacionSimplificadaResponse struct {
	// Datos de la cotización
	ID           int       `json:"id"`
	FechaCrea    time.Time `json:"fecha_crea"`
	Estado       string    `json:"estado"`
	CostoEnvio   float64   `json:"costo_envio"`
	UserID       string    `json:"user_id"`
	Nombre       string    `json:"nombre"` // Nombre del usuario
	TipoDespacho string    `json:"tipo_despacho"`
	Descripcion  *string   `json:"descripcion"`
	EstadoPago   string    `json:"estado_pago"` // Estado del pago

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
		SKU      string `json:"sku"`
		Nombre   string `json:"nombre"`
		Cantidad int    `json:"cantidad"`
	} `json:"items"`

	// Totales
	TotalItems  int     `json:"total_items"`
	TotalPrecio float64 `json:"total_precio"`
}

// PARA CHECKOUT RESUMEN DETALLE COTIZACION
// api/dtos/cotizacion_checkout_dto.go
type CheckoutCotizacionResponse struct {
	ID           int               `json:"id"`
	FechaCrea    string            `json:"fecha_crea"`
	Estado       string            `json:"estado"`
	CostoEnvio   float64           `json:"costo_envio"`
	TipoDespacho string            `json:"tipo_despacho"`
	EstadoPago   string            `json:"estado_pago"`
	Cliente      Cliente           `json:"cliente"`
	Usuario      Usuario           `json:"usuario"`
	Direccion    DireccionCliente  `json:"direccion"`
	Items        []CheckoutItemDTO `json:"items"`
	SubtotalNeto float64           `json:"subtotal_neto"`
	IVA          float64           `json:"iva"`
	Total        float64           `json:"total"`
	PreviewID    *int              `json:"preview_id,omitempty"`
}

type CheckoutItemDTO struct {
	SKU        string  `json:"sku"`
	Nombre     string  `json:"nombre"`
	Cantidad   int     `json:"cantidad"`
	PrecioUnit float64 `json:"precio_unitario"`
	Subtotal   float64 `json:"subtotal"`
	Sucursal   string  `json:"sucursal"`
}
