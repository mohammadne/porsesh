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
	ErrInsertingVote                = errors.New("")
	ErrRetrievingLastInsertedVoteID = errors.New("")
	ErrCreateVotePollNotExists      = errors.New("")
)

const queryCreateVote = `
INSERT INTO votes (user_id, poll_id, option_id, created_at)
VALUES (?, ?, ?, ?)`

func (v *votes) CreateVote(ctx context.Context, vote *Vote) (id int64, err error) {
	defer func(start time.Time) {
		if err != nil {
			v.db.Vectors.Counter.IncrementVector("votes", "create_vote", metrics.StatusFailure)
			return
		}
		v.db.Vectors.Counter.IncrementVector("votes", "create_vote", metrics.StatusSuccess)
		v.db.Vectors.Histogram.ObserveResponseTime(start, "votes", "create_vote")
	}(time.Now())

	result, err := v.db.ExecContext(ctx, queryCreateVote,
		vote.UserID, vote.PollID, vote.OptionID, time.Now())
	if err != nil {
		if pgErr, ok := err.(*pq.Error); ok && pgErr.Code == postgres.ForeignKeyViolatedCode {
			return -1, ErrCreateVotePollNotExists
		}
		return -1, errors.Join(ErrInsertingVote, err)
	}
	id, err = result.LastInsertId()
	if err != nil {
		return -1, errors.Join(ErrRetrievingLastInsertedVoteID, err)
	}

	return id, nil
}
