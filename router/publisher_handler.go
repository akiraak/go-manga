package router

import (
	"github.com/akiraak/go-manga/db"
	. "github.com/akiraak/go-manga/model"
	"github.com/gorilla/mux"
	"html/template"
	"log"
	"net/http"
	"sort"
	"strconv"
)

func publisherBooks(publisherId int64) ([]*TitleBook, int64, int) {
	books := []Book{}
	total := int64(0)
	db.ORM.Table("books").
		Where("publisher_id = ?", publisherId).
		Count(&total)
	db.ORM.Where("publisher_id = ?", publisherId).
		Limit(200).
		Order("date_publish DESC").
		Find(&books)
	tbooks := titleGroupBooks(books)
	sortedBooks := []*TitleBook{}
	for _, tbook := range tbooks {
		sortedBooks = append(sortedBooks, tbook)
	}
	sort.Slice(sortedBooks, func(i, j int) bool {
		int1, _ := strconv.Atoi(sortedBooks[i].DatePublish())
		int2, _ := strconv.Atoi(sortedBooks[j].DatePublish())
		return int1 > int2
	})
	return sortedBooks, total, len(books)
}

func PublisherHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	publisherId, _ := strconv.ParseInt(vars["id"], 10, 64)

	type Param struct {
		BaseParam
		Publisher	Publisher
		TitleBooks	[]*TitleBook
		HitTotal	int64
		AsinsCount	int
	}
	param := Param{}
	param.PageTitle = PageTitle
	param.Nav = "search"
	db.ORM.Where("id = ?", publisherId).First(&param.Publisher)
	param.TitleBooks, param.HitTotal, param.AsinsCount = publisherBooks(publisherId)

	tpl := template.Must(template.ParseFiles("template/base.html", "template/publisher.html"))
	if err := tpl.ExecuteTemplate(w, "base", param); err != nil {
		log.Println(err)
	}
}
