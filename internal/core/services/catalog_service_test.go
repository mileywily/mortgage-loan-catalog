package services_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"mortgage-loan-catalogs/internal/core/domain"
	"mortgage-loan-catalogs/internal/core/services"
)

// 1. Creamos nuestro Mock del Puerto de Salida (Repository)
type MockCatalogRepository struct {
	mock.Mock
}

func (m *MockCatalogRepository) GetCatalog(ctx context.Context, req domain.GetCatalogRequest) (interface{}, error) {
	args := m.Called(ctx, req)
	return args.Get(0), args.Error(1)
}

func TestCatalogService_GetCatalog_Success(t *testing.T) {
	// Arrange
	mockRepo := new(MockCatalogRepository)
	svc := services.NewCatalogService(mockRepo)

	req := domain.GetCatalogRequest{
		CatalogName: "Destino",
		Channel:     "WEB",
	}

	expectedResult := []domain.CatalogItem{
		{Code: "1", Description: "Vivienda Principal"},
	}

	// Configuramos el mock para que devuelva expectedResult cuando sea llamado
	mockRepo.On("GetCatalog", mock.Anything, req).Return(expectedResult, nil)

	// Act
	result, err := svc.GetCatalog(context.Background(), req)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, expectedResult, result)
	mockRepo.AssertExpectations(t) // Verifica que se llamaron a todos los métodos esperados
}

func TestCatalogService_GetCatalog_NotFound(t *testing.T) {
	// Arrange
	mockRepo := new(MockCatalogRepository)
	svc := services.NewCatalogService(mockRepo)

	req := domain.GetCatalogRequest{CatalogName: "Inexistente"}

	// Configuramos el mock para simular que el upstream devuelve NotFound
	mockRepo.On("GetCatalog", mock.Anything, req).Return(nil, services.ErrCatalogNotFound)

	// Act
	result, err := svc.GetCatalog(context.Background(), req)

	// Assert
	assert.ErrorIs(t, err, services.ErrCatalogNotFound)
	assert.Nil(t, result)
	mockRepo.AssertExpectations(t)
}
