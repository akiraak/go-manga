package db

import (
    "fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var ORM *gorm.DB

func InitDB() *gorm.DB {
	const user = "root"
	const pass = ""
	path := fmt.Sprintf("%s:%s@/manganow?charset=utf8&parseTime=True&loc=Local", user, pass)
	db, err := gorm.Open("mysql", path)
	if err != nil {
		panic("failed to connect database")
	}

	db.DB().SetMaxIdleConns(10)
	db.DB().SetMaxOpenConns(100)

	return db
}
