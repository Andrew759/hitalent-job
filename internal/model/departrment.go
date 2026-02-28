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
func GetDepartmentTree(db *gorm.DB, id int, maxDepth int) (Department, error) {
	var department Department

	// Прелоад и сортировка для сотрудников корневого департамента
	err := db.Preload("Employees", func(db *gorm.DB) *gorm.DB {
		return db.Order("created_at DESC") // DESC для новых сверху, или ASC
	}).
		//TODO: нужно ли прелоадить департаменты Employees?
		// Preload("Employees.Department").
		//По умолчанию минимальная глубина всегда будет 1
		Preload("Departments", recursivePreload(1, maxDepth)).
		First(&department, id).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return Department{}, DepartmentNotFoundErr
	}
	return department, err
}

// Вспомогательная функци ядля рекурсии
func recursivePreload(minDepth int, maxDepth int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		// Базовая настройка загрузки сотрудников с сортировкой
		employeeOrder := func(db *gorm.DB) *gorm.DB {
			return db.Order("created_at DESC")
		}

		if minDepth >= maxDepth {
			return db.Preload("Employees", employeeOrder).
				Preload("Employees.Department")
		}

		return db.Preload("Employees", employeeOrder).
			Preload("Employees.Department").
			Preload("Departments", recursivePreload(minDepth+1, maxDepth))
	}
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
