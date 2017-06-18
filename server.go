package main

import (
	"fmt"
	"log"
	"net/http"
	"github.com/akiraak/go-manga/db"
	"github.com/akiraak/go-manga/echo"
	"github.com/akiraak/go-manga/router"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"html/template"
	"io"
	"os"
	"time"
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

func getWrite(c echo.Context) error {
	cookie := new(http.Cookie)
	cookie.Name = "username"
	cookie.Value = "jon2"
	cookie.Expires = time.Now().Add(24 * time.Hour)
	c.SetCookie(cookie)

	return c.String(http.StatusOK, "Hello, World!")
}

func getRead(c echo.Context) error {
	cookie, err := c.Cookie("username")
	if err != nil {
		return err
	}
	fmt.Println(cookie.Name)
	fmt.Println(cookie.Value)

	return c.String(http.StatusOK, "Hello, World!")
}

func main() {
	logFile := initLog()
	defer logFile.Close()

	db.ORM = db.InitDB()
	defer db.ORM.Close()
	//db.ORM.LogMode(true)

	//ec.E.Debug = true

	ec.E.Use(middleware.Logger())
	ec.E.Use(middleware.Recover())
	ec.E.Use(middleware.CORS())
	ec.E.Use(middleware.Gzip())

	ec.E.Static("/static", "static")
	ec.E.GET("/", router.GetIndexHandler)
	ec.E.GET("/r18", router.GetR18Handler)
	ec.E.GET("/publisher/:id", router.GetPublisherHandler)
	ec.E.GET("/search", router.GetSearchHandler)
	ec.E.GET("/log", router.GetLogHandler)

	adminPath := os.Getenv("MANGANOW_ADMIN_PATH")
	if adminPath != "" {
		log.Println(adminPath)
		ec.E.GET(adminPath + "/publisher", router.GetAdminPublisherHandler)
		ec.E.GET(adminPath + "/publisher/:id/r18", router.GetAdminPublisherR18Handler)
	}

	ec.E.Logger.Fatal(ec.E.Start(":8000"))
}
