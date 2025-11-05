package services

import (
	"database/sql"

	"github.com/fatihesergg/go_social/internal/database"
	"github.com/fatihesergg/go_social/internal/dto"
	"github.com/fatihesergg/go_social/internal/errors"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type BaseFeedService interface {
	GetFeed(userID uuid.UUID, limit, offset, query string) ([]dto.FeedResponse, error)
}

type FeedService struct {
	storage *database.Storage
	logger  *zap.Logger
}

func NewFeedService(storage *database.Storage, logger *zap.Logger) BaseFeedService {
	return &FeedService{storage: storage, logger: logger.Named("feed_service")}
}
func (fs *FeedService) GetFeed(userID uuid.UUID, limit, offset, query string) ([]dto.FeedResponse, error) {
	search := database.NewSearch(query)
	pagination := database.NewPagination(limit, offset)
	posts, err := fs.storage.FeedStore.GetFeed(userID, pagination, search)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.NoPostsFoundError
		}
		return nil, errors.InternalServerError.Wrap(err)
	}
	response := dto.NewFeedResponse(posts)
	return response, nil
}
