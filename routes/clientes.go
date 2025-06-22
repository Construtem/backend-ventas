package routes

import (
	"net/http"

	"backend-ventas/db"
	"backend-ventas/middleware"
	"backend-ventas/models"

	"github.com/gin-gonic/gin"
)

func ClienteRoutes(r *gin.Engine) {
	r.GET("/clientes", obtenerClientes)
	r.POST("/clientes", crearCliente)
	r.PATCH("/clientes/:id", actualizarCliente)
}

func obtenerClientes(c *gin.Context) {
	var clientes []models.Cliente
	if err := db.DB.Find(&clientes).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener clientes"})
		return
	}
	c.JSON(http.StatusOK, clientes)
}

func crearCliente(c *gin.Context) {
	var cliente models.Cliente
	if err := c.BindJSON(&cliente); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "JSON inválido"})
		return
	}

	// Validaciones básicas
	if cliente.Nombre == "" || cliente.Rut == "" || cliente.Correo == "" || cliente.Telefono == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Faltan campos obligatorios"})
		return
	}

	if !middleware.VerificacionRut(cliente.Rut) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "RUT inválido"})
		return
	}

	// Verificar correo único
	var existe models.Cliente
	if err := db.DB.Where("correo = ?", cliente.Correo).First(&existe).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "El correo ya está registrado"})
		return
	}

	// Crear cliente
	if err := db.DB.Create(&cliente).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al crear cliente"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": cliente.ID})
}

func actualizarCliente(c *gin.Context) {
	id := c.Param("id")
	var cliente models.Cliente

	// Buscar cliente por ID
	if err := db.DB.First(&cliente, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Cliente no encontrado"})
		return
	}

	// Actualizar con los datos enviados
	var datosActualizados models.Cliente
	if err := c.BindJSON(&datosActualizados); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "JSON inválido"})
		return
	}

	// Validar RUT si fue enviado
	if datosActualizados.Rut != "" && !middleware.VerificacionRut(datosActualizados.Rut) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "RUT inválido"})
		return
	}

	// Actualiza todos los campos (parcial o completo)
	if err := db.DB.Model(&cliente).Updates(datosActualizados).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al actualizar cliente"})
		return
	}

	c.JSON(http.StatusOK, cliente)
}
