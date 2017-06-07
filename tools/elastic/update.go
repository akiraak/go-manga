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
		records = append(records, elastic.AsinRecord{book.Asin, elastic.AsinParam{title, publisherName, authorName}})
		if i % 1000 == 0 {
			fmt.Printf("%d%% : %d / %d\n", (i * 100  / len(books)), i, len(books))
		}
	}

	max := 20000
	for i := 0; ; i++ {
		start := (i * max)
		end := start + max
		if end >= len(records) {
			end = len(records)
		}
		fmt.Println(start, end)
		updateRecords := records[start:end]
		fmt.Println("Added index:", len(updateRecords))
		updatedIndex := elastic.BulkAsinIndex(updateRecords)
		fmt.Println("Updated index:", updatedIndex)

		if end >= len(records) {
			break
		}
	}

	fmt.Println("Success")
}
