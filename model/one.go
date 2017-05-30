package model

import (
	"github.com/akiraak/go-manga/db"
	"log"
	"time"
)

type One struct {
	LastUpdateBookPage	int
	CreatedAt	time.Time
	UpdatedAt	time.Time
}

func (One) TableName() string {
  return "one"
}

func GetOneLastUpdateBookPage() int {
	var one One
	if db.ORM.First(&one).RecordNotFound() {
		log.Panic("One table does not exist.")
	}
	return one.LastUpdateBookPage
}

func SetOneLastUpdateBookPage(page int) {
	db.ORM.Table("one").UpdateColumn("last_update_book_page", page)
}
