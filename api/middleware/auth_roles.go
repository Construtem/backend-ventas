package middleware

import (
	"backend-ventas/api/database"
	"backend-ventas/api/models"
	"backend-ventas/services"
	"context"
	"net/http"
	"slices"
	"strings"

	"github.com/gin-gonic/gin"
)

// Middleware para verificar si el usuario tiene un rol o roles especificos
func AuthRoles(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extraer token de header
		authHeader := c.GetHeader("Authorization")
		idToken := strings.TrimPrefix(authHeader, "Bearer ")

		// Verificar Token
		token, err := services.FirebaseAuth.VerifyIDToken(context.Background(), idToken)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token inválido"})
			c.Abort()
			return
		}

		// Obtener email
		email := token.Claims["email"].(string)

		// Busar usuario
		var usuario models.Usuario

		if err := database.DB.Where("email = ?", email).First(&usuario).Error; err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuario no autorizado"})
			c.Abort()
			return
		}

		// Verificar el rol del usuario en roles
		if !slices.Contains(roles, usuario.Rol.Nombre) {
			c.JSON(http.StatusForbidden, gin.H{"error": "Acceso denegado"})
			c.Abort()
			return
		}

		// Guardar información del usuario en el contexto
		c.Set("usuario", usuario)
		c.Next()
	}
}
