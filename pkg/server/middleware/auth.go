package middleware

import (
	"falcon-seed/pkg/auth/rbac"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
)

// TokenParser represents JWT token parser
type TokenParser interface {
	ParseAccessToken(string) (*jwt.Token, error)
	ParseRefreshToken(string) (*jwt.Token, error)
}

// Middleware makes JWT implement the Middleware interface.
func Auth(tokenParser TokenParser) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(context echo.Context) error {
			token, err := tokenParser.ParseAccessToken(context.Request().Header.Get("Authorization"))
			if err != nil || !token.Valid {
				return echo.NewHTTPError(http.StatusUnauthorized)
			}

			claims := token.Claims.(jwt.MapClaims)

			username := claims["user"].(string)
			email := claims["email"].(string)
			role := rbac.AccessRole(claims["role"].(float64))

			context.Set("username", username)
			context.Set("email", email)
			context.Set("role", role)

			return next(context)
		}

	}
}
