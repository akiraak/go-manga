package main

import (
	"fmt"
	"github.com/akiraak/go-manga/db"
	. "github.com/akiraak/go-manga/model"
)

const max = 1000

func updateBooks() {
	for i := 0; ; i++ {
		books := []Book{}
		db.ORM.Offset(i * max).Limit(max).Find(&books)
		for _, book := range books {
			cleanName := CleanName(book.Name)
			if cleanName != book.Name {
				fmt.Println("----", book.ID)
				fmt.Println(cleanName)
				fmt.Println(book.Name)
				db.ORM.Table("books").Where("id = ?", book.ID).UpdateColumn("name", cleanName)
			}
		}
		if len(books) == 0 {
			break
		}
	}
}

func updatePublishers() {
	for i := 0; ; i++ {
		publishers := []Publisher{}
		db.ORM.Offset(i * max).Limit(max).Find(&publishers)
		for _, publisher := range publishers {
			cleanName := CleanName(publisher.Name)
			if cleanName != publisher.Name {
				fmt.Println("----", publisher.ID)
				fmt.Println(cleanName)
				fmt.Println(publisher.Name)
				db.ORM.Table("publishers").Where("id = ?", publisher.ID).UpdateColumn("name", cleanName)
			}
		}
		if len(publishers) == 0 {
			break
		}
	}
}

func updateAuthors() {
	for i := 0; ; i++ {
		authors := []Author{}
		db.ORM.Offset(i * max).Limit(max).Find(&authors)
		for _, author := range authors {
			cleanName := CleanName(author.Name)
			if cleanName != author.Name {
				fmt.Println("----", author.ID)
				fmt.Println(cleanName)
				fmt.Println(author.Name)
				db.ORM.Table("authors").Where("id = ?", author.ID).UpdateColumn("name", cleanName)
			}
		}
		if len(authors) == 0 {
			break
		}
	}
}

func main() {
	db.ORM = db.InitDB()
	defer db.ORM.Close()
	//db.ORM.LogMode(true)

	updateBooks()
	updatePublishers()
	updateAuthors()

	fmt.Println("Finish")
}
