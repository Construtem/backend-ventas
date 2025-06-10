-- CORE

CREATE TABLE usuario_rol (
    id_usuario INT NOT NULL, -- FK
    id_lov INT NOT NULL, -- FK
    PRIMARY KEY (id_usuario, id_lov)
);

CREATE TABLE usuarios (
    id_usuario SERIAL PRIMARY KEY,
    nombre VARCHAR(100) NOT NULL,
    apellido VARCHAR(100), -- EXTRA
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    activo BOOLEAN
);

CREATE TABLE lov_universal (
    id_lov SERIAL PRIMARY KEY,
    categoria VARCHAR(100) NOT NULL,
    valor VARCHAR(100) NOT NULL,
    descripcion VARCHAR(255),
    orden INT, -- NO SÉ QUE ES
    activo BOOLEAN
);

-- VENTAS Y CORE COMERCIAL

CREATE TABLE boleta (
    id_boleta SERIAL PRIMARY KEY,
    id_venta INT NOT NULL, --FK
    fecha_emision TIMESTAMP NOT NULL,
    total DECIMAL(10, 2) NOT NULL,
    estado VARCHAR(50) NOT NULL,
    pdf_url VARCHAR(255)
);

CREATE TABLE venta_forma_entrega (
    id_venta INT NOT NULL, --FK
    id_lov INT NOT NULL, -- FK
    PRIMARY KEY (id_venta, id_lov)
);

CREATE TABLE ventas_metodo_pago (
    id_venta INT NOT NULL, --FK
    id_lov INT NOT NULL, -- FK
    PRIMARY KEY (id_venta, id_lov)
);

CREATE TABLE descuento_etiqueta (
    id_descuento INT NOT NULL, --FK
    id_lov INT NOT NULL, --FK
    PRIMARY KEY (id_descuento, id_lov)
);

CREATE TABLE despachos (
    id_despacho SERIAL PRIMARY KEY,
    id_venta INT NOT NULL, --FK
    direccion_entrega VARCHAR(255) NOT NULL,
    telefono_contacto VARCHAR(20) NOT NULL, -- Nuevo campo
    estado_envio VARCHAR(50) NOT NULL,
    volumen_total DECIMAL(10, 2),
    id_camion INT NOT NULL -- FK
);

CREATE TABLE ventas (
    id_venta SERIAL PRIMARY KEY,
    id_carrito INT NOT NULL, -- FK
    fecha TIMESTAMP NOT NULL,
    total DECIMAL(10, 2) NOT NULL,
    estado VARCHAR(50) NOT NULL
);

CREATE TABLE detalle_carrito (
    id_detalle_carrito SERIAL PRIMARY KEY,
    id_carrito INT NOT NULL, -- FK
    id_producto INT NOT NULL, -- FK
    cantidad INT NOT NULL,
    precio_unitario DECIMAL(10, 2) NOT NULL,
    descuento_aplicado DECIMAL(10, 2)
);

CREATE TABLE descuentos (
    id_descuento SERIAL PRIMARY KEY,
    nombre VARCHAR(100) NOT NULL,
    descripcion VARCHAR(255),
    porcentaje DECIMAL(5, 2), -- Ejemplo: 0.15 para 15%
    monto_fijo DECIMAL(10, 2), -- EXTRA, No es NOT NULL, si solo aplica porcentaje
    fecha_inicio TIMESTAMP NOT NULL,
    fecha_fin TIMESTAMP NOT NULL,
    activo BOOLEAN NOT NULL DEFAULT FALSE -- TIENE EXTRAS
);

CREATE TABLE carrito_compras (
    id_carrito SERIAL PRIMARY KEY,
    id_usuario INT NOT NULL, -- FK
    id_cliente INT NOT NULL, -- FK
    fecha TIMESTAMP NOT NULL,
    estado VARCHAR(50) NOT NULL
);

CREATE TABLE clientes (
    id_cliente SERIAL PRIMARY KEY,
    nombre VARCHAR(100) NOT NULL,
    apellido VARCHAR(100), -- EXTRA
    telefono VARCHAR(20),
    email VARCHAR(255) UNIQUE,
    direccion VARCHAR(255)
);

-- VENTAS POR PEDIDO

CREATE TABLE ventas_por_pedido (
    id_venta_pedido SERIAL PRIMARY KEY,
    id_cliente INT NOT NULL, -- FK
    id_proveedor INT NOT NULL, -- FK
    id_producto_prov INT NOT NULL, -- FK
    fecha TIMESTAMP,
    precio_final DECIMAL(10, 2),
    estado VARCHAR(255) NOT NULL
);

CREATE TABLE productos_proveedor (
    id_producto_prov SERIAL PRIMARY KEY,
    id_proveedor INT NOT NULL, -- FK
    nombre VARCHAR(100) NOT NULL,
    descripcion VARCHAR(255),
    precio DECIMAL(10, 2) NOT NULL,
    dimensiones VARCHAR(100),
    peso_kg DECIMAL(10, 2),
    activo BOOLEAN
);

CREATE TABLE proveedores (
    id_proveedor SERIAL PRIMARY KEY,
    nombre VARCHAR(100) NOT NULL,
    email VARCHAR(255) UNIQUE,
    telefono VARCHAR(20),
    direccion VARCHAR(255)
);

-- INVENTARIO Y DESPACHO

CREATE TABLE camiones (
    id_camion SERIAL PRIMARY KEY,
    patente VARCHAR(20) UNIQUE NOT NULL,
    capacidad_kg DECIMAL(10, 2),
    volumen_m3 DECIMAL(10, 2),
    activo BOOLEAN,
    id_lov INT NOT NULL -- FK (Ej: tipo de camión, estado, etc. de LOV_Universal)
);

CREATE TABLE producto_tipo (
    id_producto INT NOT NULL, -- FK
    id_lov INT NOT NULL, -- FK (Ej: tipo de producto de LOV_Universal)
    PRIMARY KEY (id_producto, id_lov)
);

CREATE TABLE productos_etiqueta (
    id_producto INT NOT NULL, -- FK
    id_lov INT NOT NULL, -- FK
    PRIMARY KEY (id_producto, id_lov)
);

CREATE TABLE bodegas_centrales (
    id_bodega SERIAL PRIMARY KEY,
    nombre VARCHAR(100) NOT NULL,
    direccion VARCHAR(255),
    telefono VARCHAR(20)
);

CREATE TABLE stock_sucursal (
    id_stock SERIAL PRIMARY KEY,
    id_sucursal INT NOT NULL, -- FK
    id_producto INT NOT NULL, -- FK
    cantidad INT NOT NULL
);

CREATE TABLE stock_bodega_central (
    id_stock SERIAL PRIMARY KEY,
    id_bodega INT NOT NULL, -- FK
    id_producto INT NOT NULL, -- FK
    cantidad INT NOT NULL
);

CREATE TABLE productos (
    id_producto SERIAL PRIMARY KEY,
    sku VARCHAR(50) UNIQUE NOT NULL,
    nombre VARCHAR(255) NOT NULL,
    descripcion VARCHAR(255),
    dimensiones VARCHAR(100),
    peso_kg DECIMAL(10, 2),
    precio_venta DECIMAL(10, 2),
    activo BOOLEAN
);
-- Sucursales
CREATE TABLE sucursales (
    id_sucursal SERIAL PRIMARY KEY,
    nombre VARCHAR(100) NOT NULL,
    direccion VARCHAR(255),
    telefono VARCHAR(20)
);
-- ############################################################################################################################################################################################
-- #################################################################  Creación de Claves Foráneas (FKs)  ######################################################################################
-- ############################################################################################################################################################################################

-- CORE
ALTER TABLE usuario_rol
ADD CONSTRAINT fk_usuario_rol_usuario FOREIGN KEY (id_usuario) REFERENCES usuarios(id_usuario),
ADD CONSTRAINT fk_usuario_rol_lov FOREIGN KEY (id_lov) REFERENCES lov_universal(id_lov);

-- VENTAS Y CORE COMERCIAL
ALTER TABLE boleta
ADD CONSTRAINT fk_boleta_venta FOREIGN KEY (id_venta) REFERENCES ventas(id_venta);

ALTER TABLE venta_forma_entrega
ADD CONSTRAINT fk_venta_forma_entrega_venta FOREIGN KEY (id_venta) REFERENCES ventas(id_venta),
ADD CONSTRAINT fk_venta_forma_entrega_lov FOREIGN KEY (id_lov) REFERENCES lov_universal(id_lov);

ALTER TABLE ventas_metodo_pago
ADD CONSTRAINT fk_ventas_metodo_pago_venta FOREIGN KEY (id_venta) REFERENCES ventas(id_venta),
ADD CONSTRAINT fk_ventas_metodo_pago_lov FOREIGN KEY (id_lov) REFERENCES lov_universal(id_lov);

ALTER TABLE descuento_etiqueta
ADD CONSTRAINT fk_descuento_a_etiqueta_descuento FOREIGN KEY (id_descuento) REFERENCES descuentos(id_descuento),
ADD CONSTRAINT fk_descuento_a_etiqueta_lov FOREIGN KEY (id_lov) REFERENCES lov_universal(id_lov);

ALTER TABLE despachos
ADD CONSTRAINT fk_despachos_venta FOREIGN KEY (id_venta) REFERENCES ventas(id_venta),
ADD CONSTRAINT fk_despachos_camion FOREIGN KEY (id_camion) REFERENCES camiones(id_camion);

ALTER TABLE ventas
ADD CONSTRAINT fk_ventas_carrito FOREIGN KEY (id_carrito) REFERENCES carrito_compras(id_carrito);

ALTER TABLE detalle_carrito
ADD CONSTRAINT fk_detalle_carrito_carrito FOREIGN KEY (id_carrito) REFERENCES carrito_compras(id_carrito),
ADD CONSTRAINT fk_detalle_carrito_producto FOREIGN KEY (id_producto) REFERENCES productos(id_producto);

ALTER TABLE carrito_compras
ADD CONSTRAINT fk_carrito_compras_usuario FOREIGN KEY (id_usuario) REFERENCES usuarios(id_usuario),
ADD CONSTRAINT fk_carrito_compras_cliente FOREIGN KEY (id_cliente) REFERENCES clientes(id_cliente);

-- VENTAS POR PEDIDO
ALTER TABLE ventas_por_pedido
ADD CONSTRAINT fk_ventas_por_pedido_cliente FOREIGN KEY (id_cliente) REFERENCES clientes(id_cliente),
ADD CONSTRAINT fk_ventas_por_pedido_proveedor FOREIGN KEY (id_proveedor) REFERENCES proveedores(id_proveedor),
ADD CONSTRAINT fk_ventas_por_pedido_producto_prov FOREIGN KEY (id_producto_prov) REFERENCES productos_proveedor(id_producto_prov);

ALTER TABLE productos_proveedor
ADD CONSTRAINT fk_productos_de_proveedor_proveedor FOREIGN KEY (id_proveedor) REFERENCES proveedores(id_proveedor);

-- INVENTARIO Y DESPACHO
ALTER TABLE camiones
ADD CONSTRAINT fk_camiones_lov FOREIGN KEY (id_lov) REFERENCES lov_universal(id_lov);

ALTER TABLE producto_tipo
ADD CONSTRAINT fk_producto_tipo_producto FOREIGN KEY (id_producto) REFERENCES productos(id_producto),
ADD CONSTRAINT fk_producto_tipo_lov FOREIGN KEY (id_lov) REFERENCES lov_universal(id_lov);

ALTER TABLE productos_etiqueta
ADD CONSTRAINT fk_productos_etiqueta_producto FOREIGN KEY (id_producto) REFERENCES productos(id_producto),
ADD CONSTRAINT fk_productos_etiqueta_lov FOREIGN KEY (id_lov) REFERENCES lov_universal(id_lov);

ALTER TABLE stock_sucursal
ADD CONSTRAINT fk_stock_sucursal_sucursal FOREIGN KEY (id_sucursal) REFERENCES sucursales(id_sucursal),
ADD CONSTRAINT fk_stock_sucursal_producto FOREIGN KEY (id_producto) REFERENCES productos(id_producto);

ALTER TABLE stock_bodega_central
ADD CONSTRAINT fk_stock_bodega_central_bodega FOREIGN KEY (id_bodega) REFERENCES bodegas_centrales(id_bodega),
ADD CONSTRAINT fk_stock_bodega_central_producto FOREIGN KEY (id_producto) REFERENCES productos(id_producto);