BUILDPATH=$(CURDIR)
API_NAME=mortgage-loan-catalogs
PACKAGES := $(shell go list ./internal/... | grep -vE 'mocks')

build:
	@echo "Creando Binario ..."
	@go build -ldflags '-s -w' -o $(BUILDPATH)/build/bin/dist cmd/api/main.go
	@echo "Binario generado en build/bin/dist"

test:
	@echo "Ejecutando tests..."
	@go test  $(go list ./internal/... | grep -vE 'mocks|ports|config|cmd') --coverprofile coverfile_out >> /dev/null

	@go tool cover -func coverfile_out

coverage:
	@echo "Coverfile..."
	go test $(PACKAGES) --coverprofile coverfile_out >> /dev/null
	@go tool cover -func coverfile_out
	@go tool cover -func coverfile_out | grep total | grep -o '[0-9]*\.[0-9]*' | cut -d' ' -f1 > coverage.txt

.PHONY: test build

