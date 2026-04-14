ALTER TABLE art_tiles
    DROP COLUMN size,
    DROP COLUMN price_cents;

ALTER TABLE prints
    DROP COLUMN quantity_in_stock;
