package service

import (
	"library-go/pkg/jwt"
	"library-go/pkg/logging"
)

type authService struct {
	logger *logging.Logger
}

func NewAuthService(logger *logging.Logger) AuthService {
	return &authService{
		logger: logger,
	}
}

func (a *authService) VerifyToken(token string) (bool, error) {
	return jwt.VerifyToken(token)
}
