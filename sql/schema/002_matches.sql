-- +goose Up
CREATE TABLE matches (
    id UUID PRIMARY KEY,
    player1_id UUID NOT NULL REFERENCES users (id),
    player2_id UUID NOT NULL REFERENCES users (id),
    winner_id UUID REFERENCES users (id),
    created_at TIMESTAMP NOT NULL,
    completed_at TIMESTAMP
);

-- +goose Down
DROP TABLE matches;