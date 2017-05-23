package main

import (
	"io"
	"net/http"
	"log"
	"github.com/gorilla/mux"
	"github.com/akiraak/go-manga/db"
	. "github.com/akiraak/go-manga/model"
	"time"
	"html/template"
	"os"
)

const Title = "漫画書店 ver.ω."

type BaseParam struct {
	Title	string
	Nav		string
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	type Day struct {
		Date	time.Time
		Books	[]Book
	}
	type Param struct {
		BaseParam
		Days	[]Day
	}
	jst, _ := time.LoadLocation("Asia/Tokyo")
	now := time.Now().In(jst)
	days := 5
	param := Param{BaseParam{Title, "index"}, make([]Day, days)}
	for i := 0; i < days; i++ {
		date := now.AddDate(0, 0, -i)
		param.Days[i].Date = date
		datePublish := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, jst)
		db.ORM.Where("date_publish = ?", datePublish).Find(&param.Days[i].Books)
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
	param := Param{BaseParam{Title, "log"}, []UpdateLog{}}
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
