package middleware

type contextKey string

const (
	CreateDepartmentKey        contextKey = "validateCreatedDepartment"
	ChangeDepartmentKey        contextKey = "validateChangedDepartment"
	DeleteDepartmentRequestKey contextKey = "validateDeletedDepartmentRequest"
	CreateEmployeeKey          contextKey = "validateCreatedEmployee"
)
