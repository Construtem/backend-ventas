package handlers

import (
	"net/http"
	"strconv"

	"backend-ventas/api/controllers"
	"github.com/gin-gonic/gin"
)

// Ejemplo: GET /api/productos/inventario?sucursal_id=2&page=1&limit=25
func ObtenerProductosInventario() gin.HandlerFunc {
	return func(c *gin.Context) {
		/* -------------------- 1) Validar parámetros -------------------- */
		sucursalIDStr := c.Query("sucursal_id")
		if sucursalIDStr == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Parámetro sucursal_id requerido",
			})
			return
		}
		sucursalID, err := strconv.Atoi(sucursalIDStr)
		if err != nil || sucursalID <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "sucursal_id inválido",
			})
			return
		}

		page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
		limit, _ := strconv.Atoi(c.DefaultQuery("limit", "25"))
		if page <= 0 {
			page = 1
		}
		if limit <= 0 || limit > 100 {
			limit = 25
		}

		/* -------------------- 2) Obtener datos -------------------- */
		resp, err := controllers.ListarProductosConInventario(sucursalID, page, limit)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		/* -------------------- 3) Responder -------------------- */
		c.JSON(http.StatusOK, resp)
	}
}
