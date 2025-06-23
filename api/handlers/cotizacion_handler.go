package handlers

import (
	"backend-ventas/api/controllers"
	"backend-ventas/api/dtos" // Importa tus DTOs
	"backend-ventas/api/models"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UpdateEstadoRequest struct {
	Estado    string `json:"estado" binding:"required"`
	UsuarioID *uint  `json:"usuario_id"`
}

func GetCotizacionesHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		queryParams := make(map[string]string)
		for k, v := range c.Request.URL.Query() {
			queryParams[k] = v[0]
		}

		cotizaciones, total, err := controllers.GetCotizaciones(db, queryParams)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		limit := 10
		if limitStr, ok := queryParams["limit"]; ok {
			if l, err := strconv.Atoi(limitStr); err == nil {
				limit = l
			}
		}
		page := 1
		if pageStr, ok := queryParams["page"]; ok {
			if p, err := strconv.Atoi(pageStr); err == nil {
				page = p
			}
		}

		c.JSON(http.StatusOK, dtos.CotizacionesListResponse{
			Data:  cotizaciones,
			Total: total,
			Page:  page,
			Limit: limit,
		})
	}
}

func GetCotizacionByIDHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.ParseUint(idStr, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "ID de cotización inválido"})
			return
		}

		cotizacion, err := controllers.GetCotizacionByID(db, uint(id))
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				c.JSON(http.StatusNotFound, gin.H{"error": "Cotización no encontrada"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener la cotización: " + err.Error()})
			}
			return
		}
		c.JSON(http.StatusOK, cotizacion)
	}
}

func CreateCotizacionHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var nuevaCotizacion models.Cotizacion

		if err := c.ShouldBindJSON(&nuevaCotizacion); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Solicitud inválida", "details": err.Error()})
			return
		}

		if err := controllers.CreateCotizacion(db, &nuevaCotizacion); err != nil {
			if err.Error() == "los campos 'aprobadaPorID' y 'fechaAprobacion' son requeridos si el estado es 'Aprobada'" {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo crear la cotización", "details": err.Error()})
			}
			return
		}

		createdCotizacionDTO, err := controllers.GetCotizacionByID(db, nuevaCotizacion.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Cotización creada, pero error al recargar para la respuesta", "details": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"message": "Cotización creada con éxito", "cotizacion": createdCotizacionDTO})
	}
}

func UpdateEstadoCotizacionHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.ParseUint(idStr, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "ID de cotización inválido"})
			return
		}

		var req UpdateEstadoRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		cotizacionActualizadaDTO, err := controllers.UpdateEstadoCotizacion(db, uint(id), req.Estado, req.UsuarioID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				c.JSON(http.StatusNotFound, gin.H{"error": "Cotización no encontrada"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al actualizar el estado de la cotización: " + err.Error()})
			}
			return
		}
		c.JSON(http.StatusOK, cotizacionActualizadaDTO)
	}
}

func UpdateCotizacionDetalleHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		cotizacionIDStr := c.Param("id")
		cotizacionID, err := strconv.ParseUint(cotizacionIDStr, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "ID de cotización inválido", "details": "El ID debe ser un número entero sin signo."})
			return
		}

		var nuevoDetalle []models.DetalleCotizacion
		if err := c.ShouldBindJSON(&nuevoDetalle); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Datos de detalle de cotización inválidos", "details": err.Error()})
			return
		}

		if len(nuevoDetalle) == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "La lista de detalles de cotización no puede estar vacía."})
			return
		}

		updatedCotizacionDTO, err := controllers.UpdateCotizacionDetalle(db, uint(cotizacionID), nuevoDetalle)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudieron actualizar los detalles de la cotización", "details": err.Error()})
			}
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Detalles de cotización actualizados con éxito", "cotizacion": updatedCotizacionDTO})
	}
}
