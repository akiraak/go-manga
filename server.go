package main

import (
	"github.com/akiraak/go-manga/db"
	"github.com/akiraak/go-manga/router"
	"github.com/gorilla/mux"
	"io"
	"log"
	"net/http"
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

func main() {
	logFile := initLog()
	defer logFile.Close()

	db.ORM = db.InitDB()
	defer db.ORM.Close()
	//db.ORM.LogMode(true)

	r := mux.NewRouter()
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static/"))))
	r.HandleFunc("/", router.IndexHandler)
	r.HandleFunc("/r18", router.R18Handler)
	r.HandleFunc("/publisher/{id:[0-9]+}", router.PublisherHandler)
	r.HandleFunc("/log", router.LogHandler)
	r.HandleFunc("/search", router.SearchHandler)

	adminPath := os.Getenv("MANGANOW_ADMIN_PATH")
	if adminPath != "" {
		log.Println(adminPath)
		r.HandleFunc(adminPath + "/publisher", router.AdminPublisherHandler)
		r.HandleFunc(adminPath + "/publisher/{id:[0-9]+}/r18", router.AdminPublisherR18Handler)
	}

	log.Fatal(http.ListenAndServe(":8000", r))
}
