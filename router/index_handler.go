package router

import (
	"fmt"
	"github.com/akiraak/go-manga/db"
	"github.com/akiraak/go-manga/web"
	"github.com/akiraak/go-manga/elastic"
	. "github.com/akiraak/go-manga/model"
	"github.com/labstack/echo"
	"net/http"
	"sort"
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

func searchBooks(keyword string) ([]*TitleBook, int64, int) {
	if len(keyword) > 1 {
		searchBooks, hitTotal := elastic.SearchAsins(keyword)
		ids := make([]string, len(searchBooks))
		for _, book := range searchBooks {
			ids = append(ids, book.Asin)
		}

		books := []Book{}
		db.ORM.Where("asin IN (?)", ids).Find(&books)
		tbooks := titleGroupBooks(books)
		sortedBooks := []*TitleBook{}
		for _, tbook := range tbooks {
			sortedBooks = append(sortedBooks, tbook)
		}
		sort.Slice(sortedBooks, func(i, j int) bool {
			int1, _ := strconv.Atoi(sortedBooks[i].DatePublish())
			int2, _ := strconv.Atoi(sortedBooks[j].DatePublish())
			return int1 > int2
		})
		return sortedBooks, hitTotal, len(searchBooks)
	}
	return []*TitleBook{}, 0, 0
}

func GetSearchHandler(c echo.Context) error {
	keyword := c.QueryParam("key")
	type Param struct {
		BaseParam
		TitleBooks	[]*TitleBook
		HitTotal	int64
		AsinsCount	int
	}
	param := Param{BaseParam: BaseParam{"search", keyword}}
	param.TitleBooks, param.HitTotal, param.AsinsCount = searchBooks(keyword)

	return web.RenderTemplate(
		c,
		http.StatusOK,
		[]string{"template/search.html"},
		param)
}

func PageQuery(c echo.Context) int {
	page, err := strconv.Atoi(c.QueryParam("p"))
	if page == 0 || err != nil {
		page = 1
	}
	return page
}
