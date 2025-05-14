package storage

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/lib/pq"
	"github.com/mohammadne/porsesh/pkg/databases/postgres"
	"github.com/mohammadne/porsesh/pkg/observability/metrics"
	"go.uber.org/zap"
)

type Votes interface {
	CreateVote(ctx context.Context, vote *Vote) (err error)
	GetPollOptionVotesCount(ctx context.Context, optionID int64) (result uint64, err error)
	GetCurrentDateUserVoteCount(ctx context.Context, userID int64) (result int64, err error)
}

func NewVotes(lg *zap.Logger, database *postgres.Postgres) Votes {
	return &votes{logger: lg, db: database}
}

type votes struct {
	logger *zap.Logger
	db     *postgres.Postgres
}

type Vote struct {
	UserID   int64
	PollID   int64
	OptionID sql.NullInt64
	ActedAt  time.Time
}

var (
	ErrInsertingVote           = errors.New("")
	ErrCreateVotePollNotExists = errors.New("")
)

const queryCreateVote = `
INSERT INTO votes (user_id, poll_id, option_id, acted_at)
VALUES ($1, $2, $3, $4)`

func (v *votes) CreateVote(ctx context.Context, vote *Vote) (err error) {
	defer func(start time.Time) {
		if err != nil {
			v.db.Vectors.Counter.IncrementVector("votes", "create_vote", metrics.StatusFailure)
			return
		}
		v.db.Vectors.Counter.IncrementVector("votes", "create_vote", metrics.StatusSuccess)
		v.db.Vectors.Histogram.ObserveResponseTime(start, "votes", "create_vote")
	}(time.Now())

	_, err = v.db.ExecContext(ctx, queryCreateVote,
		vote.UserID, vote.PollID, vote.OptionID, time.Now())
	if err != nil {
		if pgErr, ok := err.(*pq.Error); ok && pgErr.Code == postgres.ForeignKeyViolatedCode {
			return ErrCreateVotePollNotExists
		}
		return errors.Join(ErrInsertingVote, err)
	}

	return nil
}

var (
	errQueryGetPollOptionCount = errors.New("")
)

const queryGetPollOptionCount = `
SELECT COUNT(*) AS vote_count
FROM votes
WHERE option_id = $1
GROUP BY option_id`

func (v *votes) GetPollOptionVotesCount(ctx context.Context, optionID int64) (result uint64, err error) {
	defer func(start time.Time) {
		if err != nil {
			v.db.Vectors.Counter.IncrementVector("votes", "get_poll_option_votes_count", metrics.StatusFailure)
			return
		}
		v.db.Vectors.Counter.IncrementVector("votes", "get_poll_option_votes_count", metrics.StatusSuccess)
		v.db.Vectors.Histogram.ObserveResponseTime(start, "votes", "get_poll_option_votes_count")
	}(time.Now())

	err = v.db.QueryRowContext(ctx, queryGetPollOptionCount, optionID).Scan(&result)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, nil
		}
		return 0, errors.Join(errQueryGetPollOptionCount, err)
	}

	return result, nil
}

var (
	errQueryGetCurrentDateUserVoteCount = errors.New("")
)

const queryGetCurrentDateUserVoteCount = `
SELECT COUNT(*) 
FROM votes 
WHERE user_id = $1 AND acted_at::date = CURRENT_DATE`

func (v *votes) GetCurrentDateUserVoteCount(ctx context.Context, userID int64) (result int64, err error) {
	defer func(start time.Time) {
		if err != nil {
			v.db.Vectors.Counter.IncrementVector("votes", "get_current_date_user_vote_count", metrics.StatusFailure)
			return
		}
		v.db.Vectors.Counter.IncrementVector("votes", "get_current_date_user_vote_count", metrics.StatusSuccess)
		v.db.Vectors.Histogram.ObserveResponseTime(start, "votes", "get_current_date_user_vote_count")
	}(time.Now())

	err = v.db.QueryRowContext(ctx, queryGetCurrentDateUserVoteCount, userID).Scan(&result)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, nil
		}
		return 0, errors.Join(errQueryGetCurrentDateUserVoteCount, err)
	}

	return result, nil
}
