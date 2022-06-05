CREATE TABLE users
(
    id           SERIAL PRIMARY KEY,
    login        VARCHAR UNIQUE,
    passwordHash VARCHAR
);
