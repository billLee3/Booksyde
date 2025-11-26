-- +goose Up
CREATE TABLE users(
  id UUID PRIMARY KEY,
  created_at TIMESTAMP NOT NULL,
  updated_at TIMESTAMP NOT NULL,
  first_name TEXT NOT NULL,
  last_name TEXT NOT NULL,
  email TEXT UNIQUE NOT NULL,
  password_hash TEXT NOT NULL,
  subscribed BOOL NOT NULL,
  birth_month TEXT NOT NULL,
  birth_year INT NOT NULL
);

-- +goose Down
DROP TABLE users;
