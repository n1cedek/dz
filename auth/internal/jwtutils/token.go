package jwtu

import (
	"dz/auth/internal/model"
	"github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"
	"time"
)

func GenerateToken(info model.UserInfo, secretKey []byte, duration time.Duration) (string, error) {
	claim := model.UserClaim{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(duration).Unix(),
		},
		Username: info.Username,
		Role:     info.Role,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)

	return token.SignedString(secretKey)
}

func VerifyToken(tokenStr string, secretKey []byte) (*model.UserClaim, error) {

	token, err := jwt.ParseWithClaims(tokenStr, &model.UserClaim{},
		func(token *jwt.Token) (interface{}, error) {
			_, ok := token.Method.(*jwt.SigningMethodHMAC)
			if !ok {
				return nil, errors.New("unexpected token signing method")
			}

			return secretKey, nil
		},
	)
	if err != nil {
		return nil, errors.Errorf("invalid token: %v", err.Error())
	}

	claim, ok := token.Claims.(*model.UserClaim)
	if !ok {
		return nil, errors.New("invalid token claims")
	}

	return claim, nil

}
