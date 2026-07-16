package jwt

import (
	"time"

	"github.com/golang-jwt/jwt/v5"

	appjwt "github.com/antongolenev23/voltake-services/pkg/jwt"
	"github.com/antongolenev23/voltake-services/services/auth/internal/config"
	"github.com/antongolenev23/voltake-services/services/auth/internal/domain"
)

func GenerateToken(user domain.User, JWTcfg *config.ConfigJWT) (string, error) {
	claims := appjwt.Claims{
		UserID:  user.ID,
		Email:   user.Email,
		IsAdmin: user.IsAdmin,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(JWTcfg.TTL))),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(JWTcfg.Secret))
}
