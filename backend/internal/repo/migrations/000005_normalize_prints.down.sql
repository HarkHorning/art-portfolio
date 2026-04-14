ALTER TABLE prints
    DROP FOREIGN KEY fk_prints_art_tile,
    DROP COLUMN art_tile_id,
    DROP COLUMN archived_at,
    ADD COLUMN title VARCHAR(255) NOT NULL DEFAULT '' AFTER id,
    ADD COLUMN description TEXT AFTER title,
    ADD COLUMN portrait BOOLEAN NOT NULL DEFAULT FALSE AFTER description,
    ADD COLUMN url_low VARCHAR(512) NOT NULL DEFAULT '' AFTER portrait,
    ADD COLUMN url_high VARCHAR(512) NOT NULL DEFAULT '' AFTER url_low,
    ADD COLUMN display_order INT DEFAULT 0 AFTER url_high;
