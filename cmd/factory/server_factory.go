package factory

import (
	"hitalent/cmd/service"
	"hitalent/internal/base"
	internalService "hitalent/internal/controller"
	"hitalent/internal/middleware"
	"net/http"
)

func BuildAndServe(dbDecorator service.DBDecorator) {
	mux := BuildServer(dbDecorator)

	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		panic(err)
	}
}

func BuildServer(dbDecorator service.DBDecorator) *http.ServeMux {
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
		CreateDV: middleware.CreateDepartmentValidator{
			DBDecorator: diContainer.DBDecorator,
		},
		ChangeDV: middleware.ChangeDepartmentValidator{
			DBDecorator: diContainer.DBDecorator,
		},
	}
	employeeController := internalService.EmployeeController{
		Controller: base.Controller{
			ServeMux:     mux,
			Dependencies: diContainer,
		},
		Validator: middleware.EmployeeValidator{
			DBDecorator: diContainer.DBDecorator,
		},
	}

	departmentController.HandleRequest()
	employeeController.HandleRequest()
}
