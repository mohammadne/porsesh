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
	StartTransaction(ctx context.Context) (*sqlx.Tx, error)

	CreatePoll(ctx context.Context, tx *sqlx.Tx, poll *Poll) (id int64, err error)
	ListPollsByTag(ctx context.Context, userID, tagID int64, limit, offset int) (result []Poll, err error)

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

var (
	errListPolls               = errors.New("")
	errScanningPollInListPolls = errors.New("")
	errIteratingInListPolls    = errors.New("")
)

const (
	queryListPollsWithoutTag = `
	SELECT p.id, p.title, p.created_at
	FROM polls p
	LEFT JOIN votes v ON v.poll_id = p.id AND v.user_id = $1
	WHERE v.poll_id IS NULL
	ORDER BY p.created_at DESC
	LIMIT $2 OFFSET $3`

	queryListPollsByTag = `
	SELECT p.id, p.title, p.created_at
	FROM polls p
	JOIN poll_tags pt ON pt.poll_id = p.id
	LEFT JOIN votes v ON v.poll_id = p.id AND v.user_id = $1
	WHERE v.poll_id IS NULL AND pt.tag_id = $4
	ORDER BY p.created_at DESC
	LIMIT $2 OFFSET $3`
)

func (c *polls) ListPollsByTag(ctx context.Context, userID, tagID int64, limit, offset int) (result []Poll, err error) {
	defer func(start time.Time) {
		if err != nil {
			c.db.Vectors.Counter.IncrementVector("polls", "list_polls", metrics.StatusFailure)
			return
		}
		c.db.Vectors.Counter.IncrementVector("polls", "list_polls", metrics.StatusSuccess)
		c.db.Vectors.Histogram.ObserveResponseTime(start, "polls", "list_polls")
	}(time.Now())

	query := queryListPollsWithoutTag
	args := []any{userID, limit, offset}
	if tagID > 0 {
		query = queryListPollsByTag
		args = append(args, tagID)
	}

	rows, err := c.db.QueryContext(ctx, query, args...)
	if err != nil {
		if err == sql.ErrNoRows {
			return []Poll{}, nil
		}
		return nil, errors.Join(errListPolls, err)
	}
	defer rows.Close() // ignore error

	result = make([]Poll, 0)
	for rows.Next() {
		poll := Poll{}
		err = rows.Scan(&poll.ID, &poll.Title, &poll.CreatedAt)
		if err != nil {
			return nil, errors.Join(errScanningPollInListPolls, err)
		}
		result = append(result, poll)
	}
	if rows.Err() != nil {
		return nil, errors.Join(errIteratingInListPolls, err)
	}

	return result, nil
}
