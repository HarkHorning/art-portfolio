-- Clear existing print data (will be reseeded; art_tile_id cannot be inferred from titles safely)
DELETE FROM prints;

ALTER TABLE prints
    DROP COLUMN title,
    DROP COLUMN description,
    DROP COLUMN portrait,
    DROP COLUMN url_low,
    DROP COLUMN url_high,
    DROP COLUMN display_order,
    ADD COLUMN art_tile_id INT NOT NULL AFTER id,
    ADD COLUMN archived_at TIMESTAMP NULL DEFAULT NULL,
    ADD CONSTRAINT fk_prints_art_tile FOREIGN KEY (art_tile_id) REFERENCES art_tiles(id) ON DELETE RESTRICT;
