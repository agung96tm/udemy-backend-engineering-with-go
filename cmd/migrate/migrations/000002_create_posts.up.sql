CREATE TABLE IF NOT EXISTS posts (
    id BIGSERIAL PRIMARY KEY,
    title text NOT NULL,
    content text NOT NULL,
    user_id bigint NOT NULL,

    tags TEXT[] DEFAULT '{}',

    created_at timestamp(0) WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at timestamp(0) WITH TIME ZONE NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_posts_user
        FOREIGN KEY (user_id)
            REFERENCES users(id)
            ON DELETE CASCADE
)