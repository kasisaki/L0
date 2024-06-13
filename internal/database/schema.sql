CREATE TABLE IF NOT EXISTS orders
(
    order_uid          VARCHAR(255) PRIMARY KEY,
    track_number       VARCHAR(255),
    entry              VARCHAR(10),
    delivery_info      JSONB,
    payment_info       JSONB,
    items_info         JSONB,
    locale             VARCHAR(10),
    internal_signature VARCHAR(255),
    customer_id        VARCHAR(50),
    delivery_service   VARCHAR(50),
    shardkey           VARCHAR(10),
    sm_id              INT,
    date_created       TIMESTAMP,
    oof_shard          VARCHAR(10)
);

CREATE INDEX IF NOT EXISTS idx_orders_order_uid ON orders (order_uid);
