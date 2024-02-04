-- +goose Up
CREATE TABLE newsletters(
  id SERIAL PRIMARY KEY,

  email VARCHAR(255) NOT NULL UNIQUE,

  updated_at TIMESTAMP NOT NULL,
  created_at TIMESTAMP NOT NULL
);

CREATE INDEX newsletters_email_idx ON newsletters (email);

-- +goose Down
DROP TABLE newsletters;
