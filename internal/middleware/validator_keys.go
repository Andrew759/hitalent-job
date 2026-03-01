package middleware

type contextKey string

const (
	CreateDepartmentKey  contextKey = "validateCreatedDepartment"
	ChangeDepartmentKey  contextKey = "validateChangedDepartment"
	DeleteeDepartmentKey contextKey = "validateDeletedDepartment"
	CreateEmployeeKey    contextKey = "validateCreatedEmployee"
)
