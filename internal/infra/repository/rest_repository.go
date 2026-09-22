package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"mortgage-loan-catalogs/internal/core/domain"
	"mortgage-loan-catalogs/internal/core/ports"
	"mortgage-loan-catalogs/internal/core/services"
)

type restRepository struct {
	baseURL       string
	routeTemplate string
	isFinnflow    bool
	username      string
	password      string
	httpClient    *http.Client
}

func NewRestRepository(baseURL, routeTemplate string, isFinnflow bool, user, pass string) ports.CatalogRepository {
	return &restRepository{
		baseURL:       baseURL,
		routeTemplate: routeTemplate,
		isFinnflow:    isFinnflow,
		username:      user,
		password:      pass,
		httpClient:    &http.Client{},
	}
}

func (r *restRepository) GetCatalog(ctx context.Context, req domain.GetCatalogRequest) (interface{}, error) {
	url := fmt.Sprintf("%s%s%s", r.baseURL, r.routeTemplate, req.CatalogName)

	method := http.MethodPost
	if r.isFinnflow {
		method = http.MethodGet // Finnflow uses GET natively
	}

	httpReq, err := http.NewRequestWithContext(ctx, method, url, nil)
	if err != nil {
		return nil, err
	}

	if r.isFinnflow {
		// Java legacy uses Basic Auth (Base64 username:password) exactly as createAuthHeaders() does
		// Content-Type: application/json is set by Java's createAuthHeaders()
		httpReq.Header.Set("Content-Type", "application/json")
		if r.username != "" || r.password != "" {
			httpReq.SetBasicAuth(r.username, r.password)
		}
	} else {
		// Java Legacy Proxy or regular POST uses X-Headers
		if req.Channel != "" { httpReq.Header.Set("X-Channel", req.Channel) }
		if req.Commerce != "" { httpReq.Header.Set("X-Commerce", req.Commerce) }
		if req.TransactionID != "" { httpReq.Header.Set("X-Transaction-ID", req.TransactionID) }
	}

	resp, err := r.httpClient.Do(httpReq)
	if err != nil {
		if strings.Contains(err.Error(), "timeout") || strings.Contains(err.Error(), "connection") || strings.Contains(err.Error(), "refused") {
			return []domain.NoResultsCatalogItem{{CodRespuesta: 3, Mensaje: "Sin resultados.", Excepcion: "Ninguna"}}, nil
		}
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		switch resp.StatusCode {
		case http.StatusUnauthorized:
			return nil, fmt.Errorf("AuthenticationError: %d", resp.StatusCode)
		case http.StatusPaymentRequired, http.StatusNotFound:
			return nil, services.ErrCatalogNotFound
		default:
			return nil, fmt.Errorf("LegacySystemError: %d", resp.StatusCode)
		}
	}

	catalogLower := strings.ToLower(req.CatalogName)
	if strings.HasPrefix(catalogLower, "seguros") {
		var upstreamItems []struct {
			IdentificadorSeguro       string  `json:"IdentificadorSeguro"`
			Descripcion               string  `json:"Descripcion"`
			Poliza                    string  `json:"Poliza"`
			NombreCompania            string  `json:"NombreCompania"`
			CodigoCompania            int     `json:"CodigoCompania"`
			CodigoTipoSeguro          int     `json:"CodigoTipoSeguro"`
			CorrelativoPoliza         int     `json:"CorrelativoPoliza"`
			Tasa                      float64 `json:"Tasa"`
			Factor                    float64 `json:"Factor"`
			IndicadorPolizaIndividual int     `json:"IndicadorPolizaIndividual"`
			PorValorCuota             int     `json:"PorValorCuota"`
			IndicadorPolizaExterna    int     `json:"IndicadorPolizaExterna"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&upstreamItems); err != nil {
			if strings.Contains(err.Error(), "EOF") {
				return []domain.NoResultsCatalogItem{{CodRespuesta: 3, Mensaje: "Sin resultados.", Excepcion: "Ninguna"}}, nil
			}
			return nil, err
		}
		domainItems := make([]domain.InsuranceCatalogItem, len(upstreamItems))
		for i, item := range upstreamItems {
			domainItems[i] = domain.InsuranceCatalogItem{
				InsuranceIdentifier:  item.IdentificadorSeguro,
				Description:          item.Descripcion,
				Policy:               item.Poliza,
				CompanyName:          item.NombreCompania,
				CompanyCode:          item.CodigoCompania,
				InsuranceTypeCode:    item.CodigoTipoSeguro,
				PolicyCorrelative:    item.CorrelativoPoliza,
				Rate:                 item.Tasa,
				Factor:               item.Factor,
				IndividualPolicyFlag: item.IndicadorPolizaIndividual,
				PerQuotaValueFlag:    item.PorValorCuota,
				ExternalPolicyFlag:   item.IndicadorPolizaExterna,
			}
		}
		return domainItems, nil
	}

	if req.CatalogName == "TiposDocumentos" {
		var upstreamItems []struct {
			CodigoAdm   int    `json:"codigo_adm"`
			Descripcion string `json:"descripcion"`
			GrupoID     int    `json:"grupo_id"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&upstreamItems); err != nil {
			if strings.Contains(err.Error(), "EOF") {
				return []domain.NoResultsCatalogItem{{CodRespuesta: 3, Mensaje: "Sin resultados.", Excepcion: "Ninguna"}}, nil
			}
			return nil, err
		}
		domainItems := make([]domain.TipoDocumentoCatalogItem, len(upstreamItems))
		for i, item := range upstreamItems {
			domainItems[i] = domain.TipoDocumentoCatalogItem{
				Code:        item.CodigoAdm,
				Description: item.Descripcion,
				GroupID:     item.GrupoID,
			}
		}
		return domainItems, nil
	}

	if req.CatalogName == "Comunas" {
		var upstreamItems []struct {
			CodigoAdm   interface{} `json:"codigo_adm"`
			Descripcion string      `json:"descripcion"`
			RegionID    int         `json:"region_id"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&upstreamItems); err != nil {
			if strings.Contains(err.Error(), "EOF") {
				return []domain.NoResultsCatalogItem{{CodRespuesta: 3, Mensaje: "Sin resultados.", Excepcion: "Ninguna"}}, nil
			}
			return nil, err
		}
		domainItems := make([]domain.ComunaCatalogItem, len(upstreamItems))
		for i, item := range upstreamItems {
			domainItems[i] = domain.ComunaCatalogItem{
				Code:        fmt.Sprintf("%v", item.CodigoAdm),
				Description: item.Descripcion,
				RegionID:    item.RegionID,
			}
		}
		return domainItems, nil
	}

	var upstreamItems []struct {
		CodigoAdm    interface{} `json:"codigo_adm"`
		Descripcion  string      `json:"descripcion"`
		CodRespuesta *int        `json:"codRespuesta"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&upstreamItems); err != nil {
		if strings.Contains(err.Error(), "EOF") {
			return []domain.NoResultsCatalogItem{{CodRespuesta: 3, Mensaje: "Sin resultados.", Excepcion: "Ninguna"}}, nil
		}
		return nil, err
	}

	if len(upstreamItems) == 0 || (len(upstreamItems) > 0 && upstreamItems[0].CodRespuesta != nil) {
		return []domain.NoResultsCatalogItem{{CodRespuesta: 3, Mensaje: "Sin resultados.", Excepcion: "Ninguna"}}, nil
	}

	domainItems := make([]domain.CatalogItem, len(upstreamItems))
	for i, item := range upstreamItems {
		domainItems[i] = domain.CatalogItem{
			Code:        fmt.Sprintf("%v", item.CodigoAdm),
			Description: item.Descripcion,
		}
	}
	return domainItems, nil
}
