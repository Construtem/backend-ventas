package mappers

import (
	"backend-ventas/api/dtos"
	"backend-ventas/api/models"
)

// Helper para mapear models.Cliente a dtos.ClienteResponse
func MapClienteToDTO(cliente models.Cliente) *dtos.ClienteResponse {
	if cliente.ID == 0 { // Cliente no cargado o es cero-value
		return nil
	}
	return &dtos.ClienteResponse{
		ID:        cliente.ID,
		Nombre:    cliente.Nombre,
		Telefono:  cliente.Telefono,
		Email:     cliente.Email,
		Direccion: cliente.Direccion,
	}
}

// Helper para mapear models.Usuario a dtos.VendedorResponse
func MapUsuarioToVendedorDTO(usuario models.Usuario) *dtos.VendedorResponse {
	if usuario.ID == 0 {
		return nil
	}
	return &dtos.VendedorResponse{
		ID:     usuario.ID,
		Nombre: usuario.Nombre,
		Email:  usuario.Email,
	}
}

// Helper para mapear models.Usuario a dtos.UsuarioAprobadorResponse
func MapUsuarioToAprobadorDTO(usuario models.Usuario) *dtos.UsuarioAprobadorResponse {
	if usuario.ID == 0 {
		return nil
	}
	return &dtos.UsuarioAprobadorResponse{
		ID:     usuario.ID,
		Nombre: usuario.Nombre,
		Email:  usuario.Email,
	}
}

// Helper para mapear models.Ubicacion a dtos.UbicacionResponse
func MapUbicacionToDTO(ubicacion models.Ubicacion) *dtos.UbicacionResponse {
	if ubicacion.ID == 0 {
		return nil
	}
	return &dtos.UbicacionResponse{
		ID:        ubicacion.ID,
		Nombre:    ubicacion.Nombre,
		Tipo:      ubicacion.Tipo,
		Direccion: ubicacion.Direccion,
	}
}

// Helper para mapear models.Producto a dtos.ProductoDetalleResponse
func MapProductoToDetalleDTO(producto models.Producto) *dtos.ProductoDetalleResponse {
	if producto.ID == 0 {
		return nil
	}
	return &dtos.ProductoDetalleResponse{
		ID:          producto.ID,
		Codigo:      producto.Codigo,
		Nombre:      producto.Nombre,
		Descripcion: producto.Descripcion,
		PrecioVenta: producto.PrecioVenta,
	}
}

// Helper para mapear models.DetalleCotizacion a dtos.DetalleCotizacionResponse
func MapDetalleCotizacionToDTO(detalle models.DetalleCotizacion) dtos.DetalleCotizacionResponse {
	return dtos.DetalleCotizacionResponse{
		CotizacionID:   detalle.CotizacionID,
		ProductoID:     detalle.ProductoID,
		Cantidad:       detalle.Cantidad,
		PrecioUnitario: detalle.PrecioUnitario,
		Producto:       MapProductoToDetalleDTO(detalle.Producto), // Usa el helper para Producto
	}
}

// Helper principal para mapear models.Cotizacion a dtos.CotizacionResponse
func MapCotizacionToDTO(cotizacion models.Cotizacion) dtos.CotizacionResponse {
	detalleResponses := make([]dtos.DetalleCotizacionResponse, len(cotizacion.DetalleCotizacion))
	for i, det := range cotizacion.DetalleCotizacion {
		detalleResponses[i] = MapDetalleCotizacionToDTO(det) // Usa el helper para Detalle
	}

	var aprobadaPorDTO *dtos.UsuarioAprobadorResponse
	if cotizacion.AprobadaPor != nil {
		aprobadaPorDTO = MapUsuarioToAprobadorDTO(*cotizacion.AprobadaPor)
	}

	return dtos.CotizacionResponse{
		ID:                cotizacion.ID,
		Fecha:             cotizacion.Fecha,
		Cliente:           MapClienteToDTO(cotizacion.Cliente),
		Vendedor:          MapUsuarioToVendedorDTO(cotizacion.Vendedor),
		Ubicacion:         MapUbicacionToDTO(cotizacion.Ubicacion),
		Estado:            cotizacion.Estado,
		AprobadaPorID:     cotizacion.AprobadaPorID,
		AprobadaPor:       aprobadaPorDTO,
		FechaAprobacion:   cotizacion.FechaAprobacion,
		DetalleCotizacion: detalleResponses,
	}
}
