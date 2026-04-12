package api

import (
	"log/slog"
	"net/http"
	"strconv"

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
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func (h *Handler) ReadyCheck(c *gin.Context) {
	if err := h.sqlResource.Ping(); err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"status": "not ready",
			"reason": "database unavailable",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "ready"})
}

func (h *Handler) GetArtTiles(c *gin.Context) {
	category := c.Query("category")

	var result interface{}
	var err error

	if category != "" {
		result, err = h.sqlResource.TilesByCategory(category)
	} else {
		result, err = h.sqlResource.TopTiles(12)
	}

	if err != nil {
		slog.Error("failed to get art tiles", "error", err, "category", category)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get art"})
		return
	}

	c.JSON(http.StatusOK, result)
}

func (h *Handler) GetArtByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	art, err := h.sqlResource.ArtByID(id)
	if err != nil {
		slog.Error("failed to get art by id", "id", id, "error", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "art not found"})
		return
	}

	c.JSON(http.StatusOK, art)
}

func (h *Handler) GetCategories(c *gin.Context) {
	categories, err := h.sqlResource.AllCategories()
	if err != nil {
		slog.Error("failed to get categories", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get categories"})
		return
	}
	c.JSON(http.StatusOK, categories)
}
