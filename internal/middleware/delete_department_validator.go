package middleware

import (
	"context"
	"errors"
	"hitalent/cmd/service"
	"hitalent/internal/base"
	"hitalent/internal/model"
	"hitalent/internal/request"
	"net/http"
	"strconv"
)

type DeleteDepartmentValidator struct {
	service.DBDecorator
}

func (ddv DeleteDepartmentValidator) Validate(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var ddr request.DeleteDepartmentRequest

		bRequest := base.NewRequest(r)
		id, err := bRequest.HTTPId()
		if err != nil {
			base.NewResponse().SendError(w, err.Error(), http.StatusBadRequest)
			return
		}
		d, err := model.GetDepartmentById(ddv.DBDecorator.GormInterface, id)
		if err != nil && errors.Is(err, model.DepartmentNotFoundErr) {
			base.NewResponse().SendError(w, err.Error(), http.StatusNotFound)
			return
		}
		ddr.DepartmentId = d.Id

		query := bRequest.URL.Query()
		mode := query.Get("mode")

		//TODO: флаг не используется, поскольку каскадное удаление устанавливается на уровне БД. При этом работает
		// и reassign
		if mode == "cascade" {
			ddr.Cascade = true
		}
		if mode == "reassign" {
			ddr.Reassign = true
		}

		reassignToDepartmentIdS := query.Get("reassign_to_department_id")
		if ddr.Reassign {
			if reassignToDepartmentIdS == "" {
				base.NewResponse().SendError(w, errors.New("reassign_to_department_id required if mode reassign").Error(), http.StatusBadRequest)
				return
			}
			reassignToDepartmentId, err := strconv.Atoi(reassignToDepartmentIdS)
			if err != nil {
				base.NewResponse().SendError(w, errors.New("invalid reassign_to_department_id param").Error(), http.StatusBadRequest)
				return
			}

			_, err = model.GetDepartmentById(ddv.DBDecorator.GormInterface, reassignToDepartmentId)
			if err != nil && errors.Is(err, model.DepartmentNotFoundErr) {
				base.NewResponse().SendError(w, err.Error(), http.StatusNotFound)
				return
			}
			ddr.ReassignToDepartmentId = reassignToDepartmentId
		}

		ctx := context.WithValue(r.Context(), DeleteDepartmentRequestKey, &ddr)
		next(w, r.WithContext(ctx))
	}
}
