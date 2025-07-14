package routes

import (
	"fmt"
	"net/http"

	"backend-ventas/api/database"
	"backend-ventas/api/middleware"
	"backend-ventas/api/models"

	"github.com/gin-gonic/gin"
)

func ClienteRoutes(r *gin.Engine) {
	r.GET("/clientes", obtenerClientes)
	r.POST("/clientes", crearCliente)
	r.PATCH("/clientes/:rut", actualizarCliente)
	r.DELETE("/clientes/:rut", eliminarCliente)
}


		// FUNCION OBTENER DATOS //
func obtenerClientes(c *gin.Context) {
	var clientes []models.Cliente
	if err := database.DB.Find(&clientes).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener clientes"})
		fmt.Print("\n\t\t<<<< ERROR AL OBTERE CLIENTES >>>>\n")
		fmt.Printf("Error: %v", err)
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

	// Validaciones campos obligatorios
	if cliente.Nombre == "" || cliente.Rut == "" {
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

	// Verificar si el rut ya estan registrados
	var existe models.Cliente
	if err := database.DB.Where("rut = ?", cliente.Rut).First(&existe).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "El rut ya está registrado"})
		fmt.Print("\n\t\t<<<< EL CORREO YA ESTA REGISTRADO >>>>\n")
		return
	}
	// se pueden añadir mas, ej: Email
	// if err := database.DB.Where("correo = ?", cliente.Email).First(&existe).Error; err == nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": "El correo ya está registrado"})
	// 	fmt.Print("\n\t\t<<<< EL CORREO YA ESTA REGISTRADO >>>>\n")
	// 	return
	// }

	// Crear cliente
	if err := database.DB.Create(&cliente).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al crear cliente"})
		fmt.Print("\n\t\t<<<< ERROR AL CREAR CLIENTE >>>>\n")
		return
	}

	c.JSON(http.StatusCreated, gin.H{"cliente creado con rut": cliente.Rut})
	fmt.Print("\n\t\t<<<< CLIENTE CREADO CON EXITO >>>>\n")
}



func actualizarCliente(c *gin.Context) {
	rut_cliente := c.Param("rut")
	var cliente models.Cliente

	// Buscar cliente por RUT
	if err := database.DB.Where("rut = ?", rut_cliente).First(&cliente).Error; err != nil {
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
	if err := database.DB.Model(&cliente).Updates(datosActualizados).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al actualizar cliente"})
		fmt.Print("\n\t\t<<<< ERROR AL ACTUALIZAR CLIENTE >>>>\n")
		return
	}

	c.JSON(http.StatusOK, cliente)
}


// ELIMINAR CLIENTE
func eliminarCliente(c *gin.Context) {
    rut_cliente := c.Param("rut")
    var cliente models.Cliente

    // Buscar el cliente por RUT
    if err := database.DB.Where("rut = ?", rut_cliente).First(&cliente).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Cliente no encontrado"})
		fmt.Print("\n\t\t<<<< CLIENTE NO ENCONTRADO >>>>\n")
		return
	}

    // Eliminar el cliente
    if err := database.DB.Delete(&cliente).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al eliminar cliente"})
        fmt.Print("\n\t\t<<<< ERROR AL ELIMINAR CLIENTE >>>>\n")
        return
    }

    c.JSON(http.StatusOK, gin.H{"mensaje": "Cliente eliminado correctamente"})
    fmt.Printf("\n\t\t<<<< CLIENTE CON RUT %s ELIMINADO >>>>\n", rut_cliente)
}