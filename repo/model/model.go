package model

import (
	"github.com/jinzhu/gorm"
)

type User struct {
	gorm.Model `json:"-"`
	Email string  `json:"email" gorm:"type:varchar(100);unique_index"`
	Name string `json:"name"`
	Password string `json:"password" query:"password"`
	Role string  `json:"-" gorm:"size:255"`
	Article []Article `json:"articles"`
}

type Article struct {
	gorm.Model
	Title string
	Body string
}