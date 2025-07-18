package models

/* -------------------------------------------------------------------------- */
/*  CLIENTE                                                                   */
/* -------------------------------------------------------------------------- */

type Cliente struct {
	// ID          uint      `gorm:"primaryKey" json:"id"`
	Nombre      string    `gorm:"not null" json:"nombre"`
	Telefono    *string   `json:"telefono"`
	Email       *string   `gorm:"uniqueIndex" json:"email"`
	RazonSocial *string   `gorm:"column:razon_social" json:"razon_social"`
	Rut         string    `gorm:"primaryKey;uniqueIndex:uni_clientes_rut;not null" json:"rut"`
	TipoID      uint      `gorm:"column:tipo_id;not null" json:"tipo_id"`
	Giro 		string	  `json:"giro"`
	//CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
	//UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	// Relaciones
	TipoCliente  *TipoCliente `gorm:"foreignKey:TipoID" json:"tipo_cliente,omitempty"`
	Direcciones  []DirCliente `gorm:"foreignKey:RutCliente;references:Rut" json:"direcciones,omitempty"`
	Cotizaciones []Cotizacion `gorm:"foreignKey:RutCliente;references:Rut" json:"cotizaciones,omitempty"`
}

func (Cliente) TableName() string { return "clientes" }

/* ---------- tipo_cliente ---------- */

type TipoCliente struct {
	ID     uint   `gorm:"primaryKey" json:"id"`
	Nombre string `gorm:"not null"    json:"nombre"`
}

func (TipoCliente) TableName() string { return "tipo_cliente" }

/* ---------- dir_cliente (direcciones) ---------- */

type DirCliente struct {
	ID         uint   `gorm:"primaryKey" json:"id"`
	RutCliente string `gorm:"column:rut_cliente;not null" json:"rut_cliente"`
	Direccion  string `json:"direccion"`
	Comuna     string `json:"comuna"`
	Ciudad     string `json:"ciudad"`

	// Relación inversa por RUT
	Cliente *Cliente `gorm:"foreignKey:RutCliente;references:Rut" json:"cliente,omitempty"`
}

func (DirCliente) TableName() string { return "dir_cliente" }
