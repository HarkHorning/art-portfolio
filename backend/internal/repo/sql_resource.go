package repo

import (
	"fmt"

	"github.com/HarkHorning/portfolio-go-svelte-azure-k8/internal/models"
	"github.com/jmoiron/sqlx"
)

type Repo struct {
	db *sqlx.DB
}

func NewRepo(db *sqlx.DB) *Repo {
	return &Repo{db: db}
}

func (repo *Repo) Ping() error {
	return repo.db.Ping()
}

const displayURLSubquery = `(
	SELECT url FROM images
	WHERE art_tile_id = at.id AND variant = 'low'
	ORDER BY sort_order ASC LIMIT 1
) AS display_url`

// ArtTiles returns artwork filtered by any combination of category, size, and price range.
// Empty string means no filter for category/size; -1 means no filter for price.
func (repo *Repo) ArtTiles(category, size string, minPrice, maxPrice int) ([]models.ArtModel, error) {
	artTiles := make([]models.ArtModel, 0)
	args := make([]any, 0)

	var query string
	if category != "" {
		query = fmt.Sprintf(`
			SELECT at.id, at.title, at.description, at.portrait, %s,
			       at.made_year, at.sold, at.size, at.price_cents
			FROM art_tiles at
			JOIN art_categories ac ON at.id = ac.art_id
			JOIN categories c ON ac.category_id = c.id
			WHERE c.slug = ? AND at.archived_at IS NULL`, displayURLSubquery)
		args = append(args, category)
	} else {
		query = fmt.Sprintf(`
			SELECT at.id, at.title, at.description, at.portrait, %s,
			       at.made_year, at.sold, at.size, at.price_cents
			FROM art_tiles at
			WHERE at.archived_at IS NULL`, displayURLSubquery)
	}

	if size != "" {
		query += ` AND at.size = ?`
		args = append(args, size)
	}
	if minPrice >= 0 {
		query += ` AND at.price_cents >= ?`
		args = append(args, minPrice)
	}
	if maxPrice >= 0 {
		query += ` AND at.price_cents <= ?`
		args = append(args, maxPrice)
	}

	query += ` ORDER BY at.display_order ASC, at.id ASC`

	err := repo.db.Select(&artTiles, query, args...)
	if err != nil {
		return nil, fmt.Errorf("could not list art tiles: %w", err)
	}

	return artTiles, nil
}

func (repo *Repo) ArtSizes() ([]string, error) {
	var sizes []string
	err := repo.db.Select(&sizes, `
		SELECT DISTINCT size FROM art_tiles
		WHERE size IS NOT NULL AND archived_at IS NULL
		ORDER BY size ASC`)
	if err != nil {
		return nil, fmt.Errorf("could not list art sizes: %w", err)
	}
	return sizes, nil
}

func (repo *Repo) ArtByID(id int) (*models.ArtDetailModel, error) {
	var art models.ArtDetailModel
	err := repo.db.Get(&art, `
		SELECT id, title, description, portrait, made_year, sold, size, price_cents
		FROM art_tiles
		WHERE id = ? AND archived_at IS NULL
	`, id)
	if err != nil {
		return nil, fmt.Errorf("art not found: %w", err)
	}

	images := make([]models.ImageModel, 0)
	err = repo.db.Select(&images, `
		SELECT id, art_tile_id, variant, url, filename, sort_order
		FROM images
		WHERE art_tile_id = ?
		ORDER BY variant ASC, sort_order ASC
	`, id)
	if err != nil {
		return nil, fmt.Errorf("could not get images for art: %w", err)
	}

	categories := make([]models.CategoryModel, 0)
	err = repo.db.Select(&categories, `
		SELECT c.id, c.name, c.slug
		FROM categories c
		JOIN art_categories ac ON c.id = ac.category_id
		WHERE ac.art_id = ?
		ORDER BY c.name ASC
	`, id)
	if err != nil {
		return nil, fmt.Errorf("could not get categories for art: %w", err)
	}

	art.Images = images
	art.Categories = categories
	return &art, nil
}

func (repo *Repo) AllCategories() ([]models.CategoryModel, error) {
	categories := make([]models.CategoryModel, 0)
	err := repo.db.Select(&categories, `SELECT id, name, slug FROM categories ORDER BY name ASC`)
	if err != nil {
		return nil, fmt.Errorf("could not list categories: %w", err)
	}
	return categories, nil
}
