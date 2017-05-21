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
	"path/filepath"
	"os"
)

func month2int(m time.Month) int {
	return int(m)
}

func index(w http.ResponseWriter, r *http.Request) {
	f := template.FuncMap{
		"month2int": month2int,
	}
	files := []string{"template/index.html"}
	tname := filepath.Base(files[0])
	tpl, _ := template.New(tname).Funcs(f).ParseFiles(files...)

	type Day struct {
		Date	time.Time
		Books	[]Book
	}
	type Param struct {
		Days	[]Day
	}
	jst, _ := time.LoadLocation("Asia/Tokyo")
	now := time.Now().In(jst)
	days := 5
	param := Param{make([]Day, days)}
	for i := 0; i < days; i++ {
		date := now.AddDate(0, 0, -i)
		param.Days[i].Date = date
		datePublish := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, jst)
		db.ORM.Where("date_publish = ?", datePublish).Find(&param.Days[i].Books)
	}
	if err := tpl.Execute(w, param); err != nil {
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
	r.HandleFunc("/", index)
	log.Fatal(http.ListenAndServe(":8000", r))
}
