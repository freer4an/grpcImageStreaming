-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS images (
    id UUID PRIMARY KEY,
    original_path TEXT,
    thumbnail_path TEXT,
    width INT,
    height INT,
    uploaded_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS images;
-- +goose StatementEnd
