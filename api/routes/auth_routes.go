package routes

import (
	"backend-ventas/api/database"
	"backend-ventas/api/models"
	"backend-ventas/api/utils" // Para generar el JWT
	"log"
	"net/http"
	"strings" // Para verificar errores de unique constraint

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt" // Para comparar y HASHEAR contraseñas
	"gorm.io/gorm"
)

// LoginRequest define la estructura del cuerpo de la solicitud de login.
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// RegisterRequest define la estructura del cuerpo de la solicitud de registro.
// Aquí se incluyen todos los campos necesarios para crear un nuevo usuario.
type RegisterRequest struct {
	Nombre      string `json:"nombre" binding:"required"`
	Email       string `json:"email" binding:"required,email"`
	Contrasena  string `json:"contrasena" binding:"required,min=6"` // Mínimo 6 caracteres para la contraseña
	RolID       uint   `json:"rol_id" binding:"required"`           // El ID del rol al que se asignará el usuario
	UbicacionID *uint  `json:"ubicacion_id"`                        // Puede ser nulo si no todas las cuentas necesitan una ubicación
}

// LoginHandler (tu función existente, no necesita cambios)
func LoginHandler(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var usuario models.Usuario
	if err := database.DB.Preload("Rol").Where("email = ?", req.Email).First(&usuario).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Credenciales inválidas: usuario no encontrado."})
			return
		} else {
			log.Printf("Error de DB al buscar usuario: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error interno del servidor."})
			return
		}
	}

	if err := bcrypt.CompareHashAndPassword([]byte(usuario.Contrasena), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Credenciales inválidas: contraseña incorrecta."})
		return
	}

	if usuario.Rol == nil || usuario.Rol.Nombre == "" {
		log.Printf("Error: El rol del usuario %s (ID: %d) no pudo ser cargado o no tiene nombre.", usuario.Email, usuario.ID)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo determinar el rol del usuario."})
		return
	}

	tokenString, err := utils.GenerateJWT(usuario.ID, usuario.Rol.Nombre)
	if err != nil {
		log.Printf("Error al generar JWT para usuario %s: %v", usuario.Email, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al generar token de autenticación."})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Login exitoso!", "token": tokenString})
}

// RegisterHandler maneja el registro de nuevos usuarios.
func RegisterHandler(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Verificar si el RolID proporcionado existe
	var rol models.Rol
	if err := database.DB.First(&rol, req.RolID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusBadRequest, gin.H{"error": "El RolID proporcionado no existe."})
			return
		}
		log.Printf("Error al verificar RolID en registro: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error interno al verificar el rol."})
		return
	}

	// Verificar si la UbicacionID proporcionada existe (si se incluyó)
	if req.UbicacionID != nil {
		var ubicacion models.Ubicacion
		if err := database.DB.First(&ubicacion, *req.UbicacionID).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				c.JSON(http.StatusBadRequest, gin.H{"error": "La UbicacionID proporcionada no existe."})
				return
			}
			log.Printf("Error al verificar UbicacionID en registro: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error interno al verificar la ubicación."})
			return
		}
	}

	// Hashear la contraseña antes de guardarla
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Contrasena), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Error al hashear contraseña en registro: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al procesar la contraseña."})
		return
	}

	// Crear el objeto Usuario
	newUser := models.Usuario{
		Nombre:      req.Nombre,
		Email:       req.Email,
		Contrasena:  string(hashedPassword),
		RolID:       req.RolID,
		UbicacionID: req.UbicacionID,
	}

	// Guardar el nuevo usuario en la base de datos
	if err := database.DB.Create(&newUser).Error; err != nil {
		// Manejar error si el email ya existe (violación de unique constraint)
		if strings.Contains(err.Error(), "duplicate key value violates unique constraint") || strings.Contains(err.Error(), "unique constraint failed") {
			c.JSON(http.StatusConflict, gin.H{"error": "El email ya está registrado."})
			return
		}
		log.Printf("Error al crear nuevo usuario: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al registrar usuario."})
		return
	}

	// Generar JWT para el usuario recién registrado (lo loguea automáticamente)
	// Aseguramos que el rol del usuario se cargue para el JWT
	if err := database.DB.Preload("Rol").First(&newUser, newUser.ID).Error; err != nil {
		log.Printf("Error al recargar usuario con rol después de registro: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Usuario registrado, pero error al generar token."})
		return
	}

	if newUser.Rol == nil || newUser.Rol.Nombre == "" {
		log.Printf("Error: Rol del usuario recién registrado no disponible para token. Email: %s", newUser.Email)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Usuario registrado, pero rol no encontrado para token."})
		return
	}

	tokenString, err := utils.GenerateJWT(newUser.ID, newUser.Rol.Nombre)
	if err != nil {
		log.Printf("Error al generar JWT para usuario %s después de registro: %v", newUser.Email, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Usuario registrado, pero error al generar token."})
		return
	}

	// Cargar las relaciones para que la respuesta incluya los nombres de Rol y Ubicacion
	database.DB.Preload("Rol").Preload("Ubicacion").First(&newUser, newUser.ID)

	c.JSON(http.StatusCreated, gin.H{
		"message": "Registro exitoso!",
		"user":    newUser,
		"token":   tokenString,
	})
}

// SetupAuthRoutes registra las rutas relacionadas con la autenticación.
func SetupAuthRoutes(router *gin.Engine) {
	// Grupo de rutas para autenticación, por ejemplo, /auth/login, /auth/registrar.
	authRoutes := router.Group("/auth")
	{
		authRoutes.POST("/login", LoginHandler)
		authRoutes.POST("/registrar", RegisterHandler)
		// Puedes añadir otras rutas de autenticación aquí, como reseteo de contraseña, etc.
	}
}
