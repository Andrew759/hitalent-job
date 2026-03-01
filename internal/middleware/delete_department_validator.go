package middleware

import (
	"context"
	"errors"
	"hitalent/cmd/service"
	"hitalent/internal/base"
	"hitalent/internal/model"
	"net/http"
)

type DeleteDepartmentValidator struct {
	service.DBDecorator
}

func (ddv DeleteDepartmentValidator) Validate(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

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

		ctx := context.WithValue(r.Context(), DeleteeDepartmentKey, &d)
		next(w, r.WithContext(ctx))
	}
}
