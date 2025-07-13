# Middleware de Autenticación de Roles

## Descripción
Middleware para verificar roles de usuario usando Firebase Auth + Base de Datos.

## Uso

### Importar
```go
import "backend-ventas/api/middleware"
```

### Aplicar a Rutas
```go
// Solo gerentes
router.PATCH("/cotizaciones/:id/estado", middleware.AuthRoles("gerente"), handler)

// Múltiples roles
router.POST("/cotizaciones", middleware.AuthRoles("admin", "vendedor"), handler)

// Cualquier rol autenticado
router.GET("/perfil", middleware.AuthRoles("admin", "vendedor", "inventario"), handler)
```

## Flujo de Autenticación
1. **Extrae token** del header `Authorization: Bearer <token>`
2. **Verifica token** con Firebase Auth
3. **Extrae email** del token verificado
4. **Busca usuario** en DB por email
5. **Verifica rol** del usuario contra roles permitidos
6. **Permite/Deniega** acceso según el rol

## Respuestas de Error
- `401 Unauthorized`: Token inválido o usuario no encontrado
- `403 Forbidden`: Rol no autorizado para la ruta

## Acceso al Usuario en Handlers
```go
func MiHandler(c *gin.Context) {
    usuario := c.MustGet("usuario").(models.Usuario)
    fmt.Printf("Usuario: %s, Rol: %s", usuario.Nombre, usuario.Rol.Nombre)
}
```
