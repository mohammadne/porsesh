CREATE TABLE polls (
    id BIGSERIAL PRIMARY KEY,
    creator_id BIGINT REFERENCES users(id),
    question TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT now()
);

CREATE TABLE poll_options (
    id BIGSERIAL PRIMARY KEY,
    poll_id BIGINT REFERENCES polls(id) ON DELETE CASCADE,
    text TEXT NOT NULL
);

CREATE INDEX idx_poll_options_poll_id ON poll_options (poll_id);
