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
	twtm "hitalent/pkg/gorm_tweaks/time"
	"net/http"
	"strings"
	"time"
)

type EmployeeValidator struct {
	service.DBDecorator
}

func (ev EmployeeValidator) Validate(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var cer request.CreateEmployeeRequest

		if err := json.NewDecoder(r.Body).Decode(&cer); err != nil {
			base.NewResponse().SendError(w, "invalid JSON: "+err.Error(), http.StatusBadRequest)
			return
		}
		if errorList := ev.validateRequestRules(cer); len(errorList) > 0 {
			errResponse := base.NewResponse()
			errResponse.AddErrorsToErrorContainer(errorList)

			errResponse.Send(w, http.StatusBadRequest)
			return
		}

		rDecorator := base.NewRequest(r)
		departmentId, err := rDecorator.HTTPId()
		if err != nil {
			base.NewResponse().SendError(w, "invalid query departmentId : "+err.Error(), http.StatusBadRequest)
			return
		}

		vErr := ev.validateAndSendResponseByDBRules(departmentId)
		if vErr != nil {
			base.NewResponse().SendError(w, vErr.Message, vErr.Code)
			return
		}

		var e model.Employee
		e.DepartmentId = departmentId
		e.FullName = cer.FullName
		e.Position = cer.Position
		e.CreatedAt = twtm.TimestampWithTimeZoneMicro{Time: time.Now()}
		e.HiredAt = cer.HiredAt

		ctx := context.WithValue(r.Context(), CreateEmployeeKey, &e)
		next(w, r.WithContext(ctx))
	}
}

func (ev EmployeeValidator) validateRequestRules(cer request.CreateEmployeeRequest) []error {
	var errList []error

	if strings.TrimSpace(cer.FullName) == "" || !vService.IsHasCorrectLength(cer.FullName, 1, 200) {
		errList = append(errList, errors.New("invalid full employee name length"))
	}

	if strings.TrimSpace(cer.Position) == "" || !vService.IsHasCorrectLength(cer.Position, 1, 200) {
		errList = append(errList, errors.New("invalid position length"))
	}

	return errList
}

func (ev EmployeeValidator) validateAndSendResponseByDBRules(departmentId int) *ValidatorError {
	_, err := model.GetDepartmentById(ev.DBDecorator.GormInterface, departmentId)
	if err != nil && errors.Is(err, model.DepartmentNotFoundErr) {
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
