package main

import (
	"context"

	"github.com/AFK068/antiplagiarism/analyzer-service/internal/application/services"
	"github.com/AFK068/antiplagiarism/analyzer-service/internal/config"
	"github.com/AFK068/antiplagiarism/analyzer-service/internal/domain"
	"github.com/AFK068/antiplagiarism/analyzer-service/internal/infrastructure/clients/analyzer"
	"github.com/AFK068/antiplagiarism/analyzer-service/internal/infrastructure/httpapi/analyzerapi"
	"github.com/AFK068/antiplagiarism/analyzer-service/internal/migration"

	"github.com/AFK068/antiplagiarism/analyzer-service/internal/infrastructure/repository/repository/postgresdb"
	"github.com/AFK068/antiplagiarism/analyzer-service/internal/infrastructure/server"
	"github.com/AFK068/antiplagiarism/analyzer-service/pkg/logger"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/lib/pq"
	"go.uber.org/fx"
	"go.uber.org/zap"
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

			// Client.
			func(cfg *config.Config, log *zap.Logger) *analyzer.Client {
				return analyzer.NewClient(cfg.Analyzer.KeeperURL, log)
			},

			// Service.
			services.NewAnalyzerService,

			// Handler.
			analyzerapi.NewAnalyzerHandler,

			// Server.
			server.NewAnalyzerServer,
		),
		fx.Invoke(
			func(s *server.AnalyzerServer, lc fx.Lifecycle, log *zap.Logger) {
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
