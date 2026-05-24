CREATE TYPE notify_type AS ENUM ('post_like','comment_like','reply_like');


CREATE TABLE IF NOT EXISTS notifications(
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    notification_type notify_type NOT NULL,
    message VARCHAR(255) NOT NULL,
    payload JSONB NOT NULL DEFAULT '{}',
    is_read BOOLEAN NOT NULL DEFAULT false,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP

);