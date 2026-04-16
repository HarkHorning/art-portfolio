CREATE TABLE IF NOT EXISTS print_sizes (
    id INT AUTO_INCREMENT PRIMARY KEY,
    print_id INT NOT NULL,
    size VARCHAR(50) NOT NULL,
    price_cents INT NOT NULL DEFAULT 0,
    quantity_in_stock INT NOT NULL DEFAULT 0,
    sold BOOLEAN NOT NULL DEFAULT FALSE,
    archived_at TIMESTAMP NULL DEFAULT NULL,
    FOREIGN KEY (print_id) REFERENCES prints(id) ON DELETE CASCADE,
    INDEX idx_print_sizes_print (print_id)
);

-- Migrate existing size/price/stock data into print_sizes
INSERT INTO print_sizes (print_id, size, price_cents, quantity_in_stock, sold)
SELECT id, size, price_cents, quantity_in_stock, sold FROM prints;

-- orders table is always empty at this point (no purchase UI exists).
-- Recreate it referencing print_sizes instead of prints.
DROP TABLE orders;

CREATE TABLE orders (
    id INT AUTO_INCREMENT PRIMARY KEY,
    print_size_id INT NOT NULL,
    stripe_payment_intent_id VARCHAR(255) NULL UNIQUE,
    customer_name VARCHAR(255) NOT NULL,
    customer_email VARCHAR(255) NOT NULL,
    customer_phone VARCHAR(50) NULL,
    shipping_line1 VARCHAR(255) NOT NULL,
    shipping_line2 VARCHAR(255) NULL,
    shipping_city VARCHAR(100) NOT NULL,
    shipping_state VARCHAR(100) NOT NULL,
    shipping_zip VARCHAR(20) NOT NULL,
    shipping_country VARCHAR(100) NOT NULL DEFAULT 'US',
    quantity INT NOT NULL DEFAULT 1,
    price_paid_cents INT NOT NULL,
    status ENUM('paid', 'processing', 'shipped', 'fulfilled', 'refunded') NOT NULL DEFAULT 'paid',
    notes TEXT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (print_size_id) REFERENCES print_sizes(id),
    INDEX idx_orders_status (status),
    INDEX idx_orders_email (customer_email)
);

ALTER TABLE prints
    DROP COLUMN size,
    DROP COLUMN price_cents,
    DROP COLUMN quantity_in_stock,
    DROP COLUMN sold;
