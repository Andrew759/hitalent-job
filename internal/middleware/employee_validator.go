package middleware

import (
	"context"
	"encoding/json"
	"errors"
	"hitalent/cmd/service"
	"hitalent/internal/base"
	vService "hitalent/internal/middleware/service"
	"hitalent/internal/model"
	"net/http"
	"strings"
)

type EmployeeValidator struct {
	service.DBDecorator
}

func (ev EmployeeValidator) Validate(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var e model.Employee

		if err := json.NewDecoder(r.Body).Decode(&e); err != nil {
			base.NewResponse().SendError(w, "Invalid JSON: "+err.Error(), http.StatusBadRequest)
			return
		}
		if errorList := ev.validateRequestRules(e); len(errorList) > 0 {
			errResponse := base.NewResponse()
			errResponse.AddErrorsToErrorContainer(errorList)

			errResponse.Send(w, http.StatusBadRequest)
			return
		}

		vErr := ev.validateAndSendResponseByDBRules(e)
		if vErr != nil {
			base.NewResponse().SendError(w, vErr.Message, vErr.Code)
		}

		ctx := context.WithValue(r.Context(), EmployeeKey, &e)
		next(w, r.WithContext(ctx))
	}
}

func (ev EmployeeValidator) validateRequestRules(e model.Employee) []error {
	var errList []error

	if strings.TrimSpace(e.FullName) == "" || !vService.IsHasCorrectLength(e.FullName, 1, 200) {
		errList = append(errList, errors.New("invalid full employee name length"))
	}

	if strings.TrimSpace(e.Position) == "" || !vService.IsHasCorrectLength(e.FullName, 1, 200) {
		errList = append(errList, errors.New("invalid position length"))
	}

	return errList
}

func (ev EmployeeValidator) validateAndSendResponseByDBRules(e model.Employee) *ValidatorError {
	_, err := model.GetDepartmentById(ev.DBDecorator.GormInterface, e.DepartmentId)
	if err != nil && errors.Is(err, model.DepartmentNotFoundErr) {
		return &ValidatorError{
			Message: err.Error(),
			Code:    http.StatusConflict,
		}
	}
	if err != nil {
		return &ValidatorError{
			Message: err.Error(),
			Code:    http.StatusInternalServerError,
		}
	}

	return nil
}
