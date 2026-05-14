package controller

import (
	"errors"
	"hitalent/internal/base"
	"hitalent/internal/middleware"
	"hitalent/internal/model"
	"net/http"
	"strconv"
)

type DepartmentController struct {
	Controller base.Controller
	CreateDV   middleware.Validator
	ChangeDV   middleware.Validator
	DeleteDV   middleware.Validator
}

func (dc *DepartmentController) HandleRequest() {
	dc.Controller.ServeMux.HandleFunc("POST /departments",
		dc.CreateDV.Validate(func(w http.ResponseWriter, r *http.Request) {
			dc.CreateDepartment(w, base.NewRequest(r))
		}))

	dc.Controller.ServeMux.HandleFunc("GET /departments/{id}", func(w http.ResponseWriter, r *http.Request) {
		dc.GetDepartment(w, base.NewRequest(r))
	})

	dc.Controller.ServeMux.HandleFunc("DELETE /departments/{id}",
		dc.DeleteDV.Validate(func(w http.ResponseWriter, r *http.Request) {
			dc.DeleteDepartment(w, base.NewRequest(r))
		}))

	dc.Controller.ServeMux.HandleFunc("PATCH /departments/{id}",
		dc.ChangeDV.Validate(func(w http.ResponseWriter, r *http.Request) {
			dc.ChangeDepartmentParent(w, base.NewRequest(r))
		}))
}

func (dc *DepartmentController) CreateDepartment(w http.ResponseWriter, r *base.Request) {
	d := r.Context().Value(middleware.CreateDepartmentKey).(*model.Department)

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

	query := r.Request.URL.Query()
	depthS := query.Get("depth")
	var depth int
	if depthS == "" {
		depth = 1
	}
	depth, err = strconv.Atoi(depthS)
	if depth > 5 {
		base.NewResponse().SendError(w, errors.New("invalid depth").Error(), http.StatusBadRequest)
		return
	}
	includeEmployees := true
	includeEmployeesS := query.Get("include_employees")
	if includeEmployeesS != "" {
		includeEmployees, err = strconv.ParseBool(includeEmployeesS)
		if err != nil {
			base.NewResponse().SendError(w, errors.New("invalid include_employees param").Error(), http.StatusBadRequest)
			return
		}
	}

	d, err := model.GetDepartmentTree(dc.Controller.Dependencies.DBDecorator.GormInterface, id, depth, includeEmployees)
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

	query := r.Request.URL.Query()
	cascade := false
	reassign := false

	mode := query.Get("mode")
	if mode == "cascade" {
		cascade = true
	}
	if mode == "reassign" {
		reassign = true
	}

	reassignToDepartmentIdS := query.Get("reassign_to_department_id")
	if reassign && reassignToDepartmentIdS == "" {
		base.NewResponse().SendError(w, errors.New("reassign_to_department_id required if mode reassign").Error(), http.StatusBadRequest)
		return
	}
	reassignToDepartmentIdInt, err := strconv.Atoi(reassignToDepartmentIdS)
	if reassign && err != nil {
		base.NewResponse().SendError(w, errors.New("invalid reassign_to_department_id param").Error(), http.StatusBadRequest)
		return
	}

	err = model.DeleteAllSubDepartmentsByParentId(dc.Controller.Dependencies.DBDecorator.GormInterface, id, cascade, &reassignToDepartmentIdInt)
	if err != nil {
		base.NewResponse().SendError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (dc *DepartmentController) ChangeDepartmentParent(w http.ResponseWriter, r *base.Request) {
	d := r.Context().Value(middleware.ChangeDepartmentKey).(*model.Department)

	model.SaveDepartment(dc.Controller.Dependencies.DBDecorator.GormInterface, d)

	base.NewResponse().SendSuccess(w, d, http.StatusOK)
}
