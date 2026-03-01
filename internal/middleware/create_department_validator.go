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

type CreateDepartmentValidator struct {
	service.DBDecorator
}

func (dv CreateDepartmentValidator) Validate(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var cdr request.CreateDepartmentRequest

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

		vErr := dv.validateAndPrepareResponseByDBRules(cdr)
		if vErr != nil {
			base.NewResponse().SendError(w, vErr.Message, vErr.Code)
			return
		}

		//Секция установки параметров после валидации
		var d model.Department
		d.ParentId = cdr.ParentId
		d.Name = strings.TrimSpace(cdr.Name)
		d.CreatedAt = twtm.TimestampWithTimeZoneMicro{Time: time.Now()}

		ctx := context.WithValue(r.Context(), CreateDepartmentKey, &d)
		next(w, r.WithContext(ctx))
	}
}

// Множественная валидация параметров, в случае, если при валидации не требуется дергать БД
func (dv CreateDepartmentValidator) validateRequestRules(cdr request.CreateDepartmentRequest) []error {
	var errList []error

	if strings.TrimSpace(cdr.Name) == "" {
		errList = append(errList, errors.New("empty department name"))
	}
	if !vService.IsHasCorrectLength(cdr.Name, 1, 200) {
		errList = append(errList, errors.New("invalid department name length"))
	}

	return errList
}

// Валидация правил, которые требуют проверок в БД
func (dv CreateDepartmentValidator) validateAndPrepareResponseByDBRules(cdr request.CreateDepartmentRequest) *ValidatorError {
	_, err := model.HasSameDepartmentByParentIdAndName(dv.DBDecorator.GormInterface, cdr.Name, cdr.ParentId)
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
