package usecases

import (
	"context"

	"github.com/mohammadne/porsesh/internal/entities"
	"github.com/mohammadne/porsesh/internal/repository/storage"
	"go.uber.org/zap"
)

type Feeds interface {
	GetUserFeed(ctx context.Context, userID entities.UserID, tag string, page, limit int) (entities.Feed, error)
}

func NewFeeds(logger *zap.Logger) Feeds {
	return &feeds{logger: logger}
}

type feeds struct {
	logger *zap.Logger
	// storages
	pollsStorage storage.Polls
	tagsStorage  storage.Tags
	votesStorage storage.Votes
}

func (f *feeds) GetUserFeed(ctx context.Context, userID entities.UserID, tag string, page, limit int) (entities.Feed, error) {
	{ // validation
		if limit < 5 {
			limit = 5
		} else if limit > 20 {
			limit = 20
		}

		if page < 1 {
			page = 1
		}
	}

	var tagID int64
	if len(tag) != 0 {
		result, err := f.tagsStorage.GetTagByName(ctx, tag)
		if err != nil {
			f.logger.Error("error retrieving tag", zap.Error(err))
		} else if result != nil {
			tagID = result.ID
		}
	}

	storagePolls, err := f.pollsStorage.ListPollsByTag(ctx, int64(userID), tagID, page, limit)
	if err != nil {
		return nil, err
	}

	result := make([]entities.Poll, 0, len(storagePolls))
	for _, storagePoll := range storagePolls {
		result = append(result, entities.Poll{
			ID:        entities.PollID(storagePoll.ID),
			Title:     storagePoll.Title,
			UserID:    userID,
			CreatedAt: storagePoll.CreatedAt,
			Options:   nil, // TODO: implement
			Tags:      nil, // TODO: implement
		})
	}

	return result, nil
}
