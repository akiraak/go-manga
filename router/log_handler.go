package router

import (
	"github.com/akiraak/go-manga/db"
	"github.com/akiraak/go-manga/echo"
	. "github.com/akiraak/go-manga/model"
	"github.com/labstack/echo"
	"net/http"
)

func GetLogHandler(c echo.Context) error {
	type Param struct {
		BaseParam
		Logs []UpdateLog
	}
	param := Param{}
	param.PageTitle = PageTitle
	param.Nav = "log"
	db.ORM.Order("date desc").Find(&param.Logs)

	ec.E.Renderer = ec.CreateTemplate([]string{
		"template/base.html",
		"template/log.html"})
	return c.Render(http.StatusOK, "base", param)
}
