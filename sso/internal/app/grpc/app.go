package grpcapp

import (
	"fmt"
	"log/slog"
	"net"

	authgrpc "github.com/tizzhh/auth-grpc-service/sso/internal/delivery/grpc/auth"
	"google.golang.org/grpc"
)

type App struct {
	log        *slog.Logger
	gRPCServer *grpc.Server
	port       int
}

func New(log *slog.Logger, port int, authService authgrpc.Auth) *App {
	gRPCServer := grpc.NewServer()
	authgrpc.Register(gRPCServer, authService)

	return &App{
		log:        log,
		gRPCServer: gRPCServer,
		port:       port,
	}
}

func (a *App) MustRun() {
	if err := a.Run(); err != nil {
		panic(err)
	}
}

func (a *App) Run() error {
	const caller = "grpcapp.Run"

	log := a.log.With(slog.String("caller", caller))

	l, err := net.Listen("tcp", fmt.Sprintf(":%d", a.port))
	if err != nil {
		return fmt.Errorf("%s: %w", caller, err)
	}

	log.Info("starting grpc server", slog.String("addr", l.Addr().String()))

	if err := a.gRPCServer.Serve(l); err != nil {
		return fmt.Errorf("%s: %w", caller, err)
	}

	return nil
}

func (a *App) Stop() {
	const caller = "grpcapp.Stop"

	log := a.log.With(slog.String("caller", caller))

	log.Info("stopping gRPC server")

	a.gRPCServer.GracefulStop()
}
