package router

import (
	"fmt"
	"github.com/akiraak/go-manga/db"
	. "github.com/akiraak/go-manga/model"
	"github.com/akiraak/go-manga/pagination"
	"github.com/gorilla/mux"
	"html/template"
	"log"
	"math"
	"net/http"
	"strconv"
)

func AdminPublisherR18Handler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	publisherId, _ := strconv.ParseInt(vars["id"], 10, 64)
	publisher := Publisher{}
	if !db.ORM.Where("id = ?", publisherId).First(&publisher).RecordNotFound() {
		if publisher.R18 {
			publisher.R18 = false
		} else {
			publisher.R18 = true
		}
		db.ORM.Save(&publisher)
		log.Println(publisher, "update R18 to", publisher.R18)
	}
	fmt.Fprintf(w, "ok")
}

func AdminPublisherHandler(w http.ResponseWriter, r *http.Request) {
	type Param struct {
		BaseParam
		Publishers	[]Publisher
		Page		pagination.Page
	}
	param := Param{}
	param.PageTitle = PageTitle
	param.Nav = "admin_publisher"

	max := 100
	total := 0
	db.ORM.Table("publishers").Count(&total)
	page := PageQuery(r)
	pageMax := int(math.Ceil(float64(total) / float64(max)))
	param.Page = pagination.CreatePage(page, pageMax)

	db.ORM.Order("id desc").Offset(max * (page - 1)).Limit(100).Find(&param.Publishers)

	tpl := template.Must(template.ParseFiles("template/base.html", "template/admin_publisher.html"))
	if err := tpl.ExecuteTemplate(w, "base", param); err != nil {
		log.Println(err)
	}
}
