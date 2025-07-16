package dtos

// UbicacionResponse es lo que devolvemos al cliente para cada sucursal
type UbicacionResponse struct {
	ID        int    `json:"id"`
	Nombre    string `json:"nombre"`
	Telefono  string `json:"telefono"`
	Direccion string `json:"direccion"`
	Comuna    string `json:"comuna"`
	Ciudad    string `json:"ciudad"`
}
