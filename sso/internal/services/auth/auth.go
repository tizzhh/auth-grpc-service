package auth

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/tizzhh/auth-grpc-service/sso/domain/models"
	"github.com/tizzhh/auth-grpc-service/sso/internal/storage"
	"github.com/tizzhh/auth-grpc-service/sso/pkg/jwt"
	"github.com/tizzhh/auth-grpc-service/sso/pkg/logger/sl"
	"golang.org/x/crypto/bcrypt"
)

type UserSaver interface {
	SaveUser(ctx context.Context, email string, passHash []byte) (int64, error)
}

type UserProvider interface {
	User(ctx context.Context, email string) (models.User, error)
	IsAdmin(ctx context.Context, userID int64) (bool, error)
}

type AppProvider interface {
	App(ctx context.Context, appID int) (models.App, error)
}

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserAlreadyExists  = errors.New("user already exists")
)

// New returns a new instance of the Auth service.
func New(log *slog.Logger, userSaver UserSaver, userProvider UserProvider, appProvider AppProvider, tokenTTL time.Duration) *Auth {
	return &Auth{
		log:          log,
		userSaver:    userSaver,
		userProvider: userProvider,
		appProvider:  appProvider,
		tokenTTL:     tokenTTL,
	}
}

type Auth struct {
	log          *slog.Logger
	userSaver    UserSaver
	userProvider UserProvider
	appProvider  AppProvider
	tokenTTL     time.Duration
}

// Login checks the user with given credentials.

// If user exists, but password is incorrect, returns an error.
// If user doesn't exist, returns error.
func (a *Auth) Login(ctx context.Context, email string, password string, appID int) (string, error) {
	const caller = "services.auth.Login"

	log := a.log.With(
		slog.String("caller", caller),
		slog.String("email", email),
	)

	log.Info("attempting to login user")

	user, err := a.userProvider.User(ctx, email)
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			log.Warn("user not found", sl.Err(err))
			return "", fmt.Errorf("%s: %w", caller, ErrInvalidCredentials)
		}
		log.Error("failed to get user", sl.Err(err))
		return "", fmt.Errorf("%s: %w", caller, err)
	}

	if err := bcrypt.CompareHashAndPassword(user.PassHash, []byte(password)); err != nil {
		log.Info("invalid credentials", sl.Err(err))
		return "", fmt.Errorf("%s: %w", caller, ErrInvalidCredentials)
	}

	app, err := a.appProvider.App(ctx, appID)
	if err != nil {
		return "", fmt.Errorf("%s: %w", caller, err)
	}

	log.Info("user logged in successfully")

	token, err := jwt.NewToken(user, app, a.tokenTTL)
	if err != nil {
		log.Error("failed to generate token", sl.Err(err))
		return "", fmt.Errorf("%s: %w", caller, err)
	}

	return token, nil
}

// RegisterNewUser registers new user in the system and returns user ID.
// If user with given email already exists, returns an error.
func (a *Auth) RegisterNewUser(ctx context.Context, email string, password string) (int64, error) {
	const caller = "auth.RegisterNewUser"

	log := a.log.With(
		slog.String("caller", caller),
		slog.String("email", email),
	)

	log.Info("registering user")

	passHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Error("failed to generate password hash", sl.Err(err))
		return 0, fmt.Errorf("%s: %w", caller, err)
	}

	id, err := a.userSaver.SaveUser(ctx, email, passHash)
	if err != nil {
		if errors.Is(err, storage.ErrUserAlreadyExists) {
			log.Warn("user already exists", sl.Err(err))
			return 0, fmt.Errorf("%s: %w", caller, ErrUserAlreadyExists)
		}
		log.Error("failed to save user", sl.Err(err))
		return 0, fmt.Errorf("%s: %w", caller, err)
	}

	log.Info("user registered")

	return id, nil
}

// IsAdmin checks if user is admin.
func (a *Auth) IsAdmin(ctx context.Context, userID int64) (bool, error) {
	const caller = "auth.IsAdmin"

	log := a.log.With(
		slog.String("caller", caller),
		slog.Int64("user_id", userID),
	)

	log.Info("checking if user is admin")

	isAdmin, err := a.userProvider.IsAdmin(ctx, userID)
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			log.Warn("user not found", sl.Err(err))
			return false, fmt.Errorf("%s: %w", caller, ErrInvalidCredentials)
		}
		return false, fmt.Errorf("%s: %w", caller, err)
	}

	log.Info("checked if user is admin", slog.Bool("is_admin", isAdmin))

	return isAdmin, nil
}
