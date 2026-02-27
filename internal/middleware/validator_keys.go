package middleware

type contextKey string

const (
	DepartmentKey contextKey = "validateDepartment"
	Employee      contextKey = "validateEmployee"
)
