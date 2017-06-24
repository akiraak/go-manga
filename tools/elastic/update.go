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

	total := 0
	db.ORM.Table("books").Count(&total)

	records := []elastic.AsinRecord{}
	max := 1000
	for i := 0; ; i++ {
		books := []Book{}
		db.ORM.Offset(i * max).Limit(max).Find(&books)
		for _, book := range books {
			publisher := Publisher{}
			db.ORM.Where("id = ?", book.PublisherID).First(&publisher)
			author := Author{}
			db.ORM.Where("id = ?", book.AuthorID).First(&author)

			title := CleanName(book.Name)
			publisherName := CleanName(publisher.Name)
			authorName := CleanName(author.Name)
			records = append(
				records,
				elastic.AsinRecord{
					elastic.AsinParam{
						title,
						publisherName,
						authorName,
						"",
						book.DatePublishTime()},
					book.Asin})
		}
		fmt.Printf("%d%% : %d / %d\n", (i * max * 100  / total), i * max, total)
		if len(books) == 0 {
			break
		}
	}

	fmt.Println("Added index:", len(records))
	updatedIndex := elastic.BulkAsinIndex(records)
	fmt.Println("Updated index:", updatedIndex)

	fmt.Println("Success")
}
