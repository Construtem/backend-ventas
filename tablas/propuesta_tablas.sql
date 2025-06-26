Propuesta de tablas para la base de datos:

usuarios{
    usuario_id
    nombre_usuario
    rut_usuario // nuevo
    email_usuario
    contrasena_usuario
    rol_id
    sucursal_id // opcional
    fecha_creacion_usuario // opcional
    activo_usuario // opcional
    fecha_ultimo_acceso_usuario // opcional
}
roles{
    rol_id
    nombre_rol
    descripcion_rol
}
productos{
    sku
    nombre_producto
    descripcion_producto
    precio_producto
    stock_producto
    largo_producto
    ancho_producto
    alto_producto
    peso_producto
    unidad_medida
    categoria_id
    proveedor_id
    sucursal_id
    activo_producto // opcional
    fecha_creacion_producto // opcional
    imagen_producto // opcional
}
pedidos{
    pedido_id
    cotizacion_id // en caso de que el pedido sea de una cotizacion
    usuario_id // en caso de que el pedido sea de un cliente
    fecha_creacion_pedido // opcional
    fecha_entrega_pedido // opcional
    direccion_entrega_pedido // opcional
    estado_pedido
    metodo_pago
    total_pedido
    observaciones_pedido // opcional
}
detalle_pedido{
    pedido_id // PK
    producto_id // PK
    cantidad_producto
    precio_unitario_producto
    subtotal_producto
}
cotizaciones{
    cotizacion_id
    estado_cotizacion
    fecha_creacion_cotizacion
    fecha_vencimiento_cotizacion
    fecha_aprobacion_cotizacion
    total_cotizacion
    cliente_id
    vendedor_id
    aprobador_id
    observaciones_cotizacion // opcional
}
despachos{
    despacho_id
    pedido_id
    camion_id
    origen_despacho
    destino_despacho
    estado_despacho
    fecha_envio_despacho
    fecha_entrega_despacho
    costo_despacho
    observaciones_despacho
}
camiones{
    camion_id
    patente_camion
    tipo_camion_id
    sucursal_id // opcional
    responsable_camion_id // opcional
    descripcion_camion // opcional
    activo_camion // opcional
}
tipo_camion{
    tipo_camion_id
    capacidad_camion
    nombre_tipo_camion
    descripcion_tipo_camion
    activo_tipo_camion // opcional
}
detalle_envio{
    envio_camion_id // PK
    producto_id // PK
    cantidad_producto
}
envio_camion{
    envio_camion_id
    despacho_id
    camion_id
    estado_envio_camion
    fecha_envio_camion
    fecha_entrega_camion
    observaciones_envio_camion
}
stock_sucursal{
    stock_sucursal_id
    producto_id
    sucursal_id
    stock_producto
    stock_minimo_producto
}
