package storage

import (
	"context"
	"errors"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/mohammadne/porsesh/pkg/databases/postgres"
	"github.com/mohammadne/porsesh/pkg/observability/metrics"
	"go.uber.org/zap"
)

type Polls interface {
	StartTransaction(ctx context.Context) (*sqlx.Tx, error)

	CreatePoll(ctx context.Context, tx *sqlx.Tx, poll *Poll) (id int64, err error)

	CreatePollOptions(ctx context.Context, tx *sqlx.Tx, pollID int64, options []PollOption) (err error)
	GetPollOptionsByPollID(ctx context.Context, pollID int64) (result []PollOption, err error)

	CreatePollTags(ctx context.Context, tx *sqlx.Tx, pollID int64, tagIDs []int64) (err error)
}

func NewPools(lg *zap.Logger, database *postgres.Postgres) Polls {
	return &polls{logger: lg, db: database}
}

type polls struct {
	logger *zap.Logger
	db     *postgres.Postgres
}

func (c *polls) StartTransaction(ctx context.Context) (*sqlx.Tx, error) {
	return c.db.BeginTxx(ctx, nil)
}

type Poll struct {
	ID        int64
	UserID    int64
	Title     string
	CreatedAt time.Time
}

var (
	ErrInsertingPoll                = errors.New("")
	ErrRetrievingLastInsertedPollID = errors.New("")
)

const queryCreatePoll = `
INSERT INTO polls (user_id, title, created_at)
VALUES ($1, $2, $3)
RETURNING id`

func (c *polls) CreatePoll(ctx context.Context, tx *sqlx.Tx, poll *Poll) (id int64, err error) {
	defer func(start time.Time) {
		if err != nil {
			c.db.Vectors.Counter.IncrementVector("polls", "create_poll", metrics.StatusFailure)
			return
		}
		c.db.Vectors.Counter.IncrementVector("polls", "create_poll", metrics.StatusSuccess)
		c.db.Vectors.Histogram.ObserveResponseTime(start, "polls", "create_poll")
	}(time.Now())

	err = tx.QueryRowContext(ctx, queryCreatePoll, poll.UserID, poll.Title, time.Now()).Scan(&id)
	if err != nil {
		return -1, errors.Join(ErrInsertingPoll, err)
	}

	return id, nil
}
