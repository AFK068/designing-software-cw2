-- +goose Up
-- +goose StatementBegin
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE analysis (
    id UUID NOT NULL PRIMARY KEY DEFAULT uuid_generate_v4(),
    count_words INTEGER NOT NULL,
    count_characters INTEGER NOT NULL,
    is_plagiat BOOLEAN NOT NULL,
    hash VARCHAR(255) NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE analysis;
-- +goose StatementEnd