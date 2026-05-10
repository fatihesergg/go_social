package services

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/fatihesergg/go_social/internal/database"
	"github.com/fatihesergg/go_social/internal/errors"
	"github.com/fatihesergg/go_social/internal/model"
	"github.com/google/uuid"
)

type BaseFeedService interface {
	GetFeed(ctx context.Context, userID uuid.UUID, limit, offset, query string) ([]model.Post, error)
}

type FeedService struct {
	feedStore database.BaseFeedStore
}

func NewFeedService(feedStore database.BaseFeedStore) BaseFeedService {
	return &FeedService{feedStore: feedStore}
}
func (fs *FeedService) GetFeed(ctx context.Context, userID uuid.UUID, limit, offset, query string) ([]model.Post, error) {
	search := database.NewSearch(query)
	pagination := database.NewPagination(limit, offset)
	posts, err := fs.feedStore.GetFeed(ctx, userID, pagination, search)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.NoPostsFoundError.Wrap(fmt.Errorf("No posts found"))
		}
		return nil, errors.InternalServerError.Wrap(err)
	}
	return posts, nil
}
