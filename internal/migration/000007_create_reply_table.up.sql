CREATE TABLE IF NOT EXISTS replies (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    comment_id UUID REFERENCES comments(id) ON DELETE CASCADE,
    parent_id UUID REFERENCES replies(id) ON DELETE CASCADE,
    message TEXT NOT NULL,
    CHECK (
        (comment_id IS NOT NULL AND parent_id IS NULL) OR
        (comment_id IS NULL AND parent_id IS NOT NULL)
        )

);