package api

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/HarkHorning/portfolio-go-svelte-azure-k8/internal/models"
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
	size := c.Query("size")
	minPrice := queryInt(c, "min_price", -1)
	maxPrice := queryInt(c, "max_price", -1)

	tiles, err := h.sqlResource.ArtTiles(category, size, minPrice, maxPrice)
	if err != nil {
		slog.Error("failed to get art tiles", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get art"})
		return
	}

	c.JSON(http.StatusOK, tiles)
}

func (h *Handler) GetArtSizes(c *gin.Context) {
	sizes, err := h.sqlResource.ArtSizes()
	if err != nil {
		slog.Error("failed to get art sizes", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get sizes"})
		return
	}
	c.JSON(http.StatusOK, sizes)
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

	// Filter images by display selection; fall back to all if none selected.
	if displayIDs, err := h.sqlResource.ArtDisplayImageIDs(id); err == nil && len(displayIDs) > 0 {
		idSet := make(map[int]bool, len(displayIDs))
		for _, did := range displayIDs {
			idSet[did] = true
		}
		filtered := art.Images[:0]
		for _, img := range art.Images {
			if idSet[img.Id] {
				filtered = append(filtered, img)
			}
		}
		art.Images = filtered
	}

	// Keep the public API shape the frontend expects: top-level "url" field.
	// Prefer the first low-quality image for the tile display URL.
	type artResponse struct {
		*models.ArtDetailModel
		URL string `json:"url"`
	}
	displayURL := firstImageURL(art.Images, "low")
	if displayURL == "" {
		displayURL = firstImageURL(art.Images, "high")
	}
	c.JSON(http.StatusOK, artResponse{ArtDetailModel: art, URL: displayURL})
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

func (h *Handler) GetPrints(c *gin.Context) {
	size := c.Query("size")
	minPrice := queryInt(c, "min_price", -1)
	maxPrice := queryInt(c, "max_price", -1)

	prints, err := h.sqlResource.Prints(size, minPrice, maxPrice)
	if err != nil {
		slog.Error("failed to get prints", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get prints"})
		return
	}
	c.JSON(http.StatusOK, prints)
}

func (h *Handler) GetPrintByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	print, err := h.sqlResource.PrintByID(id)
	if err != nil {
		slog.Error("failed to get print by id", "id", id, "error", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "print not found"})
		return
	}
	c.JSON(http.StatusOK, print)
}

func (h *Handler) GetPrintSizes(c *gin.Context) {
	sizes, err := h.sqlResource.PrintSizes()
	if err != nil {
		slog.Error("failed to get print sizes", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get sizes"})
		return
	}
	c.JSON(http.StatusOK, sizes)
}

func firstImageURL(images []models.ImageModel, variant string) string {
	for _, img := range images {
		if img.Variant == variant {
			return img.URL
		}
	}
	return ""
}

func queryInt(c *gin.Context, key string, fallback int) int {
	if v := c.Query(key); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			return n
		}
	}
	return fallback
}
