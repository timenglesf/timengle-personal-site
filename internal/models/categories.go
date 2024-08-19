package models

import "gorm.io/gorm"

type Category struct {
	gorm.Model
	Name  string `gorm:"unique;not null"`
	Posts []Post `gorm:"many2many:post_categories"`
}
