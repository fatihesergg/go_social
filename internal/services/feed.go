package services

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/fatihesergg/go_social/internal/database"
	"github.com/fatihesergg/go_social/internal/dto"
	"github.com/fatihesergg/go_social/internal/errors"
	"github.com/google/uuid"
)

type BaseFeedService interface {
	GetFeed(ctx context.Context, userID uuid.UUID, limit, offset, query string) ([]dto.FeedResponse, error)
}

type FeedService struct {
	feedStore database.BaseFeedStore
}

func NewFeedService(feedStore database.BaseFeedStore) BaseFeedService {
	return &FeedService{feedStore: feedStore}
}
func (fs *FeedService) GetFeed(ctx context.Context, userID uuid.UUID, limit, offset, query string) ([]dto.FeedResponse, error) {
	search := database.NewSearch(query)
	pagination := database.NewPagination(limit, offset)
	posts, err := fs.feedStore.GetFeed(ctx, userID, pagination, search)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.NoPostsFoundError.Wrap(fmt.Errorf("No posts found"))
		}
		return nil, errors.InternalServerError.Wrap(err)
	}
	response := dto.NewFeedResponse(posts)
	return response, nil
}
