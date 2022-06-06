CREATE TABLE orders
(
    id          SERIAL PRIMARY KEY,
    number      VARCHAR UNIQUE,
    status      VARCHAR   NOT NULL,
    accrual     NUMERIC(32, 2),
    uploaded_at TIMESTAMPTZ NOT NULL,
    user_id     SERIAL    NOT NULL,
    CONSTRAINT fk_user
        FOREIGN KEY (user_id)
            REFERENCES users (id)
);
