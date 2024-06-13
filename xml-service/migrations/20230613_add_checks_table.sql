-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied

CREATE TABLE checks (
    id SERIAL PRIMARY KEY,
    service_id VARCHAR(255) NOT NULL,
    request TEXT NOT NULL,
    code INT NOT NULL,
    response_expectation TEXT NOT NULL
);

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back

DROP TABLE checks;
