CREATE TABLE orders (
    order_id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    product_sku BIGINT NOT NULL,
    timestamp BIGINT NOT NULL,
    status TEXT,
    reason TEXT,

    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);
