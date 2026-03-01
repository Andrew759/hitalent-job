package request

type CreateDepartmentRequest struct {
	Name     string `json:"name"`
	ParentId *int   `json:"parent_id"`
}
