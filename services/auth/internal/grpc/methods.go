package grpc

import (
	"context"
	"errors"
	"log/slog"
	"net/mail"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	authv1 "github.com/antongolenev23/voltake-protos/gen/go/auth/v1"
	passwordvalidator "github.com/wagslane/go-password-validator"

	"github.com/antongolenev23/voltake-services/pkg/logger"

	"github.com/antongolenev23/voltake-services/services/auth/internal/domain"
)

func (s *serverAPI) Register(
	ctx context.Context,
	req *authv1.Credentials,
) (*authv1.AuthResponse, error) {
	const op = "grpc.Register"

	log := logger.WithRequestContext(ctx, s.log, op)

	log.Debug("starting register request processing")

	if err := validateCredentials(req); err != nil {
		log.Info("invalid register credentials", slog.String("error", err.Error()))
		return nil, err
	}

	token, err := s.auth.Register(ctx, req.GetEmail(), req.GetPassword())
	if err != nil {
		if errors.Is(err, domain.ErrUserAlreadyExists) {
			log.Info("can not register", slog.String("error", err.Error()))
			return nil, status.Error(codes.AlreadyExists, "user already exists")
		}
		log.Error("can not register", slog.String("error", err.Error()))
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &authv1.AuthResponse{
		Token: token,
	}, nil
}

func (s *serverAPI) Login(
	ctx context.Context,
	req *authv1.Credentials,
) (*authv1.AuthResponse, error) {
	const op = "grpc.Login"

	log := logger.WithRequestContext(ctx, s.log, op)

	log.Debug("starting login request processing")

	token, err := s.auth.Login(ctx, req.GetEmail(), req.GetPassword())
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) || errors.Is(err, domain.ErrInvalidPassword) {
			log.Info("can not login", slog.String("error", err.Error()))
			return nil, status.Error(codes.Unauthenticated, "invalid credentials")
		}
		log.Error("can not login", slog.String("error", err.Error()))
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &authv1.AuthResponse{
		Token: token,
	}, nil
}

func validateCredentials(req *authv1.Credentials) error {
	if _, err := mail.ParseAddress(req.GetEmail()); err != nil {
		return status.Error(codes.InvalidArgument, "invalid email")
	}

	if len(req.GetPassword()) > 72 {
		return status.Error(codes.InvalidArgument, "password is too long")
	}

	minEntropyBits := 60.0
	if err := passwordvalidator.Validate(req.GetPassword(), minEntropyBits); err != nil {
		return status.Error(codes.InvalidArgument, "weak password")
	}

	return nil
}
