package repo

import (
	"fmt"

	"github.com/HarkHorning/portfolio-go-svelte-azure-k8/internal/models"
)

// ── Art ──────────────────────────────────────────────────────────────────────

func (repo *Repo) AdminAllArt() ([]models.ArtDetailModel, error) {
	rows, err := repo.db.Queryx(`
		SELECT id, title, description, portrait, made_year, sold, visible, size, price_cents
		FROM art_tiles
		WHERE archived_at IS NULL
		ORDER BY display_order ASC, id ASC
	`)
	if err != nil {
		return nil, fmt.Errorf("could not list art: %w", err)
	}
	defer rows.Close()

	var arts []models.ArtDetailModel
	for rows.Next() {
		var a models.ArtDetailModel
		if err := rows.StructScan(&a); err != nil {
			return nil, err
		}
		arts = append(arts, a)
	}

	for i := range arts {
		imgs := make([]models.ImageModel, 0)
		if err := repo.db.Select(&imgs, `
			SELECT id, art_tile_id, variant, url, filename, sort_order
			FROM images WHERE art_tile_id = ? ORDER BY variant ASC, sort_order ASC
		`, arts[i].Id); err != nil {
			return nil, err
		}
		arts[i].Images = imgs

		cats := make([]models.CategoryModel, 0)
		if err := repo.db.Select(&cats, `
			SELECT c.id, c.name, c.slug FROM categories c
			JOIN art_categories ac ON c.id = ac.category_id
			WHERE ac.art_id = ? ORDER BY c.name ASC
		`, arts[i].Id); err != nil {
			return nil, err
		}
		arts[i].Categories = cats
	}

	return arts, nil
}

func (repo *Repo) AdminCreateArt(title, description string, portrait bool, madeYear *int, size *string, priceCents *int, displayOrder int, visible bool) (int64, error) {
	res, err := repo.db.Exec(`
		INSERT INTO art_tiles (title, description, portrait, made_year, size, price_cents, display_order, visible)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`, title, description, portrait, madeYear, size, priceCents, displayOrder, visible)
	if err != nil {
		return 0, fmt.Errorf("could not create art: %w", err)
	}
	return res.LastInsertId()
}

func (repo *Repo) AdminUpdateArt(id int, title, description string, portrait bool, madeYear *int, size *string, priceCents *int, sold, visible bool) error {
	_, err := repo.db.Exec(`
		UPDATE art_tiles
		SET title=?, description=?, portrait=?, made_year=?, size=?, price_cents=?, sold=?, visible=?
		WHERE id=?
	`, title, description, portrait, madeYear, size, priceCents, sold, visible, id)
	return err
}

func (repo *Repo) AdminArchiveArt(id int) error {
	_, err := repo.db.Exec(`UPDATE art_tiles SET archived_at = NOW() WHERE id = ?`, id)
	return err
}

func (repo *Repo) AdminSetArtCategories(artID int, categoryIDs []int) error {
	tx, err := repo.db.Begin()
	if err != nil {
		return err
	}
	if _, err := tx.Exec(`DELETE FROM art_categories WHERE art_id = ?`, artID); err != nil {
		tx.Rollback()
		return err
	}
	for _, catID := range categoryIDs {
		if _, err := tx.Exec(`INSERT INTO art_categories (art_id, category_id) VALUES (?, ?)`, artID, catID); err != nil {
			tx.Rollback()
			return err
		}
	}
	return tx.Commit()
}

// ── Images ───────────────────────────────────────────────────────────────────

func (repo *Repo) AdminAddImage(artTileID int, variant, url, filename string, sortOrder int) (int64, error) {
	res, err := repo.db.Exec(`
		INSERT INTO images (art_tile_id, variant, url, filename, sort_order)
		VALUES (?, ?, ?, ?, ?)
	`, artTileID, variant, url, filename, sortOrder)
	if err != nil {
		return 0, fmt.Errorf("could not add image: %w", err)
	}
	return res.LastInsertId()
}

func (repo *Repo) AdminDeleteImage(imageID int) (string, error) {
	var filename string
	if err := repo.db.Get(&filename, `SELECT filename FROM images WHERE id = ?`, imageID); err != nil {
		return "", fmt.Errorf("image not found: %w", err)
	}
	if _, err := repo.db.Exec(`DELETE FROM images WHERE id = ?`, imageID); err != nil {
		return "", fmt.Errorf("could not delete image: %w", err)
	}
	return filename, nil
}

func (repo *Repo) AdminImagesByArtID(artID int) ([]models.ImageModel, error) {
	imgs := make([]models.ImageModel, 0)
	err := repo.db.Select(&imgs, `
		SELECT id, art_tile_id, variant, url, filename, sort_order
		FROM images WHERE art_tile_id = ?
		ORDER BY variant ASC, sort_order ASC
	`, artID)
	return imgs, err
}

// ── Prints ───────────────────────────────────────────────────────────────────

func (repo *Repo) AdminAllPrints() ([]models.PrintModel, error) {
	prints := make([]models.PrintModel, 0)
	if err := repo.db.Select(&prints, fmt.Sprintf(`
		SELECT p.id, p.art_tile_id, at.title, at.description, at.portrait, %s
		FROM prints p
		JOIN art_tiles at ON p.art_tile_id = at.id
		WHERE p.archived_at IS NULL
		ORDER BY p.id ASC
	`, printDisplayURLSubquery)); err != nil {
		return nil, err
	}
	if err := loadPrintSizes(repo, prints); err != nil {
		return nil, err
	}
	return prints, nil
}

func (repo *Repo) AdminCreatePrint(artTileID int) (int64, error) {
	res, err := repo.db.Exec(`INSERT INTO prints (art_tile_id) VALUES (?)`, artTileID)
	if err != nil {
		return 0, fmt.Errorf("could not create print: %w", err)
	}
	return res.LastInsertId()
}

func (repo *Repo) AdminArchivePrint(id int) error {
	_, err := repo.db.Exec(`UPDATE prints SET archived_at = NOW() WHERE id = ?`, id)
	return err
}

func (repo *Repo) AdminPrintSizesByPrint(printID int) ([]models.PrintSizeModel, error) {
	sizes := make([]models.PrintSizeModel, 0)
	err := repo.db.Select(&sizes, `
		SELECT id, print_id, size, price_cents, quantity_in_stock, sold
		FROM print_sizes WHERE print_id = ? AND archived_at IS NULL
		ORDER BY price_cents ASC`, printID)
	return sizes, err
}

func (repo *Repo) AdminAddPrintSize(printID int, size string, priceCents, qty int) error {
	_, err := repo.db.Exec(`
		INSERT INTO print_sizes (print_id, size, price_cents, quantity_in_stock)
		VALUES (?, ?, ?, ?)
	`, printID, size, priceCents, qty)
	return err
}

func (repo *Repo) AdminUpdatePrintSize(id int, size string, priceCents, qty int, sold bool) error {
	_, err := repo.db.Exec(`
		UPDATE print_sizes SET size=?, price_cents=?, quantity_in_stock=?, sold=? WHERE id=?
	`, size, priceCents, qty, sold, id)
	return err
}

func (repo *Repo) AdminArchivePrintSize(id int) error {
	_, err := repo.db.Exec(`UPDATE print_sizes SET archived_at = NOW() WHERE id = ?`, id)
	return err
}

// ── Categories ───────────────────────────────────────────────────────────────

func (repo *Repo) AdminCreateCategory(name, slug string) error {
	_, err := repo.db.Exec(`INSERT INTO categories (name, slug) VALUES (?, ?)`, name, slug)
	return err
}

func (repo *Repo) AdminDeleteCategory(id int) error {
	_, err := repo.db.Exec(`DELETE FROM categories WHERE id = ?`, id)
	return err
}
