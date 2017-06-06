package main

import (
	"bufio"
	"fmt"
	"github.com/akiraak/go-manga/db"
	. "github.com/akiraak/go-manga/model"
	"os"
	"strings"
)

func main() {
	db.ORM = db.InitDB()
	defer db.ORM.Close()
	//db.ORM.LogMode(true)

	buf := []string{}
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
		buf = append(buf, fmt.Sprintf(`{"index": {"_index": "asins", "_type": "asin", "_id": "%s"}}`, book.Asin))
		buf = append(buf, fmt.Sprintf(`{"title": "%s", "publisher": "%s", "author": "%s"}`, title, publisherName, authorName))
		if i % 1000 == 0 {
			fmt.Printf("%d%% : %d / %d\n", (i * 100  / len(books)), i, len(books))
		}
	}

	max := 50000
	for i := 0; (i * max) + max <= len(buf); i++ {
		start := (i * max)
		end := start + max
		if end >= len(buf) {
			end = len(buf)
		}
		f, _ := os.Create(fmt.Sprintf("requests_%d", i))
		defer f.Close()
		w := bufio.NewWriter(f)
		for _, b := range buf[start:end] {
			fmt.Fprintln(w, b)
		}
		w.Flush()
	}

	fmt.Println("Success")
}
