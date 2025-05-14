package usecases

import (
	"context"

	"github.com/mohammadne/porsesh/internal/entities"
	"go.uber.org/zap"
)

type Feeds interface {
	RetrieveUserFeed(ctx context.Context, userID entities.UserID, tag string, page, limit int) (entities.Feed, error)
}

func NewFeeds(logger *zap.Logger) Feeds {
	return nil
}

type feeds struct {
	logger *zap.Logger
}
