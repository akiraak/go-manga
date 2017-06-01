package router

import (
	"fmt"
	"github.com/akiraak/go-manga/db"
	. "github.com/akiraak/go-manga/model"
	"html/template"
	"log"
	"net/http"
	"sort"
	"strconv"
	"time"
)

const PageTitle = "漫画書店 ver.ω."

type BaseParam struct {
	PageTitle	string
	Nav			string
}

func (BaseParam)NowUnix() int64 {
	return time.Now().Unix()
}

func titleBooks(books []Book) map[int64]*TitleBook {
	resultBooks := map[int64]*TitleBook{}
	for _, book := range books {
		if !book.TitleID.Valid {
			continue
		}
		_, exist := resultBooks[book.TitleID.Int64]
		if exist {
			resultBooks[book.TitleID.Int64].AddBook(book)
		} else {
			tBook := &TitleBook{}
			tBook.AddBook(book)
			resultBooks[book.TitleID.Int64] = tBook
		}
	}
	return resultBooks
}

func publisherBooks(titleBooks map[int64]*TitleBook) map[int64]map[int64]*TitleBook {
	resultBooks := map[int64]map[int64]*TitleBook{}
	for key, tBook := range titleBooks {
		publisherID := tBook.PublisherID()
		_, exist := resultBooks[publisherID]
		if exist {
			resultBooks[publisherID][key] = tBook
		} else {
			resultBooks[publisherID] = map[int64]*TitleBook{}
			resultBooks[publisherID][key] = tBook
		}
	}
	return resultBooks
}

func dateBooks(year int, month time.Month, day int) map[int64]map[int64]*TitleBook {
	books := []Book{}
	date := fmt.Sprintf("%d%02d%02d", year, month, day)
	db.ORM.Where("date_publish = ?", date).Find(&books)
	tboos := titleBooks(books)
	pboos := publisherBooks(tboos)
	return pboos
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	type Day struct {
		Date			time.Time
		PublisherBooks	map[int64]map[int64]*TitleBook
	}
	type Param struct {
		BaseParam
		Days	[]Day
	}
	jst, _ := time.LoadLocation("Asia/Tokyo")
	now := time.Now().In(jst)
	days := 5
	param := Param{BaseParam{PageTitle, "index"}, make([]Day, days)}
	for i := 0; i < days; i++ {
		date := now.AddDate(0, 0, -i)
		param.Days[i].Date = date
		param.Days[i].PublisherBooks = dateBooks(date.Year(), date.Month(), date.Day())
	}

	tpl := template.Must(template.ParseFiles("template/base.html", "template/index.html"))
	if err := tpl.ExecuteTemplate(w, "base", param); err != nil {
		log.Println(err)
	}
}

func searchBooks(keyword string) []*TitleBook {
	if len(keyword) > 1 {
		books := []Book{}
		db.ORM.Where("name LIKE ?", "%" + keyword + "%").Find(&books)
		tbooks := titleBooks(books)
		sortedBooks := []*TitleBook{}
		for _, tbook := range tbooks {
			sortedBooks = append(sortedBooks, tbook)
		}
		sort.Slice(sortedBooks, func(i, j int) bool {
			int1, _ := strconv.Atoi(sortedBooks[i].DatePublish())
			int2, _ := strconv.Atoi(sortedBooks[j].DatePublish())
			return int1 > int2
		})
		return sortedBooks
	}
	return []*TitleBook{}
}

func SearchHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	keyword := query.Get("key")
	type Param struct {
		BaseParam
		Keyword		string
		TitleBooks	[]*TitleBook
	}
	param := Param{BaseParam{PageTitle, "search"}, keyword, []*TitleBook{}}
	param.TitleBooks = searchBooks(keyword)

	tpl := template.Must(template.ParseFiles("template/base.html", "template/search.html"))
	if err := tpl.ExecuteTemplate(w, "base", param); err != nil {
		log.Println(err)
	}
}
