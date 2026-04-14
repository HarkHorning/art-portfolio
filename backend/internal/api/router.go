package api

import (
	"log/slog"
	"os"

	"github.com/HarkHorning/portfolio-go-svelte-azure-k8/internal/admin"
	"github.com/HarkHorning/portfolio-go-svelte-azure-k8/internal/config"
	"github.com/HarkHorning/portfolio-go-svelte-azure-k8/internal/repo"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

func Routes(db *sqlx.DB, cfg config.Config) *gin.Engine {
	router := gin.Default()

	config := cors.Config{
		AllowOrigins:     getAllowedOrigins(),
		AllowMethods:     []string{"GET", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}
	router.Use(cors.New(config))
	router.Use(RateLimiter())

	sqlResource := repo.NewRepo(db)

	handle := NewHandler(*sqlResource)

	router.GET("/health", handle.HealthCheck)
	router.GET("/ready", handle.ReadyCheck)

	adminGroup := router.Group("/admin")
	if err := admin.RegisterRoutes(adminGroup, sqlResource, admin.TemplateFS, cfg); err != nil {
		slog.Error("failed to register admin routes", "error", err)
	}

	v1 := router.Group("/api/v1")
	{
		v1.GET("/art", handle.GetArtTiles)
		v1.GET("/art/:id", handle.GetArtByID)
		v1.GET("/art-sizes", handle.GetArtSizes)
		v1.GET("/categories", handle.GetCategories)
		v1.GET("/prints", handle.GetPrints)
		v1.GET("/prints/:id", handle.GetPrintByID)
		v1.GET("/print-sizes", handle.GetPrintSizes)
	}

	return router
}

func getAllowedOrigins() []string {
	if origin := os.Getenv("CORS_ORIGIN"); origin != "" {
		return []string{origin}
	}

	return []string{
		"http://localhost:3000", // Svelte dev server
		"http://localhost:5173", // Vite dev server
	}
}
