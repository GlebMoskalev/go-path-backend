package repository

import (
	"context"
	"errors"
	"time"

	"github.com/GlebMoskalev/go-path-backend/internal/model"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	UserNotFound = errors.New("user not found")
)

type UserRepository interface {
	Create(ctx context.Context, user *model.User) error
	GetByID(ctx context.Context, id uuid.UUID) (*model.User, error)
	GetByEmail(ctx context.Context, email string) (*model.User, error)
	GetByGoogleID(ctx context.Context, googleID string) (*model.User, error)
	Update(ctx context.Context, user *model.User) error
	UpdateLastLogin(ctx context.Context, userID uuid.UUID) error
	IncrementTokenVersion(ctx context.Context, userID uuid.UUID) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type userRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, user *model.User) error {
	query := `
	INSERT INTO users (id, email, name, picture, google_id, token_version, is_active, created_at, updated_at)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`

	_, err := r.db.Exec(ctx, query,
		user.ID, user.Email, user.Name, user.Picture, user.GoogleID, user.TokenVersion, user.IsActive, user.CreatedAt, user.UpdatedAt,
	)

	return err
}

func (r *userRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.User, error) {
	query := `
	SELECT id, email, name, picture, google_id, token_version, is_active, last_login_at, created_at, updated_at
	FROM users
	WHERE id = $1 and is_active = TRUE
	`

	user := &model.User{}
	err := r.db.QueryRow(ctx, query, id).Scan(
		&user.ID,
		&user.Email,
		&user.Name,
		&user.Picture,
		&user.GoogleID,
		&user.TokenVersion,
		&user.IsActive,
		&user.LastLoginAt,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, UserNotFound
		}
		return nil, err
	}
	return user, nil
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	query := `
	SELECT id, email, name, picture, google_id, token_version, is_active, last_login_at, created_at, updated_at
	FROM users
	WHERE email = $1 AND is_active = true
	`

	user := &model.User{}
	err := r.db.QueryRow(ctx, query, email).Scan(
		&user.ID,
		&user.Email,
		&user.Name,
		&user.Picture,
		&user.GoogleID,
		&user.TokenVersion,
		&user.IsActive,
		&user.LastLoginAt,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, UserNotFound
		}
		return nil, err
	}

	return user, nil
}

func (r *userRepository) GetByGoogleID(ctx context.Context, googleID string) (*model.User, error) {
	query := `
	SELECT id, email, name, picture, google_id, token_version, is_active, last_login_at, created_at, updated_at
	FROM users
	WHERE google_id = $1 AND is_active = true
	`

	user := &model.User{}
	err := r.db.QueryRow(ctx, query, googleID).Scan(
		&user.ID,
		&user.Email,
		&user.Name,
		&user.Picture,
		&user.GoogleID,
		&user.TokenVersion,
		&user.IsActive,
		&user.LastLoginAt,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, UserNotFound
		}
		return nil, err
	}

	return user, nil
}

func (r *userRepository) Update(ctx context.Context, user *model.User) error {
	query := `
    UPDATE users
    SET name = $2, picture = COALESCE(NULLIF($3, ''), picture), updated_at = $4
    WHERE id = $1
    `
	user.UpdatedAt = time.Now()

	_, err := r.db.Exec(ctx, query, user.ID, user.Name, user.Picture, user.UpdatedAt)

	return err
}

func (r *userRepository) UpdateLastLogin(ctx context.Context, userID uuid.UUID) error {
	query := `
	UPDATE users
	SET last_login_at = $2, updated_at = $3
	WHERE id = $1
	`

	now := time.Now()
	_, err := r.db.Exec(ctx, query, userID, now, now)

	return err
}

func (r *userRepository) IncrementTokenVersion(ctx context.Context, userID uuid.UUID) error {
	query := `
	UPDATE users
	SET token_version = token_version + 1, updated_at = $2
	WHERE id = $1
	`

	_, err := r.db.Exec(ctx, query, userID, time.Now())

	return err
}

func (r *userRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `
	UPDATE users
	SET is_active = false, updated_at = $2
	WHERE id = $1
	`

	_, err := r.db.Exec(ctx, query, id, time.Now())

	return err
}
