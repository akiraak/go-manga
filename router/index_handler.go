package router

import (
	"fmt"
	"github.com/akiraak/go-manga/db"
	"github.com/akiraak/go-manga/web"
	"github.com/akiraak/go-manga/elastic"
	. "github.com/akiraak/go-manga/model"
	"github.com/akiraak/go-manga/pagination"
	"github.com/labstack/echo"
	"math"
	"net/http"
	"strconv"
	"time"
)

type BaseParam struct {
	Nav			string
	SearchKey	string
}

func (BaseParam)NowUnix() int64 {
	return time.Now().Unix()
}

type Day struct {
	Date			time.Time
	PublisherBooks	map[int64]map[int64]*TitleBook
}

type DaysParam struct {
	BaseParam
	Days	[]Day
}

func titleGroupBooks(books []Book) map[int64]*TitleBook {
	resultBooks := map[int64]*TitleBook{}
	for _, book := range books {
		if !book.TitleID.Valid {
			continue
		}
		_, exist := resultBooks[book.TitleID.Int64]
		if exist {
			resultBooks[book.TitleID.Int64].AddBook(book)
		} else {
			tBook := &TitleBook{}
			tBook.AddBook(book)
			resultBooks[book.TitleID.Int64] = tBook
		}
	}
	return resultBooks
}

func publisherGroupBooks(titleBooks map[int64]*TitleBook) map[int64]map[int64]*TitleBook {
	resultBooks := map[int64]map[int64]*TitleBook{}
	for key, tBook := range titleBooks {
		publisherID := tBook.PublisherID()
		_, exist := resultBooks[publisherID]
		if exist {
			resultBooks[publisherID][key] = tBook
		} else {
			resultBooks[publisherID] = map[int64]*TitleBook{}
			resultBooks[publisherID][key] = tBook
		}
	}
	return resultBooks
}

func dateBooks(year int, month time.Month, day int, r18 bool) map[int64]map[int64]*TitleBook {
	books := []Book{}
	date := fmt.Sprintf("%d%02d%02d", year, month, day)
	r18Val := 0
	if r18 {
		r18Val = 1
	}
	db.ORM.
		Joins("left join publishers on publishers.id = books.publisher_id").
		Where("date_publish = ?", date).
		Where("publishers.r18 = ?", r18Val).
		Find(&books)
	tboos := titleGroupBooks(books)
	pboos := publisherGroupBooks(tboos)
	return pboos
}

func daysBooks(nav string, r18 bool) DaysParam {
	jst, _ := time.LoadLocation("Asia/Tokyo")
	now := time.Now().In(jst)
	days := 5
	param := DaysParam{BaseParam: BaseParam{nav, ""}}
	param.Days = make([]Day, days)
	for i := 0; i < days; i++ {
		date := now.AddDate(0, 0, -i)
		param.Days[i].Date = date
		param.Days[i].PublisherBooks = dateBooks(date.Year(), date.Month(), date.Day(), r18)
	}
	return param
}

func GetIndexHandler(c echo.Context) error {
	param := daysBooks("index", false)

	return web.RenderTemplate(
		c,
		http.StatusOK,
		[]string{
			"template/index.html",
			"template/days_books.html"},
		param)
}

func GetR18Handler(c echo.Context) error {
	param := daysBooks("r18", true)

	return web.RenderTemplate(
		c,
		http.StatusOK,
		[]string{
			"template/r18.html",
			"template/days_books.html"},
		param)
}

func searchBooks(keyword string, offset int, limit int) (TitleBooks, int64, int) {
	if len(keyword) > 1 {
		asins, hitTotal := elastic.SearchAsins([]string{keyword}, offset, limit)
		books := booksWithAsins(asins)
		return books, hitTotal, len(asins)
	}
	return TitleBooks{}, 0, 0
}

func GetSearchHandler(c echo.Context) error {
	keyword := c.QueryParam("key")
	page := PageQuery(c)
	type Param struct {
		BaseParam
		TitleBooks	TitleBooks
		Page		pagination.Page
	}
	bookMax := 10000
	limit := 200
	pageMax := bookMax / limit
	if page > pageMax {
		page = pageMax
	}
	offset := limit * (page - 1)
	param := Param{BaseParam: BaseParam{"search", keyword}}
	titleBooks, hitTotal, asinsCount := searchBooks(keyword, offset, limit)
	param.TitleBooks = titleBooks

	showPageMax := int(math.Ceil(float64(hitTotal) / float64(limit)))
	if showPageMax > pageMax {
		showPageMax = pageMax
	}

	param.Page = pagination.CreatePage(
		page,
		showPageMax,
		int(hitTotal),
		offset + 1,
		offset + asinsCount)

	return web.RenderTemplate(
		c,
		http.StatusOK,
		[]string{
			"template/search.html",
			"template/books.html",
			"template/pagination.html"},
		param)
}

func booksWithAsins(asins []string) TitleBooks {
	books := []Book{}
	db.ORM.Where("asin IN (?)", asins).Find(&books)
	tbooks := titleGroupBooks(books)
	sortedBooks := TitleBooks{}
	for _, tbook := range tbooks {
		sortedBooks.Add(tbook)
	}
	sortedBooks.SorteByDate()
	return sortedBooks
}

func GetDeveloperHandler(c echo.Context) error {
	type Param struct {
		BaseParam
		TitleBooks	TitleBooks
		Tags		[]string
		Page		pagination.Page
	}
	param := Param{BaseParam: BaseParam{"dev", ""}}
	page := PageQuery(c)

	bookMax := 10000
	limit := 200
	pageMax := bookMax / limit
	if page > pageMax {
		page = pageMax
	}
	offset := limit * (page - 1)

	books := []Book{}
	userId := 1
	asinsCount := 0
	db.ORM.
		Table("user_books").
		Where("user_id = ?", userId).
		Count(&asinsCount)
	db.ORM.
		Joins("left join user_books on user_books.book_id = books.id").
		Where("user_books.user_id = ?", userId).
		Order("books.date_publish desc").
		Offset(offset).
		Limit(limit).
		Find(&books)

	tbooks := titleGroupBooks(books)
	sortedBooks := TitleBooks{}
	for _, tbook := range tbooks {
		sortedBooks.Add(tbook)
	}
	sortedBooks.SorteByDate()
	param.TitleBooks = sortedBooks

	showPageMax := int(math.Ceil(float64(asinsCount) / float64(limit)))
	if showPageMax > pageMax {
		showPageMax = pageMax
	}

	param.Page = pagination.CreatePage(
		page,
		showPageMax,
		asinsCount,
		offset + 1,
		offset + len(param.TitleBooks))

	return web.RenderTemplate(
		c,
		http.StatusOK,
		[]string{
			"template/user_books.html",
			"template/books.html",
			"template/pagination.html"},
		param)
}

func PageQuery(c echo.Context) int {
	page, err := strconv.Atoi(c.QueryParam("p"))
	if page == 0 || err != nil {
		page = 1
	}
	return page
}
