package user

import "falcon-seed/pkg/user"

// Service represents user service
type Service struct {
}

// GetUser gets a single user
func (service Service) GetUser(string) user.User {
	return user.User{}
}

// NewService creates new user service
func NewService() Service {
	return Service{}
}
