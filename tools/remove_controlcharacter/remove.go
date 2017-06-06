package main

import (
	"fmt"
	"github.com/akiraak/go-manga/db"
	. "github.com/akiraak/go-manga/model"
)

func main() {
	db.ORM = db.InitDB()
	defer db.ORM.Close()
	//db.ORM.LogMode(true)

	books := []Book{}
	db.ORM.Find(&books)
	for _, book := range books {
		cleanName := CleanName(book.Name)
		if cleanName != book.Name {
			fmt.Println("----", book.ID)
			fmt.Println(cleanName)
			fmt.Println(book.Name)
			db.ORM.Table("books").Where("id = ?", book.ID).UpdateColumn("name", cleanName)
		}
	}

	publishers := []Publisher{}
	db.ORM.Find(&publishers)
	for _, publisher := range publishers {
		cleanName := CleanName(publisher.Name)
		if cleanName != publisher.Name {
			fmt.Println("----", publisher.ID)
			fmt.Println(cleanName)
			fmt.Println(publisher.Name)
			db.ORM.Table("publishers").Where("id = ?", publisher.ID).UpdateColumn("name", cleanName)
		}
	}

	authors := []Author{}
	db.ORM.Find(&authors)
	for _, author := range authors {
		cleanName := CleanName(author.Name)
		if cleanName != author.Name {
			fmt.Println("----", author.ID)
			fmt.Println(cleanName)
			fmt.Println(author.Name)
			db.ORM.Table("authors").Where("id = ?", author.ID).UpdateColumn("name", cleanName)
		}
	}

	fmt.Println("Finish")
}
