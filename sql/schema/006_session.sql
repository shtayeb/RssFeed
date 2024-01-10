-- +goose Up
CREATE TABLE sessions (
    id UUID PRIMARY KEY,

    user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    
    ip_address TEXT,
    user_agent TEXT,
    payload TEXT,
    expires_at   TIMESTAMP

);

-- +goose Down
DROP TABLE sessions;
