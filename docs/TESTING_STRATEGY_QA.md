# Estrategia y Guía de Testing - Catálogos de Créditos Hipotecarios (Go)

Este documento está dirigido a Desarrolladores y Analistas de Calidad (QA) para comprender la estrategia de validación, automatización y certificación del nuevo microservicio Go Hexagonal, asegurando la paridad total con el sistema legado Java.

---

## 1. Niveles de Prueba Automatizada (Pirámide de Testing en Go)

El proyecto utiliza la librería estándar `testing` combinada con `github.com/stretchr/testify` para aserciones (`assert`) y simulaciones (`mock`). Hemos dividido las pruebas en 4 capas aislando responsabilidades:

### A. Pruebas Unitarias de Dominio y Servicio (`/internal/core/services`)
- **Objetivo:** Validar la lógica de negocio pura.
- **Técnica:** Se inyecta un `MockCatalogRepository` para aislar el servicio de llamadas a bases de datos o HTTP.
- **Cobertura actual:** Escenarios de catálogo encontrado (`200 OK`) y no encontrado (`ErrCatalogNotFound`).

### B. Pruebas de Integración de Entrada / API (`/internal/app/handler`)
- **Objetivo:** Validar el enrutador HTTP (Gin Framework) y el parseo de peticiones y respuestas JSON.
- **Técnica:** Se levanta un servidor de pruebas ligero mediante `httptest.NewRecorder()` inyectando un `MockCatalogService`.
- **Cobertura actual:** Parseo de variables de entorno (`X-Backend-Env`), mapeo de JSON tags hacia el cliente, validación de que *no* se exigen headers de negocio obligatorios (paridad con Java).

### C. Pruebas de Integración de Salida / Upstream (`/internal/infra/repository`)
- **Objetivo:** Validar que el cliente HTTP realice la petición correcta hacia *Finnflow* y procese sus datos y errores.
- **Técnica:** Se interceptan las llamadas de red reales inyectando un `httptest.NewServer` que simula al proveedor externo.
- **Cobertura actual:** 
  - Verificación estricta de envío de credenciales (`Basic Auth`).
  - Mapeos de deserialización complejos (ej. Campo `Factor` manejando notación científica `2.255E-4`).
  - Captura y transformación de errores `401 Unauthorized` y `404 Not Found`.
  - Simulación de *Timeouts* / *Connection Refused* resolviendo al payload de respaldo: `Sin resultados`.

### D. Pruebas End-to-End (E2E) (`/test/e2e`)
- **Objetivo:** Simular la vida completa de un Request desde el inicio (Middleware) hasta el fin (Respuesta).
- **Técnica:** Cablea toda la aplicación real usando el repositorio en memoria (`DummyRepository`) para no depender de la red externa en pipelines CI/CD.
- **Cobertura actual:** Flujo completo de un caso genérico, validando adicionalmente que **no se haga echo** de headers de entrada en la respuesta.

---

## 2. Ejecución Local (Solución a bloqueo de Windows Defender)

**Problema:** Al ejecutar el comando estándar `go test ./...`, Go compila binarios temporales en `%TEMP%`. En ciertos entornos corporativos Windows, el Antivirus bloquea la ejecución de estos archivos temporales `.test.exe`.

**Solución para el equipo (Devs & QA):** Compilar los ejecutables de test de forma explícita dentro del proyecto (`/bin`) y ejecutarlos directamente. 

Puedes copiar y pegar este bloque en una consola **PowerShell** en la raíz del proyecto para ejecutar toda la suite:

```powershell
# Crear directorio de binarios si no existe
mkdir -Force bin

Write-Host "--- EJECUTANDO TESTS DE SERVICIO (Unitarios) ---" -ForegroundColor Cyan
go test -c -o bin/services.test.exe ./internal/core/services
cmd /c ".\bin\services.test.exe -test.v"

Write-Host "--- EJECUTANDO TESTS DE HANDLER (Integración API) ---" -ForegroundColor Cyan
go test -c -o bin/handler.test.exe ./internal/app/handler
cmd /c ".\bin\handler.test.exe -test.v"

Write-Host "--- EJECUTANDO TESTS DE REPOSITORIO (Integración Upstream) ---" -ForegroundColor Cyan
go test -c -o bin/repository.test.exe ./internal/infra/repository
cmd /c ".\bin\repository.test.exe -test.v"

Write-Host "--- EJECUTANDO TESTS END-TO-END (E2E) ---" -ForegroundColor Cyan
go test -c -o bin/e2e.test.exe ./test/e2e
cmd /c ".\bin\e2e.test.exe -test.v"
```

---

## 3. Pruebas Funcionales (Guía para QA Manual)

El equipo de QA dispone de **3 Artefactos Oficiales** alojados en la carpeta `docs/` para su proceso de certificación:

1. **`CASOS_DE_USO_LEGADO.md`**: Es la "Biblia" de producto. Detalla el árbol de decisiones de enrutamiento exacto (`equals` vs `startsWith`), mapeos de error (401, 402, 400), y las reglas de negocio base.
2. **`PAYLOADS_CASOS_DE_USO.md`**: Capturas literales (JSON) en vivo comprobando que lo que emite Java Legado es 100% igual a lo que emite Go.
3. **`Mortgage_Loan_Catalogs.postman_collection.json`**: Colección con los 14 escenarios ya configurados. 

### Pasos para validación con Postman:
1. Asegurarse de tener encendido el Simulador local (`9090`), el legado Java (`8080`) y el microservicio Go (`8082`).
2. Importar el archivo JSON en Postman.
3. Seleccionar la colección. Modificar la variable **`BASE_URL`** a:
   - `http://localhost:8080` (Para validar cómo responde el legado).
   - `http://localhost:8082` (Para validar cómo responde el nuevo código Go).
4. QA deberá certificar que ejecutando la misma petición en ambos `BASE_URL`, el código HTTP y el JSON resultante sean **idénticos** en los 14 escenarios.
