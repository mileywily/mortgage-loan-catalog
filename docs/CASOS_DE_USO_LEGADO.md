# Casos de Uso - Sistema Legado Java mortgage-loan-catalogs

> **Documento generado por auditoría automatizada del código fuente Java.**  
> Sirve como especificación oficial de comportamiento para pruebas en Postman, tests unitarios, de integración y E2E.  
> Todo comportamiento descrito aquí debe replicarse **exactamente** en el sistema Go Hexagonal.

---

## 1. Descripción General

El sistema `mortgage-loan-catalogs` es un **middleware/pasarela** que consulta catálogos relacionados con créditos hipotecarios expuestos por el proveedor externo **Finnflow**. Sus responsabilidades son:

- Enrutar la petición del cliente al upstream correcto según el nombre del catálogo.
- Transformar la respuesta del upstream al DTO de salida esperado por el cliente.
- Absorber y estandarizar errores del upstream (timeout, 4xx, 5xx) hacia respuestas predecibles.
- No valida ni rechaza peticiones por ausencia de headers de negocio — esto es responsabilidad del upstream.

---

## 2. Configuración y Variables de Entorno

| Variable de Entorno | Uso | Valor local (dev) |
|---------------------|-----|-------------------|
| `FINNFLOW_URL` | Base URL del upstream Finnflow | `http://localhost:9090` |
| `FINNFLOW_KEY` | Usuario para Basic Auth | configurado en entorno |
| `FINNFLOW_SECRET` | Contraseña para Basic Auth | configurado en entorno |
| `PORT` | Puerto del servidor Go | `8082` |
| `DEFAULT_BACKEND` | Repositorio activo por defecto (`real`, `dummy`, `java`) | `real` |
| `JAVA_LEGACY_URL` | URL del legado Java (solo para feature toggle Swagger) | `http://localhost:8080` |

> **Importante:** La URL final al upstream es: `/api/catalogo_detail/{catalog}`

---

## 3. Arquitectura de la Llamada

### 3.1 Endpoint Expuesto al Cliente

| Atributo | Valor |
|----------|-------|
| **Método HTTP** | `POST` |
| **Ruta** | `/v1/bfcl/mortgage-loan/catalogs/{catalog}` |
| **Path Parameter** | `catalog` (String, requerido) |
| **Headers de entrada** | Ninguno requerido explícitamente |
| **Content-Type respuesta** | `application/json` |

### 3.2 Llamada Upstream (hacia Finnflow)

| Atributo | Valor |
|----------|-------|
| **URL exacta** | `{FINNFLOW_URL}/api/catalogo_detail/{catalog}` |
| **Método HTTP** | `GET` |
| **Header Authorization** | `Basic {Base64(FINNFLOW_KEY:FINNFLOW_SECRET)}` |
| **Header Content-Type** | `application/json` |
| **Timeout (dev/qa)** | 30,000 ms |
| **Timeout (prod)** | 60,000 ms |
| **SSL** | verify-hostname: false · trust-all-certs: true dev/qa |

---

## 4. Lógica de Routing por Tipo de Catálogo

El `CatalogService` evalúa el parámetro `{catalog}` en **orden de prioridad**:

```
1. catalog.toLowerCase().startsWith("seguros")  → Rama Seguros     (case-INSENSITIVE)
2. "TiposDocumentos".equals(catalog)            → Rama TiposDocumentos   (case-SENSITIVE exacto)
3. "Comunas".equals(catalog)                    → Rama Comunas           (case-SENSITIVE exacto)
4. cualquier otro valor                          → Rama Genérica
```

---

## 5. Casos de Uso

---

### CU-001 — Catálogos de Seguros

**Condición de activación:** `catalog.toLowerCase().startsWith("seguros")`  
**Ejemplos:** `SegurosIncendio`, `segurosincendio`, `SegurosDesgravamen`, `SegurosCesantia`

#### Input
| Campo | Tipo | Requerido | Ejemplo |
|-------|------|-----------|---------|
| `catalog` (path) | String | Si | `SegurosIncendio` |
| Headers de negocio | - | No | Opcionales |

#### Llamada al Upstream
```
GET /api/catalogo_detail/SegurosIncendio
Authorization: Basic {Base64(KEY:SECRET)}
Content-Type: application/json
```

#### Respuesta Exitosa (200 OK)
```json
[{"IdentificadorSeguro":"PARITY-000","Descripcion":"SEGURO PARITY","Poliza":"P-000","NombreCompania":"PARITY CIA","CodigoCompania":99,"CodigoTipoSeguro":1,"CorrelativoPoliza":1,"Tasa":0.1,"Factor":2.255E-4,"IndicadorPolizaIndividual":0,"PorValorCuota":0,"IndicadorPolizaExterna":0}]
```

**⚠️ Campo Factor:** Jackson serializa `double` pequeños en notación científica. `2.255E-4` es exactamente `0.0002255`. Ambos sistemas deben emitir `2.255E-4`.

#### Escenarios
| # | Upstream responde | HTTP Cliente | Body Cliente |
|---|-------------------|-------------|--------------|
| 1 | 200 OK + datos | 200 OK | Array SeguroItem[] PascalCase |
| 2 | 200 OK + [] vacío | 200 OK | `[{"codRespuesta":3,"Mensaje":"Sin resultados.","Excepcion":"Ninguna"}]` |
| 3 | 200 OK + codRespuesta | 200 OK | `[{"codRespuesta":3,"Mensaje":"Sin resultados.","Excepcion":"Ninguna"}]` |
| 4 | Timeout / IOException | 200 OK | `[{"codRespuesta":3,"Mensaje":"Sin resultados.","Excepcion":"Ninguna"}]` |
| 5 | 401 Unauthorized | 401 | Body TokenNotValid (ver §6.1) |
| 6 | 404 o 402 | 402 | Body CatalogNotFound (ver §6.2) |
| 7 | 400 / 500 / otro | 400 | Body CatalogException (ver §6.3) |

---

### CU-002 — Catálogo TiposDocumentos

**Condición de activación:** `"TiposDocumentos".equals(catalog)` — EXACTO, case-sensitive.  
**⚠️** `tiposdocumentos` o `Tiposdocumentos` van a Rama Genérica (CU-004).

#### Respuesta Exitosa (200 OK)
```json
[{"codigo_adm":1,"descripcion":"Escritura Parity Test","grupo_id":10}]
```
**Nota:** `codigo_adm` es Integer (no String) y se incluye `grupo_id`.

#### Escenarios
| # | Upstream responde | HTTP Cliente | Body Cliente |
|---|-------------------|-------------|--------------|
| 1 | 200 OK + datos | 200 OK | Array TipoDocumentoItem[] con grupo_id |
| 2 | 200 OK + [] | 200 OK | `[{"codRespuesta":3,"Mensaje":"Sin resultados.","Excepcion":"Ninguna"}]` |
| 3 | Timeout | 200 OK | `[{"codRespuesta":3,"Mensaje":"Sin resultados.","Excepcion":"Ninguna"}]` |
| 4 | 401 | 401 | Ver §6.1 |
| 5 | 404/402 | 402 | Ver §6.2 |
| 6 | 400/500 | 400 | Ver §6.3 |

---

### CU-003 — Catálogo Comunas

**Condición de activación:** `"Comunas".equals(catalog)` — EXACTO, case-sensitive.  
**⚠️** `comunas` va a Rama Genérica (CU-004).

#### Respuesta Exitosa (200 OK)
```json
[{"codigo_adm":"13101","descripcion":"Santiago Centro","region_id":13}]
```
**Nota:** Incluye `region_id` Integer adicional. Pre-verifica raw body buscando `codRespuesta`.

#### Escenarios
| # | Upstream responde | HTTP Cliente | Body Cliente |
|---|-------------------|-------------|--------------|
| 1 | 200 OK + datos | 200 OK | Array ComunaItem[] con region_id |
| 2 | 200 OK + [] | 200 OK | `[{"codRespuesta":3,"Mensaje":"Sin resultados.","Excepcion":"Ninguna"}]` |
| 3 | 200 OK + codRespuesta | 200 OK | Payload original inalterado |
| 4 | Timeout | 200 OK | `[{"codRespuesta":3,"Mensaje":"Sin resultados.","Excepcion":"Ninguna"}]` |
| 5 | 401 | 401 | Ver §6.1 |
| 6 | 404/402 | 402 | Ver §6.2 |
| 7 | 400/500 | 400 | Ver §6.3 |

---

### CU-004 — Catálogo Genérico

**Condición de activación:** Todo lo que no capturan CU-001, CU-002, CU-003.  
**Ejemplos:** `Destino`, `Regiones`, `comunas`, `tiposdocumentos`, `SeguroItem`, `Objetivo`, `Producto`, etc.

#### Respuesta Exitosa (200 OK)
```json
[{"codigo_adm":"1","descripcion":"Vivienda Principal"}]
```
**Nota:** `codigo_adm` es **String** (a diferencia de TiposDocumentos donde es Integer).

#### Escenarios
| # | Upstream responde | HTTP Cliente | Body Cliente |
|---|-------------------|-------------|--------------|
| 1 | 200 OK + datos | 200 OK | Array CatalogItem[] con codigo_adm String |
| 2 | 200 OK + [] vacío | 200 OK | `[{"codRespuesta":3,"Mensaje":"Sin resultados.","Excepcion":"Ninguna"}]` |
| 3 | 200 OK + codRespuesta | 200 OK | `[{"codRespuesta":3,"Mensaje":"Sin resultados.","Excepcion":"Ninguna"}]` |
| 4 | Timeout / IOException | 200 OK | `[{"codRespuesta":3,"Mensaje":"Sin resultados.","Excepcion":"Ninguna"}]` |
| 5 | 401 | 401 | Ver §6.1 |
| 6 | 404/402 | 402 | Ver §6.2 |
| 7 | 400/500 | 400 | Ver §6.3 |

---

## 6. Matriz de Errores — Respuestas Exactas al Cliente

### 6.1 — 401 Unauthorized (Token Inválido)

```json
{"detail":"Given token not valid for any token type","code":"token_not_valid","messages":[{"token_class":"AccessToken","token_type":"access","message":"Token is invalid"}],"errors_detail":null}
```

### 6.2 — 402 Payment Required (Catálogo No Encontrado)

Upstream responde `404 Not Found` O `402 Payment Required`:
```json
{"errors_detail":"Catálogo no encontrado"}
```

### 6.3 — 400 Bad Request (Error Genérico del Proveedor)

Upstream responde `400`, `500`, `503` u otro código:
```json
{"errors_detail":"Error en la obtención de catálogos"}
```

---

## 7. Simulador Finnflow — Catálogos Especiales (Puerto 9090)

| Catálogo (path) | Respuesta Simulador | HTTP Cliente | Observación |
|-----------------|---------------------|-------------|-------------|
| SegurosIncendio | 200 + SeguroItem[] | 200 | Factor: 2.255E-4 |
| SegurosDesgravamen | 200 + SeguroItem[] | 200 | |
| SegurosCesantia | 200 + SeguroItem[] | 200 | |
| TiposDocumentos | 200 + TipoDocumentoItem[] | 200 | codigo_adm: Integer |
| Comunas | 200 + ComunaItem[] | 200 | region_id presente |
| Destino / Regiones / etc. | 200 + CatalogItem[] | 200 | codigo_adm: String |
| SinResultados | 200 + [{codRespuesta:3,...}] | 200 | Payload pass-through |
| Vacio | 200 + [] | 200 | Transformado a Sin resultados |
| Inexistente | 404 Not Found | 402 | CatalogNotFoundException |
| TokenInvalido | 401 Unauthorized | 401 | TokenNotValidException |
| Timeout | sleep(5s) → timeout | 200 | Transformado a Sin resultados |
| cualquier otro | 200 + [{codigo_adm:"1",...}] | 200 | Fallback genérico |

---

## 8. Colección Postman — 14 Escenarios de Prueba

### Environment Variables
```
GO_URL     = http://localhost:8082
JAVA_URL   = http://localhost:8080
BASE_PATH  = /v1/bfcl/mortgage-loan/catalogs
```

| # | Nombre | Request | HTTP Esperado | Body Esperado |
|---|--------|---------|--------------|---------------|
| 1 | Destino (estándar) | POST {{GO_URL}}{{BASE_PATH}}/Destino | 200 | CatalogItem[] genérico |
| 2 | Seguros PascalCase | POST {{GO_URL}}{{BASE_PATH}}/SegurosIncendio | 200 | SeguroItem[] con Factor:2.255E-4 |
| 3 | Seguros minúsculas | POST {{GO_URL}}{{BASE_PATH}}/segurosincendio | 200 | SeguroItem[] igual al anterior |
| 4 | TiposDocumentos exacto | POST {{GO_URL}}{{BASE_PATH}}/TiposDocumentos | 200 | Array con grupo_id |
| 5 | tiposdocumentos (genérico) | POST {{GO_URL}}{{BASE_PATH}}/tiposdocumentos | 200 | Array sin grupo_id |
| 6 | Comunas exacto | POST {{GO_URL}}{{BASE_PATH}}/Comunas | 200 | Array con region_id |
| 7 | comunas (genérico) | POST {{GO_URL}}{{BASE_PATH}}/comunas | 200 | Array sin region_id |
| 8 | Regiones (genérico) | POST {{GO_URL}}{{BASE_PATH}}/Regiones | 200 | CatalogItem[] genérico |
| 9 | SinResultados | POST {{GO_URL}}{{BASE_PATH}}/SinResultados | 200 | [{codRespuesta:3,...}] |
| 10 | Vacío | POST {{GO_URL}}{{BASE_PATH}}/Vacio | 200 | [{codRespuesta:3,...}] |
| 11 | Inexistente (402) | POST {{GO_URL}}{{BASE_PATH}}/Inexistente | 402 | {errors_detail:"Catálogo no encontrado"} |
| 12 | Token Inválido (401) | POST {{GO_URL}}{{BASE_PATH}}/TokenInvalido | 401 | {detail, code, messages} |
| 13 | Timeout (red) | POST {{GO_URL}}{{BASE_PATH}}/Timeout | 200 | [{codRespuesta:3,...}] |
| 14 | SeguroItem singular (genérico) | POST {{GO_URL}}{{BASE_PATH}}/SeguroItem | 200 | CatalogItem[] genérico |

---

## 9. Diagrama de Decisión Completo

```
POST /v1/bfcl/mortgage-loan/catalogs/{catalog}
         │
         ▼
  GET /api/catalogo_detail/{catalog}
  Authorization: Basic {Base64(KEY:SECRET)}
         │
         ├─── IOException / Timeout ──────────────► 200 [Sin resultados]
         ├─── HTTP 401 ───────────────────────────► 401 TokenNotValid
         ├─── HTTP 404 o 402 ─────────────────────► 402 CatalogNotFound  
         ├─── HTTP 400 / 500 / otro ──────────────► 400 CatalogException
         └─── HTTP 200 OK
                   │
                   ├── catalog.toLowerCase().startsWith("seguros")
                   │         └─► Parsear SeguroItem[] ───────────────────► 200 SeguroItem[]
                   │
                   ├── "TiposDocumentos".equals(catalog)  [EXACTO]
                   │         └─► Parsear TipoDocumentoItem[] ─────────────► 200 TipoDocumentoItem[]
                   │
                   ├── "Comunas".equals(catalog)  [EXACTO]
                   │         └─► Pre-check raw → Parsear ComunaItem[] ────► 200 ComunaItem[]
                   │
                   └── cualquier otro (genérico)
                             ├── body contiene codRespuesta ─────────────► 200 [Sin resultados]
                             ├── array vacío [] ──────────────────────────► 200 [Sin resultados]
                             └── array con datos ─────────────────────────► 200 CatalogItem[]
```
