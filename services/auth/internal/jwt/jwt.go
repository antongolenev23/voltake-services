package jwt

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/antongolenev23/voltake-services/services/auth/internal/config"
	"github.com/antongolenev23/voltake-services/services/auth/internal/domain/models"

	"github.com/google/uuid"

)

type Claims struct {
	UserID uuid.UUID
	IsAdmin bool
	jwt.RegisteredClaims
}

func GenerateToken(user models.User, JWTcfg *config.ConfigJWT) (string, error) {
	claims := Claims{
		UserID: user.ID,
		IsAdmin: user.IsAdmin,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(JWTcfg.TTL))),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(JWTcfg.Secret))
}
