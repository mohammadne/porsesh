package usecases

import (
	"context"
	"errors"

	"github.com/mohammadne/porsesh/internal/entities"
	"go.uber.org/zap"
)

type Polls interface {
	CreatePoll(context.Context, *entities.Poll) error
	VotePoll(ctx context.Context, v entities.PollID, u entities.UserID, index int) error
	SkipPoll(ctx context.Context, v entities.PollID, u entities.UserID) error
	Statistics(ctx context.Context, v entities.PollID) (entities.PollStatistics, error)
}

func NewPolls(logger *zap.Logger) Polls {
	return nil
}

type pools struct {
	logger *zap.Logger
}

var (
	ErrInvalidCreatePollArguments = errors.New("")
)

var (
	ErrInvalidVotePollArguments = errors.New("")
	ErrVotePollPollNotExists    = errors.New("")
)

var (
	ErrInvalidSkipPollArguments = errors.New("")
	ErrSkipPollPollNotExists    = errors.New("")
)

var (
	ErrStatisticsPollNotExists = errors.New("")
)
