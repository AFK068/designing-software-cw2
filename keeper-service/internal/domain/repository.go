package domain

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	SaveFileData(ctx context.Context, fileData *File) (uuid.UUID, error)
	GetFileData(ctx context.Context, fileID uuid.UUID) (string, error)
}
