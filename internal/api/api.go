package api

import (
	"falcon-seed/pkg/config"
	"falcon-seed/pkg/logger/zap"
	"falcon-seed/pkg/server/echo"
)

func Start(config *config.Configuration) error {
	logger, err := zap.New(config.Logging)
	if err != nil {
		return err
	}

	logger.Info("initialized logger successfully")

	server := echo.New(logger)

	server.Start(&echo.Config{
		Port:                config.Server.Port,
		ReadTimeoutSeconds:  config.Server.ReadTimeout,
		WriteTimeoutSeconds: config.Server.WriteTimeout,
		Debug:               config.Server.Debug,
	})

	return nil
}
