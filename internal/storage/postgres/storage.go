package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/lib/pq"

	"sso/internal/domain/models"
	"sso/internal/storage"
)

const (
	ConstructorDbPath = "storage.postgres.New"
	IsAdminPath       = "storage.postgres.IsAdmin"
	UserPath          = "storage.postgres.User"
	AppPath           = "storage.postgres.App"
)

type Storage struct {
	db *sql.DB
}

func New(storagePath string) (*Storage, error) {
	db, err := sql.Open("postgres", storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", ConstructorDbPath, err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) Stop() error {
	return s.db.Close()
}

// SaveUser saves user to db.
func (s *Storage) SaveUser(ctx context.Context, email string, passHash []byte) (int64, error) {
	const op = "storage.postgres.SaveUser"

	stmt, err := s.db.Prepare("INSERT INTO users(email, pass_hash) VALUES($1, $2) RETURNING id")
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	defer stmt.Close()

	var id int64
	err = stmt.QueryRowContext(ctx, email, passHash).Scan(&id)
	if err != nil {
		var pgErr *pq.Error
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return 0, fmt.Errorf("%s: %w", op, storage.ErrUserExists)
		}

		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

// User returns user by email.
func (s *Storage) User(ctx context.Context, email string) (models.User, error) {

	stmt, err := s.db.Prepare("SELECT id, email, pass_hash FROM users WHERE email = $1")
	if err != nil {
		return models.User{}, fmt.Errorf("%s: %w", UserPath, err)
	}

	defer stmt.Close()

	var user models.User
	err = stmt.QueryRowContext(ctx, email).Scan(&user.ID, &user.Email, &user.PassHash)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.User{}, fmt.Errorf("%s: %w", UserPath, storage.ErrUserNotFound)
		}

		return models.User{}, fmt.Errorf("%s: %w", UserPath, err)
	}

	return user, nil
}

func (s *Storage) App(ctx context.Context, id int) (models.App, error) {

	stmt, err := s.db.Prepare("SELECT id, name, secret FROM apps WHERE id = $1")
	if err != nil {
		return models.App{}, fmt.Errorf("%s: %w", AppPath, err)
	}

	defer stmt.Close()

	var app models.App
	err = stmt.QueryRowContext(ctx, id).Scan(&app.ID, &app.Name, &app.Secret)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.App{}, fmt.Errorf("%s: %w", AppPath, storage.ErrAppNotFound)
		}

		return models.App{}, fmt.Errorf("%s: %w", AppPath, err)
	}

	return app, nil
}

func (s *Storage) IsAdmin(ctx context.Context, userID int64) (bool, error) {

	stmt, err := s.db.Prepare("SELECT is_admin FROM users WHERE id = $1")
	if err != nil {
		return false, fmt.Errorf("%s: %w", IsAdminPath, err)
	}

	defer stmt.Close()

	var isAdmin bool
	err = stmt.QueryRowContext(ctx, userID).Scan(&isAdmin)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, fmt.Errorf("%s: %w", IsAdminPath, storage.ErrUserNotFound)
		}

		return false, fmt.Errorf("%s: %w", IsAdminPath, err)
	}

	return isAdmin, nil
}
