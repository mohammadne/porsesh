CREATE TABLE polls (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT,
    title TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT now()
);

CREATE TABLE poll_options (
    id BIGSERIAL PRIMARY KEY,
    poll_id BIGINT REFERENCES polls(id) ON DELETE CASCADE,
    content TEXT NOT NULL,
    sort INT NOT NULL
);

CREATE INDEX idx_poll_options_poll_id ON poll_options (poll_id);

CREATE TABLE poll_tags (
    poll_id BIGINT REFERENCES polls(id) ON DELETE CASCADE,
    tag_id BIGINT REFERENCES tags(id) ON DELETE CASCADE,
    PRIMARY KEY (poll_id, tag_id)
);

CREATE INDEX idx_poll_tags_tag_id ON poll_tags(tag_id);
