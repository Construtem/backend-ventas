package models

import (
	"time"
)

// // ========================================
// // MODELOS DE CLIENTE (Solo para respuestas JSON)
// // ========================================

// // Cliente representa un cliente en el sistema
// type Cliente struct {
// 	Nombre      string  `gorm:"not null" json:"nombre"`
// 	Telefono    *string `json:"telefono"`
// 	Email       *string `json:"email"`
// 	RazonSocial *string `gorm:"column:razon_social" json:"razon_social"`
// 	Rut         string  `gorm:"primaryKey" json:"rut"`
// 	TipoID      int     `gorm:"column:tipo_id;not null" json:"tipo_id"`

// 	// Relaciones
// 	TipoCliente  *TipoCliente `gorm:"foreignKey:TipoID" json:"tipo_cliente,omitempty"`
// 	Direcciones  []DirCliente `gorm:"foreignKey:RutCliente" json:"direcciones,omitempty"`
// 	Cotizaciones []Cotizacion `gorm:"foreignKey:RutCliente" json:"cotizaciones,omitempty"`
// }

// // TableName especifica el nombre de la tabla en la base de datos
// func (Cliente) TableName() string {
// 	return "clientes"
// }

// // TipoCliente representa el tipo de cliente
// type TipoCliente struct {
// 	ID     int    `gorm:"primaryKey" json:"id"`
// 	Nombre string `gorm:"not null" json:"nombre"`
// }

// func (TipoCliente) TableName() string {
// 	return "tipo_cliente"
// }

// // DirCliente representa las direcciones de un cliente
// type DirCliente struct {
// 	ID         int    `gorm:"primaryKey" json:"id"`
// 	RutCliente string `gorm:"column:rut_cliente;not null" json:"rut_cliente"`
// 	Direccion  string `gorm:"not null" json:"direccion"`
// 	Comuna     string `gorm:"not null" json:"comuna"`
// 	Ciudad     string `gorm:"not null" json:"ciudad"`

// 	// Relaciones
// 	Cliente *Cliente `gorm:"foreignKey:RutCliente" json:"cliente,omitempty"`
// }

// func (DirCliente) TableName() string {
// 	return "dir_cliente"
// }

// ========================================
// MODELOS DE COTIZACIÓN
// ========================================

type Cotizacion struct {
	ID           int       `gorm:"primaryKey" json:"id"`
	FechaCrea    time.Time `gorm:"column:fecha_crea;default:CURRENT_TIMESTAMP" json:"fecha_crea"`
	Estado       string    `gorm:"not null" json:"estado"`
	CostoEnvio   float64   `gorm:"column:costo_envio;not null" json:"costo_envio"`
	RutCliente   string    `gorm:"column:rut_cliente;not null" json:"rut_cliente"`
	UserID       string    `gorm:"column:user_id;not null" json:"user_id"`
	TipoDespacho string    `gorm:"column:tipo_despacho;not null" json:"tipo_despacho"`
	Total        *float64  `json:"total"`

	// Relaciones
	Cliente  *Cliente            `gorm:"foreignKey:RutCliente;references:Rut" json:"cliente,omitempty"`
	Usuario  *Usuario            `gorm:"foreignKey:UserID" json:"usuario,omitempty"`
	Items    []CotizacionItem    `gorm:"foreignKey:CotizacionID" json:"items,omitempty"`
	Previews []PreviewCotizacion `gorm:"foreignKey:CotizacionID" json:"previews,omitempty"`
}

func (Cotizacion) TableName() string {
	return "cotizaciones"
}

type CotizacionItem struct {
	CotizacionID int    `gorm:"primaryKey;column:cotizacion_id" json:"cotizacion_id"`
	ProductoID   string `gorm:"primaryKey;column:sku" json:"producto_id"`
	SucursalID   int    `gorm:"primaryKey;column:sucursal_id" json:"sucursal_id"`
	Cantidad     int    `gorm:"not null" json:"cantidad"`

	// Relaciones
	Producto *Producto `gorm:"foreignKey:ProductoID" json:"producto,omitempty"`
	Sucursal *Sucursal `gorm:"foreignKey:SucursalID" json:"sucursal,omitempty"`
}

func (CotizacionItem) TableName() string {
	return "cotizacion_item"
}

type PreviewCotizacion struct {
	ID                        int       `gorm:"primaryKey" json:"id"`
	CotizacionID              *int      `gorm:"column:cotizacion_id" json:"cotizacion_id"`
	IssuedAt                  time.Time `gorm:"column:issued_at;not null" json:"issued_at"`
	Subtotal                  float64   `gorm:"not null" json:"subtotal"`
	Tax                       float64   `gorm:"not null" json:"tax"`
	Total                     float64   `gorm:"not null" json:"total"`
	PaymentStatus             string    `gorm:"column:payment_status;not null;default:'pending'" json:"payment_status"`
	StatusPagado              bool      `gorm:"column:status_pagado;not null;default:false" json:"status_pagado"`
	SuccessfulPaymentIntentID *int      `gorm:"column:successful_payment_intent_id" json:"successful_payment_intent_id"`
	CreatedAt                 time.Time `gorm:"column:created_at;default:now()" json:"created_at"`
	UpdatedAt                 time.Time `gorm:"column:updated_at;default:now()" json:"updated_at"`

	// Relaciones
	Cotizacion              *Cotizacion    `gorm:"foreignKey:CotizacionID" json:"cotizacion,omitempty"`
	SuccessfulPaymentIntent *PaymentIntent `gorm:"foreignKey:SuccessfulPaymentIntentID" json:"successful_payment_intent,omitempty"`
}

func (PreviewCotizacion) TableName() string {
	return "preview_cotizacion"
}

type PaymentIntent struct {
	ID                int       `gorm:"primaryKey" json:"id"`
	QuotePreviewID    int       `gorm:"column:quote_preview_id;not null" json:"quote_preview_id"`
	PagoID            int       `gorm:"column:pago_id;not null" json:"pago_id"`
	Status            string    `gorm:"not null" json:"status"`
	TransactionAmount float64   `gorm:"column:transaction_amount;not null" json:"transaction_amount"`
	MetodoPago        *string   `gorm:"column:metodo_pago" json:"metodo_pago"`
	EventType         *string   `gorm:"column:event_type" json:"event_type"`
	Payload           *string   `gorm:"type:jsonb" json:"payload"`
	CreatedAt         time.Time `gorm:"column:created_at;default:now()" json:"created_at"`
	UpdatedAt         time.Time `gorm:"column:updated_at;default:now()" json:"updated_at"`

	// Relaciones
	QuotePreview *PreviewCotizacion `gorm:"foreignKey:QuotePreviewID" json:"quote_preview,omitempty"`
}

func (PaymentIntent) TableName() string {
	return "payment_intent"
}

// ========================================
// MODELOS DE PRODUCTOS Y CATEGORÍAS
// ========================================

type Producto struct {
	SKU         string  `gorm:"primaryKey;column:sku" json:"sku"`
	Nombre      string  `gorm:"not null" json:"nombre"`
	Descripcion string  `gorm:"not null" json:"descripcion"`
	Marca       string  `gorm:"not null" json:"marca"`
	Peso        float64 `gorm:"not null" json:"peso"`
	Largo       float64 `gorm:"not null" json:"largo"`
	Ancho       float64 `gorm:"not null" json:"ancho"`
	Alto        float64 `gorm:"not null" json:"alto"`
	Precio      float64 `gorm:"not null" json:"precio"`
	ID          int64   `gorm:"column:id" json:"id"`
	Codigo      *string `json:"codigo"`
	CategoriaID *int    `gorm:"column:categoria_id" json:"categoria_id"`
	Estado      *bool   `gorm:"default:true" json:"estado"`
	Stock       int     `gorm:"not null" json:"stock"`

	// Relaciones
	Categoria *Categoria `gorm:"foreignKey:CategoriaID" json:"categoria,omitempty"`
}

func (Producto) TableName() string {
	return "productos"
}

type Categoria struct {
	ID     int    `gorm:"primaryKey" json:"id"`
	Nombre string `gorm:"not null" json:"nombre"`
}

func (Categoria) TableName() string {
	return "categoria"
}

// ========================================
// MODELOS DE SUCURSALES
// ========================================

type Sucursal struct {
	ID        int     `gorm:"primaryKey" json:"id"`
	Nombre    string  `gorm:"not null" json:"nombre"`
	Telefono  string  `gorm:"not null" json:"telefono"`
	TipoID    int     `gorm:"column:tipo_id;not null" json:"tipo_id"`
	Direccion *string `json:"direccion"`
	Comuna    *string `json:"comuna"`
	Ciudad    *string `json:"ciudad"`

	// Relaciones
	Tipo *TipoSucursal `gorm:"foreignKey:TipoID" json:"tipo,omitempty"`
}

func (Sucursal) TableName() string {
	return "sucursales"
}

type TipoSucursal struct {
	ID     int    `gorm:"primaryKey" json:"id"`
	Nombre string `gorm:"not null" json:"nombre"`
}

func (TipoSucursal) TableName() string {
	return "tipo_sucursal"
}

// ========================================
// MODELOS DE USUARIOS Y ROLES
// ========================================

type Usuario struct {
	Email  string `gorm:"primaryKey" json:"email"`
	Nombre string `gorm:"not null" json:"nombre"`
	RolID  int    `gorm:"column:rol_id;not null" json:"rol_id"`

	// Relaciones
	Rol *Rol `gorm:"foreignKey:RolID" json:"rol,omitempty"`
}

func (Usuario) TableName() string {
	return "usuarios"
}

type Rol struct {
	ID     int    `gorm:"primaryKey" json:"id"`
	Nombre string `gorm:"not null" json:"nombre"`
}

func (Rol) TableName() string {
	return "roles"
}

// ========================================
// MODELOS DE STOCK
// ========================================

type StockSucursal struct {
	SucursalID int     `gorm:"primaryKey;column:sucursal_id" json:"sucursal_id"`
	ProductoID string  `gorm:"primaryKey;column:producto_id" json:"producto_id"`
	Cantidad   int     `gorm:"not null" json:"cantidad"`
	Descuento  float64 `gorm:"default:0" json:"descuento"`

	// Relaciones
	Sucursal *Sucursal `gorm:"foreignKey:SucursalID" json:"sucursal,omitempty"`
	Producto *Producto `gorm:"foreignKey:ProductoID" json:"producto,omitempty"`
}

func (StockSucursal) TableName() string {
	return "stock_sucursal"
}

// ========================================
// MODELOS DE FACTURACIÓN
// ========================================

type Factura struct {
	ID                  int        `gorm:"primaryKey" json:"id"`
	QuotePreviewID      int        `gorm:"column:quote_preview_id;not null" json:"quote_preview_id"`
	CotizacionID        *int       `gorm:"column:cotizacion_id" json:"cotizacion_id"`
	RutCliente          *string    `gorm:"column:rut_cliente" json:"rut_cliente"`
	TipoDocumento       *string    `gorm:"column:tipo_documento;default:'Factura Electrónica'" json:"tipo_documento"`
	Folio               string     `gorm:"not null" json:"folio"`
	FechaEmision        time.Time  `gorm:"column:fecha_emision;not null" json:"fecha_emision"`
	FechaVencimiento    *time.Time `gorm:"column:fecha_vencimiento" json:"fecha_vencimiento"`
	TimbreElectronico   *string    `gorm:"column:timbre_electronico" json:"timbre_electronico"`
	SiiIndicacion       *string    `gorm:"column:sii_indicacion" json:"sii_indicacion"`
	FraseLegal          *string    `gorm:"column:frase_legal" json:"frase_legal"`
	RutEmisor           string     `gorm:"column:rut_emisor;not null" json:"rut_emisor"`
	RazonSocialEmisor   string     `gorm:"column:razon_social_emisor;not null" json:"razon_social_emisor"`
	GiroEmisor          *string    `gorm:"column:giro_emisor" json:"giro_emisor"`
	DireccionEmisor     *string    `gorm:"column:direccion_emisor" json:"direccion_emisor"`
	ComunaEmisor        *string    `gorm:"column:comuna_emisor" json:"comuna_emisor"`
	CiudadEmisor        *string    `gorm:"column:ciudad_emisor" json:"ciudad_emisor"`
	TelefonoEmisor      *string    `gorm:"column:telefono_emisor" json:"telefono_emisor"`
	EmailEmisor         *string    `gorm:"column:email_emisor" json:"email_emisor"`
	RutReceptor         string     `gorm:"column:rut_receptor;not null" json:"rut_receptor"`
	RazonSocialReceptor string     `gorm:"column:razon_social_receptor;not null" json:"razon_social_receptor"`
	GiroReceptor        *string    `gorm:"column:giro_receptor" json:"giro_receptor"`
	DireccionReceptor   *string    `gorm:"column:direccion_receptor" json:"direccion_receptor"`
	ComunaReceptor      *string    `gorm:"column:comuna_receptor" json:"comuna_receptor"`
	CiudadReceptor      *string    `gorm:"column:ciudad_receptor" json:"ciudad_receptor"`
	ContactoReceptor    *string    `gorm:"column:contacto_receptor" json:"contacto_receptor"`
	Items               string     `gorm:"type:jsonb;not null" json:"items"`
	SubtotalNeto        float64    `gorm:"column:subtotal_neto;not null" json:"subtotal_neto"`
	Iva19               float64    `gorm:"not null" json:"iva19"`
	IvaRetenido         float64    `gorm:"column:iva_retenido;default:0" json:"iva_retenido"`
	TotalFinal          float64    `gorm:"column:total_final;not null" json:"total_final"`
	UrlPdf              *string    `gorm:"column:url_pdf" json:"url_pdf"`
	UrlVerificacion     *string    `gorm:"column:url_verificacion" json:"url_verificacion"`
	CreatedAt           time.Time  `gorm:"column:created_at;default:now()" json:"created_at"`
	UpdatedAt           time.Time  `gorm:"column:updated_at;default:now()" json:"updated_at"`

	// Relaciones
	QuotePreview *PreviewCotizacion `gorm:"foreignKey:QuotePreviewID" json:"quote_preview,omitempty"`
	Cotizacion   *Cotizacion        `gorm:"foreignKey:CotizacionID" json:"cotizacion,omitempty"`
	Cliente      *Cliente           `gorm:"foreignKey:RutCliente" json:"cliente,omitempty"`
}

func (Factura) TableName() string {
	return "facturas"
}

// ========================================
// MODELOS DE DESPACHO
// ========================================

type TipoCamion struct {
	ID         int     `gorm:"primaryKey" json:"id"`
	Volumen    float64 `gorm:"not null" json:"volumen"`
	PesoMaximo float64 `gorm:"not null" json:"peso_maximo"`
}

func (TipoCamion) TableName() string {
	return "tipo_camion"
}

type Camion struct {
	ID      int    `gorm:"primaryKey" json:"id"`
	Patente string `gorm:"uniqueIndex;not null" json:"patente"`
	TipoID  int    `gorm:"column:tipo_id;not null" json:"tipo_id"`
	Activo  bool   `gorm:"not null;default:true" json:"activo"`

	// Relaciones
	Tipo *TipoCamion `gorm:"foreignKey:TipoID" json:"tipo,omitempty"`
}

func (Camion) TableName() string {
	return "camiones"
}

type Despacho struct {
	ID            int       `gorm:"primaryKey" json:"id"`
	CotizacionID  int       `gorm:"column:cotizacion_id;not null" json:"cotizacion_id"`
	CamionID      int       `gorm:"column:camion_id;not null" json:"camion_id"`
	Origen        int       `gorm:"not null" json:"origen"`
	Destino       int       `gorm:"not null" json:"destino"`
	FechaDespacho time.Time `gorm:"column:fecha_despacho;not null" json:"fecha_despacho"`
	ValorDespacho float64   `gorm:"column:valor_despacho;not null" json:"valor_despacho"`
	Estado        *string   `json:"estado"`

	// Relaciones
	Cotizacion *Cotizacion        `gorm:"foreignKey:CotizacionID" json:"cotizacion,omitempty"`
	CamionObj  *Camion            `gorm:"foreignKey:CamionID" json:"camion_obj,omitempty"`
	OrigenObj  *Sucursal          `gorm:"foreignKey:Origen" json:"origen_obj,omitempty"`
	DestinoObj *DirCliente        `gorm:"foreignKey:Destino" json:"destino_obj,omitempty"`
	Productos  []ProductoDespacho `gorm:"foreignKey:DespachoID" json:"productos,omitempty"`
}

func (Despacho) TableName() string {
	return "despacho"
}

type ProductoDespacho struct {
	DespachoID int    `gorm:"primaryKey;column:despacho_id" json:"despacho_id"`
	ProductoID string `gorm:"primaryKey;column:producto_id" json:"producto_id"`
	Cantidad   int    `gorm:"not null" json:"cantidad"`

	// Relaciones
	Producto *Producto `gorm:"foreignKey:ProductoID" json:"producto,omitempty"`
	Despacho *Despacho `gorm:"foreignKey:DespachoID" json:"despacho,omitempty"`
}

func (ProductoDespacho) TableName() string {
	return "productos_despacho"
}

// ========================================
// MODELOS DE PROVEEDORES
// ========================================

type Proveedor struct {
	ID    int    `gorm:"primaryKey" json:"id"`
	Marca string `gorm:"not null" json:"marca"`
}

func (Proveedor) TableName() string {
	return "proveedores"
}

type StockProveedor struct {
	ProveedorID int    `gorm:"primaryKey;column:proveedor_id" json:"proveedor_id"`
	ProductoID  string `gorm:"primaryKey;column:producto_id" json:"producto_id"`
	Stock       int    `gorm:"not null" json:"stock"`

	// Relaciones
	Proveedor *Proveedor `gorm:"foreignKey:ProveedorID" json:"proveedor,omitempty"`
	Producto  *Producto  `gorm:"foreignKey:ProductoID" json:"producto,omitempty"`
}

func (StockProveedor) TableName() string {
	return "stock_proveedor"
}
