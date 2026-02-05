CREATE TABLE IF NOT EXISTS comments (
    id bigserial PRIMARY KEY,
    user_id bigint NOT NULL,
    post_id bigint NOT NULL,
    content TEXT NOT NULL,

    created_at timestamp(0) WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at timestamp(0) WITH TIME ZONE NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_comments_user
        FOREIGN KEY (user_id)
            REFERENCES users(id)
            ON DELETE CASCADE,

    CONSTRAINT fk_comments_post
        FOREIGN KEY (post_id)
            REFERENCES posts(id)
            ON DELETE CASCADE
);

CREATE INDEX idx_comments_user_id ON comments(user_id);
CREATE INDEX idx_comments_post_id ON comments(post_id);