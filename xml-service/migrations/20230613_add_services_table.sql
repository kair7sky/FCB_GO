-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied

CREATE TABLE services (
    id SERIAL PRIMARY KEY,
    serviceid VARCHAR(255) NOT NULL,
    name VARCHAR(255) NOT NULL,
    description TEXT
);

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back

DROP TABLE services;
