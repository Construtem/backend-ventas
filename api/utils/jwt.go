package utils

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
)

// GenerateJWT genera un token JWT para un usuario dado su ID y rol.
func GenerateJWT(userID uint, userRole string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,                                // Incluye el ID del usuario en los claims del token
		"rol":     userRole,                              // Incluye el rol del usuario en los claims del token
		"exp":     time.Now().Add(time.Hour * 24).Unix(), // Establece la expiración del token a 24 horas
	})

	jwtSecret := os.Getenv("JWT_SECRET_KEY")
	if jwtSecret == "" {
		return "", fmt.Errorf("JWT_SECRET_KEY no está configurada en el entorno")
	}

	return token.SignedString([]byte(jwtSecret)) // Firma el token con la clave secreta obtenida de las variables de entorno
}

// ParseJWT parsea y valida un token JWT
func ParseJWT(tokenString string) (*jwt.Token, error) {
	// Obtener la clave secreta del entorno o de alguna configuración, la misma que se usó para FIRMAR el token.
	secretKey := os.Getenv("JWT_SECRET_KEY")
	if secretKey == "" {
		return nil, fmt.Errorf("JWT_SECRET_KEY no está configurada en el entorno")
	}

	// Usar jwt.Parse con la función Keyfunc
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		// Validar que el algoritmo de firma del token sea el esperado. ejemplo: si usas HMAC-SHA256 (HS256)
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("método de firma inesperado: %v", token.Header["alg"])
		}
		// Devolve la clave secreta en formato []byte
		return []byte(secretKey), nil
	})

	if err != nil {
		return nil, err
	}

	return token, nil
}
