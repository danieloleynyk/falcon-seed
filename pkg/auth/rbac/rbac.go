package rbac

import (
	"github.com/labstack/echo/v4"
)

// RBACService represents role-based access control service interface
type RBACService interface {
	GetUser(echo.Context) User
	EnforceRole(echo.Context, AccessRole) error
}

// Service is RBAC application service
type Service struct{}

type User struct {
	ID       int
	Username string
	Email    string
	Role     AccessRole
}

// User returns user data stored in jwt token
func (service Service) GetUser(context echo.Context) User {
	id := context.Get("id").(int)
	user := context.Get("username").(string)
	email := context.Get("email").(string)
	role := context.Get("role").(AccessRole)

	return User{
		ID:       id,
		Username: user,
		Email:    email,
		Role:     role,
	}
}

// EnforceRole authorizes request by AccessRole
func (service Service) EnforceRole(context echo.Context, r AccessRole) error {
	if !(context.Get("role").(AccessRole) > r) {
		return echo.ErrForbidden
	}

	return nil
}
