# Guía para la Integración de Nuevos Catálogos

Este documento explica cómo el microservicio `mortgage-loan-catalogs` maneja los catálogos y cuál es el procedimiento técnico (paso a paso) para incorporar nuevos catálogos, manteniendo intacta la Arquitectura Hexagonal y la seguridad de tipos de Go.

---

## 1. Comportamiento por Defecto (Agnóstico)

El diseño actual del microservicio está preparado para escalar sin necesidad de modificar código ante catálogos simples.

Si se solicita un nuevo catálogo (ej. `Sucursales` o `Bancos`) y el proveedor (Finnflow) responde con una estructura genérica estándar (una lista de objetos con `codigo_adm` y `descripcion`), **no es necesario escribir código**. 

El bloque de _fallback_ en `rest_repository.go` atrapará el llamado, lo parseará de forma agnóstica a la estructura genérica `domain.CatalogItem` y lo devolverá al cliente automáticamente.

---

## 2. ¿Cuándo es necesario modificar el código?

Deberás intervenir el código del microservicio únicamente cuando un nuevo catálogo (ej. `TasasEspeciales`) cumpla al menos una de estas condiciones:
- **Respuesta no estándar:** El JSON del proveedor contiene campos obligatorios adicionales, estructuras anidadas, o llaves distintas a `codigo_adm` y `descripcion`.
- **Petición no estándar:** El proveedor externo requiere que se le envíen nuevos Query Parameters, Headers especiales o un Body específico para este catálogo.

En estos casos, no se recomienda forzar un comportamiento 100% agnóstico (usando `map[string]interface{}`), ya que esto destruye el tipado fuerte de Go y dificulta el mantenimiento. En su lugar, aplicamos la extensión hexagonal.

---

## 3. Guía Paso a Paso: Integrando un Catálogo Complejo

Sigue este orden estricto (de adentro hacia afuera) para implementar un catálogo no estándar:

### PASO 1: Capa de Dominio (`internal/core/domain/catalog.go`)
El dominio dicta las reglas. Aquí defines las estructuras de datos puras que representan el negocio.

1. Crea el `struct` para el nuevo ítem del catálogo:
```go
// Definición exacta y limpia para el negocio
type TasasEspecialesItem struct {
    TasaBase    float64  `json:"tasa_base"`
    Segmento    string   `json:"segmento"`
    Condiciones []string `json:"condiciones"`
}
```
2. *(Opcional)* Si la solicitud requiere nuevos parámetros (ej. `TipoCliente`), agrégalo al request de dominio:
```go
type GetCatalogRequest struct {
    CatalogName   string
    Channel       string
    Commerce      string
    TransactionID string
    TipoCliente   string // <--- Nuevo parámetro
}
```

### PASO 2: Capa de Repositorio (`internal/infra/repository/rest_repository.go`)
El repositorio actúa como "escudo". Su trabajo es lidiar con el desorden de las APIs externas y traducirlo al modelo puro del Dominio que creaste en el Paso 1.

Agrega un bloque `if` (antes del _fallback_ genérico final) para aislar la lógica de tu nuevo catálogo:

```go
if req.CatalogName == "TasasEspeciales" {
    // 1. Estructura temporal y sucia (Upstream)
    var upstreamResponse []struct {
        ValorTasa float64 `json:"valor_tasa_upstream"`
        Seg       string  `json:"seg_upstream"`
        Cond      string  `json:"condiciones_csv"`
    }
    
    if err := json.NewDecoder(resp.Body).Decode(&upstreamResponse); err != nil {
        return nil, err
    }

    // 2. Mapeo hacia el Dominio limpio
    domainItems := make([]domain.TasasEspecialesItem, len(upstreamResponse))
    for i, item := range upstreamResponse {
        domainItems[i] = domain.TasasEspecialesItem{
            TasaBase:    item.ValorTasa,
            Segmento:    item.Seg,
            Condiciones: strings.Split(item.Cond, ","),
        }
    }
    
    return domainItems, nil
}
```

### PASO 3: Capa de Handler / API (`internal/app/handler/catalog_handler.go`)
Si modificaste el `GetCatalogRequest` en el Paso 1 porque el cliente debe enviarte datos nuevos, este es el lugar para extraerlos de la petición HTTP (usando Gin).

```go
func (h *CatalogHandler) GetCatalog(c *gin.Context) {
    catalogName := c.Param("catalog")

    // Extraer parámetros (Headers, Query Params, Body)
    req := domain.GetCatalogRequest{
        CatalogName:   catalogName,
        Channel:       c.GetHeader("X-Channel"),
        TipoCliente:   c.Query("tipo_cliente"), // <--- Extracción de nuevo parámetro
    }
    
    // ... ejecución del servicio (sin cambios)
}
```

---

## 4. Resumen de Buenas Prácticas

- **Capa de Servicios (`catalog_service.go`):** Nota que en esta guía **nunca tocamos la capa de servicio**. Esto es porque recuperar un catálogo no suele requerir reglas de negocio complejas. Solo deberás tocar el servicio si el catálogo nuevo requiere validaciones lógicas cruzadas (ej. "Lanzar error si el canal es WEB pero pide el catálogo interno").
- **Evita `interface{}`:** Mantén el código tipado. Esto garantiza que el compilador te avise si hay errores antes de llegar a Producción.
- **Testing:** Recuerda agregar tu nuevo catálogo (ej. "TasasEspeciales") al archivo `rest_repository_test.go` para validar que el mapeo de estructuras JSON a Domain funciona correctamente.
