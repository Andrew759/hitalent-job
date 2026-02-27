package model

import "hitalent/pkg/gorm_tweaks/time"

type Department struct {
	Id          int                             `json:"id" gorm:"type:int;unique;primaryKey;autoIncrement"`
	Name        string                          `json:"name" gorm:"type:not null;size:200"`
	ParentId    int                             `json:"parentId" gorm:"type:int"`
	CreatedAt   time.TimestampWithTimeZoneMicro `json:"created_at"`
	Departments []Department                    `json:"departments" gorm:"foreignKey:ParentId;references:Id"`
	//TODO: добавить пользователей?
}
