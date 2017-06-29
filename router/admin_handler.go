package router

import (
	"github.com/akiraak/go-manga/db"
	"github.com/akiraak/go-manga/web"
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
	param := Param{BaseParam: BaseParam{"admin_publisher", ""}}

	limit := 100
	total := 0
	db.ORM.Table("publishers").Count(&total)
	pageMax := int(math.Ceil(float64(total) / float64(limit)))
	offset := limit * (page - 1)
	db.ORM.Order("id desc").Offset(offset).Limit(100).Find(&param.Publishers)
	param.Page = pagination.CreatePage(
		page,
		pageMax,
		total,
		offset + 1,
		offset + len(param.Publishers))

	return web.RenderTemplate(
		c,
		http.StatusOK,
		[]string{
			"template/admin_publisher.html",
			"template/pagination.html"},
		param)
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
