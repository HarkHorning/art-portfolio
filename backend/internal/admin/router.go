package admin

import (
	"io/fs"
	"net/http"

	"github.com/HarkHorning/portfolio-go-svelte-azure-k8/internal/config"
	"github.com/HarkHorning/portfolio-go-svelte-azure-k8/internal/repo"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(rg *gin.RouterGroup, r *repo.Repo, tmplFS fs.FS, cfg config.Config) error {
	h, err := NewHandler(r, tmplFS, cfg)
	if err != nil {
		return err
	}

	rg.GET("/login", h.GetLogin)
	rg.POST("/login", h.PostLogin)

	auth := rg.Group("/")
	auth.Use(requireAuth())
	{
		auth.POST("/logout", h.PostLogout)

		auth.GET("", func(c *gin.Context) { c.Redirect(http.StatusFound, "/admin/art") })

		auth.GET("/art", h.GetArtList)
		auth.GET("/art/new", h.GetArtNew)
		auth.POST("/art", h.PostArtCreate)
		auth.GET("/art/:id/edit", h.GetArtEdit)
		auth.POST("/art/:id", h.PostArtUpdate)
		auth.POST("/art/:id/archive", h.PostArtArchive)
		auth.POST("/art/:id/publish", h.PostArtTogglePublish)
		auth.POST("/art/:id/images", h.PostImageUpload)
		auth.POST("/art/:id/images/:imageId/delete", h.DeleteImage)

		auth.GET("/prints", h.GetPrintList)
		auth.GET("/prints/new", h.GetPrintNew)
		auth.POST("/prints", h.PostPrintCreate)
		auth.GET("/prints/:id/edit", h.GetPrintEdit)
		auth.POST("/prints/:id/archive", h.PostPrintArchive)
		auth.POST("/prints/:id/publish", h.PostPrintTogglePublish)
		auth.POST("/prints/:id/sizes", h.PostPrintSizeAdd)
		auth.POST("/prints/:id/sizes/:psid", h.PostPrintSizeUpdate)
		auth.POST("/prints/:id/sizes/:psid/delete", h.PostPrintSizeDelete)

		auth.GET("/categories", h.GetCategories)
		auth.POST("/categories", h.PostCategoryCreate)
		auth.POST("/categories/:id/delete", h.PostCategoryDelete)
	}

	return nil
}

func requireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !isAuthenticated(c.Request) {
			c.Redirect(http.StatusFound, "/admin/login")
			c.Abort()
			return
		}
		c.Next()
	}
}
