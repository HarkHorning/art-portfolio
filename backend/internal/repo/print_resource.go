package repo

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/HarkHorning/portfolio-go-svelte-azure-k8/internal/models"
)

const printDisplayURLSubquery = `COALESCE((
	SELECT url FROM images
	WHERE art_tile_id = p.art_tile_id AND variant = 'low'
	ORDER BY sort_order ASC LIMIT 1
), '') AS display_url`

func (repo *Repo) Prints(size string, minPrice, maxPrice int) ([]models.PrintModel, error) {
	prints := make([]models.PrintModel, 0)

	query := fmt.Sprintf(`
		SELECT p.id, p.art_tile_id, at.title, at.description, at.portrait, %s,
		       p.price_cents, p.size, p.sold, p.quantity_in_stock
		FROM prints p
		JOIN art_tiles at ON p.art_tile_id = at.id
		WHERE p.archived_at IS NULL AND at.archived_at IS NULL`, printDisplayURLSubquery)

	args := make([]any, 0)

	if size != "" {
		query += ` AND p.size = ?`
		args = append(args, size)
	}
	if minPrice >= 0 {
		query += ` AND p.price_cents >= ?`
		args = append(args, minPrice)
	}
	if maxPrice >= 0 {
		query += ` AND p.price_cents <= ?`
		args = append(args, maxPrice)
	}

	query += ` ORDER BY p.id ASC`

	err := repo.db.Select(&prints, query, args...)
	if err != nil {
		return nil, fmt.Errorf("could not list prints: %w", err)
	}

	return prints, nil
}

func (repo *Repo) PrintByID(id int) (*models.PrintModel, error) {
	var p models.PrintModel
	err := repo.db.Get(&p, fmt.Sprintf(`
		SELECT p.id, p.art_tile_id, at.title, at.description, at.portrait, %s,
		       p.price_cents, p.size, p.sold, p.quantity_in_stock
		FROM prints p
		JOIN art_tiles at ON p.art_tile_id = at.id
		WHERE p.id = ? AND p.archived_at IS NULL
	`, printDisplayURLSubquery), id)
	if err != nil {
		return nil, fmt.Errorf("print not found: %w", err)
	}
	return &p, nil
}

func (repo *Repo) PrintSizes() ([]string, error) {
	var sizes []string
	err := repo.db.Select(&sizes, `
		SELECT DISTINCT size FROM prints
		WHERE archived_at IS NULL
		ORDER BY size ASC`)
	if err != nil {
		return nil, fmt.Errorf("could not list print sizes: %w", err)
	}

	sort.Slice(sizes, func(i, j int) bool {
		return firstDim(sizes[i]) < firstDim(sizes[j])
	})

	return sizes, nil
}

func firstDim(size string) int {
	parts := strings.SplitN(size, "x", 2)
	if len(parts) == 0 {
		return 0
	}
	n, _ := strconv.Atoi(strings.TrimSpace(parts[0]))
	return n
}
