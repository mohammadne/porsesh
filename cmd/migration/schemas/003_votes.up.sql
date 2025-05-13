CREATE TABLE votes (
    user_id BIGINT REFERENCES users(id),
    poll_id BIGINT REFERENCES polls(id) ON DELETE CASCADE,
    option_id BIGINT REFERENCES poll_options(id) ON DELETE CASCADE, -- NULL means skipped
    acted_at TIMESTAMP DEFAULT now(),
    PRIMARY KEY (user_id, poll_id)
);

CREATE INDEX idx_votes_user_id ON votes (user_id);
CREATE INDEX idx_votes_option_id ON votes (option_id);
