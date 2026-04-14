package repo

import (
	"fmt"

	"github.com/HarkHorning/portfolio-go-svelte-azure-k8/internal/models"
)

// ── Art ──────────────────────────────────────────────────────────────────────

func (repo *Repo) AdminAllArt() ([]models.ArtDetailModel, error) {
	rows, err := repo.db.Queryx(`
		SELECT id, title, description, portrait, made_year, sold, size, price_cents
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

func (repo *Repo) AdminCreateArt(title, description string, portrait bool, madeYear *int, size *string, priceCents *int, displayOrder int) (int64, error) {
	res, err := repo.db.Exec(`
		INSERT INTO art_tiles (title, description, portrait, made_year, size, price_cents, display_order)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`, title, description, portrait, madeYear, size, priceCents, displayOrder)
	if err != nil {
		return 0, fmt.Errorf("could not create art: %w", err)
	}
	return res.LastInsertId()
}

func (repo *Repo) AdminUpdateArt(id int, title, description string, portrait bool, madeYear *int, size *string, priceCents *int, sold bool) error {
	_, err := repo.db.Exec(`
		UPDATE art_tiles
		SET title=?, description=?, portrait=?, made_year=?, size=?, price_cents=?, sold=?
		WHERE id=?
	`, title, description, portrait, madeYear, size, priceCents, sold, id)
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
	err := repo.db.Select(&prints, fmt.Sprintf(`
		SELECT p.id, p.art_tile_id, at.title, at.description, at.portrait, %s,
		       p.price_cents, p.size, p.sold, p.quantity_in_stock
		FROM prints p
		JOIN art_tiles at ON p.art_tile_id = at.id
		WHERE p.archived_at IS NULL
		ORDER BY p.id ASC
	`, printDisplayURLSubquery))
	return prints, err
}

func (repo *Repo) AdminCreatePrint(artTileID int, size string, priceCents int, quantity int) (int64, error) {
	res, err := repo.db.Exec(`
		INSERT INTO prints (art_tile_id, size, price_cents, quantity_in_stock)
		VALUES (?, ?, ?, ?)
	`, artTileID, size, priceCents, quantity)
	if err != nil {
		return 0, fmt.Errorf("could not create print: %w", err)
	}
	return res.LastInsertId()
}

func (repo *Repo) AdminUpdatePrint(id, priceCents, quantity int, size string, sold bool) error {
	_, err := repo.db.Exec(`
		UPDATE prints SET price_cents=?, quantity_in_stock=?, size=?, sold=? WHERE id=?
	`, priceCents, quantity, size, sold, id)
	return err
}

func (repo *Repo) AdminArchivePrint(id int) error {
	_, err := repo.db.Exec(`UPDATE prints SET archived_at = NOW() WHERE id = ?`, id)
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
