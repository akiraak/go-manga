package ec

import (
	. "github.com/akiraak/go-manga/model"
	"github.com/labstack/echo"
	"html/template"
	"io"
)

var E = echo.New()

type Template struct {
	Templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.Templates.ExecuteTemplate(w, name, data)
}

func CreateTemplate(files []string) *Template {
	return &Template{template.Must(template.ParseFiles(files...))}
}

type CContext struct {
	echo.Context
	User	User
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
			return h(cc)
		}
	}
}
