package request

import "encoding/json"

type ChangeDepartmentRequest struct {
	Name      *string `json:"name"`
	ParentId  *int    `json:"parent_id"`
	BodyBytes []byte
}

func (cdr ChangeDepartmentRequest) HasParentId() bool {
	var rawMap map[string]json.RawMessage
	_ = json.Unmarshal(cdr.BodyBytes, &rawMap)

	_, exists := rawMap["parent_id"]
	if exists {
		return true
	}

	return false
}
