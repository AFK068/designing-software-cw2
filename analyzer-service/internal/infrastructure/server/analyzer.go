package server

import (
	"context"
	"time"

	analyzertypes "github.com/AFK068/antiplagiarism/analyzer-service/internal/api/openapi/analyzer/v1"
	"github.com/AFK068/antiplagiarism/analyzer-service/internal/infrastructure/httpapi/analyzerapi"
	echoSwagger "github.com/swaggo/echo-swagger"

	"github.com/labstack/echo/v4"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type AnalyzerServer struct {
	Handler *analyzerapi.AnalyzerHandler
	Echo    *echo.Echo
}

func NewAnalyzerServer(handler *analyzerapi.AnalyzerHandler) *AnalyzerServer {
	return &AnalyzerServer{
		Handler: handler,
		Echo:    echo.New(),
	}
}

func (s *AnalyzerServer) Start() error {
	s.Echo.GET("/swagger/*", echoSwagger.WrapHandler)

	analyzertypes.RegisterHandlers(s.Echo, s.Handler)

	return s.Echo.Start(":" + "8082")
}

func (s *AnalyzerServer) Stop() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	return s.Echo.Shutdown(ctx)
}

func (s *AnalyzerServer) RegisterHooks(lc fx.Lifecycle, log *zap.Logger) {
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
