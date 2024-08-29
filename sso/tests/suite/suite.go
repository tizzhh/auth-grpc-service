package suite

import (
	"context"
	"net"
	"strconv"
	"testing"

	ssvo1 "github.com/tizzhh/auth-grpc-service/protos/gen/go/sso"
	"github.com/tizzhh/auth-grpc-service/sso/internal/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Suite struct {
	*testing.T
	Cfg        *config.Config
	AuthClient ssvo1.AuthClient
}

const (
	grpcHost = "localhost"
)

func New(t *testing.T) (context.Context, *Suite) {
	t.Helper()
	t.Parallel()

	cfg := config.MustLoadWithArg("../config/local.yaml")

	ctx, cancel := context.WithTimeout(context.Background(), cfg.Grpc.Timeout)

	t.Cleanup(func() {
		t.Helper()
		cancel()
	})

	grpc, err := grpc.DialContext(context.Background(), grpcAddress(cfg), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("grpc server connection failed: %v", err)
	}

	return ctx, &Suite{
		T:          t,
		Cfg:        cfg,
		AuthClient: ssvo1.NewAuthClient(grpc),
	}
}

func grpcAddress(cfg *config.Config) string {
	return net.JoinHostPort(grpcHost, strconv.Itoa(cfg.Grpc.Port))
}
