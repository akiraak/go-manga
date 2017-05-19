package main

import (
	"net/http"
	"log"
	"github.com/gorilla/mux"
	"github.com/akiraak/go-manga/db"
	. "github.com/akiraak/go-manga/model"
	"time"
	"html/template"
	"path/filepath"
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
	datePublish := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, jst)
	log.Println(datePublish)
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

func main() {
	db.ORM = db.InitDB()
	defer db.ORM.Close()

	r := mux.NewRouter()
	r.HandleFunc("/", index)
	log.Fatal(http.ListenAndServe(":8000", r))

}
