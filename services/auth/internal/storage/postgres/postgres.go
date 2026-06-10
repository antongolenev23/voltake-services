package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/antongolenev23/voltake-services/services/auth/internal/config"
	"github.com/antongolenev23/voltake-services/services/auth/internal/domain/models"
	"github.com/antongolenev23/voltake-services/services/auth/internal/storage"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

func NewPostgres(ctx context.Context, cfg *config.ConfigStorage) (*pgxpool.Pool, error) {
	const op = "storage.NewPostgres"

	connStr := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=%s",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Name,
		cfg.SSLMode,
	)

	pool, err := pgxpool.New(ctx, connStr)
	if err != nil {
		return nil, fmt.Errorf("%s: %v", op, err)
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("%s: %v", op, err)
	}

	return pool, nil
}

type Storage struct {
	db *pgxpool.Pool
}

func New(db *pgxpool.Pool) *Storage {
	return &Storage{db: db}
}

func (s *Storage) SaveUser(
	ctx context.Context,
	email string,
	passHash []byte,
) (models.User, error) {
	const op = "postgres.SaveUser"

	user := models.User{
		ID:       uuid.New(),
		Email:    email,
		PassHash: passHash,
		IsAdmin:  false,
	}

	query := `
		INSERT INTO users (id, email, pass_hash, is_admin)
		VALUES ($1, $2, $3, $4)
		RETURNING id, email, pass_hash, is_admin
	`

	err := s.db.QueryRow(
		ctx,
		query,
		user.ID,
		user.Email,
		user.PassHash,
		user.IsAdmin,
	).Scan(
		&user.ID,
		&user.Email,
		&user.PassHash,
		&user.IsAdmin,
	)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" {
				return models.User{}, storage.ErrUserAlreadyExists
			}
		}

		return models.User{}, fmt.Errorf("%s: %w", op, err)
	}

	return user, nil
}

func (s *Storage) GetUser(
	ctx context.Context,
	email string,
) (models.User, error) {
	const op = "postgres.GetUser"

	var user models.User

	query := `
		SELECT id, email, pass_hash, is_admin
		FROM users
		WHERE email = $1
	`

	err := s.db.QueryRow(ctx, query, email).Scan(
		&user.ID,
		&user.Email,
		&user.PassHash,
		&user.IsAdmin,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.User{}, storage.ErrUserNotFound
		}

		return models.User{}, fmt.Errorf("%s: %w", op, err)
	}
	return user, nil
}
