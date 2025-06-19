package dtos

import "time"

// DTOs para las relaciones anidadas
type ClienteResponse struct {
	ID        uint    `json:"id"`
	Nombre    string  `json:"nombre"`
	Telefono  *string `json:"telefono,omitempty"`
	Email     *string `json:"email,omitempty"`
	Direccion *string `json:"direccion,omitempty"`
}

type VendedorResponse struct {
	ID     uint   `json:"id"`
	Nombre string `json:"nombre"`
	Email  string `json:"email"`
	// No incluir Contrasena, RolID, etc., a menos que sea necesarios para el frontend
}

type UbicacionResponse struct {
	ID        uint    `json:"id"`
	Nombre    string  `json:"nombre"`
	Tipo      string  `json:"tipo"`
	Direccion *string `json:"direccion,omitempty"`
}

type UsuarioAprobadorResponse struct { // Usar para AprobadaPor
	ID     uint   `json:"id"`
	Nombre string `json:"nombre"`
	Email  string `json:"email"`
}

type ProductoDetalleResponse struct {
	ID          uint    `json:"id"`
	Codigo      *string `json:"codigo,omitempty"`
	Nombre      string  `json:"nombre"`
	Descripcion *string `json:"descripcion,omitempty"`
	PrecioVenta float64 `json:"precio_venta"`
	// No incluir Categoria, Proveedor, PrecioCosto, Activo, a menos que sean explícitamente requeridos
}

// DTO para el detalle de la cotización
type DetalleCotizacionResponse struct {
	CotizacionID   uint                     `json:"cotizacion_id"`
	ProductoID     uint                     `json:"producto_id"`
	Cantidad       int                      `json:"cantidad"`
	PrecioUnitario float64                  `json:"precio_unitario"`
	Producto       *ProductoDetalleResponse `json:"producto,omitempty"` // Aquí anidamos el Producto
}

// DTO principal para la Cotización
type CotizacionResponse struct {
	ID                uint                        `json:"id"`
	Fecha             time.Time                   `json:"fecha"`
	Cliente           *ClienteResponse            `json:"cliente"`   // Puntero para nil si no se carga
	Vendedor          *VendedorResponse           `json:"vendedor"`  // Puntero para nil si no se carga
	Ubicacion         *UbicacionResponse          `json:"ubicacion"` // Puntero para nil si no se carga
	Estado            string                      `json:"estado"`
	AprobadaPorID     *uint                       `json:"aprobada_por_id,omitempty"` // Mantener ID si es relevante
	AprobadaPor       *UsuarioAprobadorResponse   `json:"aprobada_por,omitempty"`    // Datos del usuario aprobador
	FechaAprobacion   *time.Time                  `json:"fecha_aprobacion,omitempty"`
	DetalleCotizacion []DetalleCotizacionResponse `json:"detalle_cotizacion,omitempty"`
}

// DTO para la respuesta de listado de cotizaciones con paginación
type CotizacionesListResponse struct {
	Data  []CotizacionResponse `json:"data"`
	Total int64                `json:"total"`
	Page  int                  `json:"page"`
	Limit int                  `json:"limit"`
}
