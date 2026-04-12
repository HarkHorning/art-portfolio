-- Migration 000001: Initial Schema
-- Creates the core tables for the portfolio application

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
);

CREATE TABLE IF NOT EXISTS categories (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE,
    slug VARCHAR(100) NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS art_categories (
    art_id INT NOT NULL,
    category_id INT NOT NULL,
    PRIMARY KEY (art_id, category_id),
    FOREIGN KEY (art_id) REFERENCES art_tiles(id) ON DELETE CASCADE,
    FOREIGN KEY (category_id) REFERENCES categories(id) ON DELETE CASCADE
);
