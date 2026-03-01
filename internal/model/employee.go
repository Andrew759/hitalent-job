package model

import (
	"hitalent/pkg/gorm_tweaks/time"

	"gorm.io/gorm"
)

type Employee struct {
	Id           int                              `json:"id" gorm:"primaryKey;autoIncrement"`
	DepartmentId int                              `json:"department_id" gorm:"not null"`
	Department   Department                       `json:"department" gorm:"foreignKey:DepartmentId;references:Id"`
	FullName     string                           `json:"full_name" gorm:"not null;size:200"`
	Position     string                           `json:"position" gorm:"not null;size:200"`
	HiredAt      *time.TimestampWithTimeZoneMicro `json:"hired_at"`
	CreatedAt    time.TimestampWithTimeZoneMicro  `json:"created_at"`
}

func ReassignEmployeeToDepById(db gorm.DB, depId int, newDepId int) error {
	result := db.Model(&Employee{}).
		Where("department_id = ?", depId).
		Update("department_id", newDepId)

	return result.Error
}

func CreateEmployee(db *gorm.DB, e *Employee) error {
	return db.Create(e).Error
}
