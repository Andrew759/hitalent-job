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

type DepartmentValidator struct {
	service.DBDecorator
}

func (dv DepartmentValidator) Validate(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var d model.Department

		if err := json.NewDecoder(r.Body).Decode(&d); err != nil {
			base.NewResponse().SendError(w, "Invalid JSON: "+err.Error(), http.StatusBadRequest)
			return
		}
		if errorList := dv.validateRequestRules(d); len(errorList) > 0 {
			errResponse := base.NewResponse()
			errResponse.AddErrorsToErrorContainer(errorList)

			errResponse.Send(w, http.StatusBadRequest)
			return
		}

		vErr := dv.validateAndSendResponseByDBRules(d)
		if vErr != nil {
			base.NewResponse().SendError(w, vErr.Message, vErr.Code)
		}

		ctx := context.WithValue(r.Context(), DepartmentKey, &d)
		next(w, r.WithContext(ctx))
	}
}

// Множественная валидация параметров, в случае, если при валидации не требуется дергать БД
func (dv DepartmentValidator) validateRequestRules(d model.Department) []error {
	var errList []error

	if strings.TrimSpace(d.Name) == "" || !vService.IsHasCorrectLength(d.Name, 1, 200) {
		errList = append(errList, errors.New("invalid department name"))
	}

	return errList
}

// Валидация правил, которые требуют праверок в БД
func (dv DepartmentValidator) validateAndSendResponseByDBRules(d model.Department) *ValidatorError {
	_, err := model.HasSameDepartmentByParentIdAndName(*dv.DBDecorator.GormInterface, d)
	if err != nil && errors.Is(err, model.DepartmentAlreadyExists) {
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
