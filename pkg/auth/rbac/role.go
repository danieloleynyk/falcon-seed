package rbac

// AccessRole represents access role type
type AccessRole int

const (
	// AdminRole has admin specific permissions
	AdminRole AccessRole = 110

	// UserRole is a standard user
	UserRole AccessRole = 200
)

// Role model
type Role struct {
	ID          int        `json:"id"`
	AccessLevel AccessRole `json:"access_level"`
	Name        string     `json:"name"`
}
