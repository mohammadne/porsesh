-- the users table
CREATE TABLE IF NOT EXISTS users (
	id BIGSERIAL PRIMARY KEY,
	username TEXT UNIQUE NOT NULL,
	created_at TIMESTAMP DEFAULT now()
);
