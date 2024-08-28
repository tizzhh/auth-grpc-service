package main

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/tizzhh/auth-grpc-service/sso/internal/app"
	"github.com/tizzhh/auth-grpc-service/sso/internal/config"
	"github.com/tizzhh/auth-grpc-service/sso/pkg/logger/sl"
)

func main() {
	cfg := config.MustLoad()
	log := sl.GetLogger()
	log.Info("starting application", slog.String("env", cfg.Env))

	application := app.New(log, cfg.Grpc.Port, cfg.StoragePath, cfg.TokenTTL)
	go application.GrpcServer.MustRun()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	<-stop

	application.GrpcServer.Stop()
	log.Info("application stopped")
}
