package router

import (
	"github.com/akiraak/go-manga/db"
	. "github.com/akiraak/go-manga/model"
	"html/template"
	"log"
	"net/http"
)

func LogHandler(w http.ResponseWriter, r *http.Request) {
	type Param struct {
		BaseParam
		Logs []UpdateLog
	}
	param := Param{}
	param.PageTitle = PageTitle
	param.Nav = "log"
	db.ORM.Order("date desc").Find(&param.Logs)

	tpl := template.Must(template.ParseFiles("template/base.html", "template/log.html"))
	if err := tpl.ExecuteTemplate(w, "base", param); err != nil {
		log.Println(err)
	}
}
