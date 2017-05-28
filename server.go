package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/akiraak/go-manga/db"
	. "github.com/akiraak/go-manga/model"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

const PageTitle = "漫画書店 ver.ω."

type BaseParam struct {
	PageTitle	string
	Nav			string
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

func indexHandler(w http.ResponseWriter, r *http.Request) {
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

func logHandler(w http.ResponseWriter, r *http.Request) {
	type Param struct {
		BaseParam
		Logs []UpdateLog
	}
	param := Param{BaseParam{PageTitle, "log"}, []UpdateLog{}}
	db.ORM.Order("date desc").Find(&param.Logs)

	tpl := template.Must(template.ParseFiles("template/base.html", "template/log.html"))
	if err := tpl.ExecuteTemplate(w, "base", param); err != nil {
		log.Println(err)
	}
}

func initLog() *os.File {
	filePath := os.Getenv("MANGANOW_LOG_FILE")
	f, err := os.OpenFile(filePath, os.O_RDWR | os.O_CREATE | os.O_APPEND, 0666)
	if err != nil {
		log.Panicf("Error opening:%v", err)
	}
	log.SetOutput(io.MultiWriter(f, os.Stdout))
	return f
}

func main() {
	logFile := initLog()
	defer logFile.Close()

	db.ORM = db.InitDB()
	defer db.ORM.Close()
	//db.ORM.LogMode(true)

	r := mux.NewRouter()
	r.HandleFunc("/", indexHandler)
	r.HandleFunc("/log", logHandler)
	log.Fatal(http.ListenAndServe(":8000", r))
}
