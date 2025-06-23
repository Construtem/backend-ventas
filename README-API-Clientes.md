# API Clientes - Documentación

## VEN-RNF-BE-01: Documentación API Clientes

### Base URL
```
http://localhost:7777
```

### Endpoints Disponibles

## 1. Obtener Todos los Clientes

### GET /clientes

Obtiene la lista completa de todos los clientes registrados en el sistema.

#### Request
```http
GET /clientes
```

#### Headers
```
Content-Type: application/json
```

#### Response

**Status: 200 OK**
```json
[
  {
    "id": 1,
    "nombre": "Juan Pérez",
    "rut": "12345678-9",
    "correo": "juan.perez@email.com",
    "telefono": "+56912345678",
    "region_comuna": "Región Metropolitana, Santiago"
  },
  {
    "id": 2,
    "nombre": "María González",
    "rut": "98765432-1",
    "correo": "maria.gonzalez@email.com",
    "telefono": "+56987654321",
    "region_comuna": "Valparaíso, Viña del Mar"
  }
]
```

**Status: 500 Internal Server Error**
```json
{
  "error": "Error al obtener clientes"
}
```

#### Posibles Errores
- **500 Internal Server Error**: Error interno del servidor o problema de conexión con la base de datos

---

## 2. Crear Nuevo Cliente

### POST /clientes

Crea un nuevo cliente en el sistema con validaciones completas.

#### Request
```http
POST /clientes
```

#### Headers
```
Content-Type: application/json
```

#### Body
```json
{
  "nombre": "Carlos Rodríguez",
  "rut": "11222333-4",
  "correo": "carlos.rodriguez@email.com",
  "telefono": "+56911223344",
  "region_comuna": "Región Metropolitana, Providencia"
}
```

#### Validaciones

1. **Campos Obligatorios**: `nombre`, `rut`, `correo`, `telefono` son requeridos
2. **Validación de RUT**: El RUT debe ser válido según el algoritmo chileno
3. **Correo Único**: El correo electrónico debe ser único en el sistema

#### Response

**Status: 201 Created**
```json
{
  "id": 3
}
```

**Status: 400 Bad Request - JSON Inválido**
```json
{
  "error": "JSON inválido"
}
```

**Status: 400 Bad Request - Campos Faltantes**
```json
{
  "error": "Faltan campos obligatorios"
}
```

**Status: 400 Bad Request - RUT Inválido**
```json
{
  "error": "RUT inválido"
}
```

**Status: 400 Bad Request - Correo Duplicado**
```json
{
  "error": "El correo ya está registrado"
}
```

**Status: 500 Internal Server Error**
```json
{
  "error": "Error al crear cliente"
}
```

#### Posibles Errores
- **400 Bad Request**: 
  - JSON malformado
  - Campos obligatorios faltantes
  - RUT inválido
  - Correo electrónico ya registrado
- **500 Internal Server Error**: Error interno del servidor

---

## 3. Actualizar Cliente

### PATCH /clientes/{id}

Actualiza parcial o completamente los datos de un cliente existente.

#### Request
```http
PATCH /clientes/1
```

#### Headers
```
Content-Type: application/json
```

#### Body (Actualización Parcial)
```json
{
  "telefono": "+56999888777",
  "region_comuna": "Región Metropolitana, Las Condes"
}
```

#### Body (Actualización Completa)
```json
{
  "nombre": "Juan Pérez Actualizado",
  "rut": "12345678-9",
  "correo": "juan.perez.nuevo@email.com",
  "telefono": "+56999888777",
  "region_comuna": "Región Metropolitana, Las Condes"
}
```

#### Validaciones

1. **Cliente Existente**: El ID del cliente debe existir en el sistema
2. **Validación de RUT**: Si se envía un RUT, debe ser válido
3. **JSON Válido**: El cuerpo de la petición debe ser JSON válido

#### Response

**Status: 200 OK**
```json
{
  "id": 1,
  "nombre": "Juan Pérez Actualizado",
  "rut": "12345678-9",
  "correo": "juan.perez.nuevo@email.com",
  "telefono": "+56999888777",
  "region_comuna": "Región Metropolitana, Las Condes"
}
```

**Status: 400 Bad Request - JSON Inválido**
```json
{
  "error": "JSON inválido"
}
```

**Status: 400 Bad Request - RUT Inválido**
```json
{
  "error": "RUT inválido"
}
```

**Status: 404 Not Found**
```json
{
  "error": "Cliente no encontrado"
}
```

**Status: 500 Internal Server Error**
```json
{
  "error": "Error al actualizar cliente"
}
```

#### Posibles Errores
- **400 Bad Request**: 
  - JSON malformado
  - RUT inválido (si se proporciona)
- **404 Not Found**: Cliente con el ID especificado no existe
- **500 Internal Server Error**: Error interno del servidor

---

## Modelo de Datos

### Estructura del Cliente

```json
{
  "id": "uint (Primary Key)",
  "nombre": "string (Required)",
  "rut": "string (Required, Valid RUT)",
  "correo": "string (Required, Unique)",
  "telefono": "string (Required)",
  "region_comuna": "string (Optional)"
}
```

### Campos Obligatorios
- **nombre**: Nombre completo del cliente
- **rut**: RUT chileno válido (formato: 12345678-9)
- **correo**: Correo electrónico único
- **telefono**: Número de teléfono

### Campos Opcionales
- **region_comuna**: Región y comuna del cliente

---

## Validaciones Específicas

### Validación de RUT
El sistema implementa el algoritmo oficial chileno para validar RUTs:
- Elimina puntos y guiones
- Calcula el dígito verificador
- Valida que el dígito verificador sea correcto

**Ejemplos de RUTs Válidos:**
- `12345678-9`
- `98765432-1`
- `11222333-4`

**Ejemplos de RUTs Inválidos:**
- `12345678-0` (dígito verificador incorrecto)
- `1234567-9` (formato incorrecto)
- `abcdefg-h` (no numérico)

### Validación de Correo Único
- Cada correo electrónico debe ser único en el sistema
- No se permiten duplicados al crear nuevos clientes
- Se permite actualizar un cliente manteniendo su correo original

---

## Ejemplos de Uso

### Ejemplo 1: Crear un Cliente Nuevo
```bash
curl -X POST http://localhost:7777/clientes \
  -H "Content-Type: application/json" \
  -d '{
    "nombre": "Ana Silva",
    "rut": "15678901-2",
    "correo": "ana.silva@email.com",
    "telefono": "+56915678901",
    "region_comuna": "Región Metropolitana, Ñuñoa"
  }'
```

### Ejemplo 2: Obtener Todos los Clientes
```bash
curl -X GET http://localhost:7777/clientes
```

### Ejemplo 3: Actualizar Teléfono de un Cliente
```bash
curl -X PATCH http://localhost:7777/clientes/1 \
  -H "Content-Type: application/json" \
  -d '{
    "telefono": "+56912345678"
  }'
```

---

## Códigos de Estado HTTP

| Código | Descripción | Uso |
|--------|-------------|-----|
| 200 | OK | Operación exitosa (GET, PATCH) |
| 201 | Created | Cliente creado exitosamente |
| 400 | Bad Request | Datos inválidos o faltantes |
| 404 | Not Found | Cliente no encontrado |
| 500 | Internal Server Error | Error interno del servidor |

---

## Notas de Implementación

- **Base de Datos**: PostgreSQL con GORM como ORM
- **Framework**: Gin para el servidor HTTP
- **Validaciones**: Implementadas en el middleware
- **Respuestas**: Formato JSON consistente
- **Puerto**: 7777 (configurable en main.go)

---

## Testing

Puedes usar el archivo `test.http` incluido en el proyecto para probar los endpoints, o usar herramientas como:
- Postman
- Insomnia
- curl (línea de comandos)
- Thunder Client (VS Code extension) 