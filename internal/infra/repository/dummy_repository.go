package repository

import (
	"context"

	"mortgage-loan-catalogs/internal/core/domain"
	"mortgage-loan-catalogs/internal/core/ports"
)

type dummyRepository struct{}

// NewDummyRepository returns a mock repository matching the legacy Java data.
func NewDummyRepository() ports.CatalogRepository {
	return &dummyRepository{}
}

func (r *dummyRepository) GetCatalog(ctx context.Context, req domain.GetCatalogRequest) (interface{}, error) {
	// Dummy response matching EXACTLY the Java legacy documentation
	if req.CatalogName == "SegurosIncendio" || req.CatalogName == "SegurosDesgravamen" || req.CatalogName == "SegurosCesantia" {
		return []domain.InsuranceCatalogItem{
			{
				InsuranceIdentifier:  "1100438-01-2023-000",
				Description:          "INCENDIO - Everest compañía de seguros generales Chile(0.2255300)",
				Policy:               "100438-01-2023-000",
				CompanyName:          "Everest compañía de seguros generales Chile",
				CompanyCode:          39,
				InsuranceTypeCode:    0,
				PolicyCorrelative:    28,
				Rate:                 0.22553,
				Factor:               0.0002255,
				IndividualPolicyFlag: 0,
				PerQuotaValueFlag:    0,
				ExternalPolicyFlag:   0,
			},
		}, nil
	}

	if req.CatalogName == "TiposDocumentos" {
		return []domain.TipoDocumentoCatalogItem{
			{
				Code:        1,
				Description: "Escritura de Propiedad",
				GroupID:     2,
			},
		}, nil
	}

	if req.CatalogName == "Vacio" || req.CatalogName == "ErrorRed" {
		// Simulates a network error or empty results from backend
		return []domain.NoResultsCatalogItem{
			{
				CodRespuesta: 3,
				Mensaje:      "Sin resultados.",
				Excepcion:    "Ninguna",
			},
		}, nil
	}

	if req.CatalogName == "Comunas" {
		return []domain.ComunaCatalogItem{
			{
				Code:        "01101",
				Description: "Iquique",
				RegionID:    1,
			},
		}, nil
	}

	// Standard Catalog exact match
	return []domain.CatalogItem{
		{
			Code:        "1",
			Description: "Vivienda Principal",
		},
	}, nil
}
