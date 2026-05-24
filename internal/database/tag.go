package database

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/fatihesergg/go_social/internal/model"
	"github.com/google/uuid"
	pq "github.com/lib/pq"
)

type BaseTagStore interface {
	CreateMultipleTags(ctx context.Context, tags []model.Tag) error
	CreatePostTag(ctx context.Context, tags []string, postID uuid.UUID) error
}

type TagStore struct {
	db *sql.DB
}

func NewTagStore(db *sql.DB) *TagStore {
	return &TagStore{
		db: db,
	}
}

func (ts *TagStore) CreateMultipleTags(ctx context.Context, tags []model.Tag) error {
	tx, err := ts.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	fmt.Println(tags)

	defer tx.Rollback()
	for _, tag := range tags {
		_, err := tx.ExecContext(ctx, "INSERT INTO tags(id,name) VALUES($1,$2) ON CONFLICT(name) DO NOTHING;", tag.ID, tag.Name)
		if err != nil {
			return fmt.Errorf("error while inserting tag: %w", err)
		}
	}
	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("error while committing transaction: %w", err)
	}
	return nil
}

func (ts *TagStore) CreatePostTag(ctx context.Context, tags []string, postID uuid.UUID) error {
	query := "INSERT INTO post_tags(post_id,tag_id) SELECT $1,id FROM tags WHERE tags.name = ANY($2);"
	result, err := ts.db.ExecContext(ctx, query, postID, pq.StringArray(tags))
	if err != nil {
		return fmt.Errorf("error while inserting post tag: %w", err)
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error getting affected rows: %w", err)
	}
	if rows == 0 {
		return sql.ErrNoRows
	}
	return nil

}
