package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/lib/pq"
)

type Post struct {
	ID      int64    `json:"id"`
	Content string   `json:"content"`
	Title   string   `json:"title"`
	UserID  int64    `json:"user_id"`
	Tags    []string `json:"tags"`
	Version int64    `json:"version"`

	Comments []*Comment `json:"comments"`

	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`

	User User `json:"user"`
}

type PostWithMeta struct {
	Post
	CommentCount int `json:"comment_count"`
}

type PostStore struct {
	db *sql.DB
}

func (s *PostStore) Create(ctx context.Context, post *Post) error {
	query := `INSERT INTO posts (title, content, user_id, tags) VALUES ($1, $2, $3, $4) RETURNING id, created_at, updated_at`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	err := s.db.QueryRowContext(
		ctx, query,
		post.Title,
		post.Content,
		post.UserID,
		pq.Array(post.Tags),
	).Scan(
		&post.ID,
		&post.CreatedAt,
		&post.UpdatedAt,
	)
	return err
}

func (s *PostStore) GetByID(ctx context.Context, id int64) (*Post, error) {
	query := `
		SELECT id, title, content, user_id, tags, version, created_at, updated_at 
		FROM posts WHERE id = $1
	`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	var post Post
	err := s.db.QueryRowContext(
		ctx, query, id,
	).Scan(
		&post.ID,
		&post.Title,
		&post.Content,
		&post.UserID,
		pq.Array(post.Tags),
		&post.Version,
		&post.CreatedAt,
		&post.UpdatedAt,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}
	return &post, err
}

func (s *PostStore) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM posts WHERE id = $1`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	result, err := s.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrNotFound
	}
	return nil
}

func (s *PostStore) Update(ctx context.Context, post *Post) error {
	query := `
		UPDATE posts 
			SET title = $1, content = $2, version = version + 1
		WHERE id = $3 AND version = $4
		RETURNING id, created_at, updated_at, version`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	err := s.db.QueryRowContext(
		ctx, query,
		post.Title, post.Content, post.ID, post.Version,
	).Scan(
		&post.ID, &post.CreatedAt, &post.UpdatedAt, &post.Version,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrNotFound
		default:
			return err
		}
	}
	return nil
}

func (s *PostStore) GetUserFeed(ctx context.Context, userID int64, fq PaginatedFeedQuery) ([]*PostWithMeta, error) {
	query := fmt.Sprintf(`
		SELECT
			p.id,
			p.title,
			p.content,
			p.user_id,
			p.tags,
			p.version,
			p.created_at,
			p.updated_at,
			u.id,
			u.username,
			COUNT(DISTINCT c.id) AS comment_count
		FROM posts p
		LEFT JOIN comments c ON c.post_id = p.id
		JOIN users u ON u.id = p.user_id
		LEFT JOIN followers f 
			ON f.follower_id = p.user_id
			AND f.user_id = $1
		WHERE
			(
				f.user_id IS NOT NULL
				OR p.user_id = $1
			)
			AND ($4 = '' OR p.title ILIKE '%' || $4 || '%' OR p.content ILIKE '%' || $4 || '%')
			AND (cardinality($5) = 0 OR p.tags @> $5)
			AND ($6 = '' OR p.created_at >= $6)
			AND ($7 = '' OR p.created_at <= $7)
		GROUP BY p.id, u.id
		ORDER BY p.created_at %s
		LIMIT $2 OFFSET $3
	`, fq.Sort)

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	rows, err := s.db.QueryContext(ctx, query, userID, fq.Limit, fq.Offset, fq.Search, pq.Array(fq.Sort), fq.Since, fq.Until)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []*PostWithMeta
	for rows.Next() {
		var post PostWithMeta
		err := rows.Scan(
			&post.ID, &post.Title, &post.Content, &post.UserID,
			&post.Tags, &post.Version, &post.CreatedAt, &post.UpdatedAt,
			&post.CommentCount,
			&post.User.ID,
			&post.User.Username,
		)
		if err != nil {
			return nil, err
		}
		posts = append(posts, &post)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}
