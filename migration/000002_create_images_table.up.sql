CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS images
(
    image_id     UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name         TEXT                     NOT NULL,
    user_id      UUID REFERENCES users (user_id),
    content_type TEXT                     NOT NULL,
    url          TEXT                     NOT NULL,
    size         INTEGER                  NOT NULL,
    created_at   TIMESTAMP WITH TIME ZONE NOT NULL,

    CONSTRAINT images_unique_name UNIQUE (user_id, name)
);

CREATE INDEX IF NOT EXISTS images_user_idx ON images (user_id);