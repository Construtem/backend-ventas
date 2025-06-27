package routes

import (
	"fmt"
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
	r.DELETE("/clientes/:id", eliminarCliente)
}


		// FUNCION OBTENER DATOS //
func obtenerClientes(c *gin.Context) {
	var clientes []models.Cliente
	if err := db.DB.Find(&clientes).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener clientes"})
		fmt.Print("\n\t\t<<<< ERROR AL OBTERE CLIENTES >>>>\n")
		return
	}
	c.JSON(http.StatusOK, clientes)
}


		// FUNCION CREAR CLIENTE //
func crearCliente(c *gin.Context) {
	var cliente models.Cliente
	if err := c.BindJSON(&cliente); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "JSON inválido"})
		fmt.Print("\n\t\t<<<< JSON INVALIDO >>>>\n")
		return
	}

	// Validaciones básicas
	if cliente.Nombre == "" || cliente.Rut == "" || cliente.Correo == "" || cliente.Telefono == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Faltan campos obligatorios"})
		fmt.Print("\n\t\t<<<< FALTAN CAMPOS >>>>\n")
		return
	}
	// Valida rut
	if !middleware.VerificacionRut(cliente.Rut) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "RUT inválido"})
		fmt.Print("\n\t\t<<<< RUT INVALIDO >>>>\n")
		return
	}

	// Verificar correo y rut
	var existe models.Cliente

	if err := db.DB.Where("correo = ?", cliente.Correo).First(&existe).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "El correo ya está registrado"})
		fmt.Print("\n\t\t<<<< EL CORREO YA ESTA REGISTRADO >>>>\n")
		return
	}
	if err := db.DB.Where("rut = ?", cliente.Rut).First(&existe).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "El rut ya está registrado"})
		fmt.Print("\n\t\t<<<< EL CORREO YA ESTA REGISTRADO >>>>\n")
		return
	}

	// Crear cliente
	if err := db.DB.Create(&cliente).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al crear cliente"})
		fmt.Print("\n\t\t<<<< ERROR AL CREAR CLIENTE >>>>\n")
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": cliente.ID})
	fmt.Print("\n\t\t<<<< CLIENTE CREADO CON EXITO >>>>\n")
}



func actualizarCliente(c *gin.Context) {
	id := c.Param("id")
	var cliente models.Cliente

	// Buscar cliente por ID
	if err := db.DB.First(&cliente, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Cliente no encontrado"})
		fmt.Print("\n\t\t<<<< CLIENTE NO ENCONTRADO >>>>\n")
		return
	}

	// Actualizar con los datos enviados
	var datosActualizados models.Cliente
	if err := c.BindJSON(&datosActualizados); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "JSON inválido"})
		fmt.Print("\n\t\t<<<< JSON INVALIDO >>>>\n")
		return
	}

	// Validar RUT si fue enviado

	// Si el usuario me está enviando un RUT (o sea, quiere actualizarlo), lo valido.
	// Si no me manda nada o lo deja vacío, lo ignoro
	if datosActualizados.Rut != "" && !middleware.VerificacionRut(datosActualizados.Rut) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "RUT inválido"})
		fmt.Print("\n\t\t<<<< RUT INVALIDO >>>>\n")
		return
	}

	// Actualiza todos los campos (parcial o completo)
	if err := db.DB.Model(&cliente).Updates(datosActualizados).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al actualizar cliente"})
		fmt.Print("\n\t\t<<<< ERROR AL ACTUALIZAR CLIENTE >>>>\n")
		return
	}

	c.JSON(http.StatusOK, cliente)
}


// ELIMINAR CLIENTE
func eliminarCliente(c *gin.Context) {
    id := c.Param("id")
    var cliente models.Cliente

    // Buscar el cliente por ID
    if err := db.DB.First(&cliente, id).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Cliente no encontrado"})
        fmt.Print("\n\t\t<<<< CLIENTE NO ENCONTRADO >>>>\n")
        return
    }

    // Eliminar el cliente
    if err := db.DB.Delete(&cliente).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al eliminar cliente"})
        fmt.Print("\n\t\t<<<< ERROR AL ELIMINAR CLIENTE >>>>\n")
        return
    }

    c.JSON(http.StatusOK, gin.H{"mensaje": "Cliente eliminado correctamente"})
    fmt.Printf("\n\t\t<<<< CLIENTE ID %s ELIMINADO >>>>\n", id)
}
