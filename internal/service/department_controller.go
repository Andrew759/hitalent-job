package service

import (
	"hitalent/internal/base"
	"hitalent/internal/middleware"
	"hitalent/internal/model"
	"net/http"
)

type DepartmentController struct {
	Controller base.Controller
	middleware.Validator
}

func (dc *DepartmentController) HandleRequest() {
	dc.Controller.ServeMux.HandleFunc("POST /departments",
		dc.Validate(func(w http.ResponseWriter, r *http.Request) {
			dc.CreateDepartment(w, base.NewRequest(r))
		}))
}

func (dc *DepartmentController) CreateDepartment(w http.ResponseWriter, r *base.Request) {
	d := r.Context().Value(middleware.DepartmentKey).(*model.Department)

	println(d)
}
