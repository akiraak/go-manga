package router

import (
	"github.com/akiraak/go-manga/db"
	"github.com/akiraak/go-manga/web"
	. "github.com/akiraak/go-manga/model"
	"github.com/labstack/echo"
	"net/http"
)

func GetLogHandler(c echo.Context) error {
	type Param struct {
		BaseParam
		Logs []UpdateLog
	}
	param := Param{BaseParam: BaseParam{"log", ""}}
	db.ORM.Order("date desc").Find(&param.Logs)

	return web.RenderTemplate(
		c,
		http.StatusOK,
		[]string{"template/log.html"},
		param)
}
