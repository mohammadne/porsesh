package storage

import (
	"github.com/mohammadne/porsesh/pkg/databases/postgres"
	"go.uber.org/zap"
)

type Users interface {
}

func NewUsers(lg *zap.Logger, database *postgres.Postgres) Users {
	return &users{logger: lg, database: database}
}

type users struct {
	logger   *zap.Logger
	database *postgres.Postgres
}
