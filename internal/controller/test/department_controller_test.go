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

	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestCreateDepartmentSuccess(t *testing.T) {
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

// по обновленной логике нельзя создать два департамента с одинаковыми в корне, если parent_id = null
func TestCreateChildDepartmentSuccess(t *testing.T) {
	tContainer := base.PrepareTestContainer(t)

	//Первый департамент
	firstDepartmentRequest := request.CreateDepartmentRequest{
		Name: "ИТ отдел",
	}
	body, _ := json.Marshal(firstDepartmentRequest)
	fResp, _ := tContainer.HTTPClient.Post(
		tContainer.HTTPServer.URL+"/departments",
		"application/json",
		bytes.NewBuffer(body),
	)

	//Второй департамент
	secondDepartmentRequest := request.CreateDepartmentRequest{
		Name: "ИТ отдел",
	}
	body, _ = json.Marshal(secondDepartmentRequest)
	sResp, _ := tContainer.HTTPClient.Post(
		tContainer.HTTPServer.URL+"/departments",
		"application/json",
		bytes.NewBuffer(body),
	)

	assert.Equal(t, http.StatusCreated, fResp.StatusCode)
	assert.Equal(t, http.StatusConflict, sResp.StatusCode)
}

// Создание дерева отделов с вложенностью 2 подотдела
func TestCreateChildDepartmentsSuccess(t *testing.T) {
	tContainer := base.PrepareTestContainer(t)

	//Первый департамент
	firstDepartmentRequest := request.CreateDepartmentRequest{
		Name: "ИТ отдел",
	}
	body, _ := json.Marshal(firstDepartmentRequest)
	fResp, _ := tContainer.HTTPClient.Post(
		tContainer.HTTPServer.URL+"/departments",
		"application/json",
		bytes.NewBuffer(body),
	)
	var fDecodedResponse appBase.Response
	json.NewDecoder(fResp.Body).Decode(&fDecodedResponse)
	var rootDepartment model.Department
	json.NewDecoder(fDecodedResponse.PayloadContainer).Decode(&rootDepartment)

	//Первый дочерний департамент
	firstChildDepartmentRequest := request.CreateDepartmentRequest{
		Name:     "Подотдел 1",
		ParentId: &rootDepartment.Id,
	}
	body, _ = json.Marshal(firstChildDepartmentRequest)
	fcResp, _ := tContainer.HTTPClient.Post(
		tContainer.HTTPServer.URL+"/departments",
		"application/json",
		bytes.NewBuffer(body),
	)

	//Второй дочерний департамент
	secondChildDepartmentRequest := request.CreateDepartmentRequest{
		Name:     "Подотдел 2",
		ParentId: &rootDepartment.Id,
	}
	body, _ = json.Marshal(secondChildDepartmentRequest)
	scResp, _ := tContainer.HTTPClient.Post(
		tContainer.HTTPServer.URL+"/departments",
		"application/json",
		bytes.NewBuffer(body),
	)
	var scDecodedResponse appBase.Response
	json.NewDecoder(scResp.Body).Decode(&scDecodedResponse)
	var scDepartment model.Department
	json.NewDecoder(scDecodedResponse.PayloadContainer).Decode(&scDepartment)

	//имя департамента совпадает, но это подотдел, поэтому ошибки нет
	firstSecondChildDepartmentRequest := request.CreateDepartmentRequest{
		Name:     "Подотдел 2",
		ParentId: &scDepartment.Id,
	}

	body, _ = json.Marshal(firstSecondChildDepartmentRequest)
	fscResp, _ := tContainer.HTTPClient.Post(
		tContainer.HTTPServer.URL+"/departments",
		"application/json",
		bytes.NewBuffer(body),
	)

	secondSecondChildDepartmentRequest := request.CreateDepartmentRequest{
		Name:     "Подотдел с уникальным именем",
		ParentId: &scDepartment.Id,
	}
	body, _ = json.Marshal(secondSecondChildDepartmentRequest)
	sscResp, _ := tContainer.HTTPClient.Post(
		tContainer.HTTPServer.URL+"/departments",
		"application/json",
		bytes.NewBuffer(body),
	)

	assert.Equal(t, http.StatusCreated, fcResp.StatusCode)
	assert.Equal(t, http.StatusCreated, scResp.StatusCode)
	assert.Equal(t, http.StatusCreated, fscResp.StatusCode)
	assert.Equal(t, http.StatusCreated, sscResp.StatusCode)
}

func TestCreateChildDepartmentWithSameNameFail(t *testing.T) {
	tContainer := base.PrepareTestContainer(t)

	//Первый департамент
	firstDepartmentRequest := request.CreateDepartmentRequest{
		Name: "ИТ отдел",
	}
	body, _ := json.Marshal(firstDepartmentRequest)
	fResp, _ := tContainer.HTTPClient.Post(
		tContainer.HTTPServer.URL+"/departments",
		"application/json",
		bytes.NewBuffer(body),
	)
	var fDecodedResponse appBase.Response
	json.NewDecoder(fResp.Body).Decode(&fDecodedResponse)
	var rootDepartment model.Department
	json.NewDecoder(fDecodedResponse.PayloadContainer).Decode(&rootDepartment)

	//Первый дочерний департамент
	firstChildDepartmentRequest := request.CreateDepartmentRequest{
		Name:     "Подотдел 1",
		ParentId: &rootDepartment.Id,
	}
	body, _ = json.Marshal(firstChildDepartmentRequest)
	fcResp, _ := tContainer.HTTPClient.Post(
		tContainer.HTTPServer.URL+"/departments",
		"application/json",
		bytes.NewBuffer(body),
	)

	//Второй дочерний департамент c тем же именем
	secondChildDepartmentRequest := request.CreateDepartmentRequest{
		Name:     "Подотдел 1",
		ParentId: &rootDepartment.Id,
	}
	body, _ = json.Marshal(secondChildDepartmentRequest)
	scResp, _ := tContainer.HTTPClient.Post(
		tContainer.HTTPServer.URL+"/departments",
		"application/json",
		bytes.NewBuffer(body),
	)

	assert.Equal(t, http.StatusCreated, fcResp.StatusCode)
	assert.Equal(t, http.StatusConflict, scResp.StatusCode)
}

func TestPlacingDepartmentInsideOwnSubtreeOnFirstLevelFail(t *testing.T) {
	tContainer := base.PrepareTestContainer(t)

	//Первый департамент
	firstDepartmentRequest := request.CreateDepartmentRequest{
		Name:     "ИТ отдел",
		ParentId: new(1),
	}
	body, _ := json.Marshal(firstDepartmentRequest)
	fResp, _ := tContainer.HTTPClient.Post(
		tContainer.HTTPServer.URL+"/departments",
		"application/json",
		bytes.NewBuffer(body),
	)
	var fDecodedResponse appBase.Response
	json.NewDecoder(fResp.Body).Decode(&fDecodedResponse)
	assert.Equal(t, "parent department not found", fDecodedResponse.ErrorContainer[0].Message)
}

func TestDepartmentNameLenFail(t *testing.T) {
	tContainer := base.PrepareTestContainer(t)

	newDepartmentRequest := request.CreateDepartmentRequest{
		Name: gofakeit.LetterN(201),
	}
	body, _ := json.Marshal(newDepartmentRequest)

	resp, _ := tContainer.HTTPClient.Post(
		tContainer.HTTPServer.URL+"/departments",
		"application/json",
		bytes.NewBuffer(body),
	)
	var decodedResponse appBase.Response
	json.NewDecoder(resp.Body).Decode(&decodedResponse)

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Equal(t, "invalid department name length", decodedResponse.ErrorContainer[0].Message)
}

func TestDepartmentEmptyNameFail(t *testing.T) {
	tContainer := base.PrepareTestContainer(t)

	newDepartmentRequest := request.CreateDepartmentRequest{
		Name: "",
	}
	body, _ := json.Marshal(newDepartmentRequest)

	resp, _ := tContainer.HTTPClient.Post(
		tContainer.HTTPServer.URL+"/departments",
		"application/json",
		bytes.NewBuffer(body),
	)
	var decodedResponse appBase.Response
	json.NewDecoder(resp.Body).Decode(&decodedResponse)

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Equal(t, "empty department name", decodedResponse.ErrorContainer[0].Message)
}
