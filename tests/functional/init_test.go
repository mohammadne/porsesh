package functional

import (
	"log"
	"testing"

	"github.com/mohammadne/porsesh/internal/config"
	"github.com/mohammadne/porsesh/internal/repository/storage"
	"github.com/mohammadne/porsesh/internal/usecases"
	"github.com/mohammadne/porsesh/pkg/databases/postgres"
	"github.com/mohammadne/porsesh/pkg/observability/logger"
	"go.uber.org/zap"
)

var (
	// usecases
	feedsUsecase usecases.Feeds
	pollsUsecase usecases.Polls

	// caches

	// storages
	pollsStorage storage.Polls
	tagsStorage  storage.Tags
	votesStorage storage.Votes
)

func TestMain(m *testing.M) {
	var err error

	cfg := struct {
		Logger   *logger.Config   `required:"true"`
		Postgres *postgres.Config `required:"true"`
	}{}

	if err := config.Load(&cfg, string(config.EnvironmentLocal)); err != nil {
		log.Panicf("failed to load config: \n%v", err)
	}

	postgres, err := postgres.Open(cfg.Postgres, config.Namespace, config.System)
	if err != nil {
		log.Fatalf("failed to initialize Postgres: \n%v", err)
	}

	{ // storages
		pollsStorage = storage.NewPools(zap.NewNop(), postgres)
		tagsStorage = storage.NewTags(zap.NewNop(), postgres)
		votesStorage = storage.NewVotes(zap.NewNop(), postgres)
	}

	{ // usecases
		feedsUsecase = usecases.NewFeeds(zap.NewNop())
		pollsUsecase = usecases.NewPolls(zap.NewNop(), pollsStorage, tagsStorage, votesStorage)
	}

	m.Run()
}
