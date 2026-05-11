package main

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/xlsmart-api/partner/handler"
	"github.com/xlsmart-api/partner/repository"
	"github.com/xlsmart-api/partner/usecase"
)

func main() {
	// Load .env file
	godotenv.Load()

	e := echo.New()

	// Global middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	// Initialize database connection
	dbDsn := os.Getenv("DATABASE_URL")
	if dbDsn != "" {
		e.Logger.Infof("Using DATABASE_URL from environment")
	} else {
		e.Logger.Warnf("DATABASE_URL env not found, using default localhost fallback")
		dbDsn = "host=localhost user=postgres password=postgres dbname=partner_db port=5434 sslmode=disable"
	}

	// Initialize repository
	repo := repository.NewRepository(repository.NewRepositoryOptions{Dsn: dbDsn})

	// Initialize use case
	uc := usecase.NewUseCase(usecase.NewUseCaseOptions{
		Repository: repo,
	})

	// Initialize handler
	h := handler.NewServer(handler.NewServerOptions{UseCase: uc})

	// Register routes (aligned with api-contract: /partners/orders, /partners/fulfillment)
	e.POST("/partners/orders", h.PostPartnerOrdersSubmit)
	e.POST("/partners/fulfillment", h.PostPartnerOrdersFulfillment)

	// Health check endpoint
	e.GET("/health", func(c echo.Context) error {
		return c.JSON(200, map[string]string{
			"service": "partner",
			"status":  "healthy",
		})
	})

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8086"
	}

	e.Logger.Info("Starting Partner service on port " + port)
	e.Logger.Fatal(e.Start(":" + port))
}
