package models

import (
	"time"
	"github.com/jinzhu/gorm"
)

// PostGRES part
type Blog struct {
	ID 			uint64 		`gorm:"primary_key" json:"id"`
	Title		string		`gorm:"size:255;not null" json:"title"`
	Content		string		`gorm:"size:255;not null" json:"content"`
	CreatedAt	time.Time	`gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt	time.Time	`gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

type Cust_User struct {
	gorm.Model
	Name string
}
