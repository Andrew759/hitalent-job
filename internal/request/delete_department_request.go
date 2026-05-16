package request

type DeleteDepartmentRequest struct {
	DepartmentId           int
	Cascade                bool
	Reassign               bool
	ReassignToDepartmentId int
}
