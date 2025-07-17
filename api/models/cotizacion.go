package models

import "time"

/* ─────────────────────────────────────────────  COTIZACIÓN  ───────────────────────────────────────────── */

type Cotizacion struct {
	ID           int       `gorm:"primaryKey"                                        json:"id"`
	FechaCrea    time.Time `gorm:"column:fecha_crea;default:CURRENT_TIMESTAMP"       json:"fecha_crea"`
	Estado       string    `gorm:"not null"                                          json:"estado"`
	CostoEnvio   float64   `gorm:"column:costo_envio;not null"                       json:"costo_envio"`
	RutCliente   string    `gorm:"column:rut_cliente;not null"                       json:"rut_cliente"`
	UserID       string    `gorm:"column:user_id;not null"                           json:"user_id"`
	TipoDespacho string    `gorm:"column:tipo_despacho;not null"                     json:"tipo_despacho"`
	Total        *float64  `json:"total,omitempty"`
	Descripcion  *string   `json:"descripcion,omitempty" binding:"max=1000"`
	EstadoPago   string    `gorm:"column:estado_pago;default:pendiente" json:"estado_pago"`

	Cliente           *Cliente           `gorm:"foreignKey:RutCliente;references:Rut" json:"cliente,omitempty"`
	Usuario           *Usuario           `gorm:"foreignKey:UserID"                    json:"usuario,omitempty"`
	Items             []CotizacionItem   `gorm:"foreignKey:CotizacionID"              json:"items,omitempty"`
	PreviewCotizacion *PreviewCotizacion `gorm:"foreignKey:CotizacionID"              json:"preview_cotizacion,omitempty"`
}

func (Cotizacion) TableName() string { return "cotizaciones" }

/* cotizacion_item */
type CotizacionItem struct {
	CotizacionID int    `gorm:"primaryKey;column:cotizacion_id" json:"cotizacion_id"`
	ProductoID   string `gorm:"primaryKey;column:sku" json:"producto_id"`
	SucursalID   int    `gorm:"primaryKey;column:sucursal_id" json:"sucursal_id"`
	Cantidad     int    `gorm:"not null" json:"cantidad"`

	Producto *Producto `gorm:"foreignKey:ProductoID;references:SKU" json:"producto,omitempty"`
	Sucursal *Sucursal `gorm:"foreignKey:SucursalID"                json:"sucursal,omitempty"`
}

func (CotizacionItem) TableName() string { return "cotizacion_item" }

/* preview_cotizacion */
type PreviewCotizacion struct {
	ID                        int       `gorm:"primaryKey"                    json:"id"`
	CotizacionID              *int      `gorm:"column:cotizacion_id"          json:"cotizacion_id,omitempty"`
	IssuedAt                  time.Time `gorm:"column:issued_at;not null"      json:"issued_at"`
	Subtotal                  float64   `gorm:"not null"                       json:"subtotal"`
	Tax                       float64   `gorm:"not null"                       json:"tax"`
	Total                     float64   `gorm:"not null"                       json:"total"`
	PaymentStatus             string    `gorm:"column:payment_status;default:'pending'" json:"payment_status"`
	StatusPagado              bool      `gorm:"column:status_pagado;default:false"      json:"status_pagado"`
	SuccessfulPaymentIntentID *int      `gorm:"column:successful_payment_intent_id"     json:"successful_payment_intent_id,omitempty"`
	CreatedAt                 time.Time `gorm:"column:created_at;default:now()"         json:"created_at"`
	UpdatedAt                 time.Time `gorm:"column:updated_at;default:now()"         json:"updated_at"`

	Cotizacion              *Cotizacion    `gorm:"foreignKey:CotizacionID"              json:"cotizacion,omitempty"`
	SuccessfulPaymentIntent *PaymentIntent `gorm:"foreignKey:SuccessfulPaymentIntentID" json:"successful_payment_intent,omitempty"`
}

func (PreviewCotizacion) TableName() string { return "preview_cotizacion" }

/* payment_intent */
type PaymentIntent struct {
	ID                int       `gorm:"primaryKey" json:"id"`
	QuotePreviewID    int       `gorm:"column:quote_preview_id;not null"`
	PagoID            int       `gorm:"column:pago_id;not null"`
	Status            string    `gorm:"not null"`
	TransactionAmount float64   `gorm:"column:transaction_amount;not null"`
	MetodoPago        *string   `gorm:"column:metodo_pago"`
	EventType         *string   `gorm:"column:event_type"`
	Payload           *string   `gorm:"type:jsonb"`
	CreatedAt         time.Time `gorm:"column:created_at;default:now()"`
	UpdatedAt         time.Time `gorm:"column:updated_at;default:now()"`

	QuotePreview *PreviewCotizacion `gorm:"foreignKey:QuotePreviewID" json:"quote_preview,omitempty"`
}

func (PaymentIntent) TableName() string { return "payment_intent" }

/* ───────────────────────────────  PRODUCTO – PK = SKU  ─────────────────────────────── */

type Producto struct {
	SKU    string  `gorm:"primaryKey;column:sku" json:"sku"`
	Nombre string  `json:"nombre"`
	Precio float64 `json:"precio"`
	// …otros campos…
}

func (Producto) TableName() string { return "productos" }

/* ───────────────────────────────  SUCURSAL  ─────────────────────────────── */

type Sucursal struct {
	ID     int    `gorm:"primaryKey" json:"id"`
	Nombre string `json:"nombre"`
}

func (Sucursal) TableName() string { return "sucursales" }

/* ───────────────────────────────  USUARIO / ROL  ─────────────────────────────── */

type Usuario struct {
	Email  string `gorm:"primaryKey" json:"email"`
	Nombre string `json:"nombre"`
	RolID  int    `gorm:"column:rol_id" json:"rol_id"`

	Rol *Rol `gorm:"foreignKey:RolID" json:"rol,omitempty"`
}

func (Usuario) TableName() string { return "usuarios" }

type Rol struct {
	ID     int    `gorm:"primaryKey" json:"id"`
	Nombre string `json:"nombre"`
}

func (Rol) TableName() string { return "roles" }

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
