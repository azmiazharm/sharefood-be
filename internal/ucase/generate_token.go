package ucase

import (
	"sharefood/internal/entity"
	"time"

	"github.com/dgrijalva/jwt-go/v4"
)

func GenerateJWT(user entity.User, secret []byte) (entity.TokenResponse, error) {
	expiredAt := time.Now().Add(time.Hour * 2).Local()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, entity.TokenClaims{
		ID: user.ID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: &jwt.Time{
				Time: expiredAt,
			},
		},
	})

	signedToken, errToken := token.SignedString(secret)
	if errToken != nil {
		return entity.TokenResponse{}, errToken
	}

	tokenResponse := entity.TokenResponse{
		Type:      "bearer",
		Token:     signedToken,
		ExpiredAt: expiredAt,
	}

	return tokenResponse, nil
}
