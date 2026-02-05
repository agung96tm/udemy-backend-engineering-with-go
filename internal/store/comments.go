package store

import (
	"context"
	"database/sql"
)

type Comment struct {
	ID      int64  `json:"id"`
	Content string `json:"content"`

	UserID int64 `json:"user_id"`
	PostID int64 `json:"post_id"`

	User User `json:"user"`

	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type CommentStore struct {
	db *sql.DB
}

func (c *CommentStore) GetAllByPostID(ctx context.Context, postID int64) ([]*Comment, error) {
	query := `
		SELECT c.id, c.content, c.user_id, c.post_id, u.id, u.username, u.email, c.created_at, c.updated_at
		FROM comments c
		INNER JOIN users u ON u.id = c.user_id
		WHERE c.post_id = $1
		`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	rows, err := c.db.QueryContext(ctx, query, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []*Comment
	for rows.Next() {
		var comment Comment
		err := rows.Scan(
			&comment.ID,
			&comment.Content,
			&comment.UserID,
			&comment.PostID,
			&comment.User.ID,
			&comment.User.Username,
			&comment.User.Email,
			&comment.CreatedAt,
			&comment.UpdatedAt)
		if err != nil {
			return nil, err
		}
		comments = append(comments, &comment)
	}

	return comments, nil
}

func (c *CommentStore) Create(ctx context.Context, comment *Comment) error {
	query := `INSERT INTO comments (content, user_id, post_id) VALUES ($1, $2, $3) RETURNING id, created_at, updated_at`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	err := c.db.QueryRowContext(
		ctx, query,
		comment.Content, comment.UserID, comment.PostID,
	).Scan(
		&comment.ID, &comment.CreatedAt, &comment.UpdatedAt,
	)
	if err != nil {
		return err
	}

	return nil
}
