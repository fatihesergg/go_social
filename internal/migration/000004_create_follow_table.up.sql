CREATE TABLE IF NOT EXISTS follows(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    follow_id UUID NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (follow_id) REFERENCES users(id) ON DELETE CASCADE
);