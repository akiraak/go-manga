package router

import (
	"github.com/akiraak/go-manga/db"
	"github.com/akiraak/go-manga/echo"
	. "github.com/akiraak/go-manga/model"
	"github.com/akiraak/go-manga/pagination"
	"github.com/labstack/echo"
	"log"
	"math"
	"net/http"
	"strconv"
)

func GetAdminPublisherR18Handler(c echo.Context) error {
	publisherId, _ := strconv.ParseInt(c.Param("id"), 10, 64)

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
	return c.String(http.StatusOK, "ok")
}

func GetAdminPublisherHandler(c echo.Context) error {
	page := PageQuery(c)

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
	pageMax := int(math.Ceil(float64(total) / float64(max)))
	param.Page = pagination.CreatePage(page, pageMax)
	db.ORM.Order("id desc").Offset(max * (page - 1)).Limit(100).Find(&param.Publishers)

	ec.E.Renderer = ec.CreateTemplate([]string{
		"template/base.html",
		"template/admin_publisher.html"})
	return c.Render(http.StatusOK, "base", param)
}

func GetAdminAddUserHandler(c echo.Context) error {
	user := User{}
	if db.ORM.First(&user).RecordNotFound() {
		user.Name = "開発者"
		db.ORM.Create(&user)
		return c.String(http.StatusOK, "create")
	}
	return c.String(http.StatusOK, "exist")
}
