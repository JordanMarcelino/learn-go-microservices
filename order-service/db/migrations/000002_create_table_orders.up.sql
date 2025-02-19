CREATE TABLE IF NOT EXISTS orders(
    id BIGSERIAL PRIMARY KEY,
    request_id VARCHAR(255) NOT NULL UNIQUE,
    customer_id BIGINT NOT NULL,
    total_amount DECIMAL(10, 2) NOT NULL,
    description TEXT DEFAULT NULL,
    status VARCHAR(50) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE TABLE IF NOT EXISTS order_items(
    id BIGSERIAL PRIMARY KEY,
    order_id BIGINT NOT NULL REFERENCES orders (id) ON DELETE CASCADE,
    product_id BIGINT NOT NULL REFERENCES products (id),
    quantity INT NOT NULL,
    price DECIMAL(10, 2) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX IF NOT EXISTS idx_order_status ON orders (status);
CREATE INDEX IF NOT EXISTS idx_order_created_at ON orders (created_at);
CREATE INDEX IF NOT EXISTS idx_fk_order_item_order_id ON order_items (order_id);
CREATE INDEX IF NOT EXISTS idx_fk_order_item_product_id ON order_items (product_id);
