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

func (repo *Repo) TopTiles(limit int) ([]models.ArtModel, error) {
	var artTiles []models.ArtModel
	query := `
		SELECT id, title, description, portrait, url_low
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

func (repo *Repo) ListMoreTiles(limit, offset int) ([]models.ArtModel, error) {
	var artTiles []models.ArtModel
	query := `
		SELECT id, title, description, portrait, url_low
		FROM art_tiles
		ORDER BY display_order ASC, id ASC
		LIMIT ? OFFSET ?
	`

	err := repo.db.Select(&artTiles, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("SERVER: Could not list more art tiles: %w", err)
	}

	return artTiles, nil
}
