package request

import "hitalent/pkg/gorm_tweaks/time"

type CreateEmployeeRequest struct {
	FullName string                           `json:"full_name"`
	Position string                           `json:"position"`
	HiredAt  *time.TimestampWithTimeZoneMicro `json:"hired_at"`
}
