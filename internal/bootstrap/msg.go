// Package bootstrap
package bootstrap

import (
	"sharefood/internal/consts"
	"sharefood/pkg/logger"
	"sharefood/pkg/msg"
)

func RegistryMessage() {
	err := msg.Setup("msg.yaml", consts.ConfigPath)
	if err != nil {
		logger.Fatal(logger.MessageFormat("file message multi language load error %s", err.Error()))
	}

}
