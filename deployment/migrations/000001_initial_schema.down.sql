-- Rollback Migration 000001: Initial Schema
-- Drops all core tables (order matters due to foreign keys)

DROP TABLE IF EXISTS art_categories;
DROP TABLE IF EXISTS art_tiles;
DROP TABLE IF EXISTS categories;
