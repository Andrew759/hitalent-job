package config

import (
	"hitalent/cmd/config/dto"
	enviromentNames "hitalent/pkg/config"

	"github.com/spf13/viper"
)

type AppConfigurationInterface interface {
	NewAppConfiguration() AppConfiguration
}

type AppConfiguration struct {
	Environment string
	dto.DatabaseConfig
}

func (c AppConfiguration) NewAppConfiguration() AppConfiguration {
	return AppConfiguration{
		DatabaseConfig: PrepareDatabaseConfig(),
	}
}

func PrepareDatabaseConfig() dto.DatabaseConfig {
	dbc := dto.DatabaseConfig{}

	dbc.SetHost(viper.GetString(enviromentNames.DbHost))
	dbc.SetPort(viper.GetInt(enviromentNames.DbPort))
	dbc.SetName(viper.GetString(enviromentNames.DbName))
	dbc.SetUser(viper.GetString(enviromentNames.DbUser))
	dbc.SetPassword(viper.GetString(enviromentNames.DbPass))
	dbc.SetTimezone(viper.GetString(enviromentNames.DbTimezone))

	return dbc
}
