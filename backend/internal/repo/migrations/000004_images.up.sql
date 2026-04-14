CREATE TABLE IF NOT EXISTS images (
    id INT AUTO_INCREMENT PRIMARY KEY,
    art_tile_id INT NOT NULL,
    variant ENUM('high', 'low') NOT NULL,
    url VARCHAR(2048) NOT NULL,
    filename VARCHAR(500) NOT NULL DEFAULT '',
    sort_order TINYINT NOT NULL DEFAULT 0,
    uploaded_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (art_tile_id) REFERENCES art_tiles(id) ON DELETE CASCADE,
    INDEX idx_images_art_tile (art_tile_id),
    INDEX idx_images_variant (variant)
);

INSERT INTO images (art_tile_id, variant, url, sort_order)
SELECT id, 'high', url_high, 0 FROM art_tiles WHERE url_high IS NOT NULL AND url_high != '';

INSERT INTO images (art_tile_id, variant, url, sort_order)
SELECT id, 'low', url_low, 0 FROM art_tiles WHERE url_low IS NOT NULL AND url_low != '';

ALTER TABLE art_tiles
    DROP COLUMN url_low,
    DROP COLUMN url_high,
    ADD COLUMN archived_at TIMESTAMP NULL DEFAULT NULL;
