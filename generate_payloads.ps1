$catalogs = @(
    "Destino", "SegurosIncendio", "segurosincendio", "TiposDocumentos", "tiposdocumentos",
    "Comunas", "comunas", "Regiones", "SinResultados", "Vacio", "Inexistente",
    "TokenInvalido", "Timeout", "SeguroItem"
)

$md = "# Casos de Uso - Payloads de Postman

"
$md += "Este documento contiene los payloads EXACTOS (request y response) ejecutados contra el servidor Java (legado, puerto 8080) y el servidor Go (nuevo, puerto 8082).

"

$summary = "| Catálogo | Java Status | Go Status | Paridad |
|---|---|---|---|
"

foreach ($cat in $catalogs) {
    Write-Host "Procesando $cat..."
    
    if ($cat -eq "Timeout") {
        $jResRaw = curl.exe -s -m 10 -w "HTTP_STATUS:%{http_code}" -X POST http://localhost:8080/v1/bfcl/mortgage-loan/catalogs/$cat
        $gResRaw = curl.exe -s -m 10 -w "HTTP_STATUS:%{http_code}" -X POST http://localhost:8082/v1/bfcl/mortgage-loan/catalogs/$cat
    } else {
        $jResRaw = curl.exe -s -w "HTTP_STATUS:%{http_code}" -X POST http://localhost:8080/v1/bfcl/mortgage-loan/catalogs/$cat
        $gResRaw = curl.exe -s -w "HTTP_STATUS:%{http_code}" -X POST http://localhost:8082/v1/bfcl/mortgage-loan/catalogs/$cat
    }

    # Parse Java
    $jStatus = $jResRaw -replace '.*HTTP_STATUS:(\d+)$', '$1'
    $jBody = $jResRaw -replace 'HTTP_STATUS:\d+$', ''
    
    # Parse Go
    $gStatus = $gResRaw -replace '.*HTTP_STATUS:(\d+)$', '$1'
    $gBody = $gResRaw -replace 'HTTP_STATUS:\d+$', ''

    $isIdentical = ($jBody -eq $gBody) -and ($jStatus -eq $gStatus)
    $parityStr = if ($isIdentical) { "✅ IDÉNTICO" } else { "❌ DIFERENTE" }
    
    $summary += "| $cat | $jStatus | $gStatus | $parityStr |
"

    $md += "### CU-XXX: $cat

"
    $md += "#### Request
"
    $md += "- **URL:** `POST /v1/bfcl/mortgage-loan/catalogs/$cat`
"
    $md += "- **Headers:** (Ninguno obligatorio)
"
    $md += "- **Body:** (vacío)

"
    
    $md += "#### Response — Java Legado (8080)
"
    $md += "**HTTP Status:** ``
"
    $md += "```json
$jBody
``` 

"

    $md += "#### Response — Go Hexagonal (8082)
"
    $md += "**HTTP Status:** ``
"
    $md += "```json
$gBody
``` 

"
    
    $md += "#### Paridad
"
    $md += "**$parityStr**

---

"
}

$md += "## Resumen de Paridad

"
$md += $summary

$md | Out-File -FilePath "C:\HEXAGONAL\mortgage-loan-catalogs\docs\PAYLOADS_CASOS_DE_USO.md" -Encoding UTF8 -Force
Write-Host "✅ Documento PAYLOADS_CASOS_DE_USO.md generado exitosamente."
