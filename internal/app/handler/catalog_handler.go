package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"mortgage-loan-catalogs/internal/core/domain"
	"mortgage-loan-catalogs/internal/core/ports"
	"mortgage-loan-catalogs/internal/core/services"
)

// catalogItemDTO represents the JSON payload for standard catalogs.
type catalogItemDTO struct {
	CodigoAdm   string `json:"codigo_adm"`
	Descripcion string `json:"descripcion"`
}

// comunaCatalogItemDTO represents the JSON payload for comunas.
type comunaCatalogItemDTO struct {
	CodigoAdm   string `json:"codigo_adm"`
	Descripcion string `json:"descripcion"`
	RegionID    int    `json:"region_id"`
}

// tipoDocumentoCatalogItemDTO represents the JSON payload for tipos documentos.
type tipoDocumentoCatalogItemDTO struct {
	CodigoAdm   int    `json:"codigo_adm"`
	Descripcion string `json:"descripcion"`
	GrupoID     int    `json:"grupo_id"`
}

// noResultsCatalogItemDTO represents the legacy empty/network error fallback
type noResultsCatalogItemDTO struct {
	CodRespuesta int    `json:"codRespuesta"`
	Mensaje      string `json:"Mensaje"`
	Excepcion    string `json:"Excepcion"`
}

// insuranceCatalogItemDTO represents the JSON payload for insurance catalogs.
type insuranceCatalogItemDTO struct {
	IdentificadorSeguro       string      `json:"IdentificadorSeguro"`
	Descripcion               string      `json:"Descripcion"`
	Poliza                    string      `json:"Poliza"`
	NombreCompania            string      `json:"NombreCompania"`
	CodigoCompania            int         `json:"CodigoCompania"`
	CodigoTipoSeguro          int         `json:"CodigoTipoSeguro"`
	CorrelativoPoliza         int         `json:"CorrelativoPoliza"`
	Tasa                      json.Number `json:"Tasa"`
	Factor                    json.Number `json:"Factor"`
	IndicadorPolizaIndividual int         `json:"IndicadorPolizaIndividual"`
	PorValorCuota             int         `json:"PorValorCuota"`
	IndicadorPolizaExterna    int         `json:"IndicadorPolizaExterna"`
}

// errorResponseDTO ensures exact JSON key ordering as Java
type errorResponseDTO struct {
	Detail       *string       `json:"detail"`
	Code         *string       `json:"code"`
	Messages     []interface{} `json:"messages"`
	ErrorsDetail *string       `json:"errors_detail"`
}

// internalServerErrorDTO represents the legacy 500 error format
type internalServerErrorDTO struct {
	Mensaje    string `json:"mensaje"`
	StatusCode int    `json:"statusCode"`
	Timestamp  string `json:"timestamp"`
}

type CatalogHandler struct {
	defaultService ports.CatalogService
	dummyService   ports.CatalogService
	realService    ports.CatalogService
	javaService    ports.CatalogService
}

func NewCatalogHandler(defaultSvc, dummy, real, java ports.CatalogService) *CatalogHandler {
	return &CatalogHandler{
		defaultService: defaultSvc,
		dummyService:   dummy,
		realService:    real,
		javaService:    java,
	}
}

// GetCatalog godoc
// @Summary Obtener Catálogo
// @Description Retorna los elementos de un catálogo dado (simulando comportamiento legado).
// @Tags catalogs
// @Accept json
// @Produce json
// @Param catalog path string true "Nombre del catálogo (ej. Destino, Comunas, SegurosIncendio)"
// @Param X-Channel header string true "Canal (ej. WEB)"
// @Param X-Commerce header string true "Comercio (ej. FALABELLA)"
// @Param X-Transaction-ID header string true "ID de transacción (ej. 12345)"
// @Param X-Backend-Env header string false "Entorno a inyectar: dummy, real, o java (por defecto: dummy)"
// @Success 200 {array} catalogItemDTO
// @Failure 401 {object} errorResponseDTO
// @Failure 402 {object} errorResponseDTO
// @Failure 400 {object} errorResponseDTO
// @Failure 500 {object} internalServerErrorDTO
// @Router /catalogs/{catalog} [post]
func (h *CatalogHandler) GetCatalog(c *gin.Context) {
	channel := c.GetHeader("X-Channel")
	commerce := c.GetHeader("X-Commerce")
	trxID := c.GetHeader("X-Transaction-ID")
	backendEnv := strings.ToLower(c.GetHeader("X-Backend-Env"))

	catalogName := c.Param("catalog")

	req := domain.GetCatalogRequest{
		CatalogName:   catalogName,
		Channel:       channel,
		Commerce:      commerce,
		TransactionID: trxID,
	}

	activeService := h.defaultService
	if backendEnv == "real" {
		activeService = h.realService
	} else if backendEnv == "java" {
		activeService = h.javaService
	} else if backendEnv == "dummy" {
		activeService = h.dummyService
	}

	result, err := activeService.GetCatalog(c.Request.Context(), req)
	if err != nil {
		handleError(c, err)
		return
	}

	// Map domain entities to DTOs
	switch data := result.(type) {
	case []domain.CatalogItem:
		dtos := make([]catalogItemDTO, len(data))
		for i, item := range data {
			dtos[i] = catalogItemDTO{
				CodigoAdm:   item.Code,
				Descripcion: item.Description,
			}
		}
		c.JSON(http.StatusOK, dtos)
	case []domain.ComunaCatalogItem:
		dtos := make([]comunaCatalogItemDTO, len(data))
		for i, item := range data {
			dtos[i] = comunaCatalogItemDTO{
				CodigoAdm:   item.Code,
				Descripcion: item.Description,
				RegionID:    item.RegionID,
			}
		}
		c.JSON(http.StatusOK, dtos)
	case []domain.TipoDocumentoCatalogItem:
		dtos := make([]tipoDocumentoCatalogItemDTO, len(data))
		for i, item := range data {
			dtos[i] = tipoDocumentoCatalogItemDTO{
				CodigoAdm:   item.Code,
				Descripcion: item.Description,
				GrupoID:     item.GroupID,
			}
		}
		c.JSON(http.StatusOK, dtos)
	case []domain.NoResultsCatalogItem:
		dtos := make([]noResultsCatalogItemDTO, len(data))
		for i, item := range data {
			dtos[i] = noResultsCatalogItemDTO{
				CodRespuesta: item.CodRespuesta,
				Mensaje:      item.Mensaje,
				Excepcion:    item.Excepcion,
			}
		}
		c.JSON(http.StatusOK, dtos)
	case []domain.InsuranceCatalogItem:
		dtos := make([]insuranceCatalogItemDTO, len(data))
		for i, item := range data {
			dtos[i] = insuranceCatalogItemDTO{
				IdentificadorSeguro:       item.InsuranceIdentifier,
				Descripcion:               item.Description,
				Poliza:                    item.Policy,
				NombreCompania:            item.CompanyName,
				CodigoCompania:            item.CompanyCode,
				CodigoTipoSeguro:          item.InsuranceTypeCode,
				CorrelativoPoliza:         item.PolicyCorrelative,
				Tasa:                      json.Number(formatFloatLikeJackson(item.Rate)),
				Factor:                    json.Number(formatFloatLikeJackson(item.Factor)),
				IndicadorPolizaIndividual: item.IndividualPolicyFlag,
				PorValorCuota:             item.PerQuotaValueFlag,
				IndicadorPolizaExterna:    item.ExternalPolicyFlag,
			}
		}
		c.JSON(http.StatusOK, dtos)
	default:
		c.JSON(http.StatusOK, result)
	}
}

func handleError(c *gin.Context, err error) {
	errMessage := err.Error()

	if strings.Contains(errMessage, "AuthenticationError") {
		tokenMessage := struct {
			TokenClass string `json:"token_class"`
			TokenType  string `json:"token_type"`
			Message    string `json:"message"`
		}{
			TokenClass: "AccessToken",
			TokenType:  "access",
			Message:    "Token is invalid",
		}

		d := "Given token not valid for any token type"
		code := "token_not_valid"
		c.JSON(http.StatusUnauthorized, errorResponseDTO{
			Detail:       &d,
			Code:         &code,
			Messages:     []interface{}{tokenMessage},
			ErrorsDetail: nil,
		})
		return
	}

	if errors.Is(err, services.ErrCatalogNotFound) {
		ed := "Catálogo no encontrado"
		c.JSON(http.StatusPaymentRequired, errorResponseDTO{
			Detail:       nil,
			Code:         nil,
			Messages:     nil,
			ErrorsDetail: &ed,
		})
		return
	}

	// Default fallback to 400 Bad Request (Legacy Parity for CatalogException)
	// Java legacy maps any other upstream error to 400 Bad Request
	ed := "Error en la obtenci\u00f3n de cat\u00e1logos"
	c.JSON(http.StatusBadRequest, errorResponseDTO{
		Detail:       nil,
		Code:         nil,
		Messages:     nil,
		ErrorsDetail: &ed,
	})
}

func formatFloatLikeJackson(f float64) string {
	s := strconv.FormatFloat(f, 'g', -1, 64)
	if f != 0 && f > -1e-3 && f < 1e-3 {
		s = strconv.FormatFloat(f, 'E', -1, 64)
		s = strings.Replace(s, "E-0", "E-", 1)
		s = strings.Replace(s, "E+0", "E", 1)
	}
	return s
}
