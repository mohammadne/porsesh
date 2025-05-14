package usecases

import (
	"context"
	"database/sql"
	"errors"
	"slices"
	"sync"

	"github.com/mohammadne/porsesh/internal/entities"
	"github.com/mohammadne/porsesh/internal/repository/storage"
	"go.uber.org/zap"
)

type Polls interface {
	CreatePoll(context.Context, *entities.Poll) error
	VotePoll(ctx context.Context, v entities.PollID, u entities.UserID, index int) error
	SkipPoll(ctx context.Context, v entities.PollID, u entities.UserID) error
	Statistics(ctx context.Context, v entities.PollID) (*entities.PollStatistics, error)
}

func NewPolls(logger *zap.Logger, ps storage.Polls, ts storage.Tags, vs storage.Votes) Polls {
	return &pools{
		logger:       logger,
		pollsStorage: ps,
		tagsStorage:  ts,
		votesStorage: vs,
	}
}

type pools struct {
	logger *zap.Logger
	// storages
	pollsStorage storage.Polls
	tagsStorage  storage.Tags
	votesStorage storage.Votes
}

var (
	ErrInvalidCreatePollArguments = errors.New("")
)

func (p *pools) CreatePoll(ctx context.Context, poll *entities.Poll) (err error) {
	{ // validation over poll
		if len(poll.Options) <= 0 || len(poll.Options) > 5 {
			return ErrInvalidCreatePollArguments
		}

		if len(poll.Tags) > 3 {
			return ErrInvalidCreatePollArguments
		}

	}

	tx, err := p.pollsStorage.StartTransaction(ctx)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()

	storagePoll := storage.Poll{
		UserID: int64(poll.UserID),
		Title:  poll.Title,
	}

	pollID, err := p.pollsStorage.CreatePoll(ctx, tx, &storagePoll)
	if err != nil {
		return err
	}

	{ // create poll options
		storagePollOptions := make([]storage.PollOption, 0, len(poll.Options))
		for _, option := range poll.Options {
			storagePollOptions = append(storagePollOptions, storage.PollOption{
				Content: option.Content,
				Sort:    option.Sort,
			})
		}

		err = p.pollsStorage.CreatePollOptions(ctx, tx, pollID, storagePollOptions)
		if err != nil {
			return err
		}
	}

	var tagIds []int64
	if len(poll.Tags) != 0 {
		storageTags := make([]storage.Tag, 0, len(poll.Tags))
		for _, tag := range poll.Tags {
			storageTags = append(storageTags, storage.Tag{Name: tag.Name})
		}
		tagIdsMap, err := p.tagsStorage.CreateTags(ctx, tx, storageTags)
		if err != nil {
			return err
		}

		names := make([]string, 0, len(poll.Tags)-len(tagIdsMap))
		for _, tag := range poll.Tags {
			id, exists := tagIdsMap[tag.Name]
			if !exists {
				names = append(names, tag.Name)
			} else {
				tagIds = append(tagIds, id)
			}
		}

		if len(tagIds) != len(poll.Tags) {
			existingTags, err := p.tagsStorage.GetTagsByNames(ctx, tx, names)
			if err != nil {
				return err
			}
			for _, tag := range existingTags {
				tagIds = append(tagIds, tag.ID)
			}
		}
	}

	err = p.pollsStorage.CreatePollTags(ctx, tx, pollID, tagIds)
	if err != nil {
		return err
	}

	return nil
}

var (
	ErrInvalidVotePollArguments = errors.New("")
	ErrDailyUserVotesLimit      = errors.New("")
	ErrVotePollPollNotExists    = errors.New("")
)

const dailyUserVoteLimits = 100

func (p *pools) VotePoll(ctx context.Context, pollID entities.PollID, u entities.UserID, index int) error {
	{ // validation
		if index < 0 || pollID < 0 {
			return ErrInvalidVotePollArguments
		}
	}

	count, err := p.votesStorage.GetCurrentDateUserVoteCount(ctx, int64(u))
	if err != nil { // log but continue
		p.logger.Error("error retrieving user daily count", zap.Int64("user_id", int64(u)), zap.Error(err))
	}

	if count >= dailyUserVoteLimits {
		return ErrDailyUserVotesLimit
	}

	storageOptions, err := p.pollsStorage.GetPollOptionsByPollID(ctx, int64(pollID))
	if err != nil {
		return err
	} else if len(storageOptions) == 0 {
		return ErrVotePollPollNotExists
	}

	if len(storageOptions)-1 < index {
		return ErrInvalidVotePollArguments
	}

	slices.SortFunc(storageOptions, func(a, b storage.PollOption) int {
		return a.Sort - b.Sort
	})

	nullableOption := sql.NullInt64{Valid: true, Int64: storageOptions[index].ID}
	storagePoll := storage.Vote{UserID: int64(u), PollID: int64(pollID), OptionID: nullableOption}
	if _, err := p.votesStorage.CreateVote(ctx, &storagePoll); err != nil {
		return err
	}

	return nil
}

var (
	ErrInvalidSkipPollArguments = errors.New("")
	ErrSkipPollPollNotExists    = errors.New("")
)

func (p *pools) SkipPoll(ctx context.Context, pollID entities.PollID, u entities.UserID) error {
	{ // validation
		if pollID < 0 {
			return ErrInvalidSkipPollArguments
		}
	}

	nullableOption := sql.NullInt64{Valid: true}
	storagePoll := storage.Vote{UserID: int64(u), PollID: int64(pollID), OptionID: nullableOption}
	if _, err := p.votesStorage.CreateVote(ctx, &storagePoll); err != nil {
		if errors.Is(err, storage.ErrCreateVotePollNotExists) {
			return ErrSkipPollPollNotExists
		}
		return err
	}

	return nil
}

var (
	ErrInvalidStatisticsArguments = errors.New("")
	ErrStatisticsPollNotExists    = errors.New("")
)

func (p *pools) Statistics(ctx context.Context, pollID entities.PollID) (result *entities.PollStatistics, err error) {
	{ // validation
		if pollID < 0 {
			return nil, ErrInvalidStatisticsArguments
		}
	}

	result = &entities.PollStatistics{PoolID: pollID}

	// todo: retrieve from cache

	storageOptions, err := p.pollsStorage.GetPollOptionsByPollID(ctx, int64(pollID))
	if err != nil {
		return nil, err
	} else if len(storageOptions) == 0 {
		return nil, ErrStatisticsPollNotExists
	}

	var wg sync.WaitGroup
	var mu sync.Mutex
	wg.Add(len(storageOptions))

	for _, so := range storageOptions {
		go func(so *storage.PollOption) {
			defer wg.Done()

			count, err := p.pollsStorage.GetPollOptionCount(ctx, so.ID)
			if err != nil {
				p.logger.Error("error counting options", zap.Int64("option_id", so.ID), zap.Error(err))
				return
			}

			mu.Lock()
			result.Votes = append(result.Votes, entities.PollStatisticsVote{
				Option: so.Content, Count: count,
			})
			mu.Unlock()
		}(&so)
	}

	wg.Wait()

	go func() {
		// todo: cache the result
	}()

	return result, nil
}
