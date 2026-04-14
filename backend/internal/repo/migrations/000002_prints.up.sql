CREATE TABLE IF NOT EXISTS prints (
    id INT AUTO_INCREMENT PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    portrait BOOLEAN NOT NULL DEFAULT FALSE,
    url_low VARCHAR(512) NOT NULL,
    url_high VARCHAR(512) NOT NULL,
    display_order INT DEFAULT 0,
    price_cents INT NOT NULL DEFAULT 0,
    size VARCHAR(20) NOT NULL,
    sold BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_prints_display_order (display_order),
    INDEX idx_prints_size (size),
    INDEX idx_prints_price (price_cents)
);
