package model

import (
	"github.com/jinzhu/gorm"
)

type User struct {
	gorm.Model `json:"-"`
	Email      string    `json:"email" gorm:"type:varchar(100);unique_index"`
	Name       string    `json:"name"`
	Password   string    `json:"password"`
	Role       string    `json:"-" gorm:"size:255"`
	Article    []Article `json:"articles" gorm:"foreignkey:Email"`
	Auth       Auth      `json:"-" gorm:"foreignkey:Email"`
}

type Article struct {
	gorm.Model `json:"-"`
	Email      string `json:"-" gorm:"type:varchar(100);unique_index"`
	Title      string `json:"title"`
	Body       string `json:"body"`
}

type Auth struct {
	gorm.Model   `json:"-"`
	Email        string `json:"-" gorm:"type:varchar(100);unique_index"`
	AccessToken  string
	RefreshToken string
}
