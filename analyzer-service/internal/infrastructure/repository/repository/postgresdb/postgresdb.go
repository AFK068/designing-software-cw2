package postgresdb

import (
	"context"
	"errors"

	"github.com/AFK068/antiplagiarism/analyzer-service/internal/domain"
	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresRepository(pool *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{
		pool: pool,
	}
}

func (r *PostgresRepository) GetAnalysis(ctx context.Context, id uuid.UUID) (*domain.Analysis, error) {
	query, args, err := squirrel.Select("count_words", "count_characters", "is_plagiat", "hash").
		From("analysis").
		Where(squirrel.Eq{"id": id}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return nil, err
	}

	var analysis domain.Analysis

	err = r.pool.QueryRow(ctx, query, args...).Scan(&analysis.WordCount, &analysis.CharacterCount, &analysis.IsPlagiat, &analysis.Hash)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}

		return nil, err
	}

	return &analysis, nil
}

func (r *PostgresRepository) ExistsByHash(ctx context.Context, hash string) (bool, error) {
	query, args, err := squirrel.Select("1").
		From("analysis").
		Where(squirrel.Eq{"hash": hash}).
		Limit(1).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return false, err
	}

	var exists int

	err = r.pool.QueryRow(ctx, query, args...).Scan(&exists)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, nil
		}

		return false, err
	}

	return true, nil
}

func (r *PostgresRepository) SaveAnalysis(ctx context.Context, id uuid.UUID, analysis *domain.Analysis) error {
	query, args, err := squirrel.Insert("analysis").
		Columns("id", "count_words", "count_characters", "is_plagiat", "hash").
		Values(id, analysis.WordCount, analysis.CharacterCount, analysis.IsPlagiat, analysis.Hash).
		Suffix("ON CONFLICT (id) DO UPDATE SET " +
			"count_words = EXCLUDED.count_words, " +
			"count_characters = EXCLUDED.count_characters, " +
			"is_plagiat = EXCLUDED.is_plagiat, " +
			"hash = EXCLUDED.hash").
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return err
	}

	_, err = r.pool.Exec(ctx, query, args...)
	if err != nil {
		return err
	}

	return nil
}
