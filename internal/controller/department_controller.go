package controller

import (
	"errors"
	"hitalent/internal/base"
	"hitalent/internal/middleware"
	"hitalent/internal/model"
	"hitalent/internal/request"
	"log/slog"
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
		slog.Error("failed to create department", "error", err, "name", d.Name, "parent_id", d.ParentId)
		base.NewResponse().SendError(w, "failed to create department. "+err.Error(), http.StatusInternalServerError)
		return
	}

	slog.Info("department created successfully", "id", d.Id, "name", d.Name)
	base.NewResponse().SendSuccess(w, d, http.StatusCreated)
}

func (dc *DepartmentController) GetDepartment(w http.ResponseWriter, r *base.Request) {
	id, err := r.HTTPId()
	if err != nil {
		slog.Warn("invalid HTTP id format", "error", err)
		base.NewResponse().SendError(w, err.Error(), http.StatusBadRequest)
		return
	}

	query := r.Request.URL.Query()
	depthS := query.Get("depth")
	var depth int
	if depthS == "" {
		depth = 1
	}
	depth, err = strconv.Atoi(depthS)
	if depth > 5 || err != nil {
		slog.Warn("invalid depth", "depth", depthS, "error", err)
		base.NewResponse().SendError(w, errors.New("invalid depth").Error(), http.StatusBadRequest)
		return
	}
	includeEmployees := true
	includeEmployeesS := query.Get("include_employees")
	if includeEmployeesS != "" {
		includeEmployees, err = strconv.ParseBool(includeEmployeesS)
		if err != nil {
			slog.Warn("invalid include_employees param", "include_employees", includeEmployeesS, "error", err)
			base.NewResponse().SendError(w, errors.New("invalid include_employees param").Error(), http.StatusBadRequest)
			return
		}
	}

	d, err := model.GetDepartmentTree(dc.Controller.Dependencies.DBDecorator.GormInterface, id, depth, includeEmployees)
	if err != nil && errors.Is(err, model.DepartmentNotFoundErr) {
		slog.Warn("department not found", "id", id)
		base.NewResponse().SendError(w, err.Error(), http.StatusNotFound)
		return
	} else if err != nil {
		slog.Error("failed to fetch department tree", "id", id, "error", err)
		base.NewResponse().SendError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	base.NewResponse().SendSuccess(w, d, http.StatusOK)
}

func (dc *DepartmentController) DeleteDepartment(w http.ResponseWriter, r *base.Request) {
	ddr := r.Context().Value(middleware.DeleteDepartmentRequestKey).(*request.DeleteDepartmentRequest)

	err := model.DeleteAllSubDepartmentsByParentId(dc.Controller.Dependencies.DBDecorator.GormInterface, ddr.DepartmentId, ddr.ReassignToDepartmentId)
	if err != nil {
		slog.Error("failed to delete department",
			"id", ddr.DepartmentId,
			"cascade", ddr.Cascade,
			"reassign_to", ddr.ReassignToDepartmentId,
			"error", err,
		)
		base.NewResponse().SendError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	slog.Info("department deleted successfully", "id", ddr.DepartmentId, "cascade", ddr.Cascade)
	w.WriteHeader(http.StatusNoContent)
}

func (dc *DepartmentController) ChangeDepartmentParent(w http.ResponseWriter, r *base.Request) {
	d := r.Context().Value(middleware.ChangeDepartmentKey).(*model.Department)

	err := model.SaveDepartment(dc.Controller.Dependencies.DBDecorator.GormInterface, d)
	if err != nil {
		slog.Error("failed to update department", "id", d.Id, "error", err)
		base.NewResponse().SendError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	slog.Info("department updated successfully", "id", d.Id, "new_parent_id", d.ParentId)
	base.NewResponse().SendSuccess(w, d, http.StatusOK)
}
