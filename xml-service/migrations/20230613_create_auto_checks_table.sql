-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied

CREATE TABLE auto_checks (
    id SERIAL PRIMARY KEY,
    file_path VARCHAR(255) NOT NULL,
    status VARCHAR(50) NOT NULL,
    result TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back

DROP TABLE auto_checks;
