CREATE TABLE feeds (
    user_id BIGINT NOT NULL REFERENCES users(id),
    poll_id BIGINT NOT NULL REFERENCES polls(id),
    added_at TIMESTAMP DEFAULT now(),
    expires_at TIMESTAMP NOT NULL,
    PRIMARY KEY (user_id, poll_id)
);

CREATE INDEX idx_feeds_expires_at ON feeds (expires_at);
