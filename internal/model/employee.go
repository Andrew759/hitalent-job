package model

import "hitalent/pkg/gorm_tweaks/time"

type Employee struct {
	Id           int `json:"id" gorm:"type:int;unique;primaryKey;autoIncrement"`
	DepartmentId int `json:"department_id" gorm:"type:int"`
	Department   `json:"department" gorm:"foreignKey:DepartmentId;references:Id"`
	FullName     string                          `json:"full_name" gorm:"type:string;not null;size:200"`
	Position     string                          `json:"position" gorm:"type:string;not null;size:200"`
	HiredAt      time.TimestampWithTimeZoneMicro `json:"hired_at"`
	CreatedAt    time.TimestampWithTimeZoneMicro `json:"created_at"`
}
