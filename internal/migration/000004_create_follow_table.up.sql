CREATE TYPE follow_request_status AS ENUM ('pending', 'accepted', 'rejected');


CREATE TABLE IF NOT EXISTS follows(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    follow_id UUID NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (follow_id) REFERENCES users(id) ON DELETE CASCADE,
    status follow_request_status NOT NULL DEFAULT 'pending',
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    UNIQUE(user_id, follow_id),
    CHECK (user_id != follow_id)
);

CREATE INDEX idx_follows_follow_id_status ON follows(follow_id, status);
CREATE INDEX idx_follows_user_id_status ON follows(user_id, status);
