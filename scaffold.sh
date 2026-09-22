#!/bin/bash

# Script para generar la estructura de carpetas (Arquitectura Hexagonal)
echo "Generando estructura de carpetas para el nuevo microservicio Go..."

# Capa de Inicio
mkdir -p cmd/api/config

# Capa de Aplicacion (Transporte)
mkdir -p internal/app/api/dto
mkdir -p internal/app/api/handler/middleware
mkdir -p internal/app/api/mapper

# Capa Core (Lógica de Negocio y Dominio)
mkdir -p internal/core/domain
mkdir -p internal/core/usecases

# Capa de Infraestructura
mkdir -p internal/infra/db/oracle
mkdir -p internal/infra/tracerMiddleware

echo "¡Estructura de carpetas generada exitosamente!"
echo "Puedes revisar las carpetas creadas en 'cmd' e 'internal'."
