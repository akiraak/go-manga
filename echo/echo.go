package ec

import (
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
