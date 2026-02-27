package factory

import (
	"hitalent/cmd/service"
	"hitalent/internal/base"
	"hitalent/internal/middleware"
	internalService "hitalent/internal/service"
	"net/http"
)

func BuildAndServe(dbDecorator service.DBDecorator) {
	mux := buildServer(dbDecorator)

	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		panic(err)
	}
}

func buildServer(dbDecorator service.DBDecorator) *http.ServeMux {
	mux := http.NewServeMux()

	diContainer := base.DIContainer{
		DBDecorator: dbDecorator,
	}
	initOrganizationServer(mux, diContainer)

	return mux
}

func initOrganizationServer(mux *http.ServeMux, diContainer base.DIContainer) {
	departmentController := internalService.DepartmentController{
		Controller: base.Controller{
			ServeMux:     mux,
			Dependencies: diContainer,
		},
		Validator: middleware.DepartmentValidator{
			DBDecorator: diContainer.DBDecorator,
		},
	}

	departmentController.HandleRequest()
}
