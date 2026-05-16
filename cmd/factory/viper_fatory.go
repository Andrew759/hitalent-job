package factory

import (
	"errors"
	"fmt"

	"github.com/spf13/viper"
)

func InitViper() {
	viper.SetConfigFile("/app/.env")
	readConfig()
}

func readConfig() {
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := errors.AsType[viper.ConfigFileNotFoundError](err); ok {
			panic(fmt.Errorf("hitalent config file not found: %w", err))
		}
		panic(fmt.Errorf("viper fatal error: %w", err))
	}
}
