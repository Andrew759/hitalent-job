package middleware

import (
	"context"
	"encoding/json"
	"errors"
	"hitalent/cmd/service"
	"hitalent/internal/base"
	vService "hitalent/internal/middleware/service"
	"hitalent/internal/model"
	"hitalent/internal/request"
	"net/http"
	"strings"
)

type ChangeDepartmentValidator struct {
	service.DBDecorator
}

func (dv ChangeDepartmentValidator) Validate(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var cdr request.ChangeDepartmentRequest

		if err := json.NewDecoder(r.Body).Decode(&cdr); err != nil {
			base.NewResponse().SendError(w, "Invalid JSON: "+err.Error(), http.StatusBadRequest)
			return
		}
		if errorList := dv.validateRequestRules(cdr); len(errorList) > 0 {
			errResponse := base.NewResponse()
			errResponse.AddErrorsToErrorContainer(errorList)

			errResponse.Send(w, http.StatusBadRequest)
			return
		}

		vErr := dv.validateAndSendResponseByDBRules(cdr)
		if vErr != nil {
			base.NewResponse().SendError(w, vErr.Message, vErr.Code)
		}

		bRequest := base.NewRequest(r)
		id, err := bRequest.HTTPId()
		if err != nil {
			base.NewResponse().SendError(w, err.Error(), http.StatusBadRequest)
		}
		d, err := model.GetDepartmentById(dv.DBDecorator.GormInterface, id)
		if err != nil && errors.Is(err, model.DepartmentNotFoundErr) {
			base.NewResponse().SendError(w, err.Error(), http.StatusNotFound)
		}

		d.ParentId = cdr.ParentId
		if cdr.Name != nil {
			d.Name = strings.TrimSpace(*cdr.Name)
		}

		ctx := context.WithValue(r.Context(), ChangeDepartmentKey, &d)
		next(w, r.WithContext(ctx))
	}
}

// Множественная валидация параметров, в случае, если при ваSaveDepartmentлидации не требуется дергать БД
func (dv ChangeDepartmentValidator) validateRequestRules(cdr request.ChangeDepartmentRequest) []error {
	var errList []error

	if cdr.Name != nil && strings.TrimSpace(*cdr.Name) == "" {
		errList = append(errList, errors.New("empty department name"))
	}
	if cdr.Name != nil && !vService.IsHasCorrectLength(*cdr.Name, 1, 200) {
		errList = append(errList, errors.New("invalid department name length"))
	}

	return errList
}

// Валидация правил, которые требуют проверок в БД
func (dv ChangeDepartmentValidator) validateAndSendResponseByDBRules(cdr request.ChangeDepartmentRequest) *ValidatorError {
	_, err := model.HasSameDepartmentByParentIdAndName(dv.DBDecorator.GormInterface, *cdr.ParentId, *cdr.Name)
	if err != nil && errors.Is(err, model.DepartmentAlreadyExists) {
		return &ValidatorError{
			Message: err.Error(),
			Code:    http.StatusConflict,
		}
	} else if err != nil {
		return &ValidatorError{
			Message: err.Error(),
			Code:    http.StatusInternalServerError,
		}
	}

	return nil
}
