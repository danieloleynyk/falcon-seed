package user

import (
	"falcon-seed/pkg/auth/rbac"
	"time"
)

// User represents user domain model
type User struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Username  string `json:"username"`
	Email     string `json:"email"`

	LastLogin          time.Time `json:"last_login,omitempty"`
	LastPasswordChange time.Time `json:"last_password_change,omitempty"`

	Token string `json:"-"`

	Role   *rbac.Role      `json:"role,omitempty"`
	RoleID rbac.AccessRole `json:"-"`
}
