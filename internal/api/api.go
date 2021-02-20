package api

import (
	"falcon-seed/pkg/config"
	"falcon-seed/pkg/logger/zap"
)

func Start(config *config.Configuration) error {
	log := zap.New(config.Logging)
	for {
		log.Info("initialized logger successfully")
	}

	return nil
}
