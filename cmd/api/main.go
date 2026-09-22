package main

import (
	"log/slog"
	"os"

	"github.com/gin-gonic/gin"
	gintrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/gin-gonic/gin"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"

	"mortgage-loan-catalogs/internal/app/handler"
	"mortgage-loan-catalogs/internal/app/middleware"
	"mortgage-loan-catalogs/internal/core/services"
	"mortgage-loan-catalogs/internal/infra/repository"

	_ "mortgage-loan-catalogs/docs"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title API de Catalogos de Creditos Hipotecarios
// @version 1.0
// @description Microservicio en Go para consultar los catalogos legados de creditos hipotecarios.
// @host localhost:8082
// @BasePath /v1/bfcl/mortgage-loan
func main() {
	// Configure JSON logger
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))
	slog.SetDefault(logger)

	// Datadog Tracer
	tracer.Start()
	defer tracer.Stop()

	// 1. Dependency Injection (Wiring)
	finnflowURL := os.Getenv("FINNFLOW_URL")
	if finnflowURL == "" { finnflowURL = "http://localhost:9090" }
	finnflowKey := os.Getenv("FINNFLOW_KEY")
	finnflowSecret := os.Getenv("FINNFLOW_SECRET")
	
	javaLegacyURL := os.Getenv("JAVA_LEGACY_URL")
	if javaLegacyURL == "" { javaLegacyURL = "http://localhost:8080" }

	dummyRepo := repository.NewDummyRepository()
	// realRepo is Finnflow (GET, client_id, client_secret)
	realRepo := repository.NewRestRepository(finnflowURL, "/api/catalogo_detail/", true, finnflowKey, finnflowSecret)
	// javaRepo is Java Proxy (POST, X-Headers)
	javaRepo := repository.NewRestRepository(javaLegacyURL, "/v1/bfcl/mortgage-loan/catalogs/", false, "", "")
	
	dummyService := services.NewCatalogService(dummyRepo)
	realService := services.NewCatalogService(realRepo)
	javaService := services.NewCatalogService(javaRepo)
	
	defaultBackend := os.Getenv("DEFAULT_BACKEND")
	defaultService := realService // Fallback by default (real URL)
	
	if defaultBackend == "dummy" {
		defaultService = dummyService
	} else if defaultBackend == "java" {
		defaultService = javaService
	}

	catalogHandler := handler.NewCatalogHandler(defaultService, dummyService, realService, javaService)

	// 2. Instantiate Gin
	r := gin.New()

	// Datadog Middleware for Gin
	r.Use(gintrace.Middleware("mortgage-loan-catalogs"))

	// 3. Setup routes
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "UP",
		})
	})

	// Legacy matching route
	r.POST("/v1/bfcl/mortgage-loan/catalogs/:catalog", middleware.LegacyHeaderValidator(), catalogHandler.GetCatalog)

	// Swagger route
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8082" // Align with parity test script
	}

	slog.Info("Starting application", "port", port)
	if err := r.Run(":" + port); err != nil {
		slog.Error("Server failed", "error", err)
		os.Exit(1)
	}
}
