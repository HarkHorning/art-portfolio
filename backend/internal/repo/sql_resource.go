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

func (repo *Repo) TopTiles(limit int) ([]models.ArtModel, error) {
	artTiles := make([]models.ArtModel, 0)
	query := `
		SELECT id, title, description, portrait, url_low, made_year, sold
		FROM art_tiles
		ORDER BY display_order ASC, id ASC
		LIMIT ?
	`

	err := repo.db.Select(&artTiles, query, limit)
	if err != nil {
		return nil, fmt.Errorf("SERVER: Could not list art tiles: %w", err)
	}

	return artTiles, nil
}

func (repo *Repo) TilesByCategory(slug string) ([]models.ArtModel, error) {
	artTiles := make([]models.ArtModel, 0)
	query := `
		SELECT at.id, at.title, at.description, at.portrait, at.url_low, at.made_year, at.sold
		FROM art_tiles at
		JOIN art_categories ac ON at.id = ac.art_id
		JOIN categories c ON ac.category_id = c.id
		WHERE c.slug = ?
		ORDER BY at.display_order ASC, at.id ASC
	`

	err := repo.db.Select(&artTiles, query, slug)
	if err != nil {
		return nil, fmt.Errorf("could not list tiles by category: %w", err)
	}

	return artTiles, nil
}

func (repo *Repo) ArtByID(id int) (*models.ArtDetailModel, error) {
	var art models.ArtModel
	err := repo.db.Get(&art, `
		SELECT id, title, description, portrait, url_low, made_year, sold
		FROM art_tiles
		WHERE id = ?
	`, id)
	if err != nil {
		return nil, fmt.Errorf("art not found: %w", err)
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

	return &models.ArtDetailModel{
		Id:          art.Id,
		Title:       art.Title,
		Description: art.Description,
		Portrait:    art.Portrait,
		URL:         art.URL,
		MadeYear:    art.MadeYear,
		Sold:        art.Sold,
		Categories:  categories,
	}, nil
}

func (repo *Repo) AllCategories() ([]models.CategoryModel, error) {
	categories := make([]models.CategoryModel, 0)
	query := `SELECT id, name, slug FROM categories ORDER BY name ASC`

	err := repo.db.Select(&categories, query)
	if err != nil {
		return nil, fmt.Errorf("could not list categories: %w", err)
	}

	return categories, nil
}
