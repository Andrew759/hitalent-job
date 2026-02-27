package middleware

import (
	"hitalent/cmd/service"
	"net/http"
)

type DepartmentValidator struct {
	service.DBDecorator
}

func (dv DepartmentValidator) Validate(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}
