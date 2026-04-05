package api

import (
	"log"
	"net/http"

	"github.com/HarkHorning/portfolio-go-svelte-azure-k8/internal/repo"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	sqlResource repo.Repo
}

func NewHandler(sqlResource repo.Repo) *Handler {
	return &Handler{
		sqlResource: sqlResource,
	}
}

func (h *Handler) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}

func (h *Handler) GetArtTiles(c *gin.Context) {
	tiles, err := h.sqlResource.TopTiles(12)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get art"})
		log.Println("SERVER: Failed to get art")
		return
	}
	c.JSON(http.StatusOK, tiles)
}
