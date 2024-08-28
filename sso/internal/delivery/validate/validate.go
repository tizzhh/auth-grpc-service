package validate

import (
	ssvo1 "github.com/tizzhh/auth-grpc-service/protos/gen/go/sso"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	emptyAppId  = 0
	emptyUserId = 0
)

func ValidateLoginRequest(req *ssvo1.LoginRequest) error {
	if req.GetEmail() == "" {
		return status.Error(codes.InvalidArgument, "email is required")
	}
	if req.GetPassword() == "" {
		return status.Error(codes.InvalidArgument, "password is required")
	}
	if req.GetAppId() == emptyAppId {
		return status.Error(codes.InvalidArgument, "app_id is required")
	}
	return nil
}

func ValidateRegisterRequest(req *ssvo1.RegisterRequest) error {
	if req.GetEmail() == "" {
		return status.Error(codes.InvalidArgument, "email is required")
	}
	if req.GetPassword() == "" {
		return status.Error(codes.InvalidArgument, "password is required")
	}
	return nil
}

func ValidateIsAdminRequest(req *ssvo1.IsAdminRequest) error {
	if req.GetUserId() == emptyUserId {
		return status.Error(codes.InvalidArgument, "user_id is required")
	}
	return nil
}
