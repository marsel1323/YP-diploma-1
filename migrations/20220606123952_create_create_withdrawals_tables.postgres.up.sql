CREATE TABLE withdrawals
(
    id           SERIAL PRIMARY KEY,
    "order"      varchar,
    sum          NUMERIC(32, 2) default 0,
    processed_at TIMESTAMPTZ   NOT NULL,
    user_id      SERIAL NOT NULL
);
