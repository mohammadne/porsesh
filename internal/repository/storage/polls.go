package storage

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/mohammadne/porsesh/pkg/databases/postgres"
	"github.com/mohammadne/porsesh/pkg/observability/metrics"
	"go.uber.org/zap"
)

type Polls interface {
	RetrieveTransaction(ctx context.Context) (*sqlx.Tx, error)
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

func (c *polls) RetrieveTransaction(ctx context.Context) (*sqlx.Tx, error) {
	return c.db.BeginTxx(ctx, nil)
}

// --------------------------------------------------------------------------> Polls

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
VALUES (?, ?, ?)`

func (c *polls) CreatePoll(ctx context.Context, tx *sqlx.Tx, poll *Poll) (id int64, err error) {
	defer func(start time.Time) {
		if err != nil {
			c.db.Vectors.Counter.IncrementVector("polls", "create_poll", metrics.StatusFailure)
			return
		}
		c.db.Vectors.Counter.IncrementVector("polls", "create_poll", metrics.StatusSuccess)
		c.db.Vectors.Histogram.ObserveResponseTime(start, "polls", "create_poll")
	}(time.Now())

	result, err := tx.ExecContext(ctx, queryCreatePoll,
		poll.UserID, poll.Title, time.Now())
	if err != nil {
		return -1, errors.Join(ErrInsertingPoll, err)
	}
	id, err = result.LastInsertId()
	if err != nil {
		return -1, errors.Join(ErrRetrievingLastInsertedPollID, err)
	}

	return id, nil
}

// --------------------------------------------------------------------------> PollOptions

var (
	ErrInsertingPollOption = errors.New("")
)

const queryCreatePollOptions = `
INSERT INTO poll_options (poll_id, content, sort)
VALUES (?, ?, ?)`

type PollOption struct {
	ID      int64
	PollID  int64
	Content string
	Sort    int
}

func (c *polls) CreatePollOptions(ctx context.Context, tx *sqlx.Tx, pollID int64, options []PollOption) (err error) {
	defer func(start time.Time) {
		if err != nil {
			c.db.Vectors.Counter.IncrementVector("polls", "create_poll_options", metrics.StatusFailure)
			return
		}
		c.db.Vectors.Counter.IncrementVector("polls", "create_poll_options", metrics.StatusSuccess)
		c.db.Vectors.Histogram.ObserveResponseTime(start, "polls", "create_poll_options")
	}(time.Now())

	for _, option := range options {
		_, err := c.db.ExecContext(ctx, queryCreatePollOptions, pollID, option.Content, option.Sort)
		if err != nil {
			return errors.Join(ErrInsertingPollOption, err)
		}
	}

	return nil
}

var (
	ErrQueryGetPollOptionsByPollID                 = errors.New("")
	errScanningPollOptionInGetPollOptionsByPollID  = errors.New("")
	errScanningPollOptionsInGetPollOptionsByPollID = errors.New("")
)

const queryGetPollOptionsByPollID = `
SELECT id, poll_id, content, sort
FROM poll_options
WHERE poll_id = ?`

func (c *polls) GetPollOptionsByPollID(ctx context.Context, pollID int64) (result []PollOption, err error) {
	defer func(start time.Time) {
		if err != nil {
			c.db.Vectors.Counter.IncrementVector("polls", "get_poll_options_by_poll_id", metrics.StatusFailure)
			return
		}
		c.db.Vectors.Counter.IncrementVector("polls", "get_poll_options_by_poll_id", metrics.StatusSuccess)
		c.db.Vectors.Histogram.ObserveResponseTime(start, "polls", "get_poll_options_by_poll_id")
	}(time.Now())

	rows, err := c.db.QueryContext(ctx, queryGetPollOptionsByPollID, pollID)
	if err != nil {
		if err == sql.ErrNoRows {
			return []PollOption{}, nil
		}
		return nil, errors.Join(ErrQueryGetPollOptionsByPollID, err)
	}
	defer rows.Close() // ignore error

	result = make([]PollOption, 0)
	for rows.Next() {
		po := PollOption{}
		err = rows.Scan(&po.ID, &po.PollID, &po.Content, &po.Sort)
		if err != nil {
			return nil, errors.Join(errScanningPollOptionInGetPollOptionsByPollID, err)
		}
		result = append(result, po)
	}
	if rows.Err() != nil {
		return nil, errors.Join(errScanningPollOptionsInGetPollOptionsByPollID, err)
	}

	return result, nil

}

// --------------------------------------------------------------------------> PollTags

var (
	ErrInsertingPollTag = errors.New("")
)

const queryCreatePollTag = `
INSERT INTO poll_options (poll_id, tag_id)
VALUES (?, ?)`

type PollTag struct {
	ID   uint64
	Name string
}

func (c *polls) CreatePollTags(ctx context.Context, tx *sqlx.Tx, pollID int64, tagIDs []int64) (err error) {
	defer func(start time.Time) {
		if err != nil {
			c.db.Vectors.Counter.IncrementVector("polls", "create_poll_tags", metrics.StatusFailure)
			return
		}
		c.db.Vectors.Counter.IncrementVector("polls", "create_poll_tags", metrics.StatusSuccess)
		c.db.Vectors.Histogram.ObserveResponseTime(start, "polls", "create_poll_tags")
	}(time.Now())

	for _, tagID := range tagIDs {
		_, err := c.db.ExecContext(ctx, queryCreatePollTag, pollID, tagID)
		if err != nil {
			return errors.Join(ErrInsertingPollTag, err)
		}
	}

	return nil
}
