package ports

import (
	"context"

	"mortgage-loan-catalogs/internal/core/domain"
)

// CatalogRepository defines the outbound port to fetch catalogs (e.g. from DB or external API)
type CatalogRepository interface {
	GetCatalog(ctx context.Context, req domain.GetCatalogRequest) (interface{}, error)
}

// CatalogService defines the inbound port (Use Case) for the application layer to consume
type CatalogService interface {
	GetCatalog(ctx context.Context, req domain.GetCatalogRequest) (interface{}, error)
}
