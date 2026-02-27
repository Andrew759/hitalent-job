package middleware

import (
	"net/http"
)

type Validator interface {
	Validate(next http.HandlerFunc) http.HandlerFunc
}
type ValidatorError struct {
	Message string
	Code    int
}
