package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"sync"

	"github.com/mattn/go-sqlite3"
	"github.com/tizzhh/auth-grpc-service/sso/domain/models"
	"github.com/tizzhh/auth-grpc-service/sso/internal/storage"
)

type Storage struct {
	db *sql.DB
}

var storageInstance *Storage
var once sync.Once

// Get returns a singleton of SQLite storage instance.
func Get(storagePath string) (*Storage, error) {
	var err error
	once.Do(func() {
		storageInstance, err = New(storagePath)
	})
	return storageInstance, err
}

// New creates a new instance of the SQLite storage.
func New(storagePath string) (*Storage, error) {
	const caller = "storage.sqlite.New"

	db, err := sql.Open("sqlite3", storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", caller, err)
	}

	return &Storage{db: db}, nil
}

// SaveUser saves a new user in db.
func (s *Storage) SaveUser(ctx context.Context, email string, passHash []byte) (int64, error) {
	const caller = "storage.sqlite.SaveUser"

	stmt, err := s.db.Prepare(`INSERT INTO users(email, pass_hash) VALUES(?, ?)`)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", caller, err)
	}

	res, err := stmt.ExecContext(ctx, email, passHash)
	if err != nil {
		var sqliteErr sqlite3.Error
		if errors.As(err, &sqliteErr) && sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
			return 0, fmt.Errorf("%s: %w", caller, storage.ErrUserAlreadyExists)
		}
		return 0, fmt.Errorf("%s: %w", caller, err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("%s: %w", caller, err)
	}

	return id, nil
}

// User returns user by email.
func (s *Storage) User(ctx context.Context, email string) (models.User, error) {
	const caller = "storage.sqlite.User"

	stmt, err := s.db.Prepare(`SELECT * FROM users WHERE email = ?`)
	if err != nil {
		return models.User{}, fmt.Errorf("%s: %w", caller, err)
	}

	row := stmt.QueryRowContext(ctx, email)

	var queryUser models.User
	err = row.Scan(&queryUser.ID, &queryUser.Email, &queryUser.PassHash, &queryUser.IsAdmin)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.User{}, fmt.Errorf("%s: %w", caller, storage.ErrUserNotFound)
		}
		return models.User{}, fmt.Errorf("%s: %w", caller, err)
	}

	return queryUser, nil
}

// IsAdmin checks if user is an admin.
func (s *Storage) IsAdmin(ctx context.Context, userID int64) (bool, error) {
	const caller = "storage.sqlite.IsAdmin"

	stmt, err := s.db.Prepare(`SELECT is_admin FROM users WHERE id = ?`)
	if err != nil {
		return false, fmt.Errorf("%s: %w", caller, err)
	}

	row := stmt.QueryRowContext(ctx, userID)

	var isAdmin bool
	err = row.Scan(&isAdmin)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, fmt.Errorf("%s: %w", caller, storage.ErrUserNotFound)
		}
		return false, fmt.Errorf("%s: %w", caller, err)
	}

	return isAdmin, nil
}

func (s *Storage) App(ctx context.Context, appID int) (models.App, error) {
	const caller = "storage.sqlite.App"

	stmt, err := s.db.Prepare(`SELECT * FROM apps WHERE id = ?`)
	if err != nil {
		return models.App{}, fmt.Errorf("%s: %w", caller, err)
	}

	row := stmt.QueryRowContext(ctx, appID)

	var queryApp models.App
	err = row.Scan(&queryApp.ID, &queryApp.Name, &queryApp.Secret)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.App{}, fmt.Errorf("%s: %w", caller, storage.ErrAppNotFound)
		}
		return models.App{}, fmt.Errorf("%s: %w", caller, err)
	}

	return queryApp, nil
}
