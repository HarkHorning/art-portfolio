ALTER TABLE prints
    ADD COLUMN size VARCHAR(50) NOT NULL DEFAULT '' AFTER art_tile_id,
    ADD COLUMN price_cents INT NOT NULL DEFAULT 0 AFTER size,
    ADD COLUMN quantity_in_stock INT NOT NULL DEFAULT 0 AFTER price_cents,
    ADD COLUMN sold BOOLEAN NOT NULL DEFAULT FALSE AFTER quantity_in_stock;

-- Restore one size per print (cheapest)
UPDATE prints p
INNER JOIN (
    SELECT print_id, size, price_cents, quantity_in_stock, sold
    FROM print_sizes
    WHERE (print_id, price_cents) IN (
        SELECT print_id, MIN(price_cents) FROM print_sizes GROUP BY print_id
    )
    GROUP BY print_id
) ps ON ps.print_id = p.id
SET p.size = ps.size, p.price_cents = ps.price_cents,
    p.quantity_in_stock = ps.quantity_in_stock, p.sold = ps.sold;

DROP TABLE orders;

CREATE TABLE orders (
    id INT AUTO_INCREMENT PRIMARY KEY,
    print_id INT NOT NULL,
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
    FOREIGN KEY (print_id) REFERENCES prints(id),
    INDEX idx_orders_status (status),
    INDEX idx_orders_email (customer_email)
);

DROP TABLE print_sizes;
