-- Script Tablas paraa PostgreSQL

-- =====================================================
-- TABLA: roles
-- =====================================================
CREATE TABLE roles (
    rol_id SERIAL PRIMARY KEY,
    nombre_rol VARCHAR(50) NOT NULL UNIQUE,
    descripcion_rol TEXT
);

-- Índices para roles
CREATE INDEX idx_roles_nombre ON roles(nombre_rol);

-- =====================================================
-- TABLA: usuarios
-- =====================================================
CREATE TABLE usuarios (
    usuario_id SERIAL PRIMARY KEY,
    nombre_usuario VARCHAR(100) NOT NULL,
    rut_usuario VARCHAR(20) UNIQUE NOT NULL,
    email_usuario VARCHAR(100) UNIQUE NOT NULL,
    contrasena_usuario VARCHAR(255) NOT NULL,
    rol_id INTEGER NOT NULL,
    sucursal_id INTEGER, -- FK a tabla sucursales
    fecha_creacion_usuario TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    activo_usuario BOOLEAN DEFAULT TRUE,
    fecha_ultimo_acceso_usuario TIMESTAMP,
    
    -- Restricciones
    CONSTRAINT fk_usuarios_rol FOREIGN KEY (rol_id) REFERENCES roles(rol_id) ON DELETE RESTRICT,
    CONSTRAINT chk_email_usuario CHECK (email_usuario ~* '^[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Za-z]{2,}$'),
    CONSTRAINT chk_rut_usuario CHECK (rut_usuario ~ '^[0-9]{7,8}-[0-9kK]$')
);

-- Índices para usuarios
CREATE INDEX idx_usuarios_email ON usuarios(email_usuario);
CREATE INDEX idx_usuarios_rut ON usuarios(rut_usuario);
CREATE INDEX idx_usuarios_rol ON usuarios(rol_id);
CREATE INDEX idx_usuarios_sucursal ON usuarios(sucursal_id);
CREATE INDEX idx_usuarios_activo ON usuarios(activo_usuario);

-- =====================================================
-- TABLA: productos
-- =====================================================
CREATE TABLE productos (
    sku VARCHAR(50) PRIMARY KEY,
    nombre_producto VARCHAR(200) NOT NULL,
    descripcion_producto TEXT,
    precio_producto DECIMAL(10,2) NOT NULL CHECK (precio_producto >= 0),
    stock_producto INTEGER DEFAULT 0 CHECK (stock_producto >= 0),
    largo_producto DECIMAL(8,2) CHECK (largo_producto >= 0),
    ancho_producto DECIMAL(8,2) CHECK (ancho_producto >= 0),
    alto_producto DECIMAL(8,2) CHECK (alto_producto >= 0),
    peso_producto DECIMAL(8,2) CHECK (peso_producto >= 0),
    unidad_medida VARCHAR(20) NOT NULL,
    categoria_id INTEGER, -- FK a tabla categorias
    proveedor_id INTEGER, -- FK a tabla proveedores
    sucursal_id INTEGER, -- FK a tabla sucursales
    activo_producto BOOLEAN DEFAULT TRUE,
    fecha_creacion_producto TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    imagen_producto VARCHAR(500)
);

-- Índices para productos
CREATE INDEX idx_productos_nombre ON productos(nombre_producto);
CREATE INDEX idx_productos_categoria ON productos(categoria_id);
CREATE INDEX idx_productos_proveedor ON productos(proveedor_id);
CREATE INDEX idx_productos_sucursal ON productos(sucursal_id);
CREATE INDEX idx_productos_activo ON productos(activo_producto);
CREATE INDEX idx_productos_precio ON productos(precio_producto);

-- =====================================================
-- TABLA: pedidos
-- =====================================================
CREATE TABLE pedidos (
    pedido_id SERIAL PRIMARY KEY,
    cotizacion_id INTEGER, -- FK a tabla cotizaciones
    usuario_id INTEGER NOT NULL, -- FK a tabla usuarios
    fecha_creacion_pedido TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    fecha_entrega_pedido DATE,
    direccion_entrega_pedido TEXT,
    estado_pedido VARCHAR(50) NOT NULL DEFAULT 'PENDIENTE',
    metodo_pago VARCHAR(50) NOT NULL,
    total_pedido DECIMAL(12,2) NOT NULL CHECK (total_pedido >= 0),
    observaciones_pedido TEXT,
    
    -- Restricciones
    CONSTRAINT fk_pedidos_usuario FOREIGN KEY (usuario_id) REFERENCES usuarios(usuario_id) ON DELETE RESTRICT,
    CONSTRAINT chk_estado_pedido CHECK (estado_pedido IN ('PENDIENTE', 'APROBADO', 'EN_PROCESO', 'ENVIADO', 'ENTREGADO', 'CANCELADO')),
    CONSTRAINT chk_fecha_entrega CHECK (fecha_entrega_pedido >= fecha_creacion_pedido::date)
);

-- Índices para pedidos
CREATE INDEX idx_pedidos_usuario ON pedidos(usuario_id);
CREATE INDEX idx_pedidos_cotizacion ON pedidos(cotizacion_id);
CREATE INDEX idx_pedidos_estado ON pedidos(estado_pedido);
CREATE INDEX idx_pedidos_fecha_creacion ON pedidos(fecha_creacion_pedido);
CREATE INDEX idx_pedidos_fecha_entrega ON pedidos(fecha_entrega_pedido);

-- =====================================================
-- TABLA: detalle_pedido
-- =====================================================
CREATE TABLE detalle_pedido (
    pedido_id INTEGER NOT NULL,
    producto_id VARCHAR(50) NOT NULL, -- FK a tabla productos (sku)
    cantidad_producto INTEGER NOT NULL CHECK (cantidad_producto > 0),
    precio_unitario_producto DECIMAL(10,2) NOT NULL CHECK (precio_unitario_producto >= 0),
    subtotal_producto DECIMAL(12,2) NOT NULL CHECK (subtotal_producto >= 0),
    
    -- Clave primaria compuesta
    PRIMARY KEY (pedido_id, producto_id),
    
    -- Restricciones
    CONSTRAINT fk_detalle_pedido_pedido FOREIGN KEY (pedido_id) REFERENCES pedidos(pedido_id) ON DELETE CASCADE,
    CONSTRAINT fk_detalle_pedido_producto FOREIGN KEY (producto_id) REFERENCES productos(sku) ON DELETE RESTRICT,
    CONSTRAINT chk_subtotal CHECK (subtotal_producto = cantidad_producto * precio_unitario_producto)
);

-- Índices para detalle_pedido
CREATE INDEX idx_detalle_pedido_producto ON detalle_pedido(producto_id);

-- =====================================================
-- TABLA: cotizaciones
-- =====================================================
CREATE TABLE cotizaciones (
    cotizacion_id SERIAL PRIMARY KEY,
    estado_cotizacion VARCHAR(50) NOT NULL DEFAULT 'PENDIENTE',
    fecha_creacion_cotizacion TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    fecha_vencimiento_cotizacion DATE NOT NULL,
    fecha_aprobacion_cotizacion TIMESTAMP,
    total_cotizacion DECIMAL(12,2) NOT NULL CHECK (total_cotizacion >= 0),
    cliente_id INTEGER, -- FK a tabla clientes
    vendedor_id INTEGER NOT NULL, -- FK a tabla usuarios
    aprobador_id INTEGER, -- FK a tabla usuarios
    observaciones_cotizacion TEXT,
    
    -- Restricciones
    CONSTRAINT fk_cotizaciones_vendedor FOREIGN KEY (vendedor_id) REFERENCES usuarios(usuario_id) ON DELETE RESTRICT,
    CONSTRAINT fk_cotizaciones_aprobador FOREIGN KEY (aprobador_id) REFERENCES usuarios(usuario_id) ON DELETE RESTRICT,
    CONSTRAINT chk_estado_cotizacion CHECK (estado_cotizacion IN ('PENDIENTE', 'APROBADA', 'RECHAZADA', 'VENCIDA')),
    CONSTRAINT chk_fecha_vencimiento CHECK (fecha_vencimiento_cotizacion >= fecha_creacion_cotizacion::date),
    CONSTRAINT chk_fecha_aprobacion CHECK (fecha_aprobacion_cotizacion IS NULL OR fecha_aprobacion_cotizacion >= fecha_creacion_cotizacion)
);

-- Índices para cotizaciones
CREATE INDEX idx_cotizaciones_cliente ON cotizaciones(cliente_id);
CREATE INDEX idx_cotizaciones_vendedor ON cotizaciones(vendedor_id);
CREATE INDEX idx_cotizaciones_estado ON cotizaciones(estado_cotizacion);
CREATE INDEX idx_cotizaciones_fecha_creacion ON cotizaciones(fecha_creacion_cotizacion);
CREATE INDEX idx_cotizaciones_fecha_vencimiento ON cotizaciones(fecha_vencimiento_cotizacion);

-- Agregar FK de pedidos a cotizaciones después de crear la tabla cotizaciones
ALTER TABLE pedidos ADD CONSTRAINT fk_pedidos_cotizacion FOREIGN KEY (cotizacion_id) REFERENCES cotizaciones(cotizacion_id) ON DELETE SET NULL;

-- =====================================================
-- TABLA: tipo_camion
-- =====================================================
CREATE TABLE tipo_camion (
    tipo_camion_id SERIAL PRIMARY KEY,
    capacidad_camion DECIMAL(8,2) NOT NULL CHECK (capacidad_camion > 0),
    nombre_tipo_camion VARCHAR(100) NOT NULL UNIQUE,
    descripcion_tipo_camion TEXT,
    activo_tipo_camion BOOLEAN DEFAULT TRUE
);

-- Índices para tipo_camion
CREATE INDEX idx_tipo_camion_activo ON tipo_camion(activo_tipo_camion);

-- =====================================================
-- TABLA: camiones
-- =====================================================
CREATE TABLE camiones (
    camion_id SERIAL PRIMARY KEY,
    patente_camion VARCHAR(10) NOT NULL UNIQUE,
    tipo_camion_id INTEGER NOT NULL,
    sucursal_id INTEGER, -- FK a tabla sucursales
    responsable_camion_id INTEGER, -- FK a tabla usuarios
    descripcion_camion TEXT,
    activo_camion BOOLEAN DEFAULT TRUE,
    
    -- Restricciones
    CONSTRAINT fk_camiones_tipo FOREIGN KEY (tipo_camion_id) REFERENCES tipo_camion(tipo_camion_id) ON DELETE RESTRICT,
    CONSTRAINT fk_camiones_responsable FOREIGN KEY (responsable_camion_id) REFERENCES usuarios(usuario_id) ON DELETE SET NULL,
    CONSTRAINT chk_patente_camion CHECK (patente_camion ~ '^[A-Z]{2}[A-Z0-9]{2}[0-9]{2}$')
);

-- Índices para camiones
CREATE INDEX idx_camiones_patente ON camiones(patente_camion);
CREATE INDEX idx_camiones_tipo ON camiones(tipo_camion_id);
CREATE INDEX idx_camiones_sucursal ON camiones(sucursal_id);
CREATE INDEX idx_camiones_activo ON camiones(activo_camion);

-- =====================================================
-- TABLA: despachos
-- =====================================================
CREATE TABLE despachos (
    despacho_id SERIAL PRIMARY KEY,
    pedido_id INTEGER NOT NULL,
    camion_id INTEGER NOT NULL,
    origen_despacho VARCHAR(200) NOT NULL,
    destino_despacho VARCHAR(200) NOT NULL,
    estado_despacho VARCHAR(50) NOT NULL DEFAULT 'PENDIENTE',
    fecha_envio_despacho TIMESTAMP,
    fecha_entrega_despacho TIMESTAMP,
    costo_despacho DECIMAL(10,2) CHECK (costo_despacho >= 0),
    observaciones_despacho TEXT,
    
    -- Restricciones
    CONSTRAINT fk_despachos_pedido FOREIGN KEY (pedido_id) REFERENCES pedidos(pedido_id) ON DELETE RESTRICT,
    CONSTRAINT fk_despachos_camion FOREIGN KEY (camion_id) REFERENCES camiones(camion_id) ON DELETE RESTRICT,
    CONSTRAINT chk_estado_despacho CHECK (estado_despacho IN ('PENDIENTE', 'EN_CAMINO', 'ENTREGADO', 'CANCELADO')),
    CONSTRAINT chk_fecha_entrega CHECK (fecha_entrega_despacho IS NULL OR fecha_entrega_despacho >= fecha_envio_despacho)
);

-- Índices para despachos
CREATE INDEX idx_despachos_pedido ON despachos(pedido_id);
CREATE INDEX idx_despachos_camion ON despachos(camion_id);
CREATE INDEX idx_despachos_estado ON despachos(estado_despacho);
CREATE INDEX idx_despachos_fecha_envio ON despachos(fecha_envio_despacho);

-- =====================================================
-- TABLA: envio_camion
-- =====================================================
CREATE TABLE envio_camion (
    envio_camion_id SERIAL PRIMARY KEY,
    despacho_id INTEGER NOT NULL,
    camion_id INTEGER NOT NULL,
    estado_envio_camion VARCHAR(50) NOT NULL DEFAULT 'PENDIENTE',
    fecha_envio_camion TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    fecha_entrega_camion TIMESTAMP,
    observaciones_envio_camion TEXT,
    
    -- Restricciones
    CONSTRAINT fk_envio_camion_despacho FOREIGN KEY (despacho_id) REFERENCES despachos(despacho_id) ON DELETE CASCADE,
    CONSTRAINT fk_envio_camion_camion FOREIGN KEY (camion_id) REFERENCES camiones(camion_id) ON DELETE RESTRICT,
    CONSTRAINT chk_estado_envio_camion CHECK (estado_envio_camion IN ('PENDIENTE', 'EN_TRANSITO', 'ENTREGADO', 'CANCELADO')),
    CONSTRAINT chk_fecha_entrega_camion CHECK (fecha_entrega_camion IS NULL OR fecha_entrega_camion >= fecha_envio_camion)
);

-- Índices para envio_camion
CREATE INDEX idx_envio_camion_despacho ON envio_camion(despacho_id);
CREATE INDEX idx_envio_camion_camion ON envio_camion(camion_id);
CREATE INDEX idx_envio_camion_estado ON envio_camion(estado_envio_camion);

-- =====================================================
-- TABLA: detalle_envio
-- =====================================================
CREATE TABLE detalle_envio (
    envio_camion_id INTEGER NOT NULL,
    producto_id VARCHAR(50) NOT NULL, -- FK a tabla productos (sku)
    cantidad_producto INTEGER NOT NULL CHECK (cantidad_producto > 0),
    
    -- Clave primaria compuesta
    PRIMARY KEY (envio_camion_id, producto_id),
    
    -- Restricciones
    CONSTRAINT fk_detalle_envio_envio FOREIGN KEY (envio_camion_id) REFERENCES envio_camion(envio_camion_id) ON DELETE CASCADE,
    CONSTRAINT fk_detalle_envio_producto FOREIGN KEY (producto_id) REFERENCES productos(sku) ON DELETE RESTRICT
);

-- Índices para detalle_envio
CREATE INDEX idx_detalle_envio_producto ON detalle_envio(producto_id);

-- =====================================================
-- TABLA: stock_sucursal
-- =====================================================
CREATE TABLE stock_sucursal (
    stock_sucursal_id SERIAL PRIMARY KEY,
    producto_id VARCHAR(50) NOT NULL, -- FK a tabla productos (sku)
    sucursal_id INTEGER NOT NULL, -- FK a tabla sucursales
    stock_producto INTEGER DEFAULT 0 CHECK (stock_producto >= 0),
    stock_minimo_producto INTEGER DEFAULT 0 CHECK (stock_minimo_producto >= 0),
    
    -- Restricciones
    CONSTRAINT fk_stock_sucursal_producto FOREIGN KEY (producto_id) REFERENCES productos(sku) ON DELETE CASCADE,
    CONSTRAINT uk_stock_sucursal UNIQUE (producto_id, sucursal_id)
);

-- Índices para stock_sucursal
CREATE INDEX idx_stock_sucursal_producto ON stock_sucursal(producto_id);
CREATE INDEX idx_stock_sucursal_sucursal ON stock_sucursal(sucursal_id);
CREATE INDEX idx_stock_sucursal_stock ON stock_sucursal(stock_producto);

-- =====================================================
-- COMENTARIOS EN TABLAS Y COLUMNAS
-- =====================================================

-- Comentarios en tablas
COMMENT ON TABLE usuarios IS 'Tabla que almacena la información de usuarios del sistema';
COMMENT ON TABLE productos IS 'Tabla que almacena la información de productos';
COMMENT ON TABLE pedidos IS 'Tabla que almacena los pedidos realizados';
COMMENT ON TABLE cotizaciones IS 'Tabla que almacena las cotizaciones generadas';
COMMENT ON TABLE despachos IS 'Tabla que almacena la información de despachos';
COMMENT ON TABLE camiones IS 'Tabla que almacena la información de camiones disponibles';

-- Comentarios en columnas importantes
COMMENT ON COLUMN usuarios.rut_usuario IS 'RUT único del usuario (formato: 12345678-9)';
COMMENT ON COLUMN productos.sku IS 'Código único del producto';
COMMENT ON COLUMN pedidos.estado_pedido IS 'Estado del pedido: PENDIENTE, APROBADO, EN_PROCESO, ENVIADO, ENTREGADO, CANCELADO';
COMMENT ON COLUMN cotizaciones.estado_cotizacion IS 'Estado de la cotización: PENDIENTE, APROBADA, RECHAZADA, VENCIDA';
