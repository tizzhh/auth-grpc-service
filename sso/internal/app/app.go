package app

import (
	"log/slog"
	"time"

	grpcapp "github.com/tizzhh/auth-grpc-service/sso/internal/app/grpc"
	"github.com/tizzhh/auth-grpc-service/sso/internal/services/auth"
	"github.com/tizzhh/auth-grpc-service/sso/internal/storage/sqlite"
)

type App struct {
	GrpcServer *grpcapp.App
}

func New(log *slog.Logger, grpcPort int, storagePath string, tokenTTL time.Duration) *App {
	storage, err := sqlite.Get(storagePath)
	if err != nil {
		panic(err)
	}

	authService := auth.New(log, storage, storage, storage, tokenTTL)

	grpcApp := grpcapp.New(log, grpcPort, authService)

	return &App{
		GrpcServer: grpcApp,
	}
}
