package controller

import (
	"hitalent/internal/base"
	"hitalent/internal/middleware"
	"hitalent/internal/model"
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
	r *base.Request,
) {
	e := r.Context().Value(middleware.CreateEmployeeKey).(*model.Employee)

	err := model.CreateEmployee(ec.Controller.Dependencies.DBDecorator.GormInterface, e)
	if err != nil {
		base.NewResponse().SendError(w, "Failed to create employee. "+err.Error(), http.StatusInternalServerError)
		return
	}

	base.NewResponse().SendSuccess(w, e, http.StatusCreated)
}
