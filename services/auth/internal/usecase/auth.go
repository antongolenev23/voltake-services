package usecase

import (
	"context"
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"

	"github.com/antongolenev23/voltake-services/services/auth/internal/config"
	"github.com/antongolenev23/voltake-services/services/auth/internal/domain"
	"github.com/antongolenev23/voltake-services/services/auth/internal/jwt"
)

type Storage interface {
	SaveUser(
		ctx context.Context,
		email string,
		passHash []byte,
	) (domain.User, error)

	GetUser(ctx context.Context, email string) (domain.User, error)
}

type Auth struct {
	jwtCfg  *config.ConfigJWT
	storage Storage
}

func New(
	jwtCfg *config.ConfigJWT,
	storage Storage,
) *Auth {
	return &Auth{
		jwtCfg:  jwtCfg,
		storage: storage,
	}
}

func (a *Auth) Register(
	ctx context.Context,
	email string,
	password string,
) (string, error) {
	const op = "service.Register"

	passHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	user, err := a.storage.SaveUser(ctx, email, passHash)
	if err != nil {
		if errors.Is(err, domain.ErrUserAlreadyExists) {
			return "", err
		}
		return "", fmt.Errorf("%s: %w", op, err)
	}

	token, err := jwt.GenerateToken(user, a.jwtCfg)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return token, nil
}

func (a *Auth) Login(
	ctx context.Context,
	email string,
	password string,
) (string, error) {
	const op = "service.Login"

	user, err := a.storage.GetUser(ctx, email)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return "", err
		}
		return "", fmt.Errorf("%s, %w", op, err)
	}

	if err := bcrypt.CompareHashAndPassword(user.PassHash, []byte(password)); err != nil {
		return "", domain.ErrInvalidPassword
	}

	token, err := jwt.GenerateToken(user, a.jwtCfg)
	if err != nil {
		return "", fmt.Errorf("%s, %w", op, err)
	}

	return token, nil
}
