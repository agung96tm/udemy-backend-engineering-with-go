package store

import (
	"context"
	"database/sql"
)

type Follower struct {
	UserID     int64  `json:"user_id"`
	FollowerID int64  `json:"follower_id"`
	CreatedAt  string `json:"created_at"`
}

type FollowerStore struct {
	db *sql.DB
}

func (s *FollowerStore) Follow(ctx context.Context, followerID, userID int64) error {
	query := `
		INSERT INTO followers (user_id, follower_id)
		SELECT $1::BIGINT, $2::BIGINT
		WHERE $1::BIGINT <> $2::BIGINT
		ON CONFLICT DO NOTHING
	`
	_, err := s.db.ExecContext(ctx, query, userID, followerID)
	return err
}

func (s *FollowerStore) Unfollow(ctx context.Context, followerID, userID int64) error {
	query := `
		DELETE FROM followers
		WHERE user_id = $1
		  AND follower_id = $2
	`
	result, err := s.db.ExecContext(ctx, query, userID, followerID)
	if err != nil {
		return err
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return ErrNotFound
	}

	return nil
}
