# API de Cotizaciones - Guía de Uso

---

## Endpoints principales

### Crear una cotización
**POST** `/api/cotizaciones`

**Body de ejemplo:**
```json
{
  "fecha": "2024-06-19T14:30:00Z",
  "cliente_id": 3,
  "vendedor_id": 5,
  "ubicacion_id": 7,
  "estado": "Pendiente",
  "detalle_cotizacion": [
    {
      "producto_id": 8,
      "cantidad": 1,
      "precio_unitario": 459990.00
    },
    {
      "producto_id": 9,
      "cantidad": 1,
      "precio_unitario": 658000.00
    }
  ]
}
```

**Respuesta exitosa:**
```json
{
  "message": "Cotización creada con éxito",
  "cotizacion": {
    "id": 2,
    "fecha": "2024-06-19T14:30:00-04:00",
    "cliente": { ... },
    "vendedor": { ... },
    "ubicacion": { ... },
    "estado": "Pendiente",
    "detalle_cotizacion": [ ... ]
  }
}
```

---

### Listar cotizaciones con filtros
**GET** `/api/cotizaciones`

Puedes usar los siguientes filtros como parámetros de query:

- `cliente_id` (int)
- `vendedor_id` (int)
- `ubicacion_id` (int)
- `estado` (string)
- `fecha_inicio` (YYYY-MM-DD)
- `fecha_fin` (YYYY-MM-DD)
- `cliente_nombre` (string, búsqueda parcial)
- `cliente_email` (string, búsqueda parcial)
- `vendedor_nombre` (string, búsqueda parcial)
- `ubicacion_nombre` (string, búsqueda parcial)
- `limit` (int, paginación)
- `page` (int, paginación)

#### Ejemplos de uso de filtros:

- Cotizaciones de un cliente:
  ```
  GET /api/cotizaciones?cliente_id=3
  ```
- Cotizaciones por vendedor:
  ```
  GET /api/cotizaciones?vendedor_id=5
  ```
- Cotizaciones pendientes:
  ```
  GET /api/cotizaciones?estado=Pendiente
  ```
- Cotizaciones entre fechas:
  ```
  GET /api/cotizaciones?fecha_inicio=2024-06-01&fecha_fin=2024-06-30
  ```
- Combinación de filtros:
  ```
  GET /api/cotizaciones?cliente_id=3&estado=Pendiente&fecha_inicio=2024-06-01&fecha_fin=2024-06-30
  ```
- Paginación:
  ```
  GET /api/cotizaciones?limit=5&page=2
  ```
- Búsqueda por nombre de cliente:
  ```
  GET /api/cotizaciones?cliente_nombre=José
  ```

**Respuesta:**
```json
{
  "data": [ ...cotizaciones... ],
  "total": 4,
  "page": 1,
  "limit": 10
}
```

---

### Obtener una cotización por ID
**GET** `/api/cotizaciones/{id}`

**Ejemplo:**
```
GET /api/cotizaciones/2
```

---

### Actualizar estado de una cotización
**PATCH** `/api/cotizaciones/{id}/estado`

**Body:**
```json
{
  "estado": "Aprobada",
  "usuario_id": 1
}
```

---

### Actualizar detalles de una cotización
**PATCH** `/api/cotizaciones/{id}/detalles`

**Body:**
```json
[
  {
    "producto_id": 8,
    "cantidad": 2,
    "precio_unitario": 450000
  }
]
```

---

## Notas importantes
- No envíes el campo `id` al crear una cotización.
- Los detalles no deben tener productos repetidos.
- Los IDs de cliente, vendedor, ubicación y producto deben existir en la base de datos.
- El endpoint responde con mensajes claros en caso de error.

---

## Recomendaciones
- Usa herramientas como **Postman** o **Insomnia** para probar la API.
- Si usas Docker, asegúrate de que la base de datos y el backend estén corriendo.

---