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
	CreateVote(ctx context.Context, vote *Vote) (id int64, err error)
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
INSERT INTO votes (user_id, poll_id, option_id, created_at)
VALUES ($1, $2, $3, $4)
RETURNING id`

func (v *votes) CreateVote(ctx context.Context, vote *Vote) (id int64, err error) {
	defer func(start time.Time) {
		if err != nil {
			v.db.Vectors.Counter.IncrementVector("votes", "create_vote", metrics.StatusFailure)
			return
		}
		v.db.Vectors.Counter.IncrementVector("votes", "create_vote", metrics.StatusSuccess)
		v.db.Vectors.Histogram.ObserveResponseTime(start, "votes", "create_vote")
	}(time.Now())

	err = v.db.QueryRowContext(ctx, queryCreateVote,
		vote.UserID, vote.PollID, vote.OptionID, time.Now()).Scan(&id)
	if err != nil {
		if pgErr, ok := err.(*pq.Error); ok && pgErr.Code == postgres.ForeignKeyViolatedCode {
			return -1, ErrCreateVotePollNotExists
		}
		return -1, errors.Join(ErrInsertingVote, err)
	}

	return id, nil
}

var (
	errQueryGetCurrentDateUserVoteCount = errors.New("")
)

const queryGetCurrentDateUserVoteCount = `
SELECT COUNT(*) 
FROM votes 
WHERE user_id = ? AND acted_at::date = CURRENT_DATE`

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
