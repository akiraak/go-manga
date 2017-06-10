package model

import (
	"github.com/akiraak/go-manga/db"
	"time"
)

type Publisher struct {
	ID			int64
	Name		string
	Ero			bool
	CreatedAt	time.Time
}

func (p *Publisher)LatestBooks(count int) []Book {
	books := []Book{}
	db.ORM.
		Where("publisher_id = ?", p.ID).
		Order("id desc").
		Limit(count).
		Find(&books)
	return books
}
