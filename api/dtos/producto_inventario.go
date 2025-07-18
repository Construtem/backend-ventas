package dtos

type SucursalStockDTO struct {
	SucursalID int     `json:"sucursal_id"`
	Nombre     string  `json:"nombre"`
	TipoID     int     `json:"tipo_id"`
	Stock      int     `json:"stock"`
	Descuento  float64 `json:"descuento"`
}

type ProductoInventarioDTO struct {
	SKU               string             `json:"sku"`
	Nombre            string             `json:"nombre"`
	Descripcion       string             `json:"descripcion"`
	Precio            float64            `json:"precio"`
	StockSucursal     int                `json:"stock_sucursal"`
	DescuentoSucursal float64            `json:"descuento_sucursal"`
	Bodegas           []SucursalStockDTO `json:"bodegas"`
	TotalStockBodegas int                `json:"total_stock_bodegas"`
}

type PaginatedProductos struct {
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	TotalItems int64 `json:"total_items"`
	TotalPages int   `json:"total_pages"`

	SucursalID int                     `json:"sucursal_id"`
	Productos  []ProductoInventarioDTO `json:"productos"`
}
