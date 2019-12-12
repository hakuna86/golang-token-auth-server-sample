package repo

import (
	"github.com/hakuna86/golang-token-auth-server-sample/repo/model"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"path/filepath"
)

type Repo struct {
	DB *gorm.DB
}

func NewRepo() (*Repo, error) {
	r := &Repo{}
	var err error
	r.DB, err = gorm.Open("sqlite3", filepath.Base("../sqlite/sample.sqlite"))
	if err != nil {
		return nil, err
	}
	r.DB.LogMode(true)
	return r, nil
}
func (r *Repo) CreateUser(user *model.User) error {
	return r.DB.Create(user).Error
}

func (r *Repo) IsUser(username, password string) bool {
	u := model.User{Email: username, Password: password}
	if r.DB.Find(&u).Error != nil {
		return false
	}
	return true
}

func (r *Repo) GetUser(user *model.User) *model.User {
	if r.DB.Find(user).Error != nil {
		return nil
	}
	return user
}

func (r *Repo) UpdateUser(user *model.User) error {
	return r.DB.Save(user).Error
}

func (r *Repo) UdateUserAuthInfo(auth *model.Auth) error {
	return r.DB.Save(auth).Error
}
