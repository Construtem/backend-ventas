package handlers

import (
	"backend-ventas/api/controllers"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ObtenerResumenVentas() gin.HandlerFunc {
	return func(c *gin.Context) {
		total, err := controllers.ContarVentasRealizadas()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "No se pudieron contar las ventas",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"total_ventas_realizadas": total,
		})
	}
}
