package e2e_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"mortgage-loan-catalogs/internal/app/handler"
	"mortgage-loan-catalogs/internal/app/middleware"
	"mortgage-loan-catalogs/internal/core/services"
	"mortgage-loan-catalogs/internal/infra/repository"
)

func setupFullApp() *gin.Engine {
	gin.SetMode(gin.TestMode)

	// Inicializamos el Dummy Repository (ideal para E2E locales sin DB ni Upstream de red real)
	// Actúa como nuestra "Base de datos en memoria"
	dummyRepo := repository.NewDummyRepository()
	dummyService := services.NewCatalogService(dummyRepo)

	// Instanciamos el controlador inyectando el dummy como servicio principal y fallbacks
	catalogHandler := handler.NewCatalogHandler(dummyService, dummyService, dummyService, dummyService)

	r := gin.New()
	r.POST("/v1/bfcl/mortgage-loan/catalogs/:catalog", middleware.LegacyHeaderValidator(), catalogHandler.GetCatalog)
	return r
}

func TestE2E_FullCatalogFlow(t *testing.T) {
	// Arrange: Inicializar aplicación completa con dependencias cableadas
	app := setupFullApp()

	// Simulamos un request E2E completo desde el exterior
	req, _ := http.NewRequest(http.MethodPost, "/v1/bfcl/mortgage-loan/catalogs/Destino", nil)
	req.Header.Set("X-Channel", "APP")
	req.Header.Set("X-Commerce", "TOTTUS")
	req.Header.Set("X-Transaction-ID", "999888")

	w := httptest.NewRecorder()

	// Act: Disparamos contra el Router
	app.ServeHTTP(w, req)

	// Assert: Verificar el pipeline completo (Middleware -> Handler -> Service -> Repository)
	assert.Equal(t, http.StatusOK, w.Code)
	
	// Validar payload
	var jsonResponse []map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &jsonResponse)
	assert.NoError(t, err)
	
	assert.Greater(t, len(jsonResponse), 0)
	assert.Equal(t, "1", jsonResponse[0]["codigo_adm"])
	assert.Equal(t, "Vivienda Principal", jsonResponse[0]["descripcion"])

	// Validar que el middleware NO hace echo de los headers (Paridad legado)
	assert.Empty(t, w.Header().Get("X-Transaction-ID"))
}
