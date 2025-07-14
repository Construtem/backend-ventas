package mappers

import (
	"backend-ventas/api/dtos"
	"backend-ventas/api/models"
)

// ---------- Cliente ----------
func ClienteToDTO(m *models.Cliente) dtos.Cliente {
	if m == nil {
		return dtos.Cliente{}
	}
	return dtos.Cliente{
		Rut:         m.Rut,
		Nombre:      m.Nombre,
		Telefono:    m.Telefono,
		Email:       m.Email,
		RazonSocial: m.RazonSocial,
	}
}

// ---------- Usuario ----------
func UsuarioToDTO(m *models.Usuario) dtos.Usuario {
	if m == nil {
		return dtos.Usuario{}
	}
	return dtos.Usuario{
		Email:  m.Email,
		Nombre: m.Nombre,
		RolID:  m.RolID,
	}
}

// ---------- Dirección ----------
func DirClienteToDTO(m *models.DirCliente) dtos.DireccionCliente {
	if m == nil {
		return dtos.DireccionCliente{}
	}
	return dtos.DireccionCliente{
		Direccion: m.Direccion,
		Comuna:    m.Comuna,
		Ciudad:    m.Ciudad,
	}
}
