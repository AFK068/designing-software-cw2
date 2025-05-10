package postgresdb

import (
	"context"
	"errors"

	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/AFK068/antiplagiarism/keeper-service/internal/domain"
)

type PostgresRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresRepository(pool *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{
		pool: pool,
	}
}

func (r *PostgresRepository) SaveFileData(ctx context.Context, fileData *domain.File) (uuid.UUID, error) {
	var fileID uuid.UUID

	query, args, err := squirrel.
		Insert("files").
		Columns("data", "name", "hash", "location").
		Values(fileData.Data, fileData.Name, fileData.Hash, fileData.Location).
		Suffix("ON CONFLICT (hash) DO UPDATE SET hash = EXCLUDED.hash RETURNING id").
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return uuid.Nil, err
	}

	err = r.pool.QueryRow(ctx, query, args...).Scan(&fileID)
	if err != nil {
		return uuid.Nil, err
	}

	return fileID, nil
}

func (r *PostgresRepository) GetFileData(ctx context.Context, fileID uuid.UUID) (string, error) {
	query, args, err := squirrel.
		Select("data").
		From("files").
		Where(squirrel.Eq{"id": fileID}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return "", err
	}

	var data string

	err = r.pool.QueryRow(ctx, query, args...).Scan(&data)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", errors.New("file not found")
		}

		return "", err
	}

	return data, nil
}
