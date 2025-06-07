package auth

import (
	"backend-ventas/api/database"
	"backend-ventas/api/models"
	"backend-ventas/api/utils"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"gorm.io/gorm"
)

// AutorizarRol verifica el rol del usuario. ej: "admin", "vendedor", etc.
func AutorizarRol(rolRequerido string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. Extraer y limpiar el token del encabezado (header)
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token de autorización no proporcionado."})
			c.Abort()
			return
		}

		// Verificar si la cadena del token comienza con "Bearer " y quitarlo
		if !strings.HasPrefix(tokenString, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Formato de token inválido. Debe ser 'Bearer <token>'."})
			c.Abort()
			return
		}
		tokenString = strings.TrimPrefix(tokenString, "Bearer ")

		// 2. Validar token y extraer claims
		token, err := utils.ParseJWT(tokenString)
		if err != nil {
			// Proporcionar mensajes de error más específicos para errores comunes de JWT
			if ve, ok := err.(*jwt.ValidationError); ok {
				if ve.Errors&jwt.ValidationErrorMalformed != 0 {
					c.JSON(http.StatusUnauthorized, gin.H{"error": "Ese no es un token válido."})
				} else if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
					c.JSON(http.StatusUnauthorized, gin.H{"error": "Tu token ha expirado o aún no es válido. Por favor, inicia sesión de nuevo."})
				} else {
					log.Printf("Error al procesar el token: %v", err) // Log del error
					c.JSON(http.StatusUnauthorized, gin.H{"error": "No se pudo procesar tu token."})
				}
			} else {
				log.Printf("Fallo al parsear el token: %v", err) // Log del error
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Fallo al parsear el token."})
			}
			c.Abort()
			return
		}
		// Asegurarse de que el token sea válido y contenga claims
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid { // Asegurarse de que el token sea válido en general
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Fallo al leer los claims del token o el token es inválido."})
			c.Abort()
			return
		}

		// 3. Obtener user_id de los claims y verificar el tipo
		// Los números en JWT suelen ser float64.
		userIDClaim, ok := claims["user_id"]
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "ID de usuario no encontrado en los claims del token."})
			c.Abort()
			return
		}
		// Verificar que el tipo de userIDClaim sea float64 (común en JWT)
		userID, ok := userIDClaim.(float64)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "ID de usuario no válido en los claims del token."})
			c.Abort()
			return
		}

		// Convertir float64 a uint si el ID de tu modelo es uint
		parsedUserID := uint(userID)

		// 4. Obtener usuario y su rol de la DB. Considerar eager loading si es posible.
		var usuario models.Usuario
		// Sugerencia: Es altamente recomendable cargar el rol de forma "eager" para reducir las consultas a la DB.
		// Asumiendo que tiene un 'Rol' definida.
		if err := database.DB.Preload("Rol").First(&usuario, parsedUserID).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				c.JSON(http.StatusNotFound, gin.H{"error": "Usuario asociado con el token no encontrado."})
			} else {
				log.Printf("Error de base de datos al obtener el usuario: %v", err) // Log del error
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Error de base de datos al obtener el usuario."})
			}
			c.Abort()
			return
		}

		// Verificar si la relación Rol se cargó correctamente
		if usuario.Rol == nil || usuario.Rol.Nombre == "" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo determinar el rol del usuario."})
			c.Abort()
			return
		}

		// 5. Verificar si el rol coincide con el requerido
		if usuario.Rol.Nombre != rolRequerido {
			c.JSON(http.StatusForbidden, gin.H{"error": "No tienes los permisos necesarios para esta acción."})
			c.Abort()
			return
		}

		// 6. Guardar datos del usuario en el contexto
		c.Set("usuario", usuario)
		c.Next()
	}
}
