-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied

CREATE TABLE notifications (
    id SERIAL PRIMARY KEY,
    message_to VARCHAR(255) NOT NULL,
    content TEXT NOT NULL,
    sent_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back

DROP TABLE notifications;
