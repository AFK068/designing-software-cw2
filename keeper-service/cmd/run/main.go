package main

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/fx"
	"go.uber.org/zap"

	_ "github.com/lib/pq"

	"github.com/AFK068/antiplagiarism/keeper-service/internal/config"
	"github.com/AFK068/antiplagiarism/keeper-service/internal/domain"
	"github.com/AFK068/antiplagiarism/keeper-service/internal/infrastructure/httpapi/keeperapi"
	"github.com/AFK068/antiplagiarism/keeper-service/internal/infrastructure/repository/postgresdb"
	"github.com/AFK068/antiplagiarism/keeper-service/internal/infrastructure/server"
	"github.com/AFK068/antiplagiarism/keeper-service/internal/migration"
	"github.com/AFK068/antiplagiarism/keeper-service/pkg/logger"
)

const (
	DevConfigPath  = "config/dev.yaml"
	MigrationsPath = "db/migrations"
)

func postgresDB(cfg *config.Config, log *zap.Logger, lc fx.Lifecycle) (domain.Repository, error) {
	dbPool, err := pgxpool.New(context.Background(), cfg.GetPostgresConnectionString())
	if err != nil {
		log.Fatal("Failed to connect to PostgreSQL", zap.Error(err))
	}

	lc.Append(fx.Hook{
		OnStop: func(_ context.Context) error {
			dbPool.Close()
			return nil
		},
	})

	return postgresdb.NewPostgresRepository(dbPool), nil
}

func main() {
	fx.New(
		fx.Provide(
			// Logger.
			logger.New,

			// Config.
			func() (*config.Config, error) {
				return config.NewConfig(DevConfigPath)
			},

			// PostgresDB.
			postgresDB,

			// Handler.
			keeperapi.NewKeeperHandler,

			// Server.
			server.NewKeeperServer,
		),
		fx.Invoke(
			func(s *server.KeeperServer, lc fx.Lifecycle, log *zap.Logger) {
				s.RegisterHooks(lc, log)
			},
			func(cfg *config.Config, log *zap.Logger) {
				if err := migration.RunMigrations(
					cfg.GetPostgresConnectionString(),
					MigrationsPath,
					log,
				); err != nil {
					log.Fatal("failed to run migrations", zap.Error(err))
				}
			},
		),
	).Run()
}
