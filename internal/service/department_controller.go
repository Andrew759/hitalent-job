package service

import (
	"errors"
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

	dc.Controller.ServeMux.HandleFunc("GET /departments/{id}", func(w http.ResponseWriter, r *http.Request) {
		dc.GetDepartment(w, base.NewRequest(r))
	})

	dc.Controller.ServeMux.HandleFunc("DELETE /departments/{id}", func(w http.ResponseWriter, r *http.Request) {
		dc.DeleteDepartment(w, base.NewRequest(r))
	})

}

func (dc *DepartmentController) CreateDepartment(w http.ResponseWriter, r *base.Request) {
	d := r.Context().Value(middleware.DepartmentKey).(*model.Department)

	err := model.CreateDepartment(dc.Controller.Dependencies.DBDecorator.GormInterface, d)
	if err != nil {
		base.NewResponse().SendError(w, "Failed to create department. "+err.Error(), http.StatusInternalServerError)
		return
	}

	base.NewResponse().SendSuccess(w, d, http.StatusCreated)
}

func (dc *DepartmentController) GetDepartment(w http.ResponseWriter, r *base.Request) {
	id, err := r.HTTPId()
	if err != nil {
		base.NewResponse().SendError(w, err.Error(), http.StatusBadRequest)
	}
	d, err := model.GetDepartmentById(dc.Controller.Dependencies.DBDecorator.GormInterface, id)
	if err != nil && errors.Is(err, model.DepartmentNotFoundErr) {
		base.NewResponse().SendError(w, err.Error(), http.StatusNotFound)
		return
	} else if err != nil {
		base.NewResponse().SendError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	base.NewResponse().SendSuccess(w, d, http.StatusOK)
}

func (dc *DepartmentController) DeleteDepartment(w http.ResponseWriter, r *base.Request) {
	id, err := r.HTTPId()
	if err != nil {
		base.NewResponse().SendError(w, err.Error(), http.StatusBadRequest)
		return
	}
	//TODO: это мок. Ид нужно достать из тела запроса
	parentId := 1
	err = model.DeleteAllSubDepartmentsByParentId(dc.Controller.Dependencies.DBDecorator.GormInterface, id, &parentId)
	if err != nil {
		base.NewResponse().SendError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
