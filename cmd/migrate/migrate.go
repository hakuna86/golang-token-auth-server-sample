package main

import (
	"github.com/hakuna86/golang-token-auth-server-sample/repo/model"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"log"
	"path/filepath"
)

func main()  {
	db, err := gorm.Open("sqlite3", filepath.Base("./sample.sqlite"))
	if err != nil {
		panic(err)
	}
	db.LogMode(true)
	db.AutoMigrate(&model.User{}, &model.Article{})
	log.Println("=============== Success Migrate ===============")
}