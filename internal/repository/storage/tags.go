package storage

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/mohammadne/porsesh/pkg/databases/postgres"
	"github.com/mohammadne/porsesh/pkg/observability/metrics"
	"go.uber.org/zap"
)

type Tags interface {
	CreateTags(ctx context.Context, tx *sqlx.Tx, tags []Tag) (ids map[string]int64, err error)
	GetTagsByNames(ctx context.Context, tx *sqlx.Tx, names []string) (result []Tag, err error)
}

func NewTags(lg *zap.Logger, database *postgres.Postgres) Tags {
	return &tags{logger: lg, db: database}
}

type tags struct {
	logger *zap.Logger
	db     *postgres.Postgres
}

type Tag struct {
	ID   int64
	Name string
}

var (
	errInsertingTag                = errors.New("error inserting url")
	ErrRetrievingLastInsertedTagID = errors.New("")
)

const queryCreateTag = `
INSERT INTO tags (name)
VALUES (?)`

func (c *tags) CreateTags(ctx context.Context, tx *sqlx.Tx, tags []Tag) (idsMap map[string]int64, err error) {
	defer func(start time.Time) {
		if err != nil {
			c.db.Vectors.Counter.IncrementVector("tags", "create_tags", metrics.StatusFailure)
			return
		}
		c.db.Vectors.Counter.IncrementVector("tags", "create_tags", metrics.StatusSuccess)
		c.db.Vectors.Histogram.ObserveResponseTime(start, "tags", "create_tags")
	}(time.Now())

	idsMap = make(map[string]int64, len(tags))
	for _, tag := range tags {
		result, err := tx.ExecContext(ctx, queryCreateTag, tag.Name)
		if err != nil {
			if pgErr, ok := err.(*pq.Error); ok && pgErr.Code == postgres.UniqueConstraintViolatedCode {
				continue
			}
			return nil, errors.Join(errInsertingTag, err)
		}
		id, err := result.LastInsertId()
		if err != nil {
			return nil, errors.Join(ErrRetrievingLastInsertedTagID, err)
		}
		idsMap[tag.Name] = id
	}

	return idsMap, nil
}

var (
	// errInsertingTag                = errors.New("error inserting url")
	// ErrRetrievingLastInsertedTagID = errors.New("")
	errPrepareGetTagsByNamesQuery   = errors.New("errPrepareGetTagsByNamesQuery")
	errGetTagsByNames               = errors.New("")
	errScanningTagInGetTagsByNames  = errors.New("")
	errScanningTagsInGetTagsByNames = errors.New("")
)

const queryGetTagsByNames = `
SELECT ID, NAME
FROM tags
WHERE NAME IN (?)`

func (c *tags) GetTagsByNames(ctx context.Context, tx *sqlx.Tx, names []string) (result []Tag, err error) {
	defer func(start time.Time) {
		if err != nil {
			c.db.Vectors.Counter.IncrementVector("tags", "get_tags_by_names", metrics.StatusFailure)
			return
		}
		c.db.Vectors.Counter.IncrementVector("tags", "get_tags_by_names", metrics.StatusSuccess)
		c.db.Vectors.Histogram.ObserveResponseTime(start, "tags", "get_tags_by_names")
	}(time.Now())

	expandedQuery, args, err := sqlx.In(queryGetTagsByNames, names)
	if err != nil {
		return nil, errors.Join(errPrepareGetTagsByNamesQuery, err)
	}
	expandedQuery = tx.Rebind(expandedQuery)

	rows, err := tx.QueryContext(ctx, expandedQuery, args...)
	if err != nil {
		if err == sql.ErrNoRows {
			return []Tag{}, nil
		}
		return nil, errors.Join(errGetTagsByNames, err)
	}
	defer rows.Close() // ignore error

	result = make([]Tag, 0)
	for rows.Next() {
		badge := Tag{}
		err = rows.Scan(&badge.ID, &badge.Name)
		if err != nil {
			return nil, errors.Join(errScanningTagInGetTagsByNames, err)
		}
		result = append(result, badge)
	}
	if rows.Err() != nil {
		return nil, errors.Join(errScanningTagsInGetTagsByNames, err)
	}

	return result, nil
}
