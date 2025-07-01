package models

import (
	"time"
)

type Cotizacion struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	FechaCrea    time.Time `gorm:"column:fecha_crea;autoCreateTime" json:"fecha_crea"`
	Estado       string    `json:"estado"`
	CostoEnvio   float64   `gorm:"column:costo_envio" json:"costo_envio"`
	ClienteID    uint      `gorm:"column:cliente_id" json:"cliente_id"`
	UserID       string    `gorm:"column:user_id" json:"user_id"`
	TipoDespacho string    `gorm:"column:tipo_despacho" json:"tipo_despacho"`

	// Relaciones
	Cliente  *Cliente            `gorm:"foreignKey:ClienteID" json:"cliente,omitempty"`
	Usuario  *Usuario            `gorm:"foreignKey:UserID" json:"usuario,omitempty"`
	Items    []CotizacionItem    `gorm:"foreignKey:CotizacionID" json:"items,omitempty"`
	Previews []PreviewCotizacion `gorm:"foreignKey:CotizacionID" json:"previews,omitempty"`
}

func (Cotizacion) TableName() string {
	return "cotizaciones"
}

type CotizacionItem struct {
	CotizacionID uint `gorm:"primaryKey;column:cotizacion_id" json:"cotizacion_id"`
	ProductoID   uint `gorm:"primaryKey;column:producto_id" json:"producto_id"`
	SucursalID   uint `gorm:"primaryKey;column:sucursal_id" json:"sucursal_id"`
	Cantidad     int  `json:"cantidad"`

	// Relaciones
	Producto *Producto `gorm:"foreignKey:ProductoID" json:"producto,omitempty"`
	Sucursal *Sucursal `gorm:"foreignKey:SucursalID" json:"sucursal,omitempty"`
}

func (CotizacionItem) TableName() string {
	return "cotizacion_item"
}

type PreviewCotizacion struct {
	ID              uint      `gorm:"primaryKey" json:"id"`
	CotizacionID    uint      `gorm:"column:cotizacion_id" json:"cotizacion_id"`
	TokenAcceso     string    `gorm:"column:token_acceso" json:"token_acceso"`
	FechaExpiracion time.Time `gorm:"column:fecha_expiracion" json:"fecha_expiracion"`
	StatusPagoID    uint      `gorm:"column:status_pago_id" json:"status_pago_id"`

	// Relaciones
	Cotizacion *Cotizacion    `gorm:"foreignKey:CotizacionID" json:"cotizacion,omitempty"`
	StatusPago *PaymentStatus `gorm:"foreignKey:StatusPagoID" json:"status_pago,omitempty"`
}

func (PreviewCotizacion) TableName() string {
	return "preview_cotizacion"
}

type PaymentStatus struct {
	ID     uint   `gorm:"primaryKey" json:"id"`
	Nombre string `json:"nombre"`
}

func (PaymentStatus) TableName() string {
	return "payment_status"
}

type Producto struct {
	SKU         uint    `gorm:"primaryKey;column:sku" json:"sku"`
	Nombre      string  `json:"nombre"`
	Descripcion string  `json:"descripcion"`
	Marca       string  `json:"marca"`
	Peso        float64 `json:"peso"`
	Largo       float64 `json:"largo"`
	Ancho       float64 `json:"ancho"`
	Alto        float64 `json:"alto"`
	Precio      float64 `json:"precio"`
}

func (Producto) TableName() string {
	return "productos"
}

type Sucursal struct {
	ID        uint   `gorm:"primaryKey" json:"id"`
	Nombre    string `json:"nombre"`
	Telefono  string `json:"telefono"`
	TipoID    uint   `gorm:"column:tipo_id" json:"tipo_id"`
	Direccion string `json:"direccion"`
	Comuna    string `json:"comuna"`
	Ciudad    string `json:"ciudad"`
}

func (Sucursal) TableName() string {
	return "sucursales"
}

type Usuario struct {
	Email  string `gorm:"primaryKey" json:"email"`
	Nombre string `gorm:"not null" json:"nombre"`
	RolID  uint   `gorm:"column:rol_id;not null" json:"rol_id"`

	// Relaciones
	Rol *Rol `gorm:"foreignKey:RolID" json:"rol,omitempty"`
}

func (Usuario) TableName() string {
	return "usuarios"
}

type StockSucursal struct {
	SucursalID uint    `gorm:"primaryKey;column:sucursal_id" json:"sucursal_id"`
	ProductoID uint    `gorm:"primaryKey;column:producto_id" json:"producto_id"`
	Cantidad   int     `json:"cantidad"`
	Descuento  float64 `json:"descuento"`

	// Relaciones
	Sucursal *Sucursal `gorm:"foreignKey:SucursalID" json:"sucursal,omitempty"`
	Producto *Producto `gorm:"foreignKey:ProductoID" json:"producto,omitempty"`
}

func (StockSucursal) TableName() string {
	return "stock_sucursal"
}

// Rol representa los roles de usuario
type Rol struct {
	ID     uint   `gorm:"primaryKey" json:"id"`
	Nombre string `gorm:"not null" json:"nombre"`
}

func (Rol) TableName() string {
	return "roles"
}

// TipoSucursal representa el tipo de sucursal
type TipoSucursal struct {
	ID     uint   `gorm:"primaryKey" json:"id"`
	Nombre string `gorm:"not null" json:"nombre"`
}

func (TipoSucursal) TableName() string {
	return "tipo_sucursal"
}

// PaymentIntent representa los intentos de pago
type PaymentIntent struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	PreviewID     uint      `gorm:"column:preview_id;not null" json:"preview_id"`
	ProveedorPago string    `gorm:"column:proveedor_pago;not null" json:"proveedor_pago"`
	Monto         float64   `gorm:"not null" json:"monto"`
	Moneda        string    `gorm:"not null;default:'CLP'" json:"moneda"`
	FechaIntento  time.Time `gorm:"column:fecha_intento;not null;default:CURRENT_TIMESTAMP" json:"fecha_intento"`

	// Relaciones
	Preview *PreviewCotizacion `gorm:"foreignKey:PreviewID" json:"preview,omitempty"`
}

func (PaymentIntent) TableName() string {
	return "payment_intent"
}

// Factura representa las facturas generadas
type Factura struct {
	ID                  uint      `gorm:"primaryKey" json:"id"`
	PreviewCotizacionID uint      `gorm:"column:preview_cotizacion_id;not null" json:"preview_cotizacion_id"`
	ClienteID           uint      `gorm:"column:cliente_id;not null" json:"cliente_id"`
	FechaEmision        time.Time `gorm:"column:fecha_emision;not null;default:CURRENT_TIMESTAMP" json:"fecha_emision"`
	Monto               float64   `gorm:"not null" json:"monto"`
	EstadoPago          string    `gorm:"column:estado_pago;not null;default:'PENDIENTE'" json:"estado_pago"`

	// Relaciones
	PreviewCotizacion *PreviewCotizacion `gorm:"foreignKey:PreviewCotizacionID" json:"preview_cotizacion,omitempty"`
	Cliente           *Cliente           `gorm:"foreignKey:ClienteID" json:"cliente,omitempty"`
}

func (Factura) TableName() string {
	return "facturas"
}

// TipoCamion representa el tipo de camión
type TipoCamion struct {
	ID         uint    `gorm:"primaryKey" json:"id"`
	Volumen    float64 `gorm:"not null" json:"volumen"`
	PesoMaximo float64 `gorm:"not null" json:"peso_maximo"`
}

func (TipoCamion) TableName() string {
	return "tipo_camion"
}

// Camion representa un camión
type Camion struct {
	ID      uint   `gorm:"primaryKey" json:"id"`
	Patente string `gorm:"uniqueIndex;not null" json:"patente"`
	TipoID  uint   `gorm:"column:tipo_id;not null" json:"tipo_id"`
	Activo  bool   `gorm:"not null;default:true" json:"activo"`

	// Relaciones
	Tipo *TipoCamion `gorm:"foreignKey:TipoID" json:"tipo,omitempty"`
}

func (Camion) TableName() string {
	return "camiones"
}

// Despacho representa un despacho
type Despacho struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	CotizacionID  uint      `gorm:"column:cotizacion_id;not null" json:"cotizacion_id"`
	CamionID      uint      `gorm:"column:camion_id;not null" json:"camion_id"`
	Origen        string    `gorm:"not null" json:"origen"`
	Destino       string    `gorm:"not null" json:"destino"`
	FechaDespacho time.Time `gorm:"column:fecha_despacho;not null" json:"fecha_despacho"`
	ValorDespacho float64   `gorm:"column:valor_despacho;not null" json:"valor_despacho"`

	// Relaciones
	Cotizacion *Cotizacion        `gorm:"foreignKey:CotizacionID" json:"cotizacion,omitempty"`
	Camion     *Camion            `gorm:"foreignKey:CamionID" json:"camion,omitempty"`
	Productos  []ProductoDespacho `gorm:"foreignKey:DespachoID" json:"productos,omitempty"`
}

func (Despacho) TableName() string {
	return "despacho"
}

// ProductoDespacho representa los productos en un despacho
type ProductoDespacho struct {
	ProductoID uint `gorm:"primaryKey;column:producto_id" json:"producto_id"`
	DespachoID uint `gorm:"primaryKey;column:despacho_id" json:"despacho_id"`
	Cantidad   int  `gorm:"not null" json:"cantidad"`

	// Relaciones
	Producto *Producto `gorm:"foreignKey:ProductoID" json:"producto,omitempty"`
	Despacho *Despacho `gorm:"foreignKey:DespachoID" json:"despacho,omitempty"`
}

func (ProductoDespacho) TableName() string {
	return "productos_despacho"
}

// Proveedor representa un proveedor
type Proveedor struct {
	ID    uint   `gorm:"primaryKey" json:"id"`
	Marca string `gorm:"not null" json:"marca"`
}

func (Proveedor) TableName() string {
	return "proveedores"
}

// StockProveedor representa el stock de un proveedor
type StockProveedor struct {
	ProveedorID uint `gorm:"primaryKey;column:proveedor_id" json:"proveedor_id"`
	ProductoID  uint `gorm:"primaryKey;column:producto_id" json:"producto_id"`
	Stock       int  `gorm:"not null;default:0" json:"stock"`

	// Relaciones
	Proveedor *Proveedor `gorm:"foreignKey:ProveedorID" json:"proveedor,omitempty"`
	Producto  *Producto  `gorm:"foreignKey:ProductoID" json:"producto,omitempty"`
}

func (StockProveedor) TableName() string {
	return "stock_proveedor"
}
