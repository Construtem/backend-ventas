package controllers

import (
	"backend-ventas/api/database"
	"backend-ventas/api/dtos"
	"backend-ventas/api/models"
	_ "backend-ventas/api/models"
)

// ListarProductosConInventario obtiene la página de productos
// e incluye stock + descuento de la TIENDA seleccionada (tipo_id = 2)
// y el stock/desc. en TODAS las bodegas (tipo_id = 1).
// Además filtra fuera los productos sin stock en ninguna parte.
func ListarProductosConInventario(
	sucursalID, page, limit int,
) (dtos.PaginatedProductos, error) {
	db := database.DB

	var resp dtos.PaginatedProductos
	resp.SucursalID = sucursalID
	resp.Page = page
	resp.Limit = limit

	// Contar total productos (sin filtro de stock)
	var total int64
	if err := db.Model(&models.Producto{}).Count(&total).Error; err != nil {
		return resp, err
	}
	resp.TotalItems = total
	resp.TotalPages = int((total + int64(limit) - 1) / int64(limit))

	// Traer la “base” de productos + stock/desc. de la tienda
	offset := (page - 1) * limit
	type tiendaRow struct {
		SKU               string
		Nombre            string
		Descripcion       string
		Precio            float64
		StockSucursal     int
		DescuentoSucursal float64
	}
	var tiendaData []tiendaRow

	if err := db.Raw(`
		SELECT 
		  p.sku,
		  p.nombre,
		  p.descripcion,
		  p.precio,
		  COALESCE(ss.cantidad, 0)  AS stock_sucursal,
		  COALESCE(ss.descuento, 0) AS descuento_sucursal
		FROM productos p
		LEFT JOIN stock_sucursal ss
		  ON ss.sku = p.sku 
		 AND ss.sucursal_id = ?
		ORDER BY p.sku
		LIMIT ? OFFSET ?`,
		sucursalID, limit, offset,
	).Scan(&tiendaData).Error; err != nil {
		return resp, err
	}

	// Si no hay productos en esta página, devolvemos vacío
	if len(tiendaData) == 0 {
		resp.Productos = []dtos.ProductoInventarioDTO{}
		return resp, nil
	}

	//  Obtener stock/desc. en bodegas para estos SKUs
	var skus []string
	for _, r := range tiendaData {
		skus = append(skus, r.SKU)
	}

	type bodegaRow struct {
		SKU        string
		SucursalID int
		Nombre     string
		Stock      int
		TipoID     int
		Descuento  float64
	}
	var bodegas []bodegaRow

	if err := db.Raw(`
		SELECT 
		  ss.sku,
		  s.id           AS sucursal_id,
		  s.nombre,
		  ss.cantidad    AS stock,
		  s.tipo_id,
		  COALESCE(ss.descuento, 0) AS descuento
		FROM stock_sucursal ss
		JOIN sucursales s 
		  ON s.id = ss.sucursal_id
		WHERE s.tipo_id = 1               -- solo bodegas
		  AND ss.sku IN (?)
	`, skus).Scan(&bodegas).Error; err != nil {
		return resp, err
	}

	// Agrupar bodegas por SKU y sumar stock
	bMap := make(map[string][]dtos.SucursalStockDTO)
	totalBod := make(map[string]int)
	for _, b := range bodegas {
		entry := dtos.SucursalStockDTO{
			SucursalID: b.SucursalID,
			Nombre:     b.Nombre,
			TipoID:     b.TipoID,
			Stock:      b.Stock,
			Descuento:  b.Descuento, // ahora incluimos el % de descuento
		}
		bMap[b.SKU] = append(bMap[b.SKU], entry)
		totalBod[b.SKU] += b.Stock
	}

	// 4) Construir el DTO final filtrando sin-stock
	for _, r := range tiendaData {
		totalStock := r.StockSucursal + totalBod[r.SKU]
		if totalStock == 0 {
			continue // omitimos producto sin stock en tienda ni bodegas
		}

		resp.Productos = append(resp.Productos, dtos.ProductoInventarioDTO{
			SKU:               r.SKU,
			Nombre:            r.Nombre,
			Descripcion:       r.Descripcion,
			Precio:            r.Precio,
			StockSucursal:     r.StockSucursal,
			DescuentoSucursal: r.DescuentoSucursal,
			Bodegas:           bMap[r.SKU],
			TotalStockBodegas: totalBod[r.SKU],
		})
	}

	return resp, nil
}
