CREATE EXTENSION IF NOT EXISTS pg_trgm;
CREATE TABLE products(
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    price DECIMAL(10, 2) NOT NULL,
    quantity INT NOT NULL DEFAULT 0 CHECK (quantity >= 0),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX IF NOT EXISTS idx_product_search ON products USING GIST (
    name gist_trgm_ops(siglen = 64),
    description gist_trgm_ops(siglen = 128)
);
CREATE INDEX IF NOT EXISTS idx_product_price ON products(price);
