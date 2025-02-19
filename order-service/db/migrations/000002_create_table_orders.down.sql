DROP TABLE IF EXISTS orders CASCADE;
DROP TABLE IF EXISTS order_items CASCADE;
DROP INDEX IF EXISTS idx_order_status;
DROP INDEX IF EXISTS idx_order_created_at;
DROP INDEX IF EXISTS idx_fk_order_item_order_id;
DROP INDEX IF EXISTS idx_fk_order_item_product_id;
