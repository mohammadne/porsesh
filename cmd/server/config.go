package main

import (
	"github.com/mohammadne/porsesh/pkg/observability/logger"
)

type Config struct {
	Logger *logger.Config `required:"true"`
}

var cfg Config
