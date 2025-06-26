# 📊 Detalles Tablas par PostgreSQL

### Tablas Principales

#### 👥 **usuarios**
Almacena la información de todos los usuarios del sistema.

| Campo                          | Tipo         | Descripción                                   |
|--------------------------------|--------------|-----------------------------------------------|
| `usuario_id`                   | SERIAL       | Clave primaria auto-incremental               |
| `nombre_usuario`               | VARCHAR(100) | Nombre completo del usuario                   |
| `rut_usuario`                  | VARCHAR(20)  | RUT único del usuario (formato: 12345678-9)   |
| `email_usuario`                | VARCHAR(100) | Email único del usuario                       |
| `contrasena_usuario`           | VARCHAR(255) | Contraseña hasheada                           |
| `rol_id`                       | INTEGER      | Referencia al rol del usuario                 |
| `sucursal_id`                  | INTEGER      | Sucursal asignada (opcional)                  |
| `fecha_creacion_usuario`       | TIMESTAMP    | Fecha de creación del usuario                 |
| `activo_usuario`               | BOOLEAN      | Estado activo/inactivo                        |
| `fecha_ultimo_acceso_usuario`  | TIMESTAMP    | Último acceso del usuario                     |

#### 🏷️ **roles**
Define los roles y permisos de los usuarios.

| Campo              | Tipo         | Descripción           |
|--------------------|--------------|-----------------------|
| `rol_id`           | SERIAL       | Clave primaria        |
| `nombre_rol`       | VARCHAR(50)  | Nombre del rol        |
| `descripcion_rol`  | TEXT         | Descripción del rol   |

#### 📦 **productos**
Catálogo de productos disponibles.

| Campo                          | Tipo           | Descripción                                   |
|--------------------------------|----------------|-----------------------------------------------|
| `sku`                          | VARCHAR(50)    | Código único del producto                     |
| `nombre_producto`              | VARCHAR(200)   | Nombre del producto                           |
| `descripcion_producto`         | TEXT           | Descripción detallada                         |
| `precio_producto`              | DECIMAL(10,2)  | Precio unitario                               |
| `stock_producto`               | INTEGER        | Stock disponible                              |
| `largo_producto`               | DECIMAL(8,2)   | Largo del producto                            |
| `ancho_producto`               | DECIMAL(8,2)   | Ancho del producto                            |
| `alto_producto`                | DECIMAL(8,2)   | Alto del producto                             |
| `peso_producto`                | DECIMAL(8,2)   | Peso del producto                             |
| `unidad_medida`                | VARCHAR(20)    | Unidad de medida                              |
| `categoria_id`                 | INTEGER        | Categoría del producto                        |
| `proveedor_id`                 | INTEGER        | Proveedor del producto                        |
| `sucursal_id`                  | INTEGER        | Sucursal donde se encuentra                   |
| `activo_producto`              | BOOLEAN        | Estado activo/inactivo                        |
| `fecha_creacion_producto`      | TIMESTAMP      | Fecha de creación                             |
| `imagen_producto`              | VARCHAR(500)   | URL de la imagen                              |

### Tablas de Transacciones

#### 🛒 **pedidos**
Registra todos los pedidos realizados.

| Campo                          | Tipo           | Descripción                                   |
|--------------------------------|----------------|-----------------------------------------------|
| `pedido_id`                    | SERIAL         | Clave primaria                                |
| `cotizacion_id`                | INTEGER        | Referencia a cotización (opcional)            |
| `usuario_id`                   | INTEGER        | Cliente que realiza el pedido                 |
| `fecha_creacion_pedido`        | TIMESTAMP      | Fecha de creación                             |
| `fecha_entrega_pedido`         | DATE           | Fecha de entrega estimada                     |
| `direccion_entrega_pedido`     | TEXT           | Dirección de entrega                          |
| `estado_pedido`                | VARCHAR(50)    | Estado: PENDIENTE, APROBADO, EN_PROCESO, ENVIADO, ENTREGADO, CANCELADO |
| `metodo_pago`                  | VARCHAR(50)    | Método de pago                                |
| `total_pedido`                 | DECIMAL(12,2)  | Total del pedido                              |
| `observaciones_pedido`         | TEXT           | Observaciones adicionales                     |

#### 📋 **detalle_pedido**
Detalle de productos en cada pedido.

| Campo                          | Tipo           | Descripción                                   |
|--------------------------------|----------------|-----------------------------------------------|
| `pedido_id`                    | INTEGER        | Referencia al pedido                          |
| `producto_id`                  | VARCHAR(50)    | Referencia al producto                        |
| `cantidad_producto`            | INTEGER        | Cantidad solicitada                           |
| `precio_unitario_producto`     | DECIMAL(10,2)  | Precio unitario al momento del pedido         |
| `subtotal_producto`            | DECIMAL(12,2)  | Subtotal calculado                            |

#### 💰 **cotizaciones**
Gestiona las cotizaciones previas a los pedidos.

| Campo                          | Tipo           | Descripción                                   |
|--------------------------------|----------------|-----------------------------------------------|
| `cotizacion_id`                | SERIAL         | Clave primaria                                |
| `estado_cotizacion`            | VARCHAR(50)    | Estado: PENDIENTE, APROBADA, RECHAZADA, VENCIDA |
| `fecha_creacion_cotizacion`    | TIMESTAMP      | Fecha de creación                             |
| `fecha_vencimiento_cotizacion` | DATE           | Fecha de vencimiento                          |
| `fecha_aprobacion_cotizacion`  | TIMESTAMP      | Fecha de aprobación                           |
| `total_cotizacion`             | DECIMAL(12,2)  | Total de la cotización                        |
| `cliente_id`                   | INTEGER        | Cliente de la cotización                      |
| `vendedor_id`                  | INTEGER        | Vendedor asignado                             |
| `aprobador_id`                 | INTEGER        | Usuario que aprueba                           |
| `observaciones_cotizacion`     | TEXT           | Observaciones                                 |

### Tablas de Logística

#### 🚚 **camiones**
Flota de camiones disponibles.

| Campo                          | Tipo           | Descripción                                   |
|--------------------------------|----------------|-----------------------------------------------|
| `camion_id`                    | SERIAL         | Clave primaria                                |
| `patente_camion`               | VARCHAR(10)    | Patente del camión                            |
| `tipo_camion_id`               | INTEGER        | Tipo de camión                                |
| `sucursal_id`                  | INTEGER        | Sucursal asignada                             |
| `responsable_camion_id`        | INTEGER        | Responsable del camión                        |
| `descripcion_camion`           | TEXT           | Descripción adicional                         |
| `activo_camion`                | BOOLEAN        | Estado activo/inactivo                        |

#### 🚛 **tipo_camion**
Tipos de camiones disponibles.

| Campo                          | Tipo           | Descripción                                   |
|--------------------------------|----------------|-----------------------------------------------|
| `tipo_camion_id`               | SERIAL         | Clave primaria                                |
| `capacidad_camion`             | DECIMAL(8,2)   | Capacidad en toneladas                        |
| `nombre_tipo_camion`           | VARCHAR(100)   | Nombre del tipo                               |
| `descripcion_tipo_camion`      | TEXT           | Descripción                                   |
| `activo_tipo_camion`           | BOOLEAN        | Estado activo/inactivo                        |

#### 📦 **despachos**
Gestiona los despachos de pedidos.

| Campo                          | Tipo           | Descripción                                   |
|--------------------------------|----------------|-----------------------------------------------|
| `despacho_id`                  | SERIAL         | Clave primaria                                |
| `pedido_id`                    | INTEGER        | Pedido asociado                               |
| `camion_id`                    | INTEGER        | Camión asignado                               |
| `origen_despacho`              | VARCHAR(200)   | Origen del despacho                           |
| `destino_despacho`             | VARCHAR(200)   | Destino del despacho                          |
| `estado_despacho`              | VARCHAR(50)    | Estado: PENDIENTE, EN_CAMINO, ENTREGADO, CANCELADO |
| `fecha_envio_despacho`         | TIMESTAMP      | Fecha de envío                                |
| `fecha_entrega_despacho`       | TIMESTAMP      | Fecha de entrega                              |
| `costo_despacho`               | DECIMAL(10,2)  | Costo del despacho                            |
| `observaciones_despacho`       | TEXT           | Observaciones                                 |

#### 🚛 **envio_camion**
Detalle de envíos por camión.

| Campo                          | Tipo           | Descripción                                   |
|--------------------------------|----------------|-----------------------------------------------|
| `envio_camion_id`              | SERIAL         | Clave primaria                                |
| `despacho_id`                  | INTEGER        | Despacho asociado                             |
| `camion_id`                    | INTEGER        | Camión asignado                               |
| `estado_envio_camion`          | VARCHAR(50)    | Estado: PENDIENTE, EN_TRANSITO, ENTREGADO, CANCELADO |
| `fecha_envio_camion`           | TIMESTAMP      | Fecha de envío                                |
| `fecha_entrega_camion`         | TIMESTAMP      | Fecha de entrega                              |
| `observaciones_envio_camion`   | TEXT           | Observaciones                                 |

#### 📦 **detalle_envio**
Productos incluidos en cada envío.

| Campo                          | Tipo           | Descripción                                   |
|--------------------------------|----------------|-----------------------------------------------|
| `envio_camion_id`              | INTEGER        | Envío asociado                                |
| `producto_id`                  | VARCHAR(50)    | Producto incluido                             |
| `cantidad_producto`            | INTEGER        | Cantidad enviada                              |

### Tablas de Inventario

#### 📊 **stock_sucursal**
Control de stock por sucursal.

| Campo                          | Tipo           | Descripción                                   |
|--------------------------------|----------------|-----------------------------------------------|
| `stock_sucursal_id`            | SERIAL         | Clave primaria                                |
| `producto_id`                  | VARCHAR(50)    | Producto                                      |
| `sucursal_id`                  | INTEGER        | Sucursal                                      |
| `stock_producto`               | INTEGER        | Stock disponible                              |
| `stock_minimo_producto`        | INTEGER        | Stock mínimo para alertas                     |

## 🔧 Características Técnicas

### **Validaciones Implementadas**
- ✅ Validación de email con regex
- ✅ Validación de RUT chileno
- ✅ Validación de patente de camión
- ✅ Restricciones de valores positivos
- ✅ Estados predefinidos
- ✅ Validación de fechas lógicas

### **Índices Optimizados**
- Búsqueda por email y RUT de usuarios
- Búsqueda por SKU y nombre de productos
- Filtrado por estados de pedidos y cotizaciones
- Consultas por fechas para reportes

### **Integridad Referencial**
- Eliminación en cascada para detalles
- Restricción de eliminación para entidades principales
- Valores NULL para relaciones opcionales

## 📋 Estados del Sistema

### Estados de Pedidos
- `PENDIENTE`  - Pedido creado, pendiente de procesamiento
- `APROBADO`   - Pedido aprobado para procesamiento
- `EN_PROCESO` - Pedido en proceso de preparación
- `ENVIADO`    - Pedido enviado
- `ENTREGADO`  - Pedido entregado al cliente
- `CANCELADO`  - Pedido cancelado

### Estados de Cotizaciones
- `PENDIENTE`  - Cotización creada, pendiente de aprobación
- `APROBADA`   - Cotización aprobada
- `RECHAZADA`  - Cotización rechazada
- `VENCIDA`    - Cotización vencida

### Estados de Despachos
- `PENDIENTE` - Despacho creado, pendiente de asignación
- `EN_CAMINO` - Despacho en ruta
- `ENTREGADO` - Despacho entregado
- `CANCELADO` - Despacho cancelado