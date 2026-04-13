package repo

import (
	"fmt"
	"log/slog"

	"github.com/jmoiron/sqlx"
)

func SeedDevData(db *sqlx.DB) error {
	slog.Info("seeding development data")

	tables := []string{"art_categories", "art_tiles", "categories"}
	for _, table := range tables {
		_, err := db.Exec(fmt.Sprintf("DELETE FROM %s", table))
		if err != nil {
			return fmt.Errorf("failed to clear %s: %w", table, err)
		}
	}

	for _, table := range []string{"art_tiles", "categories"} {
		_, _ = db.Exec(fmt.Sprintf("ALTER TABLE %s AUTO_INCREMENT = 1", table))
	}

	if err := seedCategories(db); err != nil {
		return err
	}
	if err := seedArtTiles(db); err != nil {
		return err
	}
	if err := seedArtCategories(db); err != nil {
		return err
	}

	slog.Info("development data seeded")
	return nil
}

func seedCategories(db *sqlx.DB) error {
	query := `
		INSERT INTO categories (name, slug) VALUES
		('Oil', 'oil'),
		('Acrylic', 'acrylic'),
		('Watercolor', 'watercolor'),
		('Pencil Drawing', 'pencil-drawing'),
		('Mixed', 'mixed'),
		('Pastel', 'pastel'),
		('Misc', 'misc')
	`
	_, err := db.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to seed categories: %w", err)
	}
	slog.Debug("seeded", "table", "categories")
	return nil
}

func seedArtTiles(db *sqlx.DB) error {
	query := `
		INSERT INTO art_tiles (title, description, portrait, url_low, url_high, display_order, made_year, sold) VALUES
		('Woman with Flowers', 'Acrylic on canvas', TRUE,
		 'https://storage.googleapis.com/hark-portfolio-images/art/Woman%20With%20Flowers.jpeg',
		 'https://storage.googleapis.com/hark-portfolio-images/art/Woman%20With%20Flowers.jpeg', 1, 2024, FALSE),
		('Boat on Lake', 'Oil on canvas, peaceful morning scene', FALSE,
		 'https://storage.googleapis.com/hark-portfolio-images/art/Boat%20on%20Lake.jpeg',
		 'https://storage.googleapis.com/hark-portfolio-images/art/Boat%20on%20Lake.jpeg', 2, 2023, FALSE),
		('Horse Statue', 'Watercolor', TRUE,
		 'https://storage.googleapis.com/hark-portfolio-images/art/Horse%20Statue.jpeg',
		 'https://storage.googleapis.com/hark-portfolio-images/art/Horse%20Statue.jpeg', 3, 2023, FALSE),
		('Cardinal', 'Watercolor', TRUE,
		 'https://storage.googleapis.com/hark-portfolio-images/art/Cardinal.jpeg',
		 'https://storage.googleapis.com/hark-portfolio-images/art/Cardinal.jpeg', 4, 2023, FALSE),
		('Shoebill Stork', 'Watercolor', TRUE,
		 'https://storage.googleapis.com/hark-portfolio-images/art/Shoebill.jpeg',
		 'https://storage.googleapis.com/hark-portfolio-images/art/Shoebill.jpeg', 5, 2022, FALSE),
		('Boat on Lake', 'Golden hour landscape', FALSE,
		 'https://storage.googleapis.com/hark-portfolio-images/art/Boat%20on%20Lake.jpeg',
		 'https://storage.googleapis.com/hark-portfolio-images/art/Boat%20on%20Lake.jpeg', 6, 2022, TRUE)
	`
	_, err := db.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to seed art_tiles: %w", err)
	}
	slog.Debug("seeded", "table", "art_tiles")
	return nil
}

func seedArtCategories(db *sqlx.DB) error {
	query := `
		INSERT INTO art_categories (art_id, category_id) VALUES
		(1, 2),
		(2, 1),
		(3, 3),
		(4, 3),
		(5, 3),
		(6, 1)
	`
	_, err := db.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to seed art_categories: %w", err)
	}
	slog.Debug("seeded", "table", "art_categories")
	return nil
}

