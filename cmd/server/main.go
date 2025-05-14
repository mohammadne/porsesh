package main

import (
	"context"
	"flag"
	"log"
	"os/signal"
	"sync"
	"syscall"

	"go.uber.org/zap"

	"github.com/mohammadne/porsesh/cmd"
	"github.com/mohammadne/porsesh/internal/api/http"
	"github.com/mohammadne/porsesh/internal/config"
	"github.com/mohammadne/porsesh/internal/repository/storage"
	"github.com/mohammadne/porsesh/internal/usecases"
	"github.com/mohammadne/porsesh/pkg/databases/postgres"
	"github.com/mohammadne/porsesh/pkg/observability/logger"
)

func main() {
	monitorPort := flag.Int("monitor-port", 8001, "The server port which handles monitoring endpoints (default: 8001)")
	requestPort := flag.Int("request-port", 8002, "The server port which handles http requests (default: 8002)")
	environmentRaw := flag.String("environment", "", "The environment (default: local)")
	flag.Parse() // Parse the command-line flags

	var cfg Config
	if err := config.Load(&cfg, *environmentRaw); err != nil {
		log.Panicf("failed to load config: \n%v", err)
	}

	logger, err := logger.New(cfg.Logger)
	if err != nil {
		log.Fatalf("failed to initialize logger: \n%v", err)
	}

	{ // print build informations
		fields := make([]zap.Field, 0, len(cmd.BuildInfo()))
		for k, v := range cmd.BuildInfo() {
			fields = append(fields, zap.String(k, v))
			logger.Warn("Build Information", fields...)
		}
	}

	postgres, err := postgres.Open(cfg.Postgres, config.Namespace, config.System)
	if err != nil {
		log.Fatalf("failed to initialize Postgres: \n%v", err)
	}

	pollsStorage := storage.NewPools(zap.NewNop(), postgres)
	tagsStorage := storage.NewTags(zap.NewNop(), postgres)
	votesStorage := storage.NewVotes(zap.NewNop(), postgres)

	// usecases
	feeds := usecases.NewFeeds(logger)
	pools := usecases.NewPolls(logger, pollsStorage, tagsStorage, votesStorage)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	var wg sync.WaitGroup

	wg.Add(1)
	go http.New(logger, feeds, pools).Serve(ctx, &wg, *monitorPort, *requestPort)

	<-ctx.Done()
	wg.Wait()
}

type Config struct {
	Logger   *logger.Config   `required:"true"`
	Postgres *postgres.Config `required:"true"`
}
