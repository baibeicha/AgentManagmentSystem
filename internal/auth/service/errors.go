package service

import (
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserNotFound       = errors.New("user not found")
	ErrResourceNotFound   = errors.New("resource not found")
	ErrInternalError      = errors.New("internal server error")
	ErrInvalidToken       = errors.New("invalid or expired token")
	ErrTokenReused        = errors.New("session rotation failed or token reused")
	ErrUserSuspended      = errors.New("user account is suspended")
	ErrInvalidArgument    = errors.New("invalid argument")
)

func MapError(err error) error {
	if err == nil {
		return nil
	}

	switch {
	case errors.Is(err, ErrInvalidCredentials):
		return status.Error(codes.Unauthenticated, err.Error())
	case errors.Is(err, ErrUserNotFound):
		return status.Error(codes.Unauthenticated, err.Error())
	case errors.Is(err, ErrInvalidToken):
		return status.Error(codes.Unauthenticated, err.Error())
	case errors.Is(err, ErrTokenReused):
		return status.Error(codes.Unauthenticated, err.Error())
	case errors.Is(err, ErrResourceNotFound):
		return status.Error(codes.NotFound, err.Error())
	case errors.Is(err, ErrInvalidArgument):
		return status.Error(codes.InvalidArgument, err.Error())
	case errors.Is(err, ErrInternalError):
		return status.Error(codes.Internal, err.Error())
	default:
		return status.Error(codes.Internal, err.Error())
	}
}
