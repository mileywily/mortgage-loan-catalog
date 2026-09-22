package handler_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"mortgage-loan-catalogs/internal/app/handler"
	"mortgage-loan-catalogs/internal/app/middleware"
	"mortgage-loan-catalogs/internal/core/domain"
)

// Mock del servicio de catálogos (Puerto de Entrada)
type MockCatalogService struct {
	mock.Mock
}

func (m *MockCatalogService) GetCatalog(ctx context.Context, req domain.GetCatalogRequest) (interface{}, error) {
	args := m.Called(ctx, req)
	return args.Get(0), args.Error(1)
}

func setupRouter(mockService *MockCatalogService) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	// Inyectamos el mock al handler. 
	// (Nota: para el test usamos el mismo mock para todos los fallbacks)
	h := handler.NewCatalogHandler(mockService, mockService, mockService, mockService)

	r.POST("/v1/bfcl/mortgage-loan/catalogs/:catalog", middleware.LegacyHeaderValidator(), h.GetCatalog)
	return r
}

func TestCatalogHandler_MissingHeaders_Returns200_Parity(t *testing.T) {
	// Arrange
	mockSvc := new(MockCatalogService)
	router := setupRouter(mockSvc)

	// Java legacy NO retorna 401 por headers faltantes (delega al upstream)
	req, _ := http.NewRequest(http.MethodPost, "/v1/bfcl/mortgage-loan/catalogs/Destino", nil)
	w := httptest.NewRecorder()

	expectedResult := []domain.CatalogItem{
		{Code: "1", Description: "Vivienda Principal"},
	}
	mockSvc.On("GetCatalog", mock.Anything, mock.Anything).Return(expectedResult, nil)

	// Act
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)
	
	// Validar payload
	var response []map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "1", response[0]["codigo_adm"])
}

func TestCatalogHandler_Success(t *testing.T) {
	// Arrange
	mockSvc := new(MockCatalogService)
	router := setupRouter(mockSvc)

	domainReq := domain.GetCatalogRequest{
		CatalogName:   "Destino",
		Channel:       "WEB",
		Commerce:      "FALABELLA",
		TransactionID: "123",
	}

	expectedResult := []domain.CatalogItem{
		{Code: "1", Description: "Vivienda Principal"},
	}

	mockSvc.On("GetCatalog", mock.Anything, domainReq).Return(expectedResult, nil)

	req, _ := http.NewRequest(http.MethodPost, "/v1/bfcl/mortgage-loan/catalogs/Destino", nil)
	req.Header.Set("X-Channel", "WEB")
	req.Header.Set("X-Commerce", "FALABELLA")
	req.Header.Set("X-Transaction-ID", "123")
	
	w := httptest.NewRecorder()

	// Act
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Vivienda Principal")
	assert.Contains(t, w.Body.String(), "codigo_adm") // Verifica el mapeo al DTO (JSON tags correctos)
}
