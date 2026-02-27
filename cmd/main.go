package main

import (
	"hitalent/cmd/config"
	"hitalent/cmd/factory"
	"hitalent/cmd/service"
)

func main() {
	factory.InitViper()

	appConfig := config.AppConfiguration{}.NewAppConfiguration()
	dbDecorator := service.InitORM(&appConfig.DatabaseConfig)

	defer dbDecorator.CloseDB()

	factory.BuildAndServe(dbDecorator)
}
