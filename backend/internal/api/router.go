package api

import (
	"os"

	"github.com/HarkHorning/portfolio-go-svelte-azure-k8/internal/repo"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

func Routes(db *sqlx.DB) *gin.Engine {
	router := gin.Default()

	config := cors.Config{
		AllowOrigins:     getAllowedOrigins(),
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}, // come back to this later
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}
	router.Use(cors.New(config))

	sqlResource := repo.NewRepo(db)

	handle := NewHandler(*sqlResource)

	router.GET("/health", handle.HealthCheck)
	router.GET("/ready", handle.ReadyCheck)

	v1 := router.Group("/api/v1")
	{
		v1.GET("/art", handle.GetArtTiles)
		v1.GET("/art/:id", handle.GetArtByID)
		v1.GET("/categories", handle.GetCategories)
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
