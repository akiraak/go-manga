package main

import (
	"fmt"
	"github.com/akiraak/go-manga/db"
	"github.com/akiraak/go-manga/elastic"
	. "github.com/akiraak/go-manga/model"
)

func main() {
	db.ORM = db.InitDB()
	defer db.ORM.Close()
	//db.ORM.LogMode(true)

	records := []elastic.AsinRecord{}
	books := []Book{}
	db.ORM.Find(&books)
	for i, book := range books {
		publisher := Publisher{}
		db.ORM.Where("id = ?", book.PublisherID).First(&publisher)
		author := Author{}
		db.ORM.Where("id = ?", book.AuthorID).First(&author)

		title := CleanName(book.Name)
		publisherName := CleanName(publisher.Name)
		authorName := CleanName(author.Name)
		records = append(records, elastic.AsinRecord{elastic.AsinParam{title, publisherName, authorName, ""}, book.Asin})
		if i % 1000 == 0 {
			fmt.Printf("%d%% : %d / %d\n", (i * 100  / len(books)), i, len(books))
		}
	}

	fmt.Println("Added index:", len(records))
	updatedIndex := elastic.BulkAsinIndex(records)
	fmt.Println("Updated index:", updatedIndex)

	fmt.Println("Success")
}
