-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS images (
    id UUID PRIMARY KEY,
    format VARCHAR,
    width INT,
    height INT,
    original_path TEXT,
    thumbnail_path TEXT,
    uploaded_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS images;
-- +goose StatementEnd
