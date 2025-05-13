package main

import (
	"embed"
	"flag"
	"log"
	"strings"

	"github.com/mohammadne/porsesh/internal/config"
	"github.com/mohammadne/porsesh/pkg/databases/postgres"
)

//go:embed schemas/*.sql
var files embed.FS

func main() {
	direction := flag.String("direction", "", "Either 'UP' or 'DOWN'")
	environmentRaw := flag.String("environment", "", "The environment (default: local)")
	flag.Parse() // Parse the command-line flags

	if err := config.Load(&cfg, *environmentRaw); err != nil {
		log.Panicf("failed to load config: \n%v", err)
	}

	db, err := postgres.Open(cfg.Postgres, config.Namespace, config.System)
	if err != nil {
		log.Fatalf("error connecting to postgres database\n%v", err)
	}

	migrateDirection := postgres.MigrateDirection(strings.ToUpper(*direction))
	err = db.Migrate("schemas", &files, migrateDirection)
	if err != nil {
		log.Fatalf("error migrating postgres database\n%v", err)
	}

	log.Println("database has been migrated")
}
