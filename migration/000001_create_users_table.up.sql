CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS users
(
    user_id       UUID  PRIMARY KEY DEFAULT uuid_generate_v4(),
    username      TEXT UNIQUE NOT NULL,
    password_hash BYTEA       NOT NULL,
    created_at    TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()

);

CREATE INDEX IF NOT EXISTS users_username_idx ON users(username);