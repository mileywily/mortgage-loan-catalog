# Casos de Uso - Payloads de Postman

Este documento contiene los payloads EXACTOS (request y response) ejecutados contra el servidor Java (legado, puerto 8080) y el servidor Go (nuevo, puerto 8082).

### CU-XXX: Destino

#### Request
- **URL:** POST /v1/bfcl/mortgage-loan/catalogs/Destino
- **Headers:** (Ninguno obligatorio)
- **Body:** (vacío)

#### Response — Java Legado (8080)
**HTTP Status:** `
`json
[{"codigo_adm":"1","descripcion":"Vivienda Principal"}]
` 

#### Response — Go Hexagonal (8082)
**HTTP Status:** `
`json
[{"codigo_adm":"1","descripcion":"Vivienda Principal"}]
` 

#### Paridad
**✅ IDÉNTICO**

---

### CU-XXX: SegurosIncendio

#### Request
- **URL:** POST /v1/bfcl/mortgage-loan/catalogs/SegurosIncendio
- **Headers:** (Ninguno obligatorio)
- **Body:** (vacío)

#### Response — Java Legado (8080)
**HTTP Status:** `
`json
[{"IdentificadorSeguro":"PARITY-000","Descripcion":"SEGURO PARITY","Poliza":"P-000","NombreCompania":"PARITY CIA","CodigoCompania":99,"CodigoTipoSeguro":1,"CorrelativoPoliza":1,"Tasa":0.1,"Factor":2.255E-4,"IndicadorPolizaIndividual":0,"PorValorCuota":0,"IndicadorPolizaExterna":0}]
` 

#### Response — Go Hexagonal (8082)
**HTTP Status:** `
`json
[{"IdentificadorSeguro":"PARITY-000","Descripcion":"SEGURO PARITY","Poliza":"P-000","NombreCompania":"PARITY CIA","CodigoCompania":99,"CodigoTipoSeguro":1,"CorrelativoPoliza":1,"Tasa":0.1,"Factor":2.255E-4,"IndicadorPolizaIndividual":0,"PorValorCuota":0,"IndicadorPolizaExterna":0}]
` 

#### Paridad
**✅ IDÉNTICO**

---

### CU-XXX: segurosincendio

#### Request
- **URL:** POST /v1/bfcl/mortgage-loan/catalogs/segurosincendio
- **Headers:** (Ninguno obligatorio)
- **Body:** (vacío)

#### Response — Java Legado (8080)
**HTTP Status:** `
`json
[{"IdentificadorSeguro":"PARITY-000","Descripcion":"SEGURO PARITY","Poliza":"P-000","NombreCompania":"PARITY CIA","CodigoCompania":99,"CodigoTipoSeguro":1,"CorrelativoPoliza":1,"Tasa":0.1,"Factor":2.255E-4,"IndicadorPolizaIndividual":0,"PorValorCuota":0,"IndicadorPolizaExterna":0}]
` 

#### Response — Go Hexagonal (8082)
**HTTP Status:** `
`json
[{"IdentificadorSeguro":"PARITY-000","Descripcion":"SEGURO PARITY","Poliza":"P-000","NombreCompania":"PARITY CIA","CodigoCompania":99,"CodigoTipoSeguro":1,"CorrelativoPoliza":1,"Tasa":0.1,"Factor":2.255E-4,"IndicadorPolizaIndividual":0,"PorValorCuota":0,"IndicadorPolizaExterna":0}]
` 

#### Paridad
**✅ IDÉNTICO**

---

### CU-XXX: TiposDocumentos

#### Request
- **URL:** POST /v1/bfcl/mortgage-loan/catalogs/TiposDocumentos
- **Headers:** (Ninguno obligatorio)
- **Body:** (vacío)

#### Response — Java Legado (8080)
**HTTP Status:** `
`json
[{"codigo_adm":1,"descripcion":"Escritura Parity Test","grupo_id":10}]
` 

#### Response — Go Hexagonal (8082)
**HTTP Status:** `
`json
[{"codigo_adm":1,"descripcion":"Escritura Parity Test","grupo_id":10}]
` 

#### Paridad
**✅ IDÉNTICO**

---

### CU-XXX: tiposdocumentos

#### Request
- **URL:** POST /v1/bfcl/mortgage-loan/catalogs/tiposdocumentos
- **Headers:** (Ninguno obligatorio)
- **Body:** (vacío)

#### Response — Java Legado (8080)
**HTTP Status:** `
`json
[{"codigo_adm":"1","descripcion":"Vivienda Principal"}]
` 

#### Response — Go Hexagonal (8082)
**HTTP Status:** `
`json
[{"codigo_adm":"1","descripcion":"Vivienda Principal"}]
` 

#### Paridad
**✅ IDÉNTICO**

---

### CU-XXX: Comunas

#### Request
- **URL:** POST /v1/bfcl/mortgage-loan/catalogs/Comunas
- **Headers:** (Ninguno obligatorio)
- **Body:** (vacío)

#### Response — Java Legado (8080)
**HTTP Status:** `
`json
[{"codigo_adm":"13101","descripcion":"Santiago Centro","region_id":13}]
` 

#### Response — Go Hexagonal (8082)
**HTTP Status:** `
`json
[{"codigo_adm":"13101","descripcion":"Santiago Centro","region_id":13}]
` 

#### Paridad
**✅ IDÉNTICO**

---

### CU-XXX: comunas

#### Request
- **URL:** POST /v1/bfcl/mortgage-loan/catalogs/comunas
- **Headers:** (Ninguno obligatorio)
- **Body:** (vacío)

#### Response — Java Legado (8080)
**HTTP Status:** `
`json
[{"codigo_adm":"1","descripcion":"Vivienda Principal"}]
` 

#### Response — Go Hexagonal (8082)
**HTTP Status:** `
`json
[{"codigo_adm":"1","descripcion":"Vivienda Principal"}]
` 

#### Paridad
**✅ IDÉNTICO**

---

### CU-XXX: Regiones

#### Request
- **URL:** POST /v1/bfcl/mortgage-loan/catalogs/Regiones
- **Headers:** (Ninguno obligatorio)
- **Body:** (vacío)

#### Response — Java Legado (8080)
**HTTP Status:** `
`json
[{"codigo_adm":"13","descripcion":"Metropolitana de Santiago"}]
` 

#### Response — Go Hexagonal (8082)
**HTTP Status:** `
`json
[{"codigo_adm":"13","descripcion":"Metropolitana de Santiago"}]
` 

#### Paridad
**✅ IDÉNTICO**

---

### CU-XXX: SinResultados

#### Request
- **URL:** POST /v1/bfcl/mortgage-loan/catalogs/SinResultados
- **Headers:** (Ninguno obligatorio)
- **Body:** (vacío)

#### Response — Java Legado (8080)
**HTTP Status:** `
`json
[{"codRespuesta":3,"Mensaje":"Sin resultados.","Excepcion":"Ninguna"}]
` 

#### Response — Go Hexagonal (8082)
**HTTP Status:** `
`json
[{"codRespuesta":3,"Mensaje":"Sin resultados.","Excepcion":"Ninguna"}]
` 

#### Paridad
**✅ IDÉNTICO**

---

### CU-XXX: Vacio

#### Request
- **URL:** POST /v1/bfcl/mortgage-loan/catalogs/Vacio
- **Headers:** (Ninguno obligatorio)
- **Body:** (vacío)

#### Response — Java Legado (8080)
**HTTP Status:** `
`json
[{"codRespuesta":3,"Mensaje":"Sin resultados.","Excepcion":"Ninguna"}]
` 

#### Response — Go Hexagonal (8082)
**HTTP Status:** `
`json
[{"codRespuesta":3,"Mensaje":"Sin resultados.","Excepcion":"Ninguna"}]
` 

#### Paridad
**✅ IDÉNTICO**

---

### CU-XXX: Inexistente

#### Request
- **URL:** POST /v1/bfcl/mortgage-loan/catalogs/Inexistente
- **Headers:** (Ninguno obligatorio)
- **Body:** (vacío)

#### Response — Java Legado (8080)
**HTTP Status:** `
`json
{"detail":null,"code":null,"messages":null,"errors_detail":"Cat├ílogo no encontrado"}
` 

#### Response — Go Hexagonal (8082)
**HTTP Status:** `
`json
{"detail":null,"code":null,"messages":null,"errors_detail":"Cat├ílogo no encontrado"}
` 

#### Paridad
**✅ IDÉNTICO**

---

### CU-XXX: TokenInvalido

#### Request
- **URL:** POST /v1/bfcl/mortgage-loan/catalogs/TokenInvalido
- **Headers:** (Ninguno obligatorio)
- **Body:** (vacío)

#### Response — Java Legado (8080)
**HTTP Status:** `
`json
{"detail":"Given token not valid for any token type","code":"token_not_valid","messages":[{"token_class":"AccessToken","token_type":"access","message":"Token is invalid"}],"errors_detail":null}
` 

#### Response — Go Hexagonal (8082)
**HTTP Status:** `
`json
{"detail":"Given token not valid for any token type","code":"token_not_valid","messages":[{"token_class":"AccessToken","token_type":"access","message":"Token is invalid"}],"errors_detail":null}
` 

#### Paridad
**✅ IDÉNTICO**

---

### CU-XXX: Timeout

#### Request
- **URL:** POST /v1/bfcl/mortgage-loan/catalogs/Timeout
- **Headers:** (Ninguno obligatorio)
- **Body:** (vacío)

#### Response — Java Legado (8080)
**HTTP Status:** `
`json
[{"codRespuesta":3,"Mensaje":"Sin resultados.","Excepcion":"Ninguna"}]
` 

#### Response — Go Hexagonal (8082)
**HTTP Status:** `
`json
[{"codRespuesta":3,"Mensaje":"Sin resultados.","Excepcion":"Ninguna"}]
` 

#### Paridad
**✅ IDÉNTICO**

---

### CU-XXX: SeguroItem

#### Request
- **URL:** POST /v1/bfcl/mortgage-loan/catalogs/SeguroItem
- **Headers:** (Ninguno obligatorio)
- **Body:** (vacío)

#### Response — Java Legado (8080)
**HTTP Status:** `
`json
[{"codigo_adm":"1","descripcion":"Vivienda Principal"}]
` 

#### Response — Go Hexagonal (8082)
**HTTP Status:** `
`json
[{"codigo_adm":"1","descripcion":"Vivienda Principal"}]
` 

#### Paridad
**✅ IDÉNTICO**

---

## Resumen de Paridad

| Catálogo | Java Status | Go Status | Paridad |
|---|---|---|---|
| Destino | 200 | 200 | ✅ IDÉNTICO |
| SegurosIncendio | 200 | 200 | ✅ IDÉNTICO |
| segurosincendio | 200 | 200 | ✅ IDÉNTICO |
| TiposDocumentos | 200 | 200 | ✅ IDÉNTICO |
| tiposdocumentos | 200 | 200 | ✅ IDÉNTICO |
| Comunas | 200 | 200 | ✅ IDÉNTICO |
| comunas | 200 | 200 | ✅ IDÉNTICO |
| Regiones | 200 | 200 | ✅ IDÉNTICO |
| SinResultados | 200 | 200 | ✅ IDÉNTICO |
| Vacio | 200 | 200 | ✅ IDÉNTICO |
| Inexistente | 402 | 402 | ✅ IDÉNTICO |
| TokenInvalido | 401 | 401 | ✅ IDÉNTICO |
| Timeout | 200 | 200 | ✅ IDÉNTICO |
| SeguroItem | 200 | 200 | ✅ IDÉNTICO |

