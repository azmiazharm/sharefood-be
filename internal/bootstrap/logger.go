// Package bootstrap
package bootstrap

import (
	"sharefood/internal/appctx"
	"sharefood/pkg/logger"
	"sharefood/pkg/util"
)

func RegistryLogger(cfg *appctx.Config) {
	logger.Setup(logger.Config{
		Environment: util.EnvironmentTransform(cfg.App.Env),
		Debug:       cfg.App.Debug,
		Level:       cfg.Logger.Level,
		ServiceName: cfg.Logger.Name,
	})
}
