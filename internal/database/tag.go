package database

import (
	"database/sql"
	"fmt"

	"github.com/fatihesergg/go_social/internal/model"
	"github.com/google/uuid"
	pq "github.com/lib/pq"
	"go.uber.org/zap"
)

type BaseTagStore interface {
	CreateMultipleTags(tags []model.Tag) error
	CreatePostTag(tags []string, postID uuid.UUID) error
}

type TagStore struct {
	db     *sql.DB
	logger *zap.Logger
}

func NewTagStore(db *sql.DB, logger *zap.Logger) *TagStore {
	return &TagStore{
		db:     db,
		logger: logger.Named("tag_store"),
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
			ts.logger.Error("Error while inserting tag", zap.Error(err))
			return err
		}
	}
	err = tx.Commit()
	if err != nil {
		ts.logger.Error("Error while committing transaction", zap.Error(err))
		return err
	}
	return nil
}

func (ts *TagStore) CreatePostTag(tags []string, postID uuid.UUID) error {
	query := "INSERT INTO post_tags(post_id,tag_id) SELECT $1,id FROM tags WHERE tags.name = ANY($2);"
	result, err := ts.db.Exec(query, postID, pq.StringArray(tags))
	if err != nil {
		ts.logger.Error("Error while inserting post tag", zap.Error(err))
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		ts.logger.Error("Error getting affected rows", zap.Error(err))
		return err
	}
	if rows == 0 {
		return sql.ErrNoRows
	}
	return nil

}
