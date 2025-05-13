package main

import "github.com/mohammadne/porsesh/pkg/databases/postgres"

type Config struct {
	Postgres *postgres.Config `required:"true"`
}

var cfg Config
