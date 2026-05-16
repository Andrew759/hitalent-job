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
		if cdr.Name != nil {
			cdr.Name = new(strings.TrimSpace(*cdr.Name))
		}

		if errorList := dv.validateRequestRules(cdr); len(errorList) > 0 {
			errResponse := base.NewResponse()
			errResponse.AddErrorsToErrorContainer(errorList)

			errResponse.Send(w, http.StatusBadRequest)
			return
		}

		bRequest := base.NewRequest(r)
		id, err := bRequest.HTTPId()
		if err != nil {
			base.NewResponse().SendError(w, err.Error(), http.StatusBadRequest)
			return
		}

		d, vErr := dv.validateAndPrepareResponseByDBRules(cdr, id)
		if vErr != nil {
			base.NewResponse().SendError(w, vErr.Message, vErr.Code)
			return
		}

		d.ParentId = cdr.ParentId
		if cdr.Name != nil {
			d.Name = *cdr.Name
		}

		ctx := context.WithValue(r.Context(), ChangeDepartmentKey, d)
		next(w, r.WithContext(ctx))
	}
}

// Множественная валидация параметров, в случае, если при валидации не требуется дергать БД
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
func (dv ChangeDepartmentValidator) validateAndPrepareResponseByDBRules(cdr request.ChangeDepartmentRequest, depId int) (*model.Department, *ValidatorError) {
	_, err := model.HasSameDepartmentByParentIdAndName(dv.GormInterface, *cdr.Name, cdr.ParentId)
	if err != nil && errors.Is(err, model.DepartmentAlreadyExists) {
		return nil, &ValidatorError{
			Message: err.Error(),
			Code:    http.StatusConflict,
		}
	} else if err != nil {
		return nil, &ValidatorError{
			Message: err.Error(),
			Code:    http.StatusInternalServerError,
		}
	}

	d, err := model.GetDepartmentById(dv.GormInterface, depId)
	if err != nil {
		if errors.Is(err, model.DepartmentNotFoundErr) {
			return nil, &ValidatorError{
				Message: err.Error(),
				Code:    http.StatusNotFound,
			}
		}

		return nil, &ValidatorError{
			Message: err.Error(),
			Code:    http.StatusInternalServerError,
		}
	}

	vErr := dv.validateSubDepsTreeCycle(d.Id, *cdr.ParentId)
	if vErr != nil {
		return nil, vErr
	}

	return &d, nil
}

func (dv ChangeDepartmentValidator) validateSubDepsTreeCycle(currentDepId, newParentId int) *ValidatorError {
	if currentDepId == newParentId {
		return &ValidatorError{
			Message: "cannot create a department as its own parent",
			Code:    http.StatusConflict,
		}
	}

	subDeps := model.GetSubDepartmentsByParentId(dv.DBDecorator.GormInterface, currentDepId)
	for _, sub := range subDeps {
		if sub.Id == newParentId {
			return &ValidatorError{
				Message: "circular reference: cannot move department into its own subtree",
				Code:    http.StatusConflict,
			}
		}

		if vErr := dv.validateSubDepsTreeCycle(sub.Id, newParentId); vErr != nil {
			return vErr
		}
	}

	return nil
}
