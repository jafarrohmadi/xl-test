package main

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/xlsmart-api/sf-backend/handler"
	"github.com/xlsmart-api/sf-backend/repository"
	"github.com/xlsmart-api/sf-backend/usecase"
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
		dbDsn = "host=localhost user=postgres password=postgres dbname=sf_backend_db port=5432 sslmode=disable"
	}

	// Initialize repository
	repo := repository.NewRepository(repository.NewRepositoryOptions{Dsn: dbDsn})

	// Initialize use case
	uc := usecase.NewUseCase(usecase.NewUseCaseOptions{
		Repository: repo,
	})

	// Initialize handler
	h := handler.NewServer(handler.NewServerOptions{UseCase: uc})

	// Register routes (updated endpoints)
	e.POST("/orders", h.PostOrders)
	e.POST("/orders/fulfillment/callback", h.PostOrdersFulfillmentCallback)
	e.POST("/internal/notifications", h.PostInternalNotifications)

	// Health check endpoint
	e.GET("/health", func(c echo.Context) error {
		return c.JSON(200, map[string]string{
			"service": "sf-backend",
			"status":  "healthy",
		})
	})

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8084"
	}

	e.Logger.Info("Starting SF Backend service on port " + port)
	e.Logger.Fatal(e.Start(":" + port))
}
