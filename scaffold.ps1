# Script para generar la estructura de carpetas (Arquitectura Hexagonal)
Write-Host "Generando estructura de carpetas para el nuevo microservicio Go..." -ForegroundColor Cyan

# Capa de Inicio
New-Item -ItemType Directory -Force -Path "cmd\api\config" | Out-Null

# Capa de Aplicacion (Transporte)
New-Item -ItemType Directory -Force -Path "internal\app\api\dto" | Out-Null
New-Item -ItemType Directory -Force -Path "internal\app\api\handler\middleware" | Out-Null
New-Item -ItemType Directory -Force -Path "internal\app\api\mapper" | Out-Null

# Capa Core (Lógica de Negocio y Dominio)
New-Item -ItemType Directory -Force -Path "internal\core\domain" | Out-Null
New-Item -ItemType Directory -Force -Path "internal\core\usecases" | Out-Null

# Capa de Infraestructura
New-Item -ItemType Directory -Force -Path "internal\infra\db\oracle" | Out-Null
New-Item -ItemType Directory -Force -Path "internal\infra\tracerMiddleware" | Out-Null

Write-Host "¡Estructura de carpetas generada exitosamente!" -ForegroundColor Green
Write-Host "Puedes revisar las carpetas creadas en 'cmd' e 'internal'." -ForegroundColor Yellow
