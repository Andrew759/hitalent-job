package model

import (
	"errors"
	"hitalent/pkg/gorm_tweaks/time"

	"gorm.io/gorm"
)

type Department struct {
	Id          int                             `json:"id" gorm:"primaryKey;autoIncrement"`
	Name        string                          `json:"name" gorm:"not null;size:200;uniqueIndex:idx_name_parent"`
	ParentId    *int                            `json:"parentId" gorm:"index:idx_name_parent"`
	CreatedAt   time.TimestampWithTimeZoneMicro `json:"created_at"`
	Departments []Department                    `json:"children" gorm:"foreignKey:ParentId"`
	Employees   []Employee                      `json:"employees" gorm:"foreignKey:DepartmentId"`
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

// GetDepartmentTree Функция, возвращающая вложенные результаты рекурсивно
func GetDepartmentTree(db *gorm.DB, id int, maxDepth int, includeEmployees bool) (Department, error) {
	var department Department

	if includeEmployees {
		db = db.Preload("Employees", sortedEmployees)
	}

	err := db.
		Preload("Departments", recursivePreload(1, maxDepth)).
		First(&department, id).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return Department{}, DepartmentNotFoundErr
	}
	return department, err
}

// Вспомогательная функция для рекурсии
func recursivePreload(currentDepth int, maxDepth int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		//query := db.Preload("Employees", sortedEmployees)

		//Если лимит не достигнут - рекурсия продолжается
		if currentDepth < maxDepth {
			db = db.Preload("Departments", recursivePreload(currentDepth+1, maxDepth))
		}

		return db
	}
}

// Сортировка employees по дате создания
func sortedEmployees(db *gorm.DB) *gorm.DB {
	return db.Order("created_at DESC")
}

// DeleteAllSubDepartmentsByParentId TODO: избавиться от циклов и переписать всё на нативный SQL запрос?
func DeleteAllSubDepartmentsByParentId(db *gorm.DB, parentId int, cascade bool, reassignToDepId *int) error {
	if cascade {
		subDeps := GetSubDepartmentsByParentId(db, parentId)
		for _, sub := range subDeps {
			err := DeleteAllSubDepartmentsByParentId(db, sub.Id, cascade, reassignToDepId)
			if err != nil {
				return err
			}
		}
	}

	if reassignToDepId != nil {
		err := ReassignEmployeeToDepById(*db, parentId, *reassignToDepId)
		if err != nil {
			return err
		}
	}

	if err := db.Delete(&Department{}, parentId).Error; err != nil {
		return err
	}

	return nil
}

func GetSubDepartmentsByParentId(db *gorm.DB, parentId int) []Department {
	var subDeps []Department
	db.Where("parent_id = ?", parentId).Find(&subDeps)

	return subDeps
}

func HasSameDepartmentByParentIdAndName(db *gorm.DB, dName string, parentId *int) (bool, error) {
	//TODO: если parentId = null, то сейчас код позволяет сохранить еще один департамент с тем же именем
	var count int64
	err := db.Model(&Department{}).Where("parent_id = ? AND name = ?", parentId, dName).Count(&count).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) && count == 0 {
		return false, nil
	}
	if count > 0 {
		return true, DepartmentAlreadyExists
	}

	return false, err
}

func SaveDepartment(db *gorm.DB, d *Department) {
	db.Save(d)
}
