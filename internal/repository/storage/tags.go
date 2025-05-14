package storage

import (
	"context"
	"errors"
	"fmt"

	"github.com/lib/pq"
	"github.com/mohammadne/porsesh/pkg/databases/postgres"
	"go.uber.org/zap"
)

type Tags interface {
}

func NewTags(lg *zap.Logger, database *postgres.Postgres) Tags {
	return &tags{logger: lg, database: database}
}

type tags struct {
	logger   *zap.Logger
	database *postgres.Postgres
}

var (
	errTagAlreadyExists = errors.New("error duplicate tag")
	errInsertingTag     = errors.New("error inserting url")
)

const queryCreateTag = `
	INSERT INTO tags (name)
	VALUES (?)`

func (c *tags) CreateTag(ctx context.Context, name string) (id uint64, err error) {
	result, err := c.database.ExecContext(ctx, queryCreateTag, name)
	if err != nil {
		if pgErr, ok := err.(*pq.Error); ok && pgErr.Code == postgres.UniqueConstraintViolatedCode {
			return 0, errTagAlreadyExists
		}
		return 0, errors.Join(errInsertingTag, err)
	}
	id64, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("error retrieving last inserted tag's id: %v", err)
	}

	return uint64(id64), nil
}
