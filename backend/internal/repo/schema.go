package repo

import (
	"fmt"
	"log/slog"

	"github.com/jmoiron/sqlx"
)

func SeedDevData(db *sqlx.DB) error {
	slog.Info("seeding development data")

	tables := []string{"orders", "print_sizes", "art_categories", "images", "prints", "art_tiles", "categories"}
	for _, table := range tables {
		_, err := db.Exec(fmt.Sprintf("DELETE FROM %s", table))
		if err != nil {
			return fmt.Errorf("failed to clear %s: %w", table, err)
		}
	}

	for _, table := range []string{"art_tiles", "categories", "prints", "print_sizes", "images", "orders"} {
		_, _ = db.Exec(fmt.Sprintf("ALTER TABLE %s AUTO_INCREMENT = 1", table))
	}

	if err := seedCategories(db); err != nil {
		return err
	}
	if err := seedArtTiles(db); err != nil {
		return err
	}
	if err := seedImages(db); err != nil {
		return err
	}
	if err := seedArtCategories(db); err != nil {
		return err
	}
	if err := seedPrints(db); err != nil {
		return err
	}

	slog.Info("development data seeded")
	return nil
}

func seedCategories(db *sqlx.DB) error {
	_, err := db.Exec(`
		INSERT INTO categories (name, slug) VALUES
		('Oil', 'oil'),
		('Acrylic', 'acrylic'),
		('Watercolor', 'watercolor'),
		('Pencil Drawing', 'pencil-drawing'),
		('Mixed', 'mixed'),
		('Pastel', 'pastel'),
		('Misc', 'misc')
	`)
	if err != nil {
		return fmt.Errorf("failed to seed categories: %w", err)
	}
	slog.Debug("seeded", "table", "categories")
	return nil
}

func seedArtTiles(db *sqlx.DB) error {
	_, err := db.Exec(`
		INSERT INTO art_tiles (title, description, portrait, display_order, made_year, sold) VALUES
		('Woman with Flowers',  'Acrylic on canvas',            TRUE,  1,  2024, FALSE),
		('Boat on Lake',        'Oil on canvas',                FALSE, 2,  2023, FALSE),
		('Horse Statue',        'Watercolor',                   TRUE,  3,  2023, FALSE),
		('Cardinal',            'Watercolor',                   TRUE,  4,  2023, FALSE),
		('Shoebill Stork',      'Watercolor',                   TRUE,  5,  2022, FALSE),
		('Mother with Children','Acrylic on canvas',            FALSE, 6,  2023, FALSE),
		('Bird Drawing',        'Pencil on paper',              TRUE,  7,  2022, FALSE),
		('Pencil Study',        'Pencil on paper',              TRUE,  8,  2022, FALSE),
		('Ass of an Artist',    'Acrylic on canvas',            TRUE,  9,  NULL, FALSE),
		('Yellow Bird',         'Acrylic and colored pencil',   TRUE,  10, NULL, FALSE)
	`)
	if err != nil {
		return fmt.Errorf("failed to seed art_tiles: %w", err)
	}
	slog.Debug("seeded", "table", "art_tiles")
	return nil
}

func seedImages(db *sqlx.DB) error {
	_, err := db.Exec(`
		INSERT INTO images (art_tile_id, variant, url, filename, sort_order) VALUES
		(1,  'high', 'https://storage.googleapis.com/hark-portfolio-images/art/Woman%20With%20Flowers.jpeg',  'Woman With Flowers.jpeg',  0),
		(1,  'low',  'https://storage.googleapis.com/hark-portfolio-images/art/Woman%20With%20Flowers.jpeg',  'Woman With Flowers.jpeg',  0),
		(2,  'high', 'https://storage.googleapis.com/hark-portfolio-images/art/Boat%20on%20Lake.jpeg',        'Boat on Lake.jpeg',        0),
		(2,  'low',  'https://storage.googleapis.com/hark-portfolio-images/art/Boat%20on%20Lake.jpeg',        'Boat on Lake.jpeg',        0),
		(3,  'high', 'https://storage.googleapis.com/hark-portfolio-images/art/Horse%20Statue.jpeg',          'Horse Statue.jpeg',        0),
		(3,  'low',  'https://storage.googleapis.com/hark-portfolio-images/art/Horse%20Statue.jpeg',          'Horse Statue.jpeg',        0),
		(4,  'high', 'https://storage.googleapis.com/hark-portfolio-images/art/Cardinal.jpeg',                'Cardinal.jpeg',            0),
		(4,  'low',  'https://storage.googleapis.com/hark-portfolio-images/art/Cardinal.jpeg',                'Cardinal.jpeg',            0),
		(5,  'high', 'https://storage.googleapis.com/hark-portfolio-images/art/Shoebill.jpeg',                'Shoebill.jpeg',            0),
		(5,  'low',  'https://storage.googleapis.com/hark-portfolio-images/art/Shoebill.jpeg',                'Shoebill.jpeg',            0),
		(6,  'high', 'https://storage.googleapis.com/hark-portfolio-images/art/Mother%20with%20Children.jpg', 'Mother with Children.jpg', 0),
		(6,  'low',  'https://storage.googleapis.com/hark-portfolio-images/art/Mother%20with%20Children.jpg', 'Mother with Children.jpg', 0),
		(7,  'high', 'https://storage.googleapis.com/hark-portfolio-images/art/Bird%20Drawing.jpg',           'Bird Drawing.jpg',         0),
		(7,  'low',  'https://storage.googleapis.com/hark-portfolio-images/art/Bird%20Drawing.jpg',           'Bird Drawing.jpg',         0),
		(8,  'high', 'https://storage.googleapis.com/hark-portfolio-images/art/20220609_001313.jpg',          '20220609_001313.jpg',      0),
		(8,  'low',  'https://storage.googleapis.com/hark-portfolio-images/art/20220609_001313.jpg',          '20220609_001313.jpg',      0),
		(9,  'high', 'https://storage.googleapis.com/hark-portfolio-images/art/Graf-One.jpg',                 'Graf-One.jpg',             0),
		(9,  'low',  'https://storage.googleapis.com/hark-portfolio-images/art/Graf-One.jpg',                 'Graf-One.jpg',             0),
		(10, 'high', 'https://storage.googleapis.com/hark-portfolio-images/art/Yellow%20Bird.jpg',            'Yellow Bird.jpg',          0),
		(10, 'low',  'https://storage.googleapis.com/hark-portfolio-images/art/Yellow%20Bird.jpg',            'Yellow Bird.jpg',          0)
	`)
	if err != nil {
		return fmt.Errorf("failed to seed images: %w", err)
	}
	slog.Debug("seeded", "table", "images")
	return nil
}

func seedArtCategories(db *sqlx.DB) error {
	_, err := db.Exec(`
		INSERT INTO art_categories (art_id, category_id) VALUES
		(1, 2),
		(2, 1),
		(3, 3),
		(4, 3),
		(5, 3),
		(6, 2),
		(7, 4),
		(8, 4),
		(9, 2),
		(10, 5)
	`)
	if err != nil {
		return fmt.Errorf("failed to seed art_categories: %w", err)
	}
	slog.Debug("seeded", "table", "art_categories")
	return nil
}

func seedPrints(db *sqlx.DB) error {
	_, err := db.Exec(`INSERT INTO prints (art_tile_id) VALUES (1),(2),(3),(4),(5)`)
	if err != nil {
		return fmt.Errorf("failed to seed prints: %w", err)
	}
	_, err = db.Exec(`
		INSERT INTO print_sizes (print_id, size, price_cents, quantity_in_stock) VALUES
		(1, '5x7',   2500,  5),
		(1, '8x10',  4500,  5),
		(2, '8x10',  4500,  3),
		(2, '11x14', 7500,  0),
		(3, '11x14', 7500,  3),
		(3, '16x20', 12000, 0),
		(4, '5x7',   2500,  5),
		(4, '11x14', 7500,  2),
		(5, '8x10',  4500,  5),
		(5, '16x20', 12000, 1)
	`)
	if err != nil {
		return fmt.Errorf("failed to seed print sizes: %w", err)
	}
	slog.Debug("seeded", "table", "prints+print_sizes")
	return nil
}
