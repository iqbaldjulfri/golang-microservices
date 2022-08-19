package models

type Role struct {
	Model
	Code string `gorm:"unique,uniqueIndex,not null" json:"code"`
}
