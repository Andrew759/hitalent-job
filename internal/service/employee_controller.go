package service

import (
	"hitalent/internal/base"
	"hitalent/internal/middleware"
	"net/http"
)

type EmployeeController struct {
	Controller base.Controller
	middleware.Validator
}

func (ec *EmployeeController) HandleRequest() {
	ec.Controller.ServeMux.HandleFunc("POST /departments/{id}/employees",
		ec.Validate(func(w http.ResponseWriter, r *http.Request) {
			ec.CreateEmployee(w, base.NewRequest(r))
		}))
}

func (ec *EmployeeController) CreateEmployee(
	w http.ResponseWriter,
	request *base.Request,
) {

}
