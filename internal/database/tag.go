package database

import (
	"database/sql"
	"fmt"

	"github.com/fatihesergg/go_social/internal/model"
	"github.com/google/uuid"
	pq "github.com/lib/pq"
)

type BaseTagStore interface {
	CreateMultipleTags(tags []model.Tag) error
	CreatePostTag(tags []string, postID uuid.UUID) error
}

type TagStore struct {
	db *sql.DB
}

func NewTagStore(db *sql.DB) *TagStore {
	return &TagStore{
		db: db,
	}
}

func (ts *TagStore) CreateMultipleTags(tags []model.Tag) error {
	tx, err := ts.db.Begin()
	if err != nil {
		return err
	}
	fmt.Println(tags)

	defer tx.Rollback()
	for _, tag := range tags {
		_, err := tx.Exec("INSERT INTO tags(id,name) VALUES($1,$2) ON CONFLICT(name) DO NOTHING;", tag.ID, tag.Name)
		if err != nil {
			return fmt.Errorf("Error while inserting tag: %w", err)
		}
	}
	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("Error while committing transaction: %w", err)
	}
	return nil
}

func (ts *TagStore) CreatePostTag(tags []string, postID uuid.UUID) error {
	query := "INSERT INTO post_tags(post_id,tag_id) SELECT $1,id FROM tags WHERE tags.name = ANY($2);"
	result, err := ts.db.Exec(query, postID, pq.StringArray(tags))
	if err != nil {
		return fmt.Errorf("Error while inserting post tag: %w", err)
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("Error getting affected rows: %w", err)
	}
	if rows == 0 {
		return sql.ErrNoRows
	}
	return nil

}
