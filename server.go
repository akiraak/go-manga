package main

import (
	"io"
	"fmt"
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
	type Param struct {
		Date	time.Time
		Books	[]Book
	}
	jst, _ := time.LoadLocation("Asia/Tokyo")
	now := time.Now().In(jst)
	param := Param{Date: now}
	datePublish := fmt.Sprintf("%d-%d-%d", now.Year(), now.Month(), now.Day())
	db.ORM.Where("date_publish = ?", datePublish).Find(&param.Books)

	f := template.FuncMap{
		"month2int": month2int,
	}
	files := []string{"template/index.html"}
	tname := filepath.Base(files[0])
	tpl, _ := template.New(tname).Funcs(f).ParseFiles(files...)
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