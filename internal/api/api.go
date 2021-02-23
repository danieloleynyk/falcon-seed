package api

import (
	"falcon-seed/internal/api/auth"
	"falcon-seed/internal/api/user"
	"falcon-seed/pkg/auth/jwt"
	"falcon-seed/pkg/auth/rbac"
	"falcon-seed/pkg/config"
	"falcon-seed/pkg/logger/zap"
	"falcon-seed/pkg/server"
	"falcon-seed/pkg/server/middleware"
)

func Start(config *config.Configuration) error {
	logger, err := zap.New(config.Logging)
	if err != nil {
		return err
	}

	logger.Info("initialized logger successfully")

	rbacService := rbac.Service{}

	jwtService, err := jwt.New(config.JWT)
	if err != nil {
		logger.Error(err)
		return err
	}

	s := server.New()
	v1 := s.Group("/v1", middleware.Auth(jwtService))

	auth.NewHTTP(auth.NewService(jwtService, rbacService), s.Echo)
	user.RegisterHandlers(v1)

	s.Start(&server.Config{
		Port:                config.Server.Port,
		ReadTimeoutSeconds:  config.Server.ReadTimeout,
		WriteTimeoutSeconds: config.Server.WriteTimeout,
		Debug:               config.Server.Debug,
	})

	return nil
}
