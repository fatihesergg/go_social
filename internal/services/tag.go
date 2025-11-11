package services

import (
	"regexp"
	"strings"

	"github.com/fatihesergg/go_social/internal/database"
	"github.com/fatihesergg/go_social/internal/errors"
	"github.com/fatihesergg/go_social/internal/model"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type BaseTagService interface {
	AddPostTags(postID uuid.UUID, tags []string) error
	ExtractTagStringFromContent(content string) ([]string, string)
}

type TagService struct {
	storage *database.Storage
	logger  *zap.Logger
}

func NewTagService(storage *database.Storage, logger *zap.Logger) *TagService {
	return &TagService{
		storage: storage,
		logger:  logger.Named("tag_service"),
	}
}

func (ts *TagService) ExtractTagStringFromContent(content string) ([]string, string) {
	tagExpr := regexp.MustCompile(`\#[\w]+`)
	lowerContent := strings.ToLower(content)
	var tags []string
	if tagExpr.MatchString(lowerContent) {
		matches := tagExpr.FindAllString(lowerContent, -1)
		tags = make([]string, len(matches))

		counter := 0
		for _, match := range matches {
			tag := match[1:]
			tags[counter] = tag
			counter++
		}
	}
	if tags == nil {
		return []string{}, content
	}

	cleanContent := tagExpr.ReplaceAllString(content, "")
	trimmedContent := strings.TrimSpace(cleanContent)
	return tags, trimmedContent

}

func (ts *TagService) AddPostTags(postID uuid.UUID, tags []string) error {
	tagModels := make([]model.Tag, len(tags))
	for i := 0; i < len(tags); i++ {
		tagModels[i] = model.Tag{
			ID:   uuid.New(),
			Name: tags[i],
		}
	}

	err := ts.storage.TagStore.CreateMultipleTags(tagModels)
	if err != nil {
		ts.logger.Error("Error while creating post tags", zap.Error(err))
		return errors.InternalServerError.Wrap(err)
	}

	err = ts.storage.TagStore.CreatePostTag(tags, postID)
	if err != nil {
		ts.logger.Error("Error while matching post to tags", zap.Error(err))
		return errors.InternalServerError.Wrap(err)
	}
	return nil
}
