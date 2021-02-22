package auth

import (
	"errors"
	"falcon-seed/pkg/auth/jwt"
	"falcon-seed/pkg/auth/rbac"
	"github.com/labstack/echo/v4"
)

// Service represents auth application service
type Service struct {
	tokenGenerator TokenGenerator
	rbac           RBAC
	refreshTokens  map[string]string
}

// NewService creates new auth service
func NewService(tokenGenerator TokenGenerator, rbac RBAC) Service {
	return Service{
		tokenGenerator: tokenGenerator,
		rbac:           rbac,
		refreshTokens:  map[string]string{},
	}
}

// TokenGenerator represents token generator (jwt) interface
type TokenGenerator interface {
	GenerateAccessToken(jwt.User) (string, error)
	GenerateRefreshToken(string) (string, error)
}

// RBAC represents role-based-access-control interface
type RBAC interface {
	GetUser(echo.Context) rbac.User
}

func (service Service) Authenticate(context echo.Context, user string, password string) (jwt.AuthToken, error) {
	accessToken, err := service.tokenGenerator.GenerateAccessToken(struct {
		Username    string
		Email       string
		AccessLevel rbac.AccessRole
	}{
		Username:    user,
		Email:       user + "@gmail.com",
		AccessLevel: rbac.UserRole,
	})
	if err != nil {
		return jwt.AuthToken{}, errors.New("unauthorized")
	}

	refreshToken, err := service.tokenGenerator.GenerateRefreshToken(user)
	if err != nil {
		return jwt.AuthToken{}, errors.New("unauthorized")
	}

	service.refreshTokens[refreshToken] = user

	return jwt.AuthToken{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (service Service) Refresh(context echo.Context, refreshToken string) (string, error) {
	u, exist := service.refreshTokens[refreshToken]

	if !exist {
		return "", errors.New("refresh token doesn't exist")
	}

	token, err := service.tokenGenerator.GenerateAccessToken(struct {
		Username    string
		Email       string
		AccessLevel rbac.AccessRole
	}{
		Username:    u,
		Email:       u + "@gmail.com",
		AccessLevel: rbac.UserRole,
	})

	if err != nil {
		return "", err
	}

	return token, nil
}
