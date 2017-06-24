package main

import (
	"log"
	"github.com/akiraak/go-manga/db"
	"github.com/akiraak/go-manga/web"
	"github.com/akiraak/go-manga/router"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"html/template"
	"io"
	"os"
)

func initLog() *os.File {
	filePath := os.Getenv("MANGANOW_LOG_FILE")
	f, err := os.OpenFile(filePath, os.O_RDWR | os.O_CREATE | os.O_APPEND, 0666)
	if err != nil {
		log.Panicf("Error opening:%v", err)
	}
	log.SetOutput(io.MultiWriter(f, os.Stdout))
	return f
}

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func main() {
	logFile := initLog()
	defer logFile.Close()

	db.ORM = db.InitDB()
	defer db.ORM.Close()
	//db.ORM.LogMode(true)

	e := web.Echo
	e.Debug = true

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())
	e.Use(middleware.CSRF())

	e.Static("/static", "static")
	e.Use(web.UserSessionMiddleware())

	ug := e.Group("")
	ug.GET("/", router.GetIndexHandler)
	ug.GET("/r18", router.GetR18Handler)
	ug.GET("/dev", router.GetDeveloperHandler)
	ug.GET("/search", router.GetSearchHandler)
	ug.GET("/publisher/:id", router.GetPublisherHandler)
	ug.GET("/log", router.GetLogHandler)

	adminPath := os.Getenv("MANGANOW_ADMIN_PATH")
	if adminPath != "" {
		log.Println(adminPath)
		ag := e.Group(adminPath)
		ag.GET("/publisher", router.GetAdminPublisherHandler)
		ag.GET("/publisher/:id/r18", router.GetAdminPublisherR18Handler)
		ag.GET("/adduser", router.GetAdminAddUserHandler)

		aug := ag.Group("/user")
		aug.GET("", router.GetUserHandler)
	}

	e.Logger.Fatal(e.Start(":8000"))
}
