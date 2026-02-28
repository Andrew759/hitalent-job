package main

import (
	"hitalent/cmd/config"
	"hitalent/cmd/factory"
	"hitalent/cmd/service"
	"hitalent/internal/model"
)

func main() {
	factory.InitViper()

	appConfig := config.AppConfiguration{}.NewAppConfiguration()
	dbDecorator := service.InitORM(&appConfig.DatabaseConfig)

	defer dbDecorator.CloseDB()

	factory.BuildAndServe(dbDecorator)

	//TODO: это мок миграций
	dbDecorator.GormInterface.AutoMigrate(&model.Department{}, &model.Employee{})
}
