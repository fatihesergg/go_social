package services

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/fatihesergg/go_social/internal/database"
	"github.com/fatihesergg/go_social/internal/errors"
	"github.com/fatihesergg/go_social/internal/model"
	"github.com/google/uuid"
)

var tagExpr = regexp.MustCompile(`\#[\w]+`)

type BaseTagService interface {
	AddPostTags(postID uuid.UUID, tags []string) error
	ExtractTagStringFromContent(content string) ([]string, string)
}

type TagService struct {
	storage *database.Storage
}

func NewTagService(storage *database.Storage) *TagService {
	return &TagService{
		storage: storage,
	}
}

func (ts *TagService) ExtractTagStringFromContent(content string) ([]string, string) {
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
		return errors.InternalServerError.Wrap(fmt.Errorf("Error while creating post tags: %w", err)).Wrap(err)
	}

	err = ts.storage.TagStore.CreatePostTag(tags, postID)
	if err != nil {
		return errors.InternalServerError.Wrap(fmt.Errorf("Error while matching post to tags: %w", err)).Wrap(err)
	}
	return nil
}
