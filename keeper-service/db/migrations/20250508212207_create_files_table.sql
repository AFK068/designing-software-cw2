-- +goose Up
-- +goose StatementBegin
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE files (
    id UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
    data TEXT NOT NULL,
    name VARCHAR(255) NOT NULL,
    hash VARCHAR(255) NOT NULL,
    location VARCHAR(255) NOT NULL
);

CREATE UNIQUE INDEX files_hash_uindex ON files (hash);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE files;
-- +goose StatementEnd