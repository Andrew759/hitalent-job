package model

import (
	"errors"
	"hitalent/pkg/gorm_tweaks/time"

	"gorm.io/gorm"
)

type Department struct {
	Id          int                             `json:"id" gorm:"primaryKey;autoIncrement"`
	Name        string                          `json:"name" gorm:"not null;size:200;uniqueIndex:idx_name_parent"`
	ParentId    *int                            `json:"parentId" gorm:"uniqueIndex:idx_name_parent"`
	CreatedAt   time.TimestampWithTimeZoneMicro `json:"created_at"`
	Departments []Department                    `json:"children" gorm:"foreignKey:ParentId;constraint:OnDelete:CASCADE"`
	Employees   []Employee                      `json:"employees" gorm:"foreignKey:DepartmentId;constraint:OnDelete:CASCADE"`
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

	//Прелоад для корневого уровня
	if includeEmployees {
		db = db.Preload("Employees", sortedEmployees)
	}

	err := db.
		Preload("Departments", recursivePreload(1, maxDepth, includeEmployees)).
		First(&department, id).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return Department{}, DepartmentNotFoundErr
	}
	return department, err
}

// Вспомогательная функция для рекурсии
func recursivePreload(currentDepth int, maxDepth int, includeEmployees bool) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		//Если лимит не достигнут - рекурсия продолжается
		if currentDepth < maxDepth {
			db = db.Preload("Departments", recursivePreload(currentDepth+1, maxDepth, includeEmployees))
			if includeEmployees {
				db = db.Preload("Employees", sortedEmployees)
			}
		}

		return db
	}
}

func sortedEmployees(db *gorm.DB) *gorm.DB {
	return db.Order("created_at DESC")
}

// DeleteAllSubDepartmentsByParentId TODO: избавиться от циклов и переписать всё на нативный SQL запрос?
func DeleteAllSubDepartmentsByParentId(db *gorm.DB, parentId int, cascade bool, reassignToDepId int) error {
	return db.Transaction(func(tx *gorm.DB) error {
		if reassignToDepId != 0 {
			err := ReassignEmployeeToDepById(tx, parentId, reassignToDepId)
			if err != nil {
				return err
			}
		}
		if cascade {
			subDeps := GetSubDepartmentsByParentId(db, parentId)
			for _, sub := range subDeps {
				err := DeleteAllSubDepartmentsByParentId(db, sub.Id, cascade, reassignToDepId)
				if err != nil {
					return err
				}
			}
		}

		if err := tx.Delete(&Department{}, parentId).Error; err != nil {
			return err
		}

		return nil
	})
}

func GetSubDepartmentsByParentId(db *gorm.DB, parentId int) []Department {
	var subDeps []Department
	db.Where("parent_id = ?", parentId).Find(&subDeps)
	return subDeps
}

func HasSameDepartmentByParentIdAndName(db *gorm.DB, depName string, parentId *int) (bool, error) {
	var count int64
	query := db.Model(&Department{}).Where("name = ?", depName)

	if parentId == nil {
		query = query.Where("parent_id IS NULL")
	} else {
		query = query.Where("parent_id = ?", parentId)
	}

	err := query.Count(&count).Error
	if err != nil {
		return false, err
	}
	if count > 0 {
		return true, DepartmentAlreadyExists
	}

	return false, nil
}

func SaveDepartment(db *gorm.DB, d *Department) error {
	return db.Save(d).Error
}
