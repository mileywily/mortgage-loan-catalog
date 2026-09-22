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
