# Contrato de API y Validaciones de Paridad

Este documento sirve como referencia oficial para QA, Líderes Técnicos y Consumidores sobre cómo interactuar con el microservicio `mortgage-loan-catalogs` en Go, demostrando la retrocompatibilidad exacta con el servicio heredado en Java.

## 1. Rutas
El microservicio expone su funcionalidad a través del siguiente endpoint:
*   **Ruta:** `POST /v1/bfcl/mortgage-loan/catalogs/:catalog`
*   **Ejemplo:** `POST /v1/bfcl/mortgage-loan/catalogs/Destino`
*   **Swagger:** `GET /swagger/index.html`

## 2. Variables de Entorno (Configuración)
El servicio es 100% parametrizable mediante variables de entorno inyectadas en tiempo de despliegue:
*   `FINNFLOW_URL`: URL base del upstream (ej. `http://localhost:9090`).
*   `FINNFLOW_KEY`: Client ID para autenticación hacia Finnflow (`client_id`).
*   `FINNFLOW_SECRET`: Client Secret para autenticación hacia Finnflow (`client_secret`).
*   `JAVA_LEGACY_URL`: (Solo para QA/Swagger) URL del legacy container.
*   `DEFAULT_BACKEND`: Define la estrategia de proxy. Si es `real`, el tráfico regular consume Finnflow por defecto.
*   `PORT`: Puerto donde escucha la aplicación (ej. `8082`).

## 3. Request (Headers)
Por exigencia de paridad, el consumidor debe proveer **obligatoriamente** los siguientes headers:
*   `X-Channel`: Canal del origen (ej. `WEB`).
*   `X-Commerce`: Comercio (ej. `FALABELLA`).
*   `X-Transaction-ID`: ID de traza (ej. `12345678`).

*(Opcional - Feature Toggle para pruebas)*:
*   `X-Backend-Env`: Permite, desde Swagger, sobreescribir el upstream por request (`real`, `java`, o `dummy`).

## 4. Responses (Cuerpos JSON y Status Codes)
La API garantiza respuesta idéntica (paridad byte a byte) a las respuestas del framework de Java original:
*   **200 OK**: Para catálogos exitosos (Catálogos estándar o Seguros).
*   **200 OK (Sin Resultados)**: Fallback por red/timeout o catálogos vacíos. Devuelve el formato estricto: `[{"codRespuesta": 3, "Mensaje": "Sin resultados.", "Excepcion": "Ninguna"}]`.

## 5. Manejo de Errores
Se mantienen los status codes de error exactos del legado:
*   **401 Unauthorized**: Si el Upstream devuelve un 401 por autenticación denegada. Devuelve la estructura idéntica a `TokenNotValidException`.
*   **402 Payment Required**: Si el catálogo no existe o el Upstream devuelve un `404 Not Found` o `402`. Replica la excepción `CatalogNotFoundException` del legado.
*   **400 Bad Request**: Si el Upstream devuelve **cualquier otro error HTTP** (ej. 500, 502, 504), el legado lo enmascaraba bajo una `CatalogException` con HTTP 400 y mensaje `Error en la obtención de catálogos`. Este comportamiento se replica exactamente en Go.

## 6. Evidencia de Paridad
La prueba formal `test_paridad_extendido.ps1` fue superada exitosamente en 13 de 13 reglas (100% de paridad con Java), validando:
1.  **Destino (Estándar):** Casos base OK.
2.  **Seguros PascalCase / Minúsculas:** Comportamiento idéntico (Parseo distinto al catálogo estándar).
3.  **TiposDocumentos / Comunas:** Validaciones correctas y fallbacks genéricos de JSON array.
4.  **Timeouts Upstream:** Emulación idéntica devolviendo `200 Sin resultados`.
5.  **Token Inválido:** Emulación idéntica devolviendo error JSON 401 legado.
6.  **Catálogos inexistentes:** Fallback a HTTP 402 al igual que la respuesta original del Java Legacy.
