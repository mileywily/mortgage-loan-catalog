package domain_test

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"mortgage-loan-catalogs/internal/core/domain"
)

func TestDomainStruct(t *testing.T) {
	req := domain.GetCatalogRequest{CatalogName: "test"}
	assert.Equal(t, "test", req.CatalogName)
}
