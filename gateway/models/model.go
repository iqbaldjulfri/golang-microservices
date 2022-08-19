package models

import (
	"time"

	"gorm.io/gorm"
)

type Model struct {
	ID        uint            `gorm:"primarykey,not null,autoIncrement" json:"id"`
	CreatedAt time.Time       `json:"createdAt"`
	UpdatedAt time.Time       `json:"updatedAt"`
	DeletedAt *gorm.DeletedAt `gorm:"index" json:"deletedAt"`
}

type FnCallData struct {
	Params  []interface{}
	Returns []interface{}
}

func (d *FnCallData) SetParams(ps ...interface{}) *FnCallData {
	d.Params = ps
	return d
}

func (d *FnCallData) SetReturns(rs ...interface{}) *FnCallData {
	d.Returns = rs
	return d
}
