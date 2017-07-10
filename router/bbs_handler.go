package router

import (
	"github.com/akiraak/go-manga/db"
	"github.com/akiraak/go-manga/web"
	. "github.com/akiraak/go-manga/model"
	"github.com/labstack/echo"
	"net/http"
)

func GetBbsHandler(c echo.Context) error {
	type BbsParam struct {
		BaseParam
		Comments	[]BbsComment
	}
	param := BbsParam{BaseParam: BaseParam{"bbs", ""}}
	db.ORM.Order("id desc").Find(&param.Comments)
	return web.RenderTemplate(
		c,
		http.StatusOK,
		[]string{
			"template/bbs.html"},
		param)
}

func PostBbsAddHandler(c echo.Context) error {
	name := c.FormValue("name")
	if len(name) == 0 {
		name = "名無しさん"
	}
	comment := c.FormValue("comment")
	if len(comment) > 0 {
		remoteArre := c.Request().RemoteAddr
		bbsComment := BbsComment{
			Name: name,
			Comment: comment,
			IpHash: remoteArre,
		}
		db.ORM.Create(&bbsComment)
	}
	return c.Redirect(http.StatusSeeOther, "/bbs")
}
