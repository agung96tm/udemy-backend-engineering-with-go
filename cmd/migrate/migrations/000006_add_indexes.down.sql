-- text / search indexes
DROP INDEX IF EXISTS idx_comments_content;
DROP INDEX IF EXISTS idx_posts_title;
DROP INDEX IF EXISTS idx_posts_tags;

-- lookup / relation indexes
DROP INDEX IF EXISTS idx_users_username;
DROP INDEX IF EXISTS idx_posts_user_id;
DROP INDEX IF EXISTS idx_comments_post_id;

-- optional (⚠️ hanya kalau extension ini khusus migration ini)
DROP EXTENSION IF EXISTS pg_trgm;
