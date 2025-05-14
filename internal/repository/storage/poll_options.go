package storage

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/mohammadne/porsesh/pkg/observability/metrics"
)

var (
	ErrInsertingPollOption = errors.New("")
)

const queryCreatePollOptions = `
INSERT INTO poll_options (poll_id, content, sort)
VALUES ($1, $2, $3)`

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
		_, err := tx.ExecContext(ctx, queryCreatePollOptions, pollID, option.Content, option.Sort)
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
WHERE poll_id = $1`

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
