package domain

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	GetAnalysis(ctx context.Context, id uuid.UUID) (*Analysis, error)
	SaveAnalysis(ctx context.Context, id uuid.UUID, analysis *Analysis) error
	ExistsByHash(ctx context.Context, hash string) (bool, error)
}
