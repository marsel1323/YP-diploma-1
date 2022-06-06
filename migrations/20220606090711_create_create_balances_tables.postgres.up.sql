CREATE TABLE balances
(
    id        SERIAL PRIMARY KEY,
    current   NUMERIC(32, 2) CHECK (current >= 0) default 0 ,
    withdrawn NUMERIC(32, 2) CHECK (withdrawn >= 0) default 0,
    user_id   SERIAL UNIQUE NOT NULL,
    CONSTRAINT fk_user
        FOREIGN KEY (user_id)
            REFERENCES users (id)
);
