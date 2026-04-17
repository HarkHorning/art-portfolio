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
		SELECT p.id, p.art_tile_id, at.title, at.description, at.portrait, %s, p.visible
		FROM prints p
		JOIN art_tiles at ON p.art_tile_id = at.id
		WHERE p.archived_at IS NULL AND at.archived_at IS NULL AND p.visible = TRUE`,
		printDisplayURLSubquery)

	args := make([]any, 0)

	if size != "" || minPrice >= 0 || maxPrice >= 0 {
		sub := `EXISTS (
			SELECT 1 FROM print_sizes ps
			WHERE ps.print_id = p.id AND ps.archived_at IS NULL`
		if size != "" {
			sub += ` AND ps.size = ?`
			args = append(args, size)
		}
		if minPrice >= 0 {
			sub += ` AND ps.price_cents >= ?`
			args = append(args, minPrice)
		}
		if maxPrice >= 0 {
			sub += ` AND ps.price_cents <= ?`
			args = append(args, maxPrice)
		}
		sub += `)`
		query += ` AND ` + sub
	}

	query += ` ORDER BY p.id ASC`

	if err := repo.db.Select(&prints, query, args...); err != nil {
		return nil, fmt.Errorf("could not list prints: %w", err)
	}
	if err := loadPrintSizes(repo, prints); err != nil {
		return nil, err
	}
	return prints, nil
}

func (repo *Repo) PrintByID(id int) (*models.PrintModel, error) {
	var p models.PrintModel
	err := repo.db.Get(&p, fmt.Sprintf(`
		SELECT p.id, p.art_tile_id, at.title, at.description, at.portrait, %s, p.visible
		FROM prints p
		JOIN art_tiles at ON p.art_tile_id = at.id
		WHERE p.id = ? AND p.archived_at IS NULL
	`, printDisplayURLSubquery), id)
	if err != nil {
		return nil, fmt.Errorf("print not found: %w", err)
	}
	if err := repo.db.Select(&p.Sizes, `
		SELECT id, print_id, size, price_cents, quantity_in_stock, sold
		FROM print_sizes WHERE print_id = ? AND archived_at IS NULL
		ORDER BY price_cents ASC`, id); err != nil {
		return nil, fmt.Errorf("could not load print sizes: %w", err)
	}

	// Load display images: print selection → art selection → all art images
	var printImgIDs []int
	repo.db.Select(&printImgIDs, `SELECT image_id FROM print_display_images WHERE print_id = ?`, id)
	if len(printImgIDs) > 0 {
		p.Images = loadImagesByIDs(repo, printImgIDs)
	} else {
		var artImgIDs []int
		repo.db.Select(&artImgIDs, `SELECT image_id FROM art_display_images WHERE art_tile_id = ?`, p.ArtTileId)
		if len(artImgIDs) > 0 {
			p.Images = loadImagesByIDs(repo, artImgIDs)
		} else {
			imgs := make([]models.ImageModel, 0)
			repo.db.Select(&imgs, `
				SELECT id, art_tile_id, variant, url, filename, sort_order
				FROM images WHERE art_tile_id = ? AND variant = 'low'
				ORDER BY sort_order ASC`, p.ArtTileId)
			p.Images = imgs
		}
	}

	return &p, nil
}

func loadImagesByIDs(repo *Repo, ids []int) []models.ImageModel {
	placeholders := make([]string, len(ids))
	args := make([]any, len(ids))
	for i, id := range ids {
		placeholders[i] = "?"
		args[i] = id
	}
	imgs := make([]models.ImageModel, 0)
	repo.db.Select(&imgs, fmt.Sprintf(`
		SELECT id, art_tile_id, variant, url, filename, sort_order
		FROM images WHERE id IN (%s)
		ORDER BY sort_order ASC`, strings.Join(placeholders, ",")), args...)
	return imgs
}

func (repo *Repo) PrintSizes() ([]string, error) {
	var sizes []string
	err := repo.db.Select(&sizes, `
		SELECT DISTINCT size FROM print_sizes
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

// loadPrintSizes fetches all sizes for the given prints in one query
// and assigns them in-place via index into the backing array.
func loadPrintSizes(repo *Repo, prints []models.PrintModel) error {
	if len(prints) == 0 {
		return nil
	}
	idIndex := make(map[int]int, len(prints))
	placeholders := make([]string, len(prints))
	args := make([]any, len(prints))
	for i, p := range prints {
		idIndex[p.Id] = i
		placeholders[i] = "?"
		args[i] = p.Id
		prints[i].Sizes = []models.PrintSizeModel{}
	}
	query := fmt.Sprintf(`
		SELECT id, print_id, size, price_cents, quantity_in_stock, sold
		FROM print_sizes
		WHERE print_id IN (%s) AND archived_at IS NULL
		ORDER BY price_cents ASC`, strings.Join(placeholders, ","))

	var sizes []models.PrintSizeModel
	if err := repo.db.Select(&sizes, query, args...); err != nil {
		return fmt.Errorf("could not load print sizes: %w", err)
	}
	for _, s := range sizes {
		if idx, ok := idIndex[s.PrintId]; ok {
			prints[idx].Sizes = append(prints[idx].Sizes, s)
		}
	}
	return nil
}

func firstDim(size string) int {
	parts := strings.SplitN(size, "x", 2)
	if len(parts) == 0 {
		return 0
	}
	n, _ := strconv.Atoi(strings.TrimSpace(parts[0]))
	return n
}
