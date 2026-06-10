package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/antongolenev23/voltake-services/services/auth/internal/config"
	"github.com/antongolenev23/voltake-services/services/auth/internal/domain/models"
	"github.com/antongolenev23/voltake-services/services/auth/internal/jwt"
	"github.com/antongolenev23/voltake-services/services/auth/internal/storage"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrUserNotFound = errors.New("user not found")
	ErrInvalidPassword = errors.New("invalid password")
)

type UserSaver interface {
	SaveUser(
		ctx context.Context,
		email string,
		passHash []byte,
	) (models.User, error)
}

type UserProvider interface {
	GetUser(ctx context.Context, email string) (models.User, error)
}

type Auth struct {
	jwtCfg *config.ConfigJWT
	userSaver UserSaver
	UserProvider UserProvider
}

// New returns a new instance of Auth service
func New(
	jwtCfg *config.ConfigJWT,
	userSaver UserSaver,
	userProvider UserProvider,
) *Auth {
	return &Auth{
		jwtCfg: jwtCfg,
		userSaver: userSaver,
		UserProvider: userProvider,
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

	user, err := a.userSaver.SaveUser(ctx, email, passHash)
	if err != nil {
		if errors.Is(err, storage.ErrUserAlreadyExists) {
			return "", ErrUserAlreadyExists
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

	user, err := a.UserProvider.GetUser(ctx, email)
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			return "", ErrUserNotFound
		}
		return "", fmt.Errorf("%s, %w", op, err)
	}

	if err := bcrypt.CompareHashAndPassword(user.PassHash, []byte(password)); err != nil {
		return "", ErrInvalidPassword
	}

	token, err := jwt.GenerateToken(user, a.jwtCfg)
	if err != nil {
		return "", fmt.Errorf("%s, %w", op, err)
	}

	return token, nil
}

