package repo

import (
	"fmt"
	"log/slog"

	"github.com/jmoiron/sqlx"
)

// InitSchema creates the database tables if they don't exist.
// For production, use migrations instead.
func InitSchema(db *sqlx.DB) error {
	slog.Info("initializing database schema")

	if err := createArtTilesTable(db); err != nil {
		return err
	}
	if err := createCategoriesTable(db); err != nil {
		return err
	}
	if err := createArtCategoriesTable(db); err != nil {
		return err
	}

	slog.Info("database schema initialized")
	return nil
}

func createArtTilesTable(db *sqlx.DB) error {
	query := `
		CREATE TABLE IF NOT EXISTS art_tiles (
			id INT AUTO_INCREMENT PRIMARY KEY,
			title VARCHAR(255) NOT NULL,
			description TEXT,
			portrait BOOLEAN NOT NULL DEFAULT FALSE,
			url_low VARCHAR(512) NOT NULL,
			url_high VARCHAR(512) NOT NULL,
			display_order INT DEFAULT 0,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
			INDEX idx_display_order (display_order)
		)
	`
	_, err := db.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to create art_tiles table: %w", err)
	}
	slog.Debug("table ready", "table", "art_tiles")
	return nil
}

func createCategoriesTable(db *sqlx.DB) error {
	query := `
		CREATE TABLE IF NOT EXISTS categories (
			id INT AUTO_INCREMENT PRIMARY KEY,
			name VARCHAR(100) NOT NULL UNIQUE,
			slug VARCHAR(100) NOT NULL UNIQUE
		)
	`
	_, err := db.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to create categories table: %w", err)
	}
	slog.Debug("table ready", "table", "categories")
	return nil
}

func createArtCategoriesTable(db *sqlx.DB) error {
	query := `
		CREATE TABLE IF NOT EXISTS art_categories (
			art_id INT NOT NULL,
			category_id INT NOT NULL,
			PRIMARY KEY (art_id, category_id),
			FOREIGN KEY (art_id) REFERENCES art_tiles(id) ON DELETE CASCADE,
			FOREIGN KEY (category_id) REFERENCES categories(id) ON DELETE CASCADE
		)
	`
	_, err := db.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to create art_categories table: %w", err)
	}
	slog.Debug("table ready", "table", "art_categories")
	return nil
}

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
		INSERT INTO art_tiles (title, description, portrait, url_low, url_high, display_order) VALUES
		('Woman with Flowers', 'Acrylic on canvas, 2024', TRUE,
		 'https://harkportfoliostore.blob.core.windows.net/art-images/Woman With Flowers.jpeg',
		 'https://harkportfoliostore.blob.core.windows.net/art-images/Woman With Flowers.jpeg', 1),
		('Boat on Lake', 'Oil on canvas, peaceful morning scene', FALSE,
		 'https://harkportfoliostore.blob.core.windows.net/art-images/Boat on Lake.jpeg',
		 'https://harkportfoliostore.blob.core.windows.net/art-images/Boat on Lake.jpeg', 2),
		('Horse Watercolor', 'Watercolor', TRUE,
		 'https://harkportfoliostore.blob.core.windows.net/art-images/Horse Statue.jpeg',
		 'https://harkportfoliostore.blob.core.windows.net/art-images/Horse Statue.jpeg', 3),
		('Boat on Lake', 'Golden hour landscape', FALSE,
		 'https://harkportfoliostore.blob.core.windows.net/art-images/Boat on Lake.jpeg',
		 'https://harkportfoliostore.blob.core.windows.net/art-images/Boat on Lake.jpeg', 4),
		('Woman with Flowers', 'Cubist Serialist something expressionism', TRUE,
		 'https://harkportfoliostore.blob.core.windows.net/art-images/Woman With Flowers.jpeg',
		 'https://harkportfoliostore.blob.core.windows.net/art-images/Woman With Flowers.jpeg', 5),
		('Shoebill Stork Watercolor', 'Watercolor', TRUE,
		 'https://harkportfoliostore.blob.core.windows.net/art-images/Shoebill.jpeg',
		 'https://harkportfoliostore.blob.core.windows.net/art-images/Shoebill.jpeg', 6)
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
		(4, 1),
		(5, 2),
		(6, 3)
	`
	_, err := db.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to seed art_categories: %w", err)
	}
	slog.Debug("seeded", "table", "art_categories")
	return nil
}

// DropAllTables removes all portfolio tables. Use with caution!
func DropAllTables(db *sqlx.DB) error {
	slog.Warn("dropping all tables")

	tables := []string{"art_categories", "art_tiles", "categories"}
	for _, table := range tables {
		_, err := db.Exec(fmt.Sprintf("DROP TABLE IF EXISTS %s", table))
		if err != nil {
			return fmt.Errorf("failed to drop %s: %w", table, err)
		}
		slog.Debug("dropped table", "table", table)
	}

	slog.Info("all tables dropped")
	return nil
}
