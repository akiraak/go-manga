package db

import (
    "os"
    "fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var ORM *gorm.DB

func InitDB() *gorm.DB {
	user := os.Getenv("MANGANOW_MYSQL_USER")
	pass := os.Getenv("MANGANOW_MYSQL_PASS")
	path := fmt.Sprintf("%s:%s@/manganow?charset=utf8&parseTime=True&loc=Local", user, pass)
	db, err := gorm.Open("mysql", path)
	if err != nil {
		panic("failed to connect database")
	}

	db.DB().SetMaxIdleConns(10)
	db.DB().SetMaxOpenConns(100)

	return db
}
