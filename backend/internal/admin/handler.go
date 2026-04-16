package admin

import (
	"context"
	"fmt"
	"html/template"
	"io"
	"io/fs"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"time"

	"cloud.google.com/go/storage"
	"github.com/HarkHorning/portfolio-go-svelte-azure-k8/internal/config"
	"github.com/HarkHorning/portfolio-go-svelte-azure-k8/internal/models"
	"github.com/HarkHorning/portfolio-go-svelte-azure-k8/internal/repo"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type Handler struct {
	repo          *repo.Repo
	tmplFS        fs.FS
	adminUsername string
	passwordHash  []byte
	gcsBucket     string
	gcsClient     *storage.Client
}

func NewHandler(r *repo.Repo, tmplFS fs.FS, cfg config.Config) (*Handler, error) {
	initStore(cfg.Admin.SessionSecret)

	hash, err := bcrypt.GenerateFromPassword([]byte(cfg.Admin.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("could not hash admin password: %w", err)
	}

	ctx := context.Background()
	gcsClient, err := storage.NewClient(ctx)
	if err != nil {
		slog.Warn("GCS client unavailable — image upload disabled", "error", err)
		gcsClient = nil
	}

	return &Handler{
		repo:          r,
		tmplFS:        tmplFS,
		adminUsername: cfg.Admin.Username,
		passwordHash:  hash,
		gcsBucket:     cfg.Storage.Bucket,
		gcsClient:     gcsClient,
	}, nil
}

func templateFuncs() template.FuncMap {
	return template.FuncMap{
		"dollars": func(cents int) string {
			return fmt.Sprintf("$%.2f", float64(cents)/100)
		},
		"deref": func(s *string) string {
			if s == nil {
				return ""
			}
			return *s
		},
		"derefInt": func(n *int) int {
			if n == nil {
				return 0
			}
			return *n
		},
		"not": func(b bool) bool { return !b },
		"dict": func(values ...any) (map[string]any, error) {
			if len(values)%2 != 0 {
				return nil, fmt.Errorf("dict requires even number of arguments")
			}
			m := make(map[string]any, len(values)/2)
			for i := 0; i < len(values); i += 2 {
				key, ok := values[i].(string)
				if !ok {
					return nil, fmt.Errorf("dict keys must be strings")
				}
				m[key] = values[i+1]
			}
			return m, nil
		},
		"hasCategory": func(cats []models.CategoryModel, id int) bool {
			for _, c := range cats {
				if c.Id == id {
					return true
				}
			}
			return false
		},
	}
}

func (h *Handler) render(c *gin.Context, name string, data any) {
	c.Header("Content-Type", "text/html; charset=utf-8")

	isPartial := name == "image_list.html" || name == "print_size_list.html" || name == "login.html"
	var files []string
	if isPartial {
		files = []string{"templates/" + name}
	} else {
		files = []string{"templates/layout.html", "templates/" + name}
		if name == "art_form.html" {
			files = append(files, "templates/image_list.html")
		}
		if name == "print_form.html" {
			files = append(files, "templates/print_size_list.html")
		}
	}

	tmpl, err := template.New("").Funcs(templateFuncs()).ParseFS(h.tmplFS, files...)
	if err != nil {
		slog.Error("template parse error", "template", name, "error", err)
		c.String(http.StatusInternalServerError, "template error")
		return
	}

	entry := "layout"
	if isPartial {
		entry = name
	}
	if err := tmpl.ExecuteTemplate(c.Writer, entry, data); err != nil {
		slog.Error("template render error", "template", name, "error", err)
		c.String(http.StatusInternalServerError, "render error")
	}
}

// ── Auth ─────────────────────────────────────────────────────────────────────

func (h *Handler) GetLogin(c *gin.Context) {
	h.render(c, "login.html", nil)
}

func (h *Handler) PostLogin(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")

	if username != h.adminUsername || bcrypt.CompareHashAndPassword(h.passwordHash, []byte(password)) != nil {
		h.render(c, "login.html", gin.H{"error": "Invalid credentials"})
		return
	}

	if err := setAuthenticated(c.Writer, c.Request); err != nil {
		c.String(http.StatusInternalServerError, "session error")
		return
	}
	c.Redirect(http.StatusSeeOther, "/admin/art")
}

func (h *Handler) PostLogout(c *gin.Context) {
	clearSession(c.Writer, c.Request)
	c.Redirect(http.StatusSeeOther, "/admin/login")
}

// ── Art ──────────────────────────────────────────────────────────────────────

func (h *Handler) GetArtList(c *gin.Context) {
	arts, err := h.repo.AdminAllArt()
	if err != nil {
		slog.Error("admin: list art", "error", err)
		c.String(http.StatusInternalServerError, "error loading art")
		return
	}
	cats, _ := h.repo.AllCategories()
	h.render(c, "art_list.html", gin.H{"Arts": arts, "Categories": cats})
}

func (h *Handler) GetArtNew(c *gin.Context) {
	cats, _ := h.repo.AllCategories()
	h.render(c, "art_form.html", gin.H{"Art": nil, "Categories": cats, "Action": "/admin/art"})
}

func (h *Handler) PostArtCreate(c *gin.Context) {
	title := strings.TrimSpace(c.PostForm("title"))
	description := strings.TrimSpace(c.PostForm("description"))
	portrait := c.PostForm("portrait") == "true"
	madeYear := optInt(c.PostForm("made_year"))
	size := optString(c.PostForm("size"))
	priceCents := optInt(c.PostForm("price_cents"))
	displayOrder := parseInt(c.PostForm("display_order"), 0)
	categoryIDs := parseIDs(c.PostFormArray("category_ids"))

	if title == "" {
		cats, _ := h.repo.AllCategories()
		h.render(c, "art_form.html", gin.H{"Error": "Title is required", "Categories": cats, "Action": "/admin/art"})
		return
	}

	visible := c.PostForm("visible") != "false"
	id, err := h.repo.AdminCreateArt(title, description, portrait, madeYear, size, priceCents, displayOrder, visible)
	if err != nil {
		slog.Error("admin: create art", "error", err)
		c.String(http.StatusInternalServerError, "create failed")
		return
	}

	if len(categoryIDs) > 0 {
		h.repo.AdminSetArtCategories(int(id), categoryIDs)
	}

	c.Redirect(http.StatusSeeOther, "/admin/art")
}

func (h *Handler) GetArtEdit(c *gin.Context) {
	id := paramInt(c, "id")
	art, err := h.repo.ArtByID(id)
	if err != nil {
		c.String(http.StatusNotFound, "not found")
		return
	}
	cats, _ := h.repo.AllCategories()
	h.render(c, "art_form.html", gin.H{
		"Art":        art,
		"Categories": cats,
		"Action":     fmt.Sprintf("/admin/art/%d", id),
	})
}

func (h *Handler) PostArtUpdate(c *gin.Context) {
	id := paramInt(c, "id")
	title := strings.TrimSpace(c.PostForm("title"))
	description := strings.TrimSpace(c.PostForm("description"))
	portrait := c.PostForm("portrait") == "true"
	madeYear := optInt(c.PostForm("made_year"))
	size := optString(c.PostForm("size"))
	priceCents := optInt(c.PostForm("price_cents"))
	sold := c.PostForm("sold") == "true"
	visible := c.PostForm("visible") != "false"
	categoryIDs := parseIDs(c.PostFormArray("category_ids"))

	if err := h.repo.AdminUpdateArt(id, title, description, portrait, madeYear, size, priceCents, sold, visible); err != nil {
		slog.Error("admin: update art", "error", err)
		c.String(http.StatusInternalServerError, "update failed")
		return
	}
	h.repo.AdminSetArtCategories(id, categoryIDs)
	c.Redirect(http.StatusSeeOther, "/admin/art")
}

func (h *Handler) PostArtArchive(c *gin.Context) {
	id := paramInt(c, "id")
	if err := h.repo.AdminArchiveArt(id); err != nil {
		c.String(http.StatusInternalServerError, "archive failed")
		return
	}
	c.Redirect(http.StatusSeeOther, "/admin/art")
}

// ── Images ───────────────────────────────────────────────────────────────────

func (h *Handler) PostImageUpload(c *gin.Context) {
	artID := paramInt(c, "id")
	variant := c.PostForm("variant")
	if variant != "high" && variant != "low" {
		c.String(http.StatusBadRequest, "variant must be high or low")
		return
	}

	file, header, err := c.Request.FormFile("image")
	if err != nil {
		c.String(http.StatusBadRequest, "no file uploaded")
		return
	}
	defer file.Close()

	if err := validateImageFile(file, header.Size); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	file.Seek(0, io.SeekStart)

	if h.gcsClient == nil {
		c.String(http.StatusServiceUnavailable, "image storage not configured")
		return
	}

	objectName := fmt.Sprintf("art/%s/%d-%d-%s", variant, artID, time.Now().UnixMilli(), header.Filename)
	url, err := h.uploadToGCS(c.Request.Context(), objectName, file, header.Header.Get("Content-Type"))
	if err != nil {
		slog.Error("admin: gcs upload", "error", err)
		c.String(http.StatusInternalServerError, "upload failed")
		return
	}

	sortOrder := 0
	existing, _ := h.repo.AdminImagesByArtID(artID)
	for _, img := range existing {
		if img.Variant == variant && img.SortOrder >= sortOrder {
			sortOrder = img.SortOrder + 1
		}
	}

	if _, err := h.repo.AdminAddImage(artID, variant, url, header.Filename, sortOrder); err != nil {
		slog.Error("admin: save image record", "error", err)
		c.String(http.StatusInternalServerError, "save failed")
		return
	}

	imgs, _ := h.repo.AdminImagesByArtID(artID)
	h.render(c, "image_list.html", gin.H{"Images": imgs, "ArtID": artID})
}

func (h *Handler) DeleteImage(c *gin.Context) {
	imageID := paramInt(c, "imageId")
	artID := paramInt(c, "id")

	filename, err := h.repo.AdminDeleteImage(imageID)
	if err != nil {
		c.String(http.StatusInternalServerError, "delete failed")
		return
	}

	if h.gcsClient != nil && filename != "" {
		_ = h.gcsClient.Bucket(h.gcsBucket).Object(filename).Delete(c.Request.Context())
	}

	imgs, _ := h.repo.AdminImagesByArtID(artID)
	h.render(c, "image_list.html", gin.H{"Images": imgs, "ArtID": artID})
}

func (h *Handler) uploadToGCS(ctx context.Context, objectName string, r io.Reader, contentType string) (string, error) {
	wc := h.gcsClient.Bucket(h.gcsBucket).Object(objectName).NewWriter(ctx)
	wc.ContentType = contentType
	if _, err := io.Copy(wc, r); err != nil {
		wc.Close()
		return "", err
	}
	if err := wc.Close(); err != nil {
		return "", err
	}
	return fmt.Sprintf("https://storage.googleapis.com/%s/%s", h.gcsBucket, objectName), nil
}

// ── Prints ───────────────────────────────────────────────────────────────────

func (h *Handler) GetPrintList(c *gin.Context) {
	prints, err := h.repo.AdminAllPrints()
	if err != nil {
		slog.Error("admin: list prints", "error", err)
		c.String(http.StatusInternalServerError, "error loading prints")
		return
	}
	h.render(c, "prints_list.html", gin.H{"Prints": prints})
}

func (h *Handler) GetPrintNew(c *gin.Context) {
	arts, _ := h.repo.AdminAllArt()
	h.render(c, "print_form.html", gin.H{"Print": nil, "Arts": arts})
}

func (h *Handler) PostPrintCreate(c *gin.Context) {
	artTileID := parseInt(c.PostForm("art_tile_id"), 0)
	if artTileID == 0 {
		arts, _ := h.repo.AdminAllArt()
		h.render(c, "print_form.html", gin.H{"Error": "Select an artwork", "Arts": arts})
		return
	}
	id, err := h.repo.AdminCreatePrint(artTileID)
	if err != nil {
		slog.Error("admin: create print", "error", err)
		c.String(http.StatusInternalServerError, "create failed")
		return
	}
	c.Redirect(http.StatusSeeOther, fmt.Sprintf("/admin/prints/%d/edit", id))
}

func (h *Handler) GetPrintEdit(c *gin.Context) {
	id := paramInt(c, "id")
	p, err := h.repo.PrintByID(id)
	if err != nil {
		c.String(http.StatusNotFound, "not found")
		return
	}
	h.render(c, "print_form.html", gin.H{"Print": p})
}

func (h *Handler) PostPrintArchive(c *gin.Context) {
	id := paramInt(c, "id")
	if err := h.repo.AdminArchivePrint(id); err != nil {
		c.String(http.StatusInternalServerError, "archive failed")
		return
	}
	c.Redirect(http.StatusSeeOther, "/admin/prints")
}

func (h *Handler) PostPrintSizeAdd(c *gin.Context) {
	printID := paramInt(c, "id")
	size := strings.TrimSpace(c.PostForm("size"))
	priceCents := parseInt(c.PostForm("price_cents"), 0)
	qty := parseInt(c.PostForm("quantity_in_stock"), 0)

	if size != "" {
		if err := h.repo.AdminAddPrintSize(printID, size, priceCents, qty); err != nil {
			slog.Error("admin: add print size", "error", err)
		}
	}
	h.renderPrintSizeList(c, printID)
}

func (h *Handler) PostPrintSizeDelete(c *gin.Context) {
	printID := paramInt(c, "id")
	sizeID := paramInt(c, "psid")
	if err := h.repo.AdminArchivePrintSize(sizeID); err != nil {
		slog.Error("admin: archive print size", "error", err)
	}
	h.renderPrintSizeList(c, printID)
}

func (h *Handler) PostPrintSizeUpdate(c *gin.Context) {
	printID := paramInt(c, "id")
	sizeID := paramInt(c, "psid")
	size := strings.TrimSpace(c.PostForm("size"))
	priceCents := parseInt(c.PostForm("price_cents"), 0)
	qty := parseInt(c.PostForm("quantity_in_stock"), 0)
	sold := c.PostForm("sold") == "true"
	if err := h.repo.AdminUpdatePrintSize(sizeID, size, priceCents, qty, sold); err != nil {
		slog.Error("admin: update print size", "error", err)
	}
	h.renderPrintSizeList(c, printID)
}

func (h *Handler) renderPrintSizeList(c *gin.Context, printID int) {
	sizes, _ := h.repo.AdminPrintSizesByPrint(printID)
	h.render(c, "print_size_list.html", gin.H{"PrintID": printID, "Sizes": sizes})
}

// ── Categories ───────────────────────────────────────────────────────────────

func (h *Handler) GetCategories(c *gin.Context) {
	cats, err := h.repo.AllCategories()
	if err != nil {
		c.String(http.StatusInternalServerError, "error loading categories")
		return
	}
	h.render(c, "categories.html", gin.H{"Categories": cats})
}

func (h *Handler) PostCategoryCreate(c *gin.Context) {
	name := strings.TrimSpace(c.PostForm("name"))
	slug := strings.TrimSpace(c.PostForm("slug"))
	if name == "" || slug == "" {
		cats, _ := h.repo.AllCategories()
		h.render(c, "categories.html", gin.H{"Categories": cats, "Error": "Name and slug are required"})
		return
	}
	if err := h.repo.AdminCreateCategory(name, slug); err != nil {
		c.String(http.StatusInternalServerError, "create failed")
		return
	}
	c.Redirect(http.StatusSeeOther, "/admin/categories")
}

func (h *Handler) PostCategoryDelete(c *gin.Context) {
	id := paramInt(c, "id")
	if err := h.repo.AdminDeleteCategory(id); err != nil {
		c.String(http.StatusInternalServerError, "delete failed")
		return
	}
	c.Redirect(http.StatusSeeOther, "/admin/categories")
}

// ── Helpers ──────────────────────────────────────────────────────────────────

func paramInt(c *gin.Context, key string) int {
	n, _ := strconv.Atoi(c.Param(key))
	return n
}

func parseInt(s string, fallback int) int {
	n, err := strconv.Atoi(strings.TrimSpace(s))
	if err != nil {
		return fallback
	}
	return n
}

func optInt(s string) *int {
	s = strings.TrimSpace(s)
	if s == "" {
		return nil
	}
	n, err := strconv.Atoi(s)
	if err != nil {
		return nil
	}
	return &n
}

func optString(s string) *string {
	s = strings.TrimSpace(s)
	if s == "" {
		return nil
	}
	return &s
}

func parseIDs(ss []string) []int {
	var ids []int
	for _, s := range ss {
		if n, err := strconv.Atoi(s); err == nil {
			ids = append(ids, n)
		}
	}
	return ids
}

func validateImageFile(r io.Reader, size int64) error {
	const maxSize = 20 << 20 // 20 MB
	if size > maxSize {
		return fmt.Errorf("file too large (max 20 MB)")
	}
	magic := make([]byte, 12)
	if _, err := io.ReadFull(r, magic); err != nil {
		return fmt.Errorf("could not read file")
	}
	if isJPEG(magic) || isPNG(magic) || isWEBP(magic) {
		return nil
	}
	return fmt.Errorf("unsupported file type: only JPEG, PNG, and WebP are allowed")
}

func isJPEG(b []byte) bool { return len(b) >= 3 && b[0] == 0xFF && b[1] == 0xD8 && b[2] == 0xFF }
func isPNG(b []byte) bool {
	return len(b) >= 8 && b[0] == 0x89 && b[1] == 0x50 && b[2] == 0x4E && b[3] == 0x47
}
func isWEBP(b []byte) bool {
	return len(b) >= 12 && string(b[0:4]) == "RIFF" && string(b[8:12]) == "WEBP"
}

func gcsObjectName(url, bucket string) string {
	prefix := fmt.Sprintf("https://storage.googleapis.com/%s/", bucket)
	return strings.TrimPrefix(url, prefix)
}
