# Definiciones del Dominio y Contratos Legados (JSON)

Este documento fue extraído del repositorio `mortgage-loan-catalogs-go-clean` y contiene la información necesaria para replicar la paridad estricta en el nuevo microservicio Hexagonal.

## 1. Entidades del Dominio

Las entidades de dominio puras (sin tags JSON) que deben ubicarse en `internal/core/domain/`:

```go
package domain

// CatalogItem represents a generic catalog item.
type CatalogItem struct {
	Code        string
	Description string
}

type ComunaCatalogItem struct {
	Code        string
	Description string
	RegionID    int
}

type TipoDocumentoCatalogItem struct {
	Code        int
	Description string
	GroupID     int
}

type NoResultsCatalogItem struct {
	CodRespuesta int
	Mensaje      string
	Excepcion    string
}

// InsuranceCatalogItem represents an insurance catalog item.
type InsuranceCatalogItem struct {
	InsuranceIdentifier  string
	Description          string
	Policy               string
	CompanyName          string
	CompanyCode          int
	InsuranceTypeCode    int
	PolicyCorrelative    int
	Rate                 float64
	Factor               float64
	IndividualPolicyFlag int
	PerQuotaValueFlag    int
	ExternalPolicyFlag   int
}

// GetCatalogRequest represents the pure domain request to get a catalog.
type GetCatalogRequest struct {
	CatalogName   string
	Channel       string
	Commerce      string
	TransactionID string
}
```

## 2. Contratos Legados DTO (JSON Paridad Estricta)

Estos son los DTOs que se deben mapear en `internal/app/` (o los handlers) para mantener la paridad exacta con el sistema legado (Java):

```go
// catalogItemDTO represents the JSON payload for standard catalogs.
type catalogItemDTO struct {
	CodigoAdm   string `json:"codigo_adm"`
	Descripcion string `json:"descripcion"`
}

// comunaCatalogItemDTO represents the JSON payload for comunas.
type comunaCatalogItemDTO struct {
	CodigoAdm   string `json:"codigo_adm"`
	Descripcion string `json:"descripcion"`
	RegionID    int    `json:"region_id"`
}

// tipoDocumentoCatalogItemDTO represents the JSON payload for tipos documentos.
type tipoDocumentoCatalogItemDTO struct {
	CodigoAdm   int    `json:"codigo_adm"`
	Descripcion string `json:"descripcion"`
	GrupoID     int    `json:"grupo_id"`
}

// noResultsCatalogItemDTO represents the legacy empty/network error fallback
type noResultsCatalogItemDTO struct {
	CodRespuesta int    `json:"codRespuesta"`
	Mensaje      string `json:"Mensaje"`
	Excepcion    string `json:"Excepcion"`
}

// insuranceCatalogItemDTO represents the JSON payload for insurance catalogs.
type insuranceCatalogItemDTO struct {
	IdentificadorSeguro       string      `json:"IdentificadorSeguro"`
	Descripcion               string      `json:"Descripcion"`
	Poliza                    string      `json:"Poliza"`
	NombreCompania            string      `json:"NombreCompania"`
	CodigoCompania            int         `json:"CodigoCompania"`
	CodigoTipoSeguro          int         `json:"CodigoTipoSeguro"`
	CorrelativoPoliza         int         `json:"CorrelativoPoliza"`
	Tasa                      json.Number `json:"Tasa"`
	Factor                    json.Number `json:"Factor"`
	IndicadorPolizaIndividual int         `json:"IndicadorPolizaIndividual"`
	PorValorCuota             int         `json:"PorValorCuota"`
	IndicadorPolizaExterna    int         `json:"IndicadorPolizaExterna"`
}
```

### Regla Especial de Formateo de Float (Notación Científica Jackson)
Para mantener compatibilidad binaria con cómo Java de-serializa los `Double`, se requiere esta función auxiliar en el mapeo de Seguros:
```go
func formatFloatLikeJackson(f float64) string {
	s := strconv.FormatFloat(f, 'g', -1, 64)
	if f != 0 && f > -1e-3 && f < 1e-3 {
		s = strconv.FormatFloat(f, 'E', -1, 64)
		s = strings.Replace(s, "E-0", "E-", 1)
		s = strings.Replace(s, "E+0", "E", 1)
	}
	return s
}
```

## 3. Respuestas de Error Legadas (Paridad de Códigos HTTP)

Las respuestas de error exactas que se deben devolver para mantener la retrocompatibilidad:

```go
// errorResponseDTO ensures exact JSON key ordering as Java
type errorResponseDTO struct {
	Detail       *string       `json:"detail"`
	Code         *string       `json:"code"`
	Messages     []interface{} `json:"messages"`
	ErrorsDetail *string       `json:"errors_detail"`
}
```

### Escenarios de Error:

1.  **Error 401 Unauthorized (Tokens/Headers faltantes):**
    ```json
    {
      "detail": "Given token not valid for any token type",
      "code": "token_not_valid",
      "messages": [
        {
          "token_class": "AccessToken",
          "token_type": "access",
          "message": "Token is invalid"
        }
      ],
      "errors_detail": null
    }
    ```

2.  **Error 402 Payment Required (Catálogo No Encontrado):**
    ```json
    {
      "detail": null,
      "code": null,
      "messages": null,
      "errors_detail": "Catálogo no encontrado"
    }
    ```

3.  **Error 400 Bad Request (Excepciones Generales):**
    ```json
    {
      "detail": null,
      "code": null,
      "messages": null,
      "errors_detail": "Mensaje de error interno"
    }
    ```

4.  **Error 500 Internal Server Error (Fallo no manejado / Legacy):**
    ```json
    {
      "mensaje": "Error interno del servidor",
      "statusCode": 500,
      "timestamp": "2026-09-21T19:15:33.985099286"
    }
    ```

## 4. Endpoint Legacy

-   **Ruta:** `POST /v1/bfcl/mortgage-loan/catalogs/{catalog}`
-   **Headers Obligatorios:** `X-Channel`, `X-Commerce`, `X-Transaction-ID`
