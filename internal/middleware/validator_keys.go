package middleware

type contextKey string

const (
	DepartmentKey contextKey = "validateDepartment"
	EmployeeKey   contextKey = "validateEmployee"
)
