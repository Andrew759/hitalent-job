package request

type ChangeDepartmentRequest struct {
	Name     *string `json:"name"`
	ParentId *int    `json:"parent_id"`
}
