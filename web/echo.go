package web

import (
	"bytes"
	"github.com/akiraak/go-manga/db"
	. "github.com/akiraak/go-manga/model"
	"github.com/labstack/echo"
	"html/template"
)

var Echo = echo.New()

type CContext struct {
	echo.Context
	User	User
}

const PageTitle = "漫画書店 ver.ω."

func RenderTemplate(c echo.Context, code int, files []string, data interface{}) error {
	type Param struct {
		C		echo.Context
		CC		*CContext
		Title	string
		Data	interface{}
	}
	param := Param{c, c.(*CContext), PageTitle, data}
	files = append(files, "template/base.html")
	t := template.Must(template.ParseFiles(files...))
	buf := new(bytes.Buffer)
	if err := t.ExecuteTemplate(buf, "base", param); err != nil {
		return err
	}
	return c.HTMLBlob(code, buf.Bytes())
}

func UserSessionMiddleware() echo.MiddlewareFunc {
	return func(h echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cc := &CContext{c, User{}}
			/*
			userId, err := cc.Cookie("userId")
			if err == nil {
				//fmt.Println(cookie.Name)
				fmt.Println("User exist")
			} else {
				fmt.Println("No User")
			}
			*/
			user := User{}
			if !db.ORM.First(&user).RecordNotFound() {
				cc.User = user
			}
			return h(cc)
		}
	}
}
