package auth

import (
	"context"
	"errors"

	ssvo1 "github.com/tizzhh/auth-grpc-service/protos/gen/go/sso"
	"github.com/tizzhh/auth-grpc-service/sso/internal/delivery/validate"
	"github.com/tizzhh/auth-grpc-service/sso/internal/services/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type serverAPI struct {
	ssvo1.UnimplementedAuthServer
	auth Auth
}

func Register(gRPC *grpc.Server, auth Auth) {
	ssvo1.RegisterAuthServer(gRPC, &serverAPI{auth: auth})
}

type Auth interface {
	Login(ctx context.Context, email string, password string, appID int) (string, error)
	RegisterNewUser(ctx context.Context, email string, password string) (int64, error)
	IsAdmin(ctx context.Context, userID int64) (bool, error)
}

func (s *serverAPI) Login(ctx context.Context, req *ssvo1.LoginRequest) (*ssvo1.LoginResponse, error) {
	if err := validate.ValidateLoginRequest(req); err != nil {
		return nil, err
	}

	token, err := s.auth.Login(ctx, req.GetEmail(), req.GetPassword(), int(req.GetAppId()))
	if err != nil {
		if errors.Is(err, auth.ErrInvalidCredentials) {
			return nil, status.Error(codes.InvalidArgument, "invalid email and/or password")
		}
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &ssvo1.LoginResponse{
		Token: token,
	}, nil
}

func (s *serverAPI) Register(ctx context.Context, req *ssvo1.RegisterRequest) (*ssvo1.RegisterResponse, error) {
	if err := validate.ValidateRegisterRequest(req); err != nil {
		return nil, err
	}

	userID, err := s.auth.RegisterNewUser(ctx, req.GetEmail(), req.GetPassword())
	if err != nil {
		if errors.Is(err, auth.ErrUserAlreadyExists) {
			return nil, status.Error(codes.AlreadyExists, "user with this email already exists")
		}
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &ssvo1.RegisterResponse{
		UserId: userID,
	}, nil
}

func (s *serverAPI) IsAdmin(ctx context.Context, req *ssvo1.IsAdminRequest) (*ssvo1.IsAdminResponse, error) {
	if err := validate.ValidateIsAdminRequest(req); err != nil {
		return nil, err
	}

	isAdmin, err := s.auth.IsAdmin(ctx, req.GetUserId())
	if err != nil {
		if errors.Is(err, auth.ErrInvalidCredentials) {
			return nil, status.Error(codes.NotFound, "user not found")
		}
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &ssvo1.IsAdminResponse{
		IsAdmin: isAdmin,
	}, nil
}
