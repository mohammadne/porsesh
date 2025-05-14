package storage

import (
	"context"
	"errors"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/mohammadne/porsesh/pkg/observability/metrics"
)

var (
	ErrInsertingPollTag = errors.New("")
)

const queryCreatePollTag = `
INSERT INTO poll_tags (poll_id, tag_id)
VALUES ($1, $2)`

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
		_, err := tx.ExecContext(ctx, queryCreatePollTag, pollID, tagID)
		if err != nil {
			return errors.Join(ErrInsertingPollTag, err)
		}
	}

	return nil
}
