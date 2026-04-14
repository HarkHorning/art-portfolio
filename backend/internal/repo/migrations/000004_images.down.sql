ALTER TABLE art_tiles
    DROP COLUMN archived_at,
    ADD COLUMN url_low VARCHAR(512) NOT NULL DEFAULT '',
    ADD COLUMN url_high VARCHAR(512) NOT NULL DEFAULT '';

UPDATE art_tiles at
JOIN images i ON i.art_tile_id = at.id AND i.variant = 'low' AND i.sort_order = 0
SET at.url_low = i.url;

UPDATE art_tiles at
JOIN images i ON i.art_tile_id = at.id AND i.variant = 'high' AND i.sort_order = 0
SET at.url_high = i.url;

DROP TABLE IF EXISTS images;
