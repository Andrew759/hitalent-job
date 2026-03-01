package middleware

type contextKey string

const (
	CreateDepartmentKey contextKey = "validateCreatedDepartment"
	ChangeDepartmentKey contextKey = "validateChangedDepartment"
	CreateEmployeeKey   contextKey = "validateCreatedEmployee"
)
