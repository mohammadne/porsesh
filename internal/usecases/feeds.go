package usecases

import "go.uber.org/zap"

type Feeds interface{}

func NewFeeds(logger *zap.Logger) Feeds {
	return nil
}

type feeds struct {
	logger *zap.Logger
}
