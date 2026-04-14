ALTER TABLE art_tiles
    ADD COLUMN size VARCHAR(20) NULL,
    ADD COLUMN price_cents INT NULL;

ALTER TABLE prints
    ADD COLUMN quantity_in_stock INT NOT NULL DEFAULT 0;
