package services

import (
	"context"
	"errors"

	"mortgage-loan-catalogs/internal/core/domain"
	"mortgage-loan-catalogs/internal/core/ports"
)

var (
	ErrCatalogNotFound = errors.New("CatalogNotFoundError")
)

type catalogService struct {
	repo ports.CatalogRepository
}

// NewCatalogService creates a new instance of CatalogService injecting the repository dependency
func NewCatalogService(repo ports.CatalogRepository) ports.CatalogService {
	return &catalogService{
		repo: repo,
	}
}

// GetCatalog executes the core business logic.
// For now it delegates to the repository, but any future domain validation can be added here.
func (s *catalogService) GetCatalog(ctx context.Context, req domain.GetCatalogRequest) (interface{}, error) {
	// Delegate to infra (DB / Mock / Upstream)
	// Upstream system is the source of truth for whether a catalog exists
	return s.repo.GetCatalog(ctx, req)
}
