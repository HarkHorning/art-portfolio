package repo

import (
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
)

// This file is just for development when I spin up a bunch of different things.
// Also, I want to keep any ideas for creating the db here for when I make a more perminant solution.
func InitSchema(db *sqlx.DB) error {
	log.Println("Initializing database schema...")

	if err := createArtTilesTable(db); err != nil {
		return err
	}
	if err := createCategoriesTable(db); err != nil {
		return err
	}
	if err := createArtCategoriesTable(db); err != nil {
		return err
	}

	log.Println("Database schema initialized successfully")
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
	log.Println("  - art_tiles table ready")
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
	log.Println("  - categories table ready")
	return nil
}

func createArtCategoriesTable(db *sqlx.DB) error { // Eventually add genre?
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
	log.Println("  - art_categories table ready")
	return nil
}

func SeedDevData(db *sqlx.DB) error { // For development
	log.Println("Seeding development data...")

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

	log.Println("Development data seeded successfully")
	return nil
}

func seedCategories(db *sqlx.DB) error {
	query := `
		INSERT INTO categories (name, slug) VALUES
		('Watercolor', 'watercolor'),
		('Oil Painting', 'oil'),
		('Portrait', 'portrait'),
		('Landscape', 'landscape'),
		('Acrylic', 'acrylic')
	`
	_, err := db.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to seed categories: %w", err)
	}
	log.Println("  - categories seeded")
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
	log.Println("  - art_tiles seeded")
	return nil
}

func seedArtCategories(db *sqlx.DB) error {
	query := `
		INSERT INTO art_categories (art_id, category_id) VALUES
		(1, 3),  -- Portrait -> Acrylic
		(2, 2),  -- Landscape -> Acrylic
		(2, 4),  -- Boat on Lake -> Acrylic
		(3, 1),  -- Horse Watercolor -> Watercolor
		(3, 3),  -- Horse Watercolor -> Portrait
		(4, 2),  -- Evening Light -> Acrylic
		(4, 4),  -- Evening Light -> Acrylic
		(5, 5),  -- Study in Blue -> Acrylic
		(6, 2),  -- Mountain Range -> Acrylic
		(6, 4)   -- Mountain Range -> Acrylic
	`
	_, err := db.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to seed art_categories: %w", err)
	}
	log.Println("  - art_categories seeded")
	return nil
}

// DropAllTables removes all portfolio tables. Use with caution!
func DropAllTables(db *sqlx.DB) error {
	log.Println("WARNING: Dropping all tables...")

	// Order matters - drop dependent tables first
	tables := []string{"art_categories", "art_tiles", "categories"}
	for _, table := range tables {
		_, err := db.Exec(fmt.Sprintf("DROP TABLE IF EXISTS %s", table))
		if err != nil {
			return fmt.Errorf("failed to drop %s: %w", table, err)
		}
		log.Printf("  - dropped %s", table)
	}

	log.Println("All tables dropped")
	return nil
}
