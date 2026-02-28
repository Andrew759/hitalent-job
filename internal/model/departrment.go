package model

import (
	"errors"
	"hitalent/pkg/gorm_tweaks/time"

	"gorm.io/gorm"
)

type Department struct {
	Id          int                             `json:"id" gorm:"type:int;unique;primaryKey;autoIncrement"`
	Name        string                          `json:"name" gorm:"type:not null;size:200;uniqueIndex:idx_name_parent"`
	ParentId    *int                            `json:"parentId" gorm:"type:int;uniqueIndex:idx_name_parent"`
	CreatedAt   time.TimestampWithTimeZoneMicro `json:"created_at"`
	Departments []Department                    `json:"departments" gorm:"foreignKey:ParentId;references:Id;constraint:OnDelete:CASCADE"`
	Employees   []Employee                      `json:"employess" gorm:"foreignKey:DepartmentId;constraint:OnDelete:CASCADE"`
}

var DepartmentNotFoundErr = errors.New("department not found")

var DepartmentAlreadyExists = errors.New("department already exist")

func CreateDepartment(db *gorm.DB, d *Department) error {
	return db.Create(d).Error
}

func GetDepartmentById(db *gorm.DB, id int) (Department, error) {
	var department Department
	result := db.First(&department, id)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return Department{}, DepartmentNotFoundErr
	}

	return department, result.Error
}

// DeleteAllSubDepartmentsByParentId TODO: избавиться от циклов и переписать всё на нативный SQL запрос?
func DeleteAllSubDepartmentsByParentId(db *gorm.DB, parentId int, reassignToDepId *int) error {
	if reassignToDepId != nil {
		err := ReassignEmployeeToDepById(*db, parentId, *reassignToDepId)
		if err != nil {
			return err
		}
	}
	subDeps := GetSubDepartmentsByParentId(*db, parentId)
	var subDepIds []int
	for _, sub := range subDeps {
		subDepIds = append(subDepIds, sub.Id)
		if err := db.Delete(&sub).Error; err != nil {
			return err
		}
	}

	for _, id := range subDepIds {
		err := DeleteAllSubDepartmentsByParentId(db, id, reassignToDepId)
		if err != nil {
			return err
		}
	}

	return nil
}

func GetSubDepartmentsByParentId(db gorm.DB, parentId int) []Department {
	var subDeps []Department
	db.Where("parent_id = ?", parentId).Find(&subDeps)

	return subDeps
}

func HasSameDepartmentByParentIdAndName(db gorm.DB, d Department) (bool, error) {
	var count int64
	err := db.Model(&Department{}).Where("parent_id = ? AND name = ?", d.ParentId, d.Name).Count(&count).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) && count == 0 {
		return false, nil
	}
	if count > 0 {
		return true, DepartmentAlreadyExists
	}

	return false, err
}
