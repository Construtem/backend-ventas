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
    "detalle_cotizacion": [
      {
        "producto_id": 8,
        "cantidad": 1,
        "precio_unitario": 459990,
        "producto": {
          "codigo": "753698569",
          "nombre": "Cocina a Gas...",
          "descripcion": "...",
          "precio_venta": 459990
        }
      }
    ],
    "total": 1612990
  }
}

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
  "data": [
    {
      "id": 2,
      "fecha": "2024-06-19T14:30:00-04:00",
      "cliente": {
        "id": 3,
        "nombre": "Juan Pérez",
        "telefono": "+56912345678",
        "email": "juan@email.com",
        "direccion": "Av. Principal 123"
      },
      "vendedor": {
        "id": 5,
        "nombre": "María González",
        "email": "maria@empresa.com"
      },
      "ubicacion": {
        "id": 7,
        "nombre": "Sucursal Centro",
        "tipo": "Tienda",
        "direccion": "Calle Central 456"
      },
      "estado": "Pendiente",
      "productos": [
        {
          "id": 8,
          "nombre": "Cocina a Gas 4 Quemadores",
          "cantidad": 1,
          "precio_unitario": 459990
        },
        {
          "id": 9,
          "nombre": "Refrigerador Side by Side",
          "cantidad": 1,
          "precio_unitario": 658000
        }
      ],
      "total": 1117990
    }
  ],
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

**Respuesta:**
```json
{
  "id": 2,
  "fecha": "2024-06-19T14:30:00-04:00",
  "cliente": {
    "id": 3,
    "nombre": "Juan Pérez",
    "telefono": "+56912345678",
    "email": "juan@email.com",
    "direccion": "Av. Principal 123"
  },
  "vendedor": {
    "id": 5,
    "nombre": "María González",
    "email": "maria@empresa.com"
  },
  "ubicacion": {
    "id": 7,
    "nombre": "Sucursal Centro",
    "tipo": "Tienda",
    "direccion": "Calle Central 456"
  },
  "estado": "Pendiente",
  "detalle_cotizacion": [
    {
      "producto_id": 8,
      "cantidad": 1,
      "precio_unitario": 459990,
      "producto": {
        "codigo": "753698569",
        "nombre": "Cocina a Gas 4 Quemadores",
        "descripcion": "Cocina a gas con 4 quemadores y horno",
        "precio_venta": 459990
      }
    },
    {
      "producto_id": 9,
      "cantidad": 1,
      "precio_unitario": 658000,
      "producto": {
        "codigo": "753698570",
        "nombre": "Refrigerador Side by Side",
        "descripcion": "Refrigerador de dos puertas",
        "precio_venta": 658000
      }
    }
  ],
  "total": 1117990
}
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

### Historial de una cotización
**GET** `/api/cotizaciones/{id}/historial`

Devuelve los eventos asociados a la cotización.

**Ejemplo:**
```
GET /api/cotizaciones/2/historial
```

**Respuesta:**
```json
[
  {
    "id": 1,
    "fecha": "2024-06-19T14:30:00-04:00",
    "accion": "Creada"
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