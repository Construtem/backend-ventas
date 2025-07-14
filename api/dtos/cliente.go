package dtos

type Cliente struct {
	Rut         string  `json:"rut"`
	Nombre      string  `json:"nombre"`
	Telefono    *string `json:"telefono,omitempty"`
	Email       *string `json:"email,omitempty"`
	RazonSocial *string `json:"razon_social,omitempty"`
}

type Usuario struct {
	Nombre string `json:"nombre"`
	Email  string `json:"email"`
	RolID  int    `json:"rol_id"`
}

type DireccionCliente struct {
	Direccion string `json:"direccion"`
	Comuna    string `json:"comuna"`
	Ciudad    string `json:"ciudad"`
}
