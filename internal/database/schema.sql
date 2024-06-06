CREATE TABLE IF NOT EXISTS deliveries
(
    id      SERIAL PRIMARY KEY,
    name    VARCHAR(255),
    phone   VARCHAR(20),
    zip     VARCHAR(20),
    city    VARCHAR(255),
    address VARCHAR(255),
    region  VARCHAR(255),
    email   VARCHAR(255)
);

CREATE TABLE IF NOT EXISTS payments
(
    id            SERIAL PRIMARY KEY,
    transaction   VARCHAR(255),
    request_id    VARCHAR(255),
    currency      VARCHAR(3),
    provider      VARCHAR(50),
    amount        INT,
    payment_dt    TIMESTAMP,
    bank          VARCHAR(50),
    delivery_cost INT,
    goods_total   INT,
    custom_fee    INT
);

CREATE TABLE IF NOT EXISTS items
(
    id           SERIAL PRIMARY KEY,
    chrt_id      INT,
    track_number VARCHAR(255),
    price        INT,
    rid          VARCHAR(255),
    name         VARCHAR(255),
    sale         INT,
    size         VARCHAR(50),
    total_price  INT,
    nm_id        INT,
    brand        VARCHAR(255),
    status       INT
);

CREATE TABLE IF NOT EXISTS orders
(
    id                 SERIAL PRIMARY KEY,
    order_uid          VARCHAR(255),
    track_number       VARCHAR(255),
    entry              VARCHAR(10),
    delivery_id        INT REFERENCES deliveries (id),
    payment_id         INT REFERENCES payments (id),
    items_id           INT[] REFERENCES items (id),
    locale             VARCHAR(10),
    internal_signature VARCHAR(255),
    customer_id        VARCHAR(50),
    delivery_service   VARCHAR(50),
    shardkey           VARCHAR(10),
    sm_id              INT,
    date_created       TIMESTAMP,
    oof_shard          VARCHAR(10)
);