# Guía de Distribución y Pruebas Locales (Compañeros / QA)

Esta guía explica paso a paso cómo empaquetar el microservicio con los **últimos cambios aplicados** (donde se eliminaron los valores hardcodeados) para entregárselo a un compañero (usuario de Mac o Windows) para que pueda probarlo sin necesidad de instalar Go ni descargar código fuente.

---

## PASO 1: Compilar los binarios desde tu máquina (Windows)

Abre una consola **PowerShell** en la raíz de tu proyecto (`C:\HEXAGONAL\mortgage-loan-catalogs`). Vamos a generar los binarios estáticos optimizados.

Si tu compañero usa una **Mac Moderna (Procesador M1, M2, M3 - Apple Silicon)**:
```powershell
$env:GOOS="darwin"; $env:GOARCH="arm64"; go build -ldflags="-s -w" -o mortgage-loan-catalogs-mac cmd/api/main.go
```

Si tu compañero usa una **Mac Antigua (Procesador Intel)**:
```powershell
$env:GOOS="darwin"; $env:GOARCH="amd64"; go build -ldflags="-s -w" -o mortgage-loan-catalogs-mac cmd/api/main.go
```

*(Si necesitas uno para otro compañero con Windows)*:
```powershell
$env:GOOS="windows"; $env:GOARCH="amd64"; go build -ldflags="-s -w" -o mortgage-loan-catalogs.exe cmd/api/main.go
```

---

## PASO 2: Preparar los Scripts de Arranque (Estrictos)

Como eliminamos el código *hardcodeado* por seguridad de arquitectura, el servicio ahora **exige** que se le pasen las variables `PORT` y `FINNFLOW_URL`, o de lo contrario no encenderá.

Crea un archivo de texto llamado **`iniciar_servicio.sh`** (para Mac/Linux) e incluye exactamente esto:

```bash
#!/bin/bash
echo "Iniciando Microservicio de Catalogos..."

# 1. Variables de entorno OBLIGATORIAS
export PORT=8082
export FINNFLOW_URL="http://localhost:9090"

# (Opcional) Si en el futuro prueban el legado Java, descomentar esta linea:
# export JAVA_LEGACY_URL="http://localhost:8080"

# 2. Permisos de ejecución
chmod +x ./mortgage-loan-catalogs-mac

# 3. Ejecutar binario
./mortgage-loan-catalogs-mac
```

---

## PASO 3: Empaquetar y Enviar

Crea un archivo `.zip` que contenga:
1. El binario `mortgage-loan-catalogs-mac` (generado en el paso 1).
2. El script `iniciar_servicio.sh` (creado en el paso 2).
3. *(Si tu compañero no tiene forma de acceder al Finnflow real)*, envíale también el simulador `simulador-finnflow.exe` (o su equivalente para Mac).

---

## PASO 4: Instrucciones para enviarle a tu compañero (Copy & Paste)

Cópiale el siguiente texto a tu compañero junto con el `.zip`:

> **Hola! Aquí tienes el microservicio listo para probar.**
> 
> **Instrucciones (Mac):**
> 1. Descomprime el `.zip` en una carpeta.
> 2. Abre la terminal, arrastra la carpeta hacia la ventana de la terminal y presiona Enter para ubicarte allí.
> 3. Ejecuta estos comandos:
>    `chmod +x iniciar_servicio.sh`
>    `./iniciar_servicio.sh`
> 
> **⚠️ IMPORTANTE (Seguridad de Mac / Gatekeeper):**
> Como te estoy enviando este archivo directamente, la primera vez que lo corras tu Mac puede lanzar una alerta diciendo: *"no se puede abrir porque el desarrollador no está verificado"*. 
> Si esto ocurre, ve a **Preferencias del Sistema > Privacidad y Seguridad**, baja hasta el final y haz clic en **"Permitir de todos modos"** (Allow anyway) junto al nombre del archivo, y vuelve a correr el script. ¡Y listo! Estará escuchando en el puerto 8082.
