package server

import (
	"context"
	"time"

	"github.com/labstack/echo/v4"
	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/AFK068/antiplagiarism/keeper-service/internal/infrastructure/httpapi/keeperapi"

	keepertypes "github.com/AFK068/antiplagiarism/keeper-service/internal/api/openapi/keeper/v1"
)

type KeeperServer struct {
	Handler *keeperapi.KeeperHandler
	Echo    *echo.Echo
}

func NewKeeperServer(handler *keeperapi.KeeperHandler) *KeeperServer {
	return &KeeperServer{
		Handler: handler,
		Echo:    echo.New(),
	}
}

func (s *KeeperServer) Start() error {
	keepertypes.RegisterHandlers(s.Echo, s.Handler)

	return s.Echo.Start(":" + "8081")
}

func (s *KeeperServer) Stop() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	return s.Echo.Shutdown(ctx)
}

func (s *KeeperServer) RegisterHooks(lc fx.Lifecycle, log *zap.Logger) {
	lc.Append(fx.Hook{
		OnStart: func(context.Context) error {
			log.Info("Starting compressor server")

			go func() {
				if err := s.Start(); err != nil {
					log.Error("Failed to start compressor server", zap.Error(err))
				}
			}()

			return nil
		},
		OnStop: func(context.Context) error {
			log.Info("Stopping compressor server")

			if err := s.Stop(); err != nil {
				log.Error("Failed to stop compressor server", zap.Error(err))
			}

			return nil
		},
	})
}
