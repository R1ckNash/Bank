CREATE TABLE IF NOT EXISTS users
(
    id         uuid PRIMARY KEY,
    username   text      NOT NULL UNIQUE,
    email      text      NOT NULL UNIQUE,
    password   text      NOT NULL,
    created_at timestamp NOT NULL DEFAULT now()
);