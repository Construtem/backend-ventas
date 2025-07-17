package models

// Sucursales representa una fila de la tabla de ubicaciones/sucursales
type Sucursales struct {
	ID        int    `db:"id" json:"id"`
	Nombre    string `db:"nombre" json:"nombre"`
	Telefono  string `db:"telefono" json:"telefono"`
	TipoID    int    `db:"tipo_id" json:"tipo_id"`
	Direccion string `db:"direccion" json:"direccion"`
	Comuna    string `db:"comuna" json:"comuna"`
	Ciudad    string `db:"ciudad" json:"ciudad"`
}
