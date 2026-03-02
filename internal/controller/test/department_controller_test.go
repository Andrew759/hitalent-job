package test

import (
	"bytes"
	"encoding/json"
	appBase "hitalent/internal/base"
	"hitalent/internal/controller/test/base"
	"hitalent/internal/model"
	"hitalent/internal/request"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCreateDepartment(t *testing.T) {
	tContainer := base.PrepareTestContainer(t)

	newDepartmentRequest := request.CreateDepartmentRequest{
		Name: "ИТ отдел",
	}
	body, _ := json.Marshal(newDepartmentRequest)

	resp, _ := tContainer.HTTPClient.Post(
		tContainer.HTTPServer.URL+"/departments",
		"application/json",
		bytes.NewBuffer(body),
	)
	var decodedResponse appBase.Response
	json.NewDecoder(resp.Body).Decode(&decodedResponse)

	var createdDepartment model.Department
	json.NewDecoder(decodedResponse.PayloadContainer).Decode(&createdDepartment)

	assert.Equal(t, resp.StatusCode, http.StatusCreated)
	assert.Equal(t, "ИТ отдел", createdDepartment.Name)
	assert.Empty(t, createdDepartment.ParentId)
	assert.Empty(t, createdDepartment.Employees)
	assert.Empty(t, createdDepartment.Departments)

	truncatedExpectedTime := time.Now().Truncate(time.Hour).Local()
	truncatedActualTime := createdDepartment.CreatedAt.Truncate(time.Hour)
	assert.Equal(t, truncatedExpectedTime, truncatedActualTime)
}
