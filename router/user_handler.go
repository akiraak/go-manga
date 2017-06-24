package router

import (
	"github.com/akiraak/go-manga/web"
	"github.com/labstack/echo"
	"net/http"
)

func GetUserHandler(c echo.Context) error {
	type Param struct {
		BaseParam
	}
	param := Param{BaseParam: BaseParam{"user", ""}}

	return web.RenderTemplate(
		c,
		http.StatusOK,
		[]string{"template/user.html"},
		param)
}

// TODO: タグ一覧
// TODO: タグ編集
// TODO: タグ削除
