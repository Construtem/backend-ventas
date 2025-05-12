# Backend - Ventas

Este repositorio contiene el microservicio backend para el módulo de **ventas** en Construtem.

## 🛠️ Tecnologías
- Go
- REST API
- PostgreSQL (compartida con facturación)
- JWT para autenticación
- Docker

## 🚀 Funcionalidades
- Crear, actualizar y consultar órdenes de venta.
- Consultar disponibilidad de productos.
- Validaciones de negocio para ventas.

## 📦 Instalación y ejecución
```bash
docker-compose up --build
```

## 📁 Estructura del proyecto
- `/controllers`
- `/routes`
- `/models`
- `/services`

## 🧩 Integraciones
- Comunicación con `backend-facturacion` para emisión de boletas/facturas.
- Almacenamiento en base de datos compartida con facturación.
