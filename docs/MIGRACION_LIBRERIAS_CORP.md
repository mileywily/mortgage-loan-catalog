# Guía de Migración: Librerías Corporativas y Variables de Entorno

Actualmente, el microservicio `mortgage-loan-catalogs` está construido usando la Arquitectura Hexagonal y librerías estándar de Go (`gin`, `net/http`, `slog`). Esto permite un desarrollo ágil y sin dependencias de red internas.

Sin embargo, para cumplir al 100% con los estándares del banco (plantillas FIF) antes del paso a Producción, será necesario integrar las **librerías corporativas privadas**. Esta guía detalla cómo realizar esa transición.

---

## 1. Configuración de Acceso a Repositorios Privados (GOPRIVATE)

Las librerías requeridas están alojadas en GitLab interno (`gitlab.falabella.tech`). Para que Go pueda descargarlas sin fallar por temas de proxy público:

1. Configurar la variable de entorno local:
   ```bash
   go env -w GOPRIVATE="gitlab.falabella.tech/*"
   ```
2. Configurar la autenticación de Git (usando un Personal Access Token o SSH):
   ```bash
   git config --global url."https://<TU_USUARIO>:<TU_TOKEN>@gitlab.falabella.tech/".insteadOf "https://gitlab.falabella.tech/"
   ```

## 2. Mapa de Migración de Librerías (`go.mod`)

Una vez que haya acceso, se deben restaurar los bloques `replace` y `require` originales de la plantilla y refactorizar el código así:

| Componente Actual | Librería Corporativa (FIF) | Descripción del Refactor Requerido |
| :--- | :--- | :--- |
| `github.com/gin-gonic/gin` | `gin-fif` | Sustituir `gin.New()` y middleware estándar por el inicializador de `gin-fif`, que usualmente ya incluye trazabilidad, healthchecks y métricas corporativas. |
| `log/slog` | `logger-fif` | Reemplazar las instancias de log estándar. Cambiar `slog.Error()` por las funciones equivalentes en `logger-fif`. |
| `os.Getenv` | `config-fif` | Migrar la carga de variables (`PORT`, `FINNFLOW_URL`) para usar el gestor de configuración corporativo (útil si se conectará a Consul/Vault). |
| Estructuras de Error Custom | `error-fif` | Actualmente retornamos errores estandarizados a mano. Se debe acoplar el `catalog_handler.go` para retornar los códigos usando las estructuras de `error-fif`. |
| `&http.Client{}` | `http-fif` | En `rest_repository.go`, se debe inyectar el cliente HTTP corporativo en lugar del estándar para asegurar trazabilidad (DataDog APM) en peticiones de salida hacia Finnflow. |

## 3. Variables de Entorno para el Pipeline (GitLab CI)

El archivo `.gitlab-ci.yml` importado desde `fif/dx/configurable-pipelines/golang-build-pipeline` requiere ciertas variables a nivel de GitLab CI/CD (Settings > CI/CD > Variables):

- `GCR_SERVICE_ACCOUNT`: Credenciales en Base64 para publicar la imagen en el Registry (Google Container Registry).
- `TEST_AUTH_GCR`: Token/Password adicional usado durante los escaneos de vulnerabilidades.
- `DS_SCANNER_IMAGE`: Imagen de escaneo estático provista por seguridad.
- `SCAN_IMAGE`: Nombre y tag de la imagen resultante (opcional, el pipeline infiere por defecto a partir del commit).

## 4. Variables de Entorno de Ejecución (ConfigMap / Secrets)

Al desplegar en Kubernetes, el servicio (y por ende `config-fif`) requerirá las siguientes variables inyectadas:

- `PORT`: Generalmente `8080` (en nuestro desarrollo usamos 8082).
- `DEFAULT_BACKEND`: `java` o `real` (usado temporalmente para decidir comportamiento de proxy o simulación, puede retirarse en prod).
- `FINNFLOW_URL`: URL base real del microservicio Finnflow interno.
- `FINNFLOW_KEY` / `FINNFLOW_SECRET`: Credenciales para Basic Auth hacia Finnflow (Inyectadas mediante Vault o Kubernetes Secrets, **nunca en código duro**).

## 5. Actualización de Catálogo Backstage (`catalog-info.yaml`)

Antes de realizar el merge definitivo, validar con los Arquitectos los siguientes campos en `catalog-info.yaml` para asegurar que las métricas y los dashboards de DataDog queden asociados al equipo correcto:

- `owner: group:fiftech_integraciones_chile` (Asegurar que sea el grupo correcto).
- `system: fraud-and-aml-risks-chile` (Actualizar al dominio de **Créditos Hipotecarios / Mortgage** si aplica).
- `links: APM DataDog` (Sustituir el link de la plantilla por el dashboard real del nuevo microservicio).
