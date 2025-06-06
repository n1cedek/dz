package auth

import (
	jwtu "dz/auth/internal/jwtutils"
	"dz/auth/internal/model"
	"dz/auth/internal/service/jwt_service"
	"time"
)

const (
	refreshTokenSecretKey  = "W4/X+LLjehdxptt4YgGFCvMpq5ewptpZZYRHY6A72g0="
	accessTokenSecretKey   = "VqvguGiffXILza1f44TWXowDT4zwf03dtXmqWW4SYyE="
	refreshTokenExpiration = 60 * time.Minute
	accessTokenExpiration  = 1 * time.Minute
)

type authService struct {
}

func NewAuthService() jwt_service.AuthService {
	return &authService{}
}

func (a *authService) Login(username string) (string, error) {
	return jwtu.GenerateToken(model.UserInfo{
		Username: username,
		Role:     "admin",
	}, []byte(refreshTokenSecretKey), refreshTokenExpiration)
}

func (a *authService) Refresh(refreshToken string) (string, error) {
	claim, err := jwtu.VerifyToken(refreshToken, []byte(refreshTokenSecretKey))
	if err != nil {
		return "", err
	}

	return jwtu.GenerateToken(model.UserInfo{
		Username: claim.Username,
		Role:     "admin",
	}, []byte(refreshTokenSecretKey), refreshTokenExpiration)
}

func (a *authService) GenerateAccess(refreshToken string) (string, error) {
	claim, err := jwtu.VerifyToken(refreshToken, []byte(refreshTokenSecretKey))
	if err != nil {
		return "", err
	}
	return jwtu.GenerateToken(model.UserInfo{
		Username: claim.Username,
		Role:     "admin",
	}, []byte(accessTokenSecretKey), accessTokenExpiration)
}
