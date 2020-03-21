package models

import 	"github.com/jinzhu/gorm"

// Gallery is an our image container resource
type Gallery struct {
	gorm.Model
	UserID uint `gorm:"not_null;index"`
	Title string `gorm:"not_null"`
}
