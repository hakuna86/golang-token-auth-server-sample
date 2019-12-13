package model

import (
	"encoding/json"
	"github.com/jinzhu/gorm"
)

type User struct {
	gorm.Model
	Email      string    `json:"email" gorm:"type:varchar(100);unique_index"`
	Name       string    `json:"name"`
	Password   string    `json:"password"`
	Role       string    `gorm:"size:255"`
	Article    []Article `json:"articles"`
	Auth       *Auth      `gorm:"foreignkey:Email"`
}

func (u User) ToString() string {
	data, err := json.MarshalIndent(&u, "", "	")
	if err != nil {
		return err.Error()
	}
	return string(data)
}

type Article struct {
	gorm.Model
	Email      string `gorm:"type:varchar(100)"`
	Title      string `json:"title"`
	Body       string `json:"body"`
}

type Auth struct {
	gorm.Model
	Email        string `gorm:"type:varchar(100)"`
	AccessToken  string
	RefreshToken string
}
